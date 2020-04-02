package injector

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
)

// Injector interface that starts a web server listening on MutatingAdmissionWebHooks
// of Kubernetes (MutatingAdmissionController)
type Injector interface {
	Init(cfg Config, kubeClient *kubernetes.Clientset) error
	Run(ctx context.Context) error
}

type injector struct {
	server       *http.Server
	config       Config
	deserializer runtime.Decoder
	kubeClient   *kubernetes.Clientset
}

// NewInjector creates a new Injector
func NewInjector() Injector {
	mux := http.NewServeMux()
	i := &injector{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", 8443),
			Handler: mux,
		},
		deserializer: serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer(),
	}

	mux.HandleFunc("/mutate", i.handleRequest)
	return i
}

func (i *injector) Init(cfg Config, kubeClient *kubernetes.Clientset) error {

	i.config = cfg
	i.kubeClient = kubeClient

	return nil
}

func (i *injector) Run(ctx context.Context) error {
	doneChannel := make(chan struct{})

	go func() {
		select {
		case <-ctx.Done():
			log.Println("injector: shutting down")
			shwdctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			i.server.Shutdown(shwdctx) // nolint: errcheck
		case <-doneChannel:
		}
	}()

	log.Printf("injector: starting server on port %v\n", i.server.Addr)

	if err := i.server.ListenAndServeTLS(i.config.TLSCertFile, i.config.TLSKeyFile); err != http.ErrServerClosed {
		log.Printf("injector: can't start server: %s\n", err)
		close(doneChannel)
		return err
	}

	close(doneChannel)
	return nil
}

func (i *injector) handleRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data []byte
	var admissionResponse *v1beta1.AdmissionResponse
	var patchOps []PatchOperation
	var err error
	admissionReview := v1beta1.AdmissionReview{}

	// read and check body
	if r.Body != nil {
		if d, err := ioutil.ReadAll(r.Body); err == nil {
			data = d
		}
	}

	if len(data) == 0 {
		log.Println("injector: empty request body received")
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	// check Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Printf("injector: request Content-Type=%s, expect application/json\n", contentType)
		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	// deserialize review
	if _, _, err = i.deserializer.Decode(data, nil, &admissionReview); err != nil {
		log.Printf("injector: Can't decode body: %v\n", err)
		respondWithError(nil, w, err)
		return
	}

	if admissionReview.Request.Kind.Kind != "Pod" {
		log.Printf("injector: invalid kind for review: %s", admissionReview.Kind)
		respondWithError(&admissionReview, w, fmt.Errorf("invalid kind of review: %s", admissionReview.Kind))
		return
	}

	pod, err := deserializePod(admissionReview.Request)
	if err != nil {
		respondWithError(&admissionReview, w, err)
		return
	}

	patchOps = i.patchPod(pod)
	if len(patchOps) == 0 {
		admissionResponse = &v1beta1.AdmissionResponse{
			Allowed: true,
		}
		admissionResponse.Result = &metav1.Status{
			Status: "Success",
		}
	} else {
		patchType := v1beta1.PatchTypeJSONPatch
		jsonPatch, err := json.Marshal(patchOps)

		if err != nil {
			respondWithError(&admissionReview, w, err)
			return
		}

		log.Printf("Containers patched: %v", string(jsonPatch))
		admissionResponse = &v1beta1.AdmissionResponse{
			Allowed:   true,
			PatchType: &patchType,
			Patch:     jsonPatch,
		}
		admissionResponse.Result = &metav1.Status{
			Status: "Success",
		}
	}

	arResponse := v1beta1.AdmissionReview{}
	arResponse.Response = admissionResponse
	if admissionReview.Request != nil {
		arResponse.Response.UID = admissionReview.Request.UID
	}

	writeAdmissionResponse(&arResponse, w, http.StatusOK)
}

func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

func respondWithError(ar *v1beta1.AdmissionReview, w http.ResponseWriter, err error) {
	response := toAdmissionResponse(err)
	review := v1beta1.AdmissionReview{}
	review.Response = response
	if nil != ar && ar.Request != nil {
		review.Response.UID = ar.Request.UID
	}

	writeAdmissionResponse(&review, w, http.StatusInternalServerError)
}

func writeAdmissionResponse(ar *v1beta1.AdmissionReview, w http.ResponseWriter, statusCode int) {
	json, err := json.Marshal(ar)
	if err != nil {
		log.Printf("injector: can't serialize response AdmissionReview: %s", err)
		http.Error(w, fmt.Sprintf("can't serialize response: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(json)
	if err != nil {
		log.Printf("injector: cant write AdmissionReview response: %s", err)
		http.Error(w, fmt.Sprintf("can't write AdmissionReview response %s", err), http.StatusInternalServerError)
	}
}

func deserializePod(req *v1beta1.AdmissionRequest) (*corev1.Pod, error) {
	pod := &corev1.Pod{}

	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		log.Printf("can't deserialize pod from json: %v", err)
		return nil, err
	}

	return pod, nil
}
