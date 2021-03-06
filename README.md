![Publish Docker image](https://github.com/ToucanSoftware/spa-reloader/workflows/Publish%20Docker%20image/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/ToucanSoftware/spa-reloader)](https://goreportcard.com/report/github.com/ToucanSoftware/spa-reloader) [![GoDoc](https://godoc.org/github.com/ToucanSoftware/spa-reloader?status.svg)](https://godoc.org/github.com/ToucanSoftware/spa-reloader)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FToucanSoftware%2Fspa-reloader.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FToucanSoftware%2Fspa-reloader?ref=badge_shield)

# SPA Reloader

SPA Reloader provides a Kubernetes Controller that listens to changes in Deployments container image changes,
and informes Single Page Applications via WebSockets.

**Problem**: Single Page Application (SPA) are great technology, they run on the client browser and do not need round trip to the server every time we move across the application. That means we have to load the entire (or part of the) application in the browser, which will tipically communicate with a REST API in the server. But if we upgrade the SPA, unless we reload the page, the browser has not standard way of knowing if the application has changed and needs to get reloaded.

**Solution**: The solution proposed includes a Server Side component and a Client Side component.

- Server Side: Use a Kubernetes Controller that listen to changes in the Deployment configuration specified by the user and send a WebSocket message to clients.
- Client Side: Use a websocket client that subscribes to the server and stores the current image version. If the client recieves a message with a different image version it would fire an event to inform that a new version of the applition it's running. Click [here](https://github.com/ToucanSoftware/spa-reloader-vue) to get to the Client Side Repo.

You can take a look at the Demo Project [here](https://github.com/ToucanSoftware/spa-reloader-demo).

## Configuration

The controller takes the following environment variables in order to configure its behavior.

- **SPA_NAMESPACE**: is the name of the environment variable used to watch for changes in a namespace.
- **SPA_NAME**: is the name of the environment variable used to watch for changes in deployment name.
- **SPA_RESYNC_SEC**: is the number of seconds to resync.
- **SPA_WEBSOCKET_PORT**: is the port number of the websocket. Change this value to bind the websocket server to another port. Be aware that the current Dockerfile exposes port `8080` for the websocket server.

## Change detection

Every time a change is detected in the deployment that the controller is listening to, a message is dispatched to the clients connected to the web socket.

Here is an example of a message for the Deployment `spa` in Namespace `toucan`.

```json
{
  "created_at": "2020-11-10T11:08:47.073626-03:00",
  "namespace": "toucan",
  "name": "spa",
  "image": "nginx:1.14.1",
  "sha256": "32fdf92b4e986e109e4db0865758020cb0c3b70d6ba80d02fe87bad5cc3dc228"
}
```


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FToucanSoftware%2Fspa-reloader.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FToucanSoftware%2Fspa-reloader?ref=badge_large)