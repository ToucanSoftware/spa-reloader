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

package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "ws")))
)

// WebSockerServer websocket server
type WebSockerServer struct {
	BindAddress string
	upgrader    websocket.Upgrader
	hub         *Hub
}

// NewWebSockerServer creates a new web socket server
func NewWebSockerServer() *WebSockerServer {
	return &WebSockerServer{
		BindAddress: "0.0.0.0:8080",
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		hub: newHub(),
	}
}

// Run runs the websocket server
func (s *WebSockerServer) Run() error {
	logger.Info(fmt.Sprintf("Starting WebSocket server at %s", s.BindAddress))
	go s.hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error(fmt.Sprintf("Error creating upgrader: %v", err))
			return
		}
		serveWs(s.hub, conn)
	})
	return http.ListenAndServe(s.BindAddress, nil)
}
