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

package message

import (
	"time"
)

// ImageDescription description a Docker Image
type ImageDescription struct {
	// Image is the name of the image
	ImageName string `json:"name,omitempty"`
	// SHA-256 of the Docker image
	ImageSHA256 string `json:"sha256,omitempty"`
}

// ImageChangeMessage is the message for informing the clients that a images has changed
type ImageChangeMessage struct {
	CreatedAt time.Time `json:"created_at"`
	// Namespace namespace to listen to.
	Namespace string `json:"namespace,omitempty"`
	// Name name of the deployment to listen to.
	Name string `json:"name,omitempty"`
	// ImageSHA256 is the SHA-256 of the container Image ID
	CurrentImage *ImageDescription `json:"current_image,omitempty"`
	// ImageSHA256 is the SHA-256 of the container Image ID
	PreviousImage *ImageDescription `json:"previous_image,omitempty"`
}

// NewImageChangeMessage creates a new change image message
func NewImageChangeMessage(namespace string, name string,
	currentImage *ImageDescription,
	previousImage *ImageDescription) *ImageChangeMessage {
	now := time.Now()
	return &ImageChangeMessage{
		CreatedAt:     now,
		Namespace:     namespace,
		Name:          name,
		CurrentImage:  currentImage,
		PreviousImage: previousImage,
	}
}
