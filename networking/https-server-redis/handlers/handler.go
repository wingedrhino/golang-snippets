// Package handlers contains HTTP handlers required for app
package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wingedrhino/golang-snippets/networking/https-server-redis/internal/platform/stack"
)

// Handler is a http.Handler exposing business logic
type handler struct {
	st stack.Stack
}

// NewHandler returns a http.Handler
func NewHandler(s stack.Stack) http.Handler {
	return handler{st: s}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flush, flushable := w.(http.Flusher)
	timestamp := time.Now()
	reqDetails := fmt.Sprintf("Time: %v, Req IP Address: %s, Path: %s", timestamp, r.RemoteAddr, r.RequestURI)
	fmt.Println(reqDetails)

	err := h.st.Push(reqDetails)

	if err != nil {
		w.Write([]byte("Error writing to DB!"))
		return
	}

	resList, err := h.st.Read()

	if err != nil {
		w.Write([]byte("Error reading from DB!"))
		return
	}

	w.Write([]byte(fmt.Sprintf("Previous %d Requests:\n\n", len(resList))))
	// Attempt to write the response to the client immediately
	if flushable {
		flush.Flush()
	}

	for _, val := range resList {
		_, err = w.Write([]byte(val + "\n"))
		if err != nil {
			return
		}
		// Attempt to write the response to the client immediately
		if flushable {
			flush.Flush()
		}
		//  Sleep for a bit so that client can see page loading
		time.Sleep(300 * time.Millisecond)
	}

	return
}
