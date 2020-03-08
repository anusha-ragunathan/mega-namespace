/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ExecAction is inspired from api_v1's ExecAction which executes inside the container via
// the kubelet. Here, we execute on the master.
type ExecAction struct {
	// Command is the command line to execute on the host, the working directory for the
	// command  is root ('/') in the host's filesystem. The command is simply exec'd, it is
	// not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
	// a shell, you need to explicitly call out to that shell.
	// Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
	// +optional
	Command []string `json:"command,omitempty"`
}

// This struct contains a static list of k8s resources that can be provisioned by the controller.
// To implement a dynamic list, we need to use the generic runtime.Object, the dynamic
// k8s client and receive into a unstructured.unstructured object. However, having static types
// is pretty useful for other fields of the NamespaceTemplate.
// Explore the dynamic client later.
// https://stackoverflow.com/questions/53341727/how-to-submit-generic-runtime-object-to-kubernetes-api-using-client-go
// Another drawback is only one instance can be provisioned per resource type.
// XXX: Fix this later.
type AdditionalResources struct {
	PodSpec        v1.PodSpec    `json:"podspec"`
	SecretSpec     v1.Secret     `json:"secretspec"`
	LimitRangeSpec v1.LimitRange `json:"limitrangespec"`
}

// NamespaceTemplateSpec defines the desired state of NamespaceTemplate
type NamespaceTemplateSpec struct {
	// XXX: could be omitempty. which of these is optional?
	Options map[string]string `json:"options,omitempty"`

	PreCreateHook ExecAction `json:"precreatehook,omitempty"`

	PostCreateHook ExecAction `json:"postcreatehook,omitempty"`

	AddResources AdditionalResources `json:"addresources,omitempty"`
}

// NamespaceTemplateStatus defines the observed state of NamespaceTemplate
type NamespaceTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// NamespaceTemplate is the Schema for the namespacetemplates API
type NamespaceTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NamespaceTemplateSpec   `json:"spec,omitempty"`
	Status NamespaceTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NamespaceTemplateList contains a list of NamespaceTemplate
type NamespaceTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NamespaceTemplate{}, &NamespaceTemplateList{})
}
