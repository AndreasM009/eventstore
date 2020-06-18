package http

import (
	"fmt"

	cors "github.com/AdhityaRamadhanus/fasthttpcors"
	"github.com/AndreasM009/eventstore-impl/store"
	registry "github.com/AndreasM009/eventstore/pkg/eventstored/eventstore"
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
	registry registry.Registry
}

// NewServer creates a new API Server
func NewServer(port int, eventStores map[string]store.EventStore, registry registry.Registry) Server {
	return &server{
		port:     port,
		evtstore: eventStores,
		api:      NewAPI(eventStores, registry),
		registry: registry,
	}
}

func (s *server) StartNonBlocking() {
	handler := s.useCors(
		s.useRouter())

	go func() {
		err := fasthttp.ListenAndServe(fmt.Sprintf(":%v", s.port), handler)
		if err != nil {
			fmt.Println(err)
		}
	}()
}

func (s *server) useRouter() fasthttp.RequestHandler {
	s.router = routing.New()
	s.api.RegisterRoutes(s.router)
	return s.router.HandleRequest
}

func (s *server) useCors(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	h := cors.NewCorsHandler(cors.Options{
		AllowedOrigins: []string{"*"},
		Debug:          false,
	})

	return h.CorsMiddleware(next)
}
