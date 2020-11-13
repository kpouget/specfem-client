package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"

)

// GitSpec defines the path to the git repository
type GitSpec struct {
	// +kubebuilder:validation:Pattern=[a-zA-Z0-9\.\-\/]+
	Uri string `json:"uri"`

	// +kubebuilder:validation:Pattern=[a-zA-Z0-9\-]+
	Ref string `json:"ref"`
}

// MpiSpec defines the execution conditions
type ExecSpec struct {
	Nproc int32 `json:"nproc"`
	Ncore int32 `json:"ncore"`
	SlotsPerWorker int32 `json:"slotsPerWorker"`
}

// MpiSpec defines the execution conditions
type SpecfemSpec struct {
	Nex int32 `json:"nex"`
	GpuPlatform string `json:"gpuPlatform",omitempty`
	GpuDevice string `json:"gpuDevice",omitempty`
}

type NetworkType string

const (
    NetworkTypeDefault     NetworkType = "Default"
    NetworkTypeMultus      NetworkType = "Multus"
    NetworkTypeHostNetwork NetworkType = "HostNetwork"
)

// ResourcesSpec defines the OpenShift resource usage
type ResourcesSpec struct {
	StorageClassName string `json:"storageClassName"`
	WorkerNodeSelector map[string]string `json:"workerNodeSelector"`
	UseUbiImage bool `json:"useUbiImage"`
	RelyOnSharedFS bool `json:"relayOnSharedFS"`
	NetworkType NetworkType `json:"networkType"`
	Requests map[string]resource.Quantity `json:"requests",omitempty`
}

// SpecfemAppSpec defines the desired state of SpecfemApp
type SpecfemAppSpec struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	Git GitSpec `json:"git"`
	Exec ExecSpec `json:"exec"`

	Specfem SpecfemSpec `json:"specfem"`

	Resources ResourcesSpec `json:"resources"`
}

// SpecfemAppStatus defines the observed state of SpecfemApp
type SpecfemAppStatus struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	State string `json:"state"`
	
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SpecfemApp is the Schema for the specfemapps API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=specfemapps,scope=Namespaced
type SpecfemApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpecfemAppSpec   `json:"spec,omitempty"`
	Status SpecfemAppStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SpecfemAppList contains a list of SpecfemApp
type SpecfemAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SpecfemApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SpecfemApp{}, &SpecfemAppList{})
}
