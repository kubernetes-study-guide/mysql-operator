package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EnvironmentVariables mirrors - https://hub.docker.com/_/mysql#Environment_Variables (sher-start)
type EnvironmentVariables struct {
	MysqlRootPassword string `json:"mysql_root_password"`
	MysqlDatabase     string `json:"mysql_database"`
	MysqlUser         string `json:"mysql_user"`
	MysqlPassword     string `json:"mysql_password"`
}

// VolumeSettings - info needed to create PVC.
type VolumeSettings struct {
	VolumeSize   string `json:"volume_size"`
	StorageClass string `json:"storage_class"`
}

// sher-end

// MySQLSpec defines the desired state of MySQL
type MySQLSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// sher-start
	Environment EnvironmentVariables `json:"environment"`
	// here's a more sophisticated example:
	// https://github.com/Sher-Chowdhury/prometheus-jmx-exporter-operator/blob/master/pkg/apis/banzaicloud/v1alpha1/types.go#L18-L32

	Volume VolumeSettings `json:"volume"`
	// sher-end
}

// MySQLStatus defines the observed state of MySQL
type MySQLStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MySQL is the Schema for the mysqls API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=mysqls,scope=Namespaced
type MySQL struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MySQLSpec   `json:"spec,omitempty"`
	Status MySQLStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MySQLList contains a list of MySQL
type MySQLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MySQL `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MySQL{}, &MySQLList{})
}
