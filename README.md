# cretag

Retag containerd images

A standalone implementation of https://github.com/containerd/containerd/pull/3388


## installation

```
go get -u github.com/seemethere/cretag
```

## usage

```
cretag docker.io/library/alpine:latest docker.io/library/alpine:my-new-tag
```
