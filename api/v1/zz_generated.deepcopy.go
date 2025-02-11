// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AdditionalResources) DeepCopyInto(out *AdditionalResources) {
	*out = *in
	in.Pod.DeepCopyInto(&out.Pod)
	in.Secret.DeepCopyInto(&out.Secret)
	in.LimitRange.DeepCopyInto(&out.LimitRange)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdditionalResources.
func (in *AdditionalResources) DeepCopy() *AdditionalResources {
	if in == nil {
		return nil
	}
	out := new(AdditionalResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExecAction) DeepCopyInto(out *ExecAction) {
	*out = *in
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExecAction.
func (in *ExecAction) DeepCopy() *ExecAction {
	if in == nil {
		return nil
	}
	out := new(ExecAction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceTemplate) DeepCopyInto(out *NamespaceTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceTemplate.
func (in *NamespaceTemplate) DeepCopy() *NamespaceTemplate {
	if in == nil {
		return nil
	}
	out := new(NamespaceTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NamespaceTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceTemplateList) DeepCopyInto(out *NamespaceTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NamespaceTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceTemplateList.
func (in *NamespaceTemplateList) DeepCopy() *NamespaceTemplateList {
	if in == nil {
		return nil
	}
	out := new(NamespaceTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NamespaceTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceTemplateSpec) DeepCopyInto(out *NamespaceTemplateSpec) {
	*out = *in
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.PreCreateHook.DeepCopyInto(&out.PreCreateHook)
	in.PostCreateHook.DeepCopyInto(&out.PostCreateHook)
	in.AddResources.DeepCopyInto(&out.AddResources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceTemplateSpec.
func (in *NamespaceTemplateSpec) DeepCopy() *NamespaceTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(NamespaceTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceTemplateStatus) DeepCopyInto(out *NamespaceTemplateStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceTemplateStatus.
func (in *NamespaceTemplateStatus) DeepCopy() *NamespaceTemplateStatus {
	if in == nil {
		return nil
	}
	out := new(NamespaceTemplateStatus)
	in.DeepCopyInto(out)
	return out
}
