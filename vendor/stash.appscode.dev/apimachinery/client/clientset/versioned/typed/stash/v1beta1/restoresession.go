/*
Copyright The Stash Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	"time"

	v1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	scheme "stash.appscode.dev/apimachinery/client/clientset/versioned/scheme"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// RestoreSessionsGetter has a method to return a RestoreSessionInterface.
// A group's client should implement this interface.
type RestoreSessionsGetter interface {
	RestoreSessions(namespace string) RestoreSessionInterface
}

// RestoreSessionInterface has methods to work with RestoreSession resources.
type RestoreSessionInterface interface {
	Create(ctx context.Context, restoreSession *v1beta1.RestoreSession, opts v1.CreateOptions) (*v1beta1.RestoreSession, error)
	Update(ctx context.Context, restoreSession *v1beta1.RestoreSession, opts v1.UpdateOptions) (*v1beta1.RestoreSession, error)
	UpdateStatus(ctx context.Context, restoreSession *v1beta1.RestoreSession, opts v1.UpdateOptions) (*v1beta1.RestoreSession, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.RestoreSession, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.RestoreSessionList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.RestoreSession, err error)
	RestoreSessionExpansion
}

// restoreSessions implements RestoreSessionInterface
type restoreSessions struct {
	client rest.Interface
	ns     string
}

// newRestoreSessions returns a RestoreSessions
func newRestoreSessions(c *StashV1beta1Client, namespace string) *restoreSessions {
	return &restoreSessions{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the restoreSession, and returns the corresponding restoreSession object, and an error if there is any.
func (c *restoreSessions) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.RestoreSession, err error) {
	result = &v1beta1.RestoreSession{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("restoresessions").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of RestoreSessions that match those selectors.
func (c *restoreSessions) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.RestoreSessionList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.RestoreSessionList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("restoresessions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested restoreSessions.
func (c *restoreSessions) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("restoresessions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a restoreSession and creates it.  Returns the server's representation of the restoreSession, and an error, if there is any.
func (c *restoreSessions) Create(ctx context.Context, restoreSession *v1beta1.RestoreSession, opts v1.CreateOptions) (result *v1beta1.RestoreSession, err error) {
	result = &v1beta1.RestoreSession{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("restoresessions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(restoreSession).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a restoreSession and updates it. Returns the server's representation of the restoreSession, and an error, if there is any.
func (c *restoreSessions) Update(ctx context.Context, restoreSession *v1beta1.RestoreSession, opts v1.UpdateOptions) (result *v1beta1.RestoreSession, err error) {
	result = &v1beta1.RestoreSession{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("restoresessions").
		Name(restoreSession.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(restoreSession).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *restoreSessions) UpdateStatus(ctx context.Context, restoreSession *v1beta1.RestoreSession, opts v1.UpdateOptions) (result *v1beta1.RestoreSession, err error) {
	result = &v1beta1.RestoreSession{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("restoresessions").
		Name(restoreSession.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(restoreSession).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the restoreSession and deletes it. Returns an error if one occurs.
func (c *restoreSessions) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("restoresessions").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *restoreSessions) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("restoresessions").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched restoreSession.
func (c *restoreSessions) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.RestoreSession, err error) {
	result = &v1beta1.RestoreSession{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("restoresessions").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
