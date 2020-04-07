package injector

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

const (
	argMode             = "-mode"
	modeKubernetes      = "kubernetes"
	argPort             = "-port"
	argEventstores      = "-eventstores"
	argOperatorEndpoint = "-operatorendpoint"
	sidecarName         = "eventstored"
	sidecarImage        = "m009/eventstored:latest"
	httpPortName        = "http"
	operatorService     = "eventstore-operator"
	operatorServicePort = int(80)
)

func (i *injector) patchPod(pod *corev1.Pod, controlPlaneNamespace string) []PatchOperation {
	patchOperations := []PatchOperation{}

	if enabled := i.isEventstoreEnabled(pod); !enabled {
		return patchOperations
	}

	names := i.getEvenstoreNames(pod)
	if names == "" {
		return patchOperations
	}

	port := i.getEvenstorePort(pod)

	sidecar := createSidecarContainer(port, names, controlPlaneNamespace)

	if len(pod.Spec.Containers) == 0 {
		return append(patchOperations, PatchOperation{
			Path:  "spec/containers",
			Value: []corev1.Container{sidecar},
			Op:    "add",
		})
	}

	return append(patchOperations, PatchOperation{
		Path:  "/spec/containers/-",
		Value: sidecar,
		Op:    "add",
	})
}

func createSidecarContainer(port int, evtsNames, controlPlaneNamespace string) corev1.Container {
	cntr := corev1.Container{
		Name:            sidecarName,
		Image:           sidecarImage,
		ImagePullPolicy: corev1.PullAlways,
		Ports: []corev1.ContainerPort{
			{
				Name:          httpPortName,
				ContainerPort: int32(port),
			},
		},
		Command: []string{"./eventstored"},
		Args: []string{
			fmt.Sprintf("%s=%s", argMode, modeKubernetes),
			fmt.Sprintf("%s=%d", argPort, port),
			fmt.Sprintf("%s='%s'", argEventstores, evtsNames),
			fmt.Sprintf("%s=http://%s.%s.svc.cluster.local:%d", argOperatorEndpoint, operatorService, controlPlaneNamespace, operatorServicePort),
		},
	}

	return cntr
}
