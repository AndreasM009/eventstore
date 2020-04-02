package injector

import (
	"log"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

const (
	eventstoreEnabledKey = "eventstore/enabled"
	eventstorePortKey    = "eventstore/port"
	eventstoreNames      = "eventstore/names"
	evenstoreDefaultPort = 5600
)

func (i *injector) isEventstoreEnabled(pod *corev1.Pod) bool {
	val, ok := pod.Annotations[eventstoreEnabledKey]

	if !ok {
		return false
	}

	str := strings.ToLower(val)

	switch str {
	case "y", "yes", "true", "1", "on":
		return true
	default:
		return false
	}
}

func (i *injector) getEvenstorePort(pod *corev1.Pod) int {
	val, ok := pod.Annotations[eventstorePortKey]
	if !ok {
		return evenstoreDefaultPort
	}

	port, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("value of port annotation '%s' can't be converted to integer, using %d default port", val, evenstoreDefaultPort)
		return evenstoreDefaultPort
	}

	return port
}

func (i *injector) getEvenstoreNames(pod *corev1.Pod) string {
	val, ok := pod.Annotations[eventstoreNames]
	if !ok {
		log.Println("no evenstore names specified in pod annotation")
		return ""
	}

	return val
}
