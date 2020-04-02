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
	evtstore store.EventStore
	router   *routing.Router
}

// NewServer creates a new API Server
func NewServer(port int, eventStore store.EventStore) Server {
	return &server{
		port:     port,
		evtstore: eventStore,
		api:      NewAPI(eventStore),
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
