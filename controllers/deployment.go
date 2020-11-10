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

package controllers

import (
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"

	"go.uber.org/zap"

	"github.com/ToucanSoftware/spa-reloader/pkg/ws"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "controllers.deployment")))
)

// DeploymentManager is a controller for Deployments
type DeploymentManager struct {
	// Namespace namespace to listen to.
	Namespace string
	// Name name of the deployment to listen to.
	Name string
	// informer created
	informer cache.SharedIndexInformer
	// WebSocket server
	WSServer *ws.WebSockerServer
}

// NewSPAManager creates a new DeploymentController
func NewSPAManager(namespace string, name string) (*DeploymentManager, error) {
	return &DeploymentManager{
		Namespace: namespace,
		Name:      name,
		WSServer:  ws.NewWebSockerServer(),
	}, nil
}

// Start starts listesting to changes in the deployments
func (r *DeploymentManager) Start(done <-chan struct{}) error {
	logger.Info("Starting WebSocket Informer")
	go func() {
		if err := r.WSServer.Run(); err != nil {
			logger.Fatal(fmt.Sprintf("Error Starting WebSocket Informer: %v", err))
		}
	}()
	logger.Info("WebSocket Informer Started")
	logger.Info("Starting Deployment Controller")
	go func() {
		if err := r.listenAndServe(); err != nil {
			logger.Fatal(fmt.Sprintf("listen: %v", err))
		}
	}()
	logger.Info("Deployment Controller Started")

	<-done

	return nil
}

// listenAndServe implements the logic of the controller
func (r *DeploymentManager) listenAndServe() error {
	clientset, err := getClient()
	if err != nil {
		return err
	}

	var selector fields.Selector = fields.Everything()
	if r.Name != "" {
		selector, err = fields.ParseSelector(fmt.Sprintf("metadata.name=%s", r.Name))
		if err != nil {
			return err
		}
	}

	logger.Info(fmt.Sprintf("Creating Watchlist for Deployment %s at Namespace %s", r.Name, r.Namespace))
	watchList := cache.NewListWatchFromClient(clientset.AppsV1().RESTClient(), "deployments", r.Namespace, selector)

	logger.Info(fmt.Sprintf("Creating Informer for Deployment %s at Namespace %s", r.Name, r.Namespace))

	// Shared informer example
	informer := cache.NewSharedIndexInformer(
		watchList,
		&appsv1.Deployment{},
		time.Second*10,
		cache.Indexers{},
	)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.handleDeploymentAdd,
		UpdateFunc: r.handleDeploymentUpdate,
	})

	r.informer = informer

	stop := make(chan struct{})
	go informer.Run(stop)

	return nil
}

func (r *DeploymentManager) handleDeploymentAdd(obj interface{}) {
	deploy := obj.(*appsv1.Deployment)
	logger.Info(fmt.Sprintf("Deployment [%s] is added", deploy.Name))
}

func (r *DeploymentManager) handleDeploymentUpdate(old, current interface{}) {
	oldDeploy := old.(*appsv1.Deployment)
	currentDeploy := current.(*appsv1.Deployment)

	if oldDeploy.Spec.Template.Spec.Containers[0].Image != currentDeploy.Spec.Template.Spec.Containers[0].Image {
		logger.Info(fmt.Sprintf("Image changed from %s to %s", oldDeploy.Spec.Template.Spec.Containers[0].Image, currentDeploy.Spec.Template.Spec.Containers[0].Image))
		// Find Pods for the deployment and get status.ImageID for
	}

	logger.Info(fmt.Sprintf("Deployment [%s] is update", oldDeploy.Name))
}

func getClient() (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(ctrl.GetConfigOrDie())
}
