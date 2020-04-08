package http

import (
	"fmt"

	"github.com/AndreasM009/eventstore-go/store"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// Server interface for API server
type Server interface {
	StartNonBlocking()
}

type server struct {
	api      APIRoutes
	port     int
	evtstore map[string]store.EventStore
	router   *routing.Router
}

// NewServer creates a new API Server
func NewServer(port int, eventStores map[string]store.EventStore) Server {
	return &server{
		port:     port,
		evtstore: eventStores,
		api:      NewAPI(eventStores),
	}
}

func (s *server) StartNonBlocking() {
	s.router = routing.New()
	s.api.RegisterRoutes(s.router)

	go func() {
		err := fasthttp.ListenAndServe(fmt.Sprintf(":%v", s.port), s.router.HandleRequest)
		if err != nil {
			fmt.Println(err)
		}
	}()
}
