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

package fake

import (
	"context"

	v1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeBackupSessions implements BackupSessionInterface
type FakeBackupSessions struct {
	Fake *FakeStashV1beta1
	ns   string
}

var backupsessionsResource = schema.GroupVersionResource{Group: "stash.appscode.com", Version: "v1beta1", Resource: "backupsessions"}

var backupsessionsKind = schema.GroupVersionKind{Group: "stash.appscode.com", Version: "v1beta1", Kind: "BackupSession"}

// Get takes name of the backupSession, and returns the corresponding backupSession object, and an error if there is any.
func (c *FakeBackupSessions) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.BackupSession, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(backupsessionsResource, c.ns, name), &v1beta1.BackupSession{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.BackupSession), err
}

// List takes label and field selectors, and returns the list of BackupSessions that match those selectors.
func (c *FakeBackupSessions) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.BackupSessionList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(backupsessionsResource, backupsessionsKind, c.ns, opts), &v1beta1.BackupSessionList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.BackupSessionList{ListMeta: obj.(*v1beta1.BackupSessionList).ListMeta}
	for _, item := range obj.(*v1beta1.BackupSessionList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested backupSessions.
func (c *FakeBackupSessions) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(backupsessionsResource, c.ns, opts))

}

// Create takes the representation of a backupSession and creates it.  Returns the server's representation of the backupSession, and an error, if there is any.
func (c *FakeBackupSessions) Create(ctx context.Context, backupSession *v1beta1.BackupSession, opts v1.CreateOptions) (result *v1beta1.BackupSession, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(backupsessionsResource, c.ns, backupSession), &v1beta1.BackupSession{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.BackupSession), err
}

// Update takes the representation of a backupSession and updates it. Returns the server's representation of the backupSession, and an error, if there is any.
func (c *FakeBackupSessions) Update(ctx context.Context, backupSession *v1beta1.BackupSession, opts v1.UpdateOptions) (result *v1beta1.BackupSession, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(backupsessionsResource, c.ns, backupSession), &v1beta1.BackupSession{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.BackupSession), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBackupSessions) UpdateStatus(ctx context.Context, backupSession *v1beta1.BackupSession, opts v1.UpdateOptions) (*v1beta1.BackupSession, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(backupsessionsResource, "status", c.ns, backupSession), &v1beta1.BackupSession{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.BackupSession), err
}

// Delete takes name of the backupSession and deletes it. Returns an error if one occurs.
func (c *FakeBackupSessions) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(backupsessionsResource, c.ns, name), &v1beta1.BackupSession{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBackupSessions) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(backupsessionsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.BackupSessionList{})
	return err
}

// Patch applies the patch and returns the patched backupSession.
func (c *FakeBackupSessions) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.BackupSession, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(backupsessionsResource, c.ns, name, pt, data, subresources...), &v1beta1.BackupSession{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.BackupSession), err
}
