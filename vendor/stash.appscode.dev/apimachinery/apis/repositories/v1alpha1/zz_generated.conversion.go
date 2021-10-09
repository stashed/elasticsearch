//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright AppsCode Inc. and Contributors

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

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	repositories "stash.appscode.dev/apimachinery/apis/repositories"

	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*Snapshot)(nil), (*repositories.Snapshot)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Snapshot_To_repositories_Snapshot(a.(*Snapshot), b.(*repositories.Snapshot), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*repositories.Snapshot)(nil), (*Snapshot)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_repositories_Snapshot_To_v1alpha1_Snapshot(a.(*repositories.Snapshot), b.(*Snapshot), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*SnapshotList)(nil), (*repositories.SnapshotList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_SnapshotList_To_repositories_SnapshotList(a.(*SnapshotList), b.(*repositories.SnapshotList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*repositories.SnapshotList)(nil), (*SnapshotList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_repositories_SnapshotList_To_v1alpha1_SnapshotList(a.(*repositories.SnapshotList), b.(*SnapshotList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*SnapshotStatus)(nil), (*repositories.SnapshotStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_SnapshotStatus_To_repositories_SnapshotStatus(a.(*SnapshotStatus), b.(*repositories.SnapshotStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*repositories.SnapshotStatus)(nil), (*SnapshotStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_repositories_SnapshotStatus_To_v1alpha1_SnapshotStatus(a.(*repositories.SnapshotStatus), b.(*SnapshotStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_Snapshot_To_repositories_Snapshot(in *Snapshot, out *repositories.Snapshot, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_SnapshotStatus_To_repositories_SnapshotStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_Snapshot_To_repositories_Snapshot is an autogenerated conversion function.
func Convert_v1alpha1_Snapshot_To_repositories_Snapshot(in *Snapshot, out *repositories.Snapshot, s conversion.Scope) error {
	return autoConvert_v1alpha1_Snapshot_To_repositories_Snapshot(in, out, s)
}

func autoConvert_repositories_Snapshot_To_v1alpha1_Snapshot(in *repositories.Snapshot, out *Snapshot, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_repositories_SnapshotStatus_To_v1alpha1_SnapshotStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_repositories_Snapshot_To_v1alpha1_Snapshot is an autogenerated conversion function.
func Convert_repositories_Snapshot_To_v1alpha1_Snapshot(in *repositories.Snapshot, out *Snapshot, s conversion.Scope) error {
	return autoConvert_repositories_Snapshot_To_v1alpha1_Snapshot(in, out, s)
}

func autoConvert_v1alpha1_SnapshotList_To_repositories_SnapshotList(in *SnapshotList, out *repositories.SnapshotList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]repositories.Snapshot, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_Snapshot_To_repositories_Snapshot(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1alpha1_SnapshotList_To_repositories_SnapshotList is an autogenerated conversion function.
func Convert_v1alpha1_SnapshotList_To_repositories_SnapshotList(in *SnapshotList, out *repositories.SnapshotList, s conversion.Scope) error {
	return autoConvert_v1alpha1_SnapshotList_To_repositories_SnapshotList(in, out, s)
}

func autoConvert_repositories_SnapshotList_To_v1alpha1_SnapshotList(in *repositories.SnapshotList, out *SnapshotList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Snapshot, len(*in))
		for i := range *in {
			if err := Convert_repositories_Snapshot_To_v1alpha1_Snapshot(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_repositories_SnapshotList_To_v1alpha1_SnapshotList is an autogenerated conversion function.
func Convert_repositories_SnapshotList_To_v1alpha1_SnapshotList(in *repositories.SnapshotList, out *SnapshotList, s conversion.Scope) error {
	return autoConvert_repositories_SnapshotList_To_v1alpha1_SnapshotList(in, out, s)
}

func autoConvert_v1alpha1_SnapshotStatus_To_repositories_SnapshotStatus(in *SnapshotStatus, out *repositories.SnapshotStatus, s conversion.Scope) error {
	out.Tree = in.Tree
	out.Paths = *(*[]string)(unsafe.Pointer(&in.Paths))
	out.Hostname = in.Hostname
	out.Username = in.Username
	out.UID = int(in.UID)
	out.Gid = int(in.Gid)
	out.Tags = *(*[]string)(unsafe.Pointer(&in.Tags))
	out.Repository = in.Repository
	return nil
}

// Convert_v1alpha1_SnapshotStatus_To_repositories_SnapshotStatus is an autogenerated conversion function.
func Convert_v1alpha1_SnapshotStatus_To_repositories_SnapshotStatus(in *SnapshotStatus, out *repositories.SnapshotStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_SnapshotStatus_To_repositories_SnapshotStatus(in, out, s)
}

func autoConvert_repositories_SnapshotStatus_To_v1alpha1_SnapshotStatus(in *repositories.SnapshotStatus, out *SnapshotStatus, s conversion.Scope) error {
	out.Tree = in.Tree
	out.Paths = *(*[]string)(unsafe.Pointer(&in.Paths))
	out.Hostname = in.Hostname
	out.Username = in.Username
	out.UID = int32(in.UID)
	out.Gid = int32(in.Gid)
	out.Tags = *(*[]string)(unsafe.Pointer(&in.Tags))
	out.Repository = in.Repository
	return nil
}

// Convert_repositories_SnapshotStatus_To_v1alpha1_SnapshotStatus is an autogenerated conversion function.
func Convert_repositories_SnapshotStatus_To_v1alpha1_SnapshotStatus(in *repositories.SnapshotStatus, out *SnapshotStatus, s conversion.Scope) error {
	return autoConvert_repositories_SnapshotStatus_To_v1alpha1_SnapshotStatus(in, out, s)
}
