// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	redis *redis.Client

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub(redis *redis.Client) *Hub {
	return &Hub{
		redis:      redis,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	// subscribe redis channel
	ctx := context.Background()
	ps := h.redis.Subscribe(ctx, "general_chat")
	_, err := ps.Receive(ctx)
	if err != nil {
		panic(err)
	}

	// Go channel which receives messages.
	ch := ps.Channel()

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-ch:
			for client := range h.clients {
				select {
				case client.send <- []byte(message.Payload):
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
