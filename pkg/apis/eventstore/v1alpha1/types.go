package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MetadataItem is a name/value pair for a metadata
type MetadataItem struct {
	Name         string       `json:"name"`
	Value        string       `json:"value"`
	SecretKeyRef SecretKeyRef `json:"secretKeyRef,omitempty"`
}

// SecretKeyRef is a reference to a secret holding the value for the metadata item.
// Name is the secret name, and key is the field in the secret.
type SecretKeyRef struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// EventstoreSpec defines the desired state of Eventstore
type EventstoreSpec struct {
	Type     string         `json:"type"`
	Metadata []MetadataItem `json:"metadata"`
}

// EventstoreStatus defines the observed state of Eventstore
type EventstoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// Eventstore is the Schema for the eventstores API
// +genclient
// +genclient:noStatus
// +resource:path=eventstore
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Eventstore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EventstoreSpec   `json:"spec,omitempty"`
	Status EventstoreStatus `json:"status,omitempty"`
}

// EventstoreList contains a list of Eventstore
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resourcepath=eventstore
type EventstoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Eventstore `json:"items"`
}
