package main

import (
	"context"
	"errors"
	"os"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/defaults"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/namespaces"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	Version = "0.0.1"
)

func doCommand(c *cli.Context) (err error) {
	if c.NArg() != 2 {
		err = errors.New("Must have 2 arguments")
		return
	}
	client, err := containerd.New(c.GlobalString("address"), containerd.WithDefaultNamespace(c.GlobalString("namespace")))
	if err != nil {
		return
	}
	defer func() {
		err := client.Close()
		if err != nil {
			return
		}
	}()

	imageService := client.ImageService()
	if err != nil {
		return
	}

	originalTag := c.Args()[0]
	newTag := c.Args()[1]

	var image images.Image
	image, err = imageService.Get(context.Background(), originalTag)
	if errdefs.IsNotFound(err) {
		log.Debugf("%s not found, attempting to pull", originalTag)
		_, err = client.Pull(context.Background(), originalTag)
		if err != nil {
			return
		}
		image, err = imageService.Get(context.Background(), originalTag)

	}
	if err != nil {
		return
	}
	image.Name = newTag
	_, err = imageService.Create(context.Background(), image)
	if err != nil {
		return
	}
	log.Infof("Created %s", newTag)
	return
}

func main() {
	var namespace string
	app := cli.NewApp()
	app.Name = "cretag"
	app.Version = Version

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in logs",
		},
		cli.StringFlag{
			Name:  "address, a",
			Usage: "address for containerd's GRPC server",
			Value: defaults.DefaultAddress,
		},
		cli.StringFlag{
			Name:        "namespace, n",
			Usage:       "namespace to use with commands",
			Value:       namespaces.Default,
			EnvVar:      namespaces.NamespaceEnvVar,
			Destination: &namespace,
		},
	}
	app.Action = doCommand
	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}
