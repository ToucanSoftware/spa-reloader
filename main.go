/*
Copyright © 2020 ToucanSoftware

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"strconv"

	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/ToucanSoftware/spa-reloader/controllers"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "main")))
)

const (
	// spaNamepapce is the name of the environment variable used to watch for changes in a namespace.
	spaNamepapce string = "SPA_NAMESPACE"
	// spaName is the name of the environment variable used to watch for changes in deployment name.
	spaName string = "SPA_NAME"
	// spaResyncSeconds is the number of seconds to resync
	spaResyncSeconds string = "SPA_RESYNC_SEC"
	// spaWebSocketPort port to bind the websocket server
	spaWebSocketPort string = "SPA_WEBSOCKET_PORT"
)

const (
	// defaultNamespace by default, only listen to changes in `default` namespace.
	defaultNamespace string = "default"
	// defaultName by default, only listen to all deployments.
	defaultName string = ""
	//defaultResyncSeconds by default we resync every 30 sec.
	defaultResyncSeconds string = "30"
	// defaultResyncSecondsValue by default we resync every 30 sec.
	defaultResyncSecondsValue int = 10
	// defaultWebSocketPort by default we listent in port 8080
	defaultWebSocketPort string = "8080"
	// defaultWebSocketPortValue by default we listent in port 8080
	defaultWebSocketPortValue int = 8080
)

func main() {
	namespace := getenv(spaNamepapce, defaultNamespace)
	name := getenv(spaName, defaultName)

	resync := getenv(spaResyncSeconds, defaultResyncSeconds)
	resyncVal, err := strconv.Atoi(resync)
	if err != nil {
		resyncVal = defaultResyncSecondsValue
	}

	websocketPort := getenv(spaWebSocketPort, defaultWebSocketPort)
	websocketPortVal, err := strconv.Atoi(websocketPort)
	if err != nil {
		websocketPortVal = defaultWebSocketPortValue
	}

	mgr, err := controllers.NewSPAManager(namespace, name, resyncVal, websocketPortVal)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to start manager: %v", err))
		os.Exit(1)
	}

	logger.Info("Starting SPA Reloader")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		logger.Error(fmt.Sprintf("Error running manager: %v", err))
		os.Exit(1)
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
