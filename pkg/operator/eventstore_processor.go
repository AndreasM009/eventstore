package operator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"k8s.io/apimachinery/pkg/labels"

	v1alpha1 "github.com/AndreasM009/eventstore/pkg/apis/eventstore/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type eventstoreProcessor struct {
	kubeClient kubernetes.Interface
}

func newEventStoreProcessor(kubeClient kubernetes.Interface) Processor {
	return &eventstoreProcessor{
		kubeClient: kubeClient,
	}
}

func (p *eventstoreProcessor) ProcessChanged(obj interface{}) error {
	log.Println("Eventstore Object changed")

	eventstore := obj.(*v1alpha1.Eventstore)

	services, err := p.getEventstoreServices()
	if err != nil {
		log.Println("can't get evenstore services")
		return nil
	}

	endpoints := p.getEndpoints(services)

	for _, e := range endpoints {
		p.updateSidecar(e, eventstore)
	}

	return nil
}

func (p *eventstoreProcessor) ProcessDeleted(obj interface{}) error {
	log.Println("Eventstore Object deleted")
	return nil
}

func (p *eventstoreProcessor) getEventstoreServices() (*corev1.ServiceList, error) {
	services, err := p.kubeClient.CoreV1().Services(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{eventstoreEnabledKey: "true"}).String(),
	})

	if err != nil {
		return nil, err
	}

	return services, nil
}

func (p *eventstoreProcessor) getEndpoints(services *corev1.ServiceList) []*corev1.Endpoints {
	result := []*corev1.Endpoints{}
	for _, v := range services.Items {
		endpoint, err := p.kubeClient.CoreV1().Endpoints(v.GetNamespace()).Get(context.TODO(), v.GetName(), metav1.GetOptions{})
		if err != nil {
			log.Printf("failed to get endpint for service %s: %s", v.GetName(), err)
		} else {
			result = append(result, endpoint)
		}
	}

	return result
}

func (p *eventstoreProcessor) updateSidecar(endpoint *corev1.Endpoints, settings *v1alpha1.Eventstore) {
	if endpoint == nil || len(endpoint.Subsets) <= 0 {
		return
	}

	payload, err := json.Marshal(settings)

	if err != nil {
		log.Printf("can't serialize Eventstore to json: %s\n", err)
		return
	}

	for _, a := range endpoint.Subsets[0].Addresses {
		address := fmt.Sprintf("%s:%d", a.IP, evenstoreDefaultPort)

		go func(address string) {
			url := fmt.Sprintf("http://%s/configurations/%s", address, settings.ObjectMeta.Name)
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				log.Printf("failed to send update to sidecar: %s\n", err)
				return
			}

			if resp.StatusCode != http.StatusOK {
				log.Printf("update sidecar config returned %d\n", resp.StatusCode)
				return
			}

			log.Println("sidecar updated with new config")
		}(address)
	}
}
