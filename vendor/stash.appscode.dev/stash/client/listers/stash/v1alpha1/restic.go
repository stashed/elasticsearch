/*
Copyright 2019 The Stash Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "stash.appscode.dev/stash/apis/stash/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ResticLister helps list Restics.
type ResticLister interface {
	// List lists all Restics in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.Restic, err error)
	// Restics returns an object that can list and get Restics.
	Restics(namespace string) ResticNamespaceLister
	ResticListerExpansion
}

// resticLister implements the ResticLister interface.
type resticLister struct {
	indexer cache.Indexer
}

// NewResticLister returns a new ResticLister.
func NewResticLister(indexer cache.Indexer) ResticLister {
	return &resticLister{indexer: indexer}
}

// List lists all Restics in the indexer.
func (s *resticLister) List(selector labels.Selector) (ret []*v1alpha1.Restic, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Restic))
	})
	return ret, err
}

// Restics returns an object that can list and get Restics.
func (s *resticLister) Restics(namespace string) ResticNamespaceLister {
	return resticNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ResticNamespaceLister helps list and get Restics.
type ResticNamespaceLister interface {
	// List lists all Restics in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.Restic, err error)
	// Get retrieves the Restic from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.Restic, error)
	ResticNamespaceListerExpansion
}

// resticNamespaceLister implements the ResticNamespaceLister
// interface.
type resticNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Restics in the indexer for a given namespace.
func (s resticNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Restic, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Restic))
	})
	return ret, err
}

// Get retrieves the Restic from the indexer for a given namespace and name.
func (s resticNamespaceLister) Get(name string) (*v1alpha1.Restic, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("restic"), name)
	}
	return obj.(*v1alpha1.Restic), nil
}
