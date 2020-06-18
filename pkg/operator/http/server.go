package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	eventstoreclient "github.com/AndreasM009/eventstore/pkg/client/clientset/versioned"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Server http server interface
type Server interface {
	StartNonBlocking()
}

type server struct {
	port             int
	eventstoreClient *eventstoreclient.Clientset
	router           *routing.Router
}

// NewServer creates a new server
func NewServer(port int, eventstoreClient *eventstoreclient.Clientset) Server {
	s := &server{
		port:             port,
		eventstoreClient: eventstoreClient,
		router:           routing.New(),
	}

	s.router.Get("/eventstores", s.onGetComponents)
	return s
}

func (s *server) StartNonBlocking() {
	go func() {
		err := fasthttp.ListenAndServe(fmt.Sprintf(":%v", s.port), s.router.HandleRequest)
		if err != nil {
			log.Println(err)
		}
	}()
}

func (s *server) onGetComponents(c *routing.Context) error {
	stores, err := s.eventstoreClient.EventstoreV1alpha1().
		Eventstores(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		msg := NewErrorResponse("ERR_GETTING_EVENTSTORES", fmt.Sprintf("can't get EventStores %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
		return nil
	}

	data, err := json.Marshal(stores.Items)
	if err != nil {
		msg := NewErrorResponse("ERR_SERIALIZE_EVENTSTORES", fmt.Sprintf("can't serialize EventStores %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
		return nil
	}

	respondWithJSON(c.RequestCtx, fasthttp.StatusOK, data)
	return nil
}
