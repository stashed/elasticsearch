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

package v1alpha1

import (
	"fmt"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/catalog"
	"kubedb.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
)

func (f FerretDBVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralFerretDBVersion))
}

var _ apis.ResourceInfo = &FerretDBVersion{}

func (f FerretDBVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralFerretDBVersion, catalog.GroupName)
}

func (f FerretDBVersion) ResourceShortCode() string {
	return ResourceCodeFerretDBVersion
}

func (f FerretDBVersion) ResourceKind() string {
	return ResourceKindFerretDBVersion
}

func (f FerretDBVersion) ResourceSingular() string {
	return ResourceSingularFerretDBVersion
}

func (f FerretDBVersion) ResourcePlural() string {
	return ResourcePluralFerretDBVersion
}

func (f FerretDBVersion) ValidateSpecs() error {
	if f.Spec.Version == "" ||
		f.Spec.DB.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for ferretdbVersion "%v":
spec.version,
spec.db.image`, f.Name)
	}
	return nil
}
