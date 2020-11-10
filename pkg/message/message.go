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

// ImageChangeMessage is the message for informing the clients that a images has changed
type ImageChangeMessage struct {
	CreatedAt time.Time `json:"created_at"`
	// Namespace namespace to listen to.
	Namespace string `json:"namespace,omitempty"`
	// Name name of the deployment to listen to.
	Name string `json:"name,omitempty"`
	// Image is the name of the image
	Image string `json:"image,omitempty"`
	// ImageSHA256 is the SHA-256 of the container Image ID
	ImageSHA256 string `json:"sha256,omitempty"`
}

// NewImageChangeMessage creates a new change image message
func NewImageChangeMessage(namespace string, name string, image string, imageSHA256 string) *ImageChangeMessage {
	now := time.Now()
	return &ImageChangeMessage{
		CreatedAt:   now,
		Namespace:   namespace,
		Name:        name,
		Image:       image,
		ImageSHA256: imageSHA256,
	}
}
