package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,categories={fleet-networking},shortName=tme
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=`.spec.profile.name`,name="Profile",type=string
// +kubebuilder:printcolumn:JSONPath=`.spec.profile.namespace`,name="Profile-Namespace",type=string
// +kubebuilder:printcolumn:JSONPath=`.spec.endpointRef.name`,name="Backend",type=string
// +kubebuilder:printcolumn:JSONPath=`.status.conditions[?(@.type=='Accepted')].status`,name="Is-Accepted",type=string
// +kubebuilder:printcolumn:JSONPath=`.metadata.creationTimestamp`,name="Age",type=date

// TrafficManagerBackend is used to manage the Azure Traffic Manager Endpoints using cloud native way.
// A backend contains one or more endpoints. Therefore, the controller may create multiple endpoints under the Traffic
// Manager Profile.
// https://learn.microsoft.com/en-us/azure/traffic-manager/traffic-manager-endpoint-types
type TrafficManagerBackend struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// The desired state of TrafficManagerBackend.
	Spec TrafficManagerBackendSpec `json:"spec"`

	// The observed status of TrafficManagerBackend.
	// +optional
	Status TrafficManagerBackendStatus `json:"status,omitempty"`
}

type TrafficManagerBackendSpec struct {
	// +required
	// immutable
	// Which TrafficManagerProfile the backend should be attached to.
	Profile TrafficManagerProfileRef `json:"profile"`

	// The reference to a backend.
	// immutable
	// +required
	Backend TrafficManagerBackendRef `json:"backend"`

	// The total weight of endpoints behind the serviceImport when using the 'Weighted' traffic routing method.
	// Possible values are from 1 to 1000.
	// By default, the routing method is 'Weighted', so that it is required for now.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000
	// For example, if there are two clusters exporting the service via public ip, each public ip will be configured
	// as "Weight"/2.
	Weight *int64 `json:"weight,omitempty"`
}

// TrafficManagerProfileRef is a reference to a trafficManagerProfile object in the same namespace as the TrafficManagerBackend object.
type TrafficManagerProfileRef struct {
	// Name is the name of the referenced trafficManagerProfile.
	Name string `json:"name"`
}

// TrafficManagerBackendRef is the reference to a backend.
// Currently, we only support one backend type: ServiceImport.
type TrafficManagerBackendRef struct {
	// Name is the reference to the ServiceImport in the same namespace as the TrafficManagerBackend object.
	// +required
	Name string `json:"name"`
}

// TrafficManagerEndpointStatus is the status of Azure Traffic Manager endpoint which is successfully accepted under the traffic
// manager Profile.
type TrafficManagerEndpointStatus struct {
	// Name of the endpoint.
	// +required
	Name string `json:"name"`

	// The weight of this endpoint when using the 'Weighted' traffic routing method.
	// Possible values are from 1 to 1000.
	// +optional
	Weight *int64 `json:"weight,omitempty"`

	// The fully-qualified DNS name or IP address of the endpoint.
	// +optional
	Target *string `json:"target,omitempty"`

	// Cluster is where the endpoint is exported from.
	// +optional
	Cluster *ClusterStatus `json:"cluster,omitempty"`
}

type TrafficManagerBackendStatus struct {
	// Endpoints contains a list of accepted Azure endpoints which are created or updated under the traffic manager Profile.
	// +optional
	Endpoints []TrafficManagerEndpointStatus `json:"endpoints,omitempty"`

	// Current backend status.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// TrafficManagerBackendConditionType is a type of condition associated with a TrafficManagerBackendStatus. This type
// should be used within the TrafficManagerBackendStatus.Conditions field.
type TrafficManagerBackendConditionType string

// TrafficManagerBackendConditionReason defines the set of reasons that explain why a particular backend has been raised.
type TrafficManagerBackendConditionReason string

const (
	// TrafficManagerBackendConditionAccepted condition indicates whether endpoints have been created or updated for the profile.
	// This does not indicate whether or not the configuration has been propagated to the data plane.
	//
	// Possible reasons for this condition to be True are:
	//
	// * "Accepted"
	//
	// Possible reasons for this condition to be False are:
	//
	// * "Invalid"
	// * "Pending"
	TrafficManagerBackendConditionAccepted TrafficManagerBackendConditionReason = "Accepted"

	// TrafficManagerBackendReasonAccepted is used with the "Accepted" condition when the condition is True.
	TrafficManagerBackendReasonAccepted TrafficManagerBackendConditionReason = "Accepted"

	// TrafficManagerBackendReasonInvalid is used with the "Accepted" condition when one or
	// more endpoint references have an invalid or unsupported configuration
	// and cannot be configured on the Profile with more details in the message.
	TrafficManagerBackendReasonInvalid TrafficManagerBackendConditionReason = "Invalid"

	// TrafficManagerBackendReasonPending is used with the "Accepted" when creating or updating endpoint hits an internal error with
	// more details in the message and the controller will keep retry.
	TrafficManagerBackendReasonPending TrafficManagerBackendConditionReason = "Pending"
)

//+kubebuilder:object:root=true

// TrafficManagerBackendList contains a list of TrafficManagerBackend.
type TrafficManagerBackendList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	// +listType=set
	Items []TrafficManagerBackend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TrafficManagerBackend{}, &TrafficManagerBackendList{})
}
