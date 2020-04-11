package operator

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"k8s.io/client-go/kubernetes"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type deploymentProcessor struct {
	kubeClient kubernetes.Interface
}

func newDeploymentProcessor(kubeClient kubernetes.Interface) Processor {
	return &deploymentProcessor{
		kubeClient: kubeClient,
	}
}

func (p *deploymentProcessor) ProcessChanged(obj interface{}) error {
	deployment := obj.(*appsv1.Deployment)

	if enabled := p.isEventstoreEnabled(deployment); !enabled {
		return nil
	}

	appID := p.getAppID(deployment)
	port := p.getSidecarPort(deployment)
	servicename := fmt.Sprintf("%s-eventstore", appID)

	if appID == "" {
		log.Printf("Skipping creation of service for deployment %s, appid is empty\n", deployment.GetName())
		return nil
	}

	if p.serviceExists(servicename, deployment.GetNamespace()) {
		log.Printf("Service %s already exists\n", servicename)
		return nil
	}

	// service definition
	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   servicename,
			Labels: map[string]string{eventstoreEnabledKey: "true"},
		},
		Spec: corev1.ServiceSpec{
			Selector: deployment.Spec.Selector.MatchLabels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromInt(port),
					Name:       httpPortName,
				},
			},
		},
	}

	// create the service
	_, err := p.kubeClient.CoreV1().Services(deployment.GetNamespace()).Create(context.TODO(), &service, metav1.CreateOptions{})
	if err != nil {
		log.Printf("creating service for deployment %s failed: %s", deployment.GetName(), err)
	}

	log.Printf("Service '%s' for deployment '%s' in namespace '%s' created\n", servicename, deployment.GetName(), deployment.GetNamespace())

	return nil
}

func (p *deploymentProcessor) ProcessDeleted(obj interface{}) error {
	deployment := obj.(*appsv1.Deployment)

	if enabled := p.isEventstoreEnabled(deployment); !enabled {
		return nil
	}

	appID := p.getAppID(deployment)
	servicename := fmt.Sprintf("%s-eventstore", appID)

	if !p.serviceExists(servicename, deployment.GetNamespace()) {
		return nil
	}

	err := p.kubeClient.CoreV1().Services(deployment.GetNamespace()).Delete(context.TODO(), servicename, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("failed deleteing service for deployment %s: %s\n", deployment.GetName(), servicename)
	}

	log.Printf("Service '%s' for deployment '%s' in namespace '%s' deleted", servicename, deployment.GetName(), deployment.GetNamespace())

	return nil
}

func (p *deploymentProcessor) isEventstoreEnabled(deployment *appsv1.Deployment) bool {
	annotations := deployment.Spec.Template.Annotations

	enabled, ok := annotations[eventstoreEnabledKey]

	if !ok {
		return false
	}

	switch strings.ToLower(enabled) {
	case "y", "yes", "true", "on", "1":
		return true
	default:
		return false
	}
}

func (p *deploymentProcessor) getAppID(deployment *appsv1.Deployment) string {
	annotations := deployment.Spec.Template.Annotations

	if id, ok := annotations[eventstoreAppID]; ok {
		return id
	}

	return ""
}

func (p *deploymentProcessor) getSidecarPort(deployment *appsv1.Deployment) int {
	annotations := deployment.Spec.Template.Annotations

	if port, ok := annotations[eventstorePortKey]; ok {
		portnumber, err := strconv.Atoi(port)
		if err != nil {
			return evenstoreDefaultPort
		}

		return portnumber
	}

	return evenstoreDefaultPort
}

func (p *deploymentProcessor) serviceExists(name, namespace string) bool {
	_, err := p.kubeClient.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	return err == nil
}
