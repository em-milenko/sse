//main.go
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type DashBoard struct {
	Event string `json:"event"`
	Data  string `json:"data"`
	Id    int    `json:"id"`
}

type DashBoardHandler struct{}

func main() {
	dashBoardHandler := NewHandler()
	server := &http.Server{Addr: ":3000", Handler: dashBoardHandler}
	if err := server.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}

func NewHandler() *DashBoardHandler {
	return &DashBoardHandler{}
}

var channelMap = make(map[string]chan *DashBoard, 10)
var lock = sync.RWMutex{}

func (dashBoardHandler *DashBoardHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	id := request.Form.Get("id")
	lock.Lock()
	clientChan, ok := channelMap[id]
	if !ok {
		clientChan = make(chan *DashBoard, 1)
		channelMap[id] = clientChan
		go updateDashboard(id)
		go cleanDashboard(id)
	}
	lock.Unlock()

	responseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	responseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	responseWriter.Header().Set("Content-Type", "text/event-stream")
	responseWriter.Header().Set("Cache-Control", "no-cache")
	responseWriter.Header().Set("Connection", "keep-alive")

	done := request.Context().Done()

	select {
	case ev := <-clientChan:
		response := fmt.Sprintf("event:%v\ndata:%v\n\n", ev.Event, ev.Data)
		_, err = fmt.Fprint(responseWriter, response)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Print(response)
		close(clientChan)
		delete(channelMap, id)
	case <-done:
		_, err := fmt.Fprintf(responseWriter, ": nothing to sent\n\n")
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
	if f, ok := responseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func updateDashboard(id string) {
	time.Sleep(5 * time.Second)
	lock.Lock()
	defer lock.Unlock()
	events, ok := channelMap[id]
	if ok {
		events <- &DashBoard{
			Event: "update",
			Data:  "some payload: " + id,
			Id:    rand.Int(),
		}
	}
}

func cleanDashboard(id string) {
	time.Sleep(10 * time.Second)
	lock.Lock()
	defer lock.Unlock()
	events, ok := channelMap[id]
	if ok {
		close(events)
		delete(channelMap, id)
	}
}
