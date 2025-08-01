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

package v1beta2

import (
	ocpv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
	lifecycleapi "kubevirt.io/controller-lifecycle-operator-sdk/api"
)

const (
	OperatorPausedAnnotation = "kubevirt.io/operator.paused"
)

type TemplateValidator struct {
	// Replicas is the number of replicas of the template validator pod
	//+kubebuilder:validation:Minimum=0
	//+kubebuilder:default=2
	Replicas *int32 `json:"replicas,omitempty"`

	// Placement describes the node scheduling configuration
	Placement *lifecycleapi.NodePlacement `json:"placement,omitempty"`
}

type CommonTemplates struct {
	// Namespace is the k8s namespace where CommonTemplates should be installed
	//+kubebuilder:validation:MaxLength=63
	//+kubebuilder:validation:Pattern=^[a-z0-9]([-a-z0-9]*[a-z0-9])?$
	Namespace string `json:"namespace"`

	// DataImportCronTemplates defines a list of DataImportCrons managed by the SSP
	// Operator. This is intended for images used by CommonTemplates.
	DataImportCronTemplates []DataImportCronTemplate `json:"dataImportCronTemplates,omitempty"`
}

type CommonInstancetypes struct {
	// URL of a remote Kustomize target from which to generate and deploy resources.
	//
	// The following caveats apply to the provided URL:
	//
	// * Only 'https://' and 'git://' URLs are supported.
	//
	// * The URL must include '?ref=$ref' or '?version=$ref' pinning it to a specific
	//   reference. It is recommended that the reference be a specific commit or tag
	//   to ensure the generated contents does not change over time. As such it is
	//   recommended not to use branches as the ref for the time being.
	//
	// * Only VirtualMachineClusterPreference and VirtualMachineClusterInstancetype
	//   resources generated from the URL are deployed by the operand.
	//
	// See the following Kustomize documentation for more details:
	//
	// remote targets
	// https://github.com/kubernetes-sigs/kustomize/blob/master/examples/remoteBuild.md
	URL *string `json:"url,omitempty"`
}

// SSPSpec defines the desired state of SSP
type SSPSpec struct {
	// TemplateValidator is configuration of the template validator operand
	TemplateValidator *TemplateValidator `json:"templateValidator,omitempty"`

	// CommonTemplates is the configuration of the common templates operand
	CommonTemplates CommonTemplates `json:"commonTemplates"`

	// TLSSecurityProfile is a configuration for the TLS.
	TLSSecurityProfile *ocpv1.TLSSecurityProfile `json:"tlsSecurityProfile,omitempty"`

	// TokenGenerationService configures the service for generating tokens to access VNC for a VM.
	TokenGenerationService *TokenGenerationService `json:"tokenGenerationService,omitempty"`

	// CommonInstancetypes is ignored.
	// +kubebuilder:deprecatedversion:warning="commonInstancetypes filed is ignored"
	// Deprecated: This field is ignored.
	CommonInstancetypes *CommonInstancetypes `json:"commonInstancetypes,omitempty"`

	// TektonPipelines is the configuration of the tekton-pipelines operand
	// +kubebuilder:deprecatedversion:warning="tektonPipelines filed is ignored"
	// Deprecated: This field is ignored.
	TektonPipelines *TektonPipelines `json:"tektonPipelines,omitempty"`

	// TektonTasks is the configuration of the tekton-tasks operand
	// +kubebuilder:deprecatedversion:warning="tektonTasks filed is ignored"
	// Deprecated: This field is ignored.
	TektonTasks *TektonTasks `json:"tektonTasks,omitempty"`

	// FeatureGates for SSP
	FeatureGates *FeatureGates `json:"featureGates,omitempty"`
}

// TektonPipelines defines the desired state of pipelines
type TektonPipelines struct {
	Namespace string `json:"namespace,omitempty"`
}

// TektonTasks defines variables for configuration of tasks
type TektonTasks struct {
	Namespace string `json:"namespace,omitempty"`
}

// FeatureGates for SSP
type FeatureGates struct {
	// +kubebuilder:deprecatedversion:warning="tekton task resources are no longer deployed by SSP"
	// Deprecated: This field is ignored.
	DeployTektonTaskResources bool `json:"deployTektonTaskResources,omitempty"`

	// +kubebuilder:deprecatedversion:warning="use the .spec.tokenGenerationService to enable vm-console-proxy"
	// Deprecated: This field is ignored.
	DeployVmConsoleProxy bool `json:"deployVmConsoleProxy,omitempty"`

	// +kubebuilder:deprecatedversion:warning="instance types and preferences are no longer deployed by SSP"
	// Deprecated: This field is ignored.
	DeployCommonInstancetypes *bool `json:"deployCommonInstancetypes,omitempty"`
}

// DataImportCronTemplate defines the template type for DataImportCrons.
// It requires metadata.name to be specified while leaving namespace as optional.
type DataImportCronTemplate struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec cdiv1beta1.DataImportCronSpec `json:"spec"`
}

// AsDataImportCron converts the DataImportCronTemplate to a cdiv1beta1.DataImportCron
func (t *DataImportCronTemplate) AsDataImportCron() cdiv1beta1.DataImportCron {
	return cdiv1beta1.DataImportCron{
		ObjectMeta: t.ObjectMeta,
		Spec:       t.Spec,
	}
}

// TokenGenerationService configures the service for generating tokens to access VNC for a VM.
type TokenGenerationService struct {
	Enabled bool `json:"enabled,omitempty"`
}

// SSPStatus defines the observed state of SSP
type SSPStatus struct {
	lifecycleapi.Status `json:",inline"`

	// Paused is true when the operator notices paused annotation.
	Paused bool `json:"paused,omitempty"`

	// ObservedGeneration is the latest generation observed by the operator.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// SSP is the Schema for the ssps API
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:deprecatedversion:warning="SSP v1beta2 is deprecated and will not be updated with new functionality. It will stop being served in the next version."
type SSP struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SSPSpec   `json:"spec,omitempty"`
	Status SSPStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SSPList contains a list of SSP
type SSPList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SSP `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SSP{}, &SSPList{})
}
