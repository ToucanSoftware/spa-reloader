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
	"context"
	"fmt"
	"github.com/ToucanSoftware/spa-reloader/pkg/message"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
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
	// Resync is the number of sec to execute resync
	Resync int
	// CurrentImageSHA256 stores the current SHA-256 of the container Image ID
	CurrentImageSHA256 string
	// informer created
	informer cache.SharedIndexInformer
	// Kubernetes Client Set
	client *kubernetes.Clientset
	// WebSocket server
	WSServer *ws.WebSockerServer
}

// NewSPAManager creates a new DeploymentController
func NewSPAManager(namespace string, name string, resync int, websocketPort int) (*DeploymentManager, error) {
	client, err := kubernetes.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		return nil, err
	}
	return &DeploymentManager{
		Namespace:          namespace,
		Name:               name,
		Resync:             resync,
		CurrentImageSHA256: "",
		client:             client,
		WSServer:           ws.NewWebSockerServer(websocketPort),
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
	var err error
	var selector fields.Selector = fields.Everything()
	if r.Name != "" {
		selector, err = fields.ParseSelector(fmt.Sprintf("metadata.name=%s", r.Name))
		if err != nil {
			return err
		}
	}
	logger.Info(fmt.Sprintf("Creating Watchlist for Deployment %s at Namespace %s", r.Name, r.Namespace))
	watchList := cache.NewListWatchFromClient(r.client.AppsV1().RESTClient(), "deployments", r.Namespace, selector)

	logger.Info(fmt.Sprintf("Creating Informer for Deployment %s at Namespace %s", r.Name, r.Namespace))

	// Shared informer example
	informer := cache.NewSharedIndexInformer(
		watchList,
		&appsv1.Deployment{},
		time.Second*time.Duration(r.Resync),
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

	var oldImage = oldDeploy.Spec.Template.Spec.Containers[0].Image
	var currImage = currentDeploy.Spec.Template.Spec.Containers[0].Image

	if currImage == "" {
		// wait for next resync cycle
		return
	}
	pods, err := r.getPodsForDeploy(currentDeploy)
	if err != nil {
		logger.Error(fmt.Sprintf("error finding pods por deploy %s: %v\n", currentDeploy.GetName(), err))
	} else {
		for _, pod := range pods.Items {
			if len(pod.Status.ContainerStatuses) > 0 {
				var image = pod.Status.ContainerStatuses[0].Image
				// Work on the the current immage
				if currImage == image {
					var imageID = pod.Status.ContainerStatuses[0].ImageID
					// Handle the case when pod is pending
					if imageID != "" {
						var imageSHA256 = sha256FromImageID(imageID)
						if imageSHA256 != "" && imageSHA256 != r.CurrentImageSHA256 {
							logger.Info(fmt.Sprintf("Detected Pod Image ID Change from %s to %s", oldImage, currImage))
							changeImageMessage := message.NewImageChangeMessage(r.Namespace, r.Name, currImage, imageSHA256, r.CurrentImageSHA256)
							r.CurrentImageSHA256 = imageSHA256
							err = r.WSServer.BroadcastMessage(changeImageMessage)
							if err != nil {
								logger.Error(fmt.Sprintf("error sending broadcast message: %v\n", err))
							}
						}
					}
				}
			}
		}
	}
}

func (r *DeploymentManager) getPodsForDeploy(deploy *appsv1.Deployment) (*corev1.PodList, error) {
	set := labels.Set(deploy.Spec.Selector.MatchLabels)
	listOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	pods, err := r.client.CoreV1().Pods(deploy.Namespace).List(context.TODO(), listOptions)
	return pods, err
}

func sha256FromImageID(imageID string) string {
	sep := strings.Split(imageID, "@sha256:")
	if len(sep) > 1 {
		return sep[1]
	}
	return ""
}
