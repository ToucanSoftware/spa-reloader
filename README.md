[![Build Status](https://travis-ci.org/ToucanSoftware/spa-reloader.svg?branch=main)](https://travis-ci.org/ToucanSoftware/spa-reloader) [![Go Report Card](https://goreportcard.com/badge/github.com/ToucanSoftware/spa-reloader)](https://goreportcard.com/report/github.com/ToucanSoftware/spa-reloader) [![GoDoc](https://godoc.org/github.com/ToucanSoftware/spa-reloader?status.svg)](https://godoc.org/github.com/ToucanSoftware/spa-reloader)

# SPA Reloader

SPA Reloader provides a Kubernetes Controller that listens to changes in Deployments container image changes,
and informes Single Page Applications via WebSockets.

## Configuration

The controller takes the following environment variables in order to configure its behavior.

- SPA_NAMESPACE: is the name of the environment variable used to watch for changes in a namespace.
- SPA_NAME: is the name of the environment variable used to watch for changes in deployment name.
- SPA_RESYNC_SEC: is the number of seconds to resync.

## Change detection

Every time a change in a deploy that controller is listening to is detected, a message like this one is dispatch
to the clients connected to the web socket.

```json
{
  "created_at": "2020-11-10T11:08:47.073626-03:00",
  "namespace": "toucan",
  "name": "spa",
  "image": "nginx:1.14.1",
  "sha256": "32fdf92b4e986e109e4db0865758020cb0c3b70d6ba80d02fe87bad5cc3dc228"
}
```
