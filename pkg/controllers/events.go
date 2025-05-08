package controllers

import (
	"context"
	"fmt"
	"net/http"
)

func ApiEvents(writer http.ResponseWriter, request *http.Request) {
	// Set headers for SSE
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")

	// Create a context for handling client disconnection
	_, cancel := context.WithCancel(request.Context())
	defer cancel()

	flusher, ok := writer.(http.Flusher)
	if !ok {
		http.Error(writer, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// subscribe to msgbus
	updateCh := updateBus.Subscribe()

	// Send data to the client when we receive an event
	for {
		select {
		case data := <-updateCh:
			fmt.Fprintf(writer, "data: %s\n\n", data)
			flusher.Flush()
		}
	}

}
