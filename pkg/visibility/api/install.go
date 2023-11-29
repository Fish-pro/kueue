// Copyright 2023 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"

	v1alpha1 "sigs.k8s.io/kueue/apis/visibility/v1alpha1"
	"sigs.k8s.io/kueue/pkg/queue"
	apirest "sigs.k8s.io/kueue/pkg/visibility/api/rest"
)

var (
	Scheme = runtime.NewScheme()

	Codecs = serializer.NewCodecFactory(Scheme)

	ParameterCodec = runtime.NewParameterCodec(Scheme)
)

func init() {
	utilruntime.Must(v1alpha1.AddToScheme(Scheme))
	utilruntime.Must(Scheme.SetVersionPriority(v1alpha1.GroupVersion))
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
}

// Install installs API scheme defined in apis/v1alpha1 and registers storage
func Install(server *genericapiserver.GenericAPIServer, kueueMgr *queue.Manager) error {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(v1alpha1.GroupVersion.Group, Scheme, ParameterCodec, Codecs)
	pendingWorkloadsInCqREST := apirest.NewPendingWorkloadsInCqREST(kueueMgr)
	cqREST := apirest.NewCqREST()
	pendingWorkloadsInLqREST := apirest.NewPendingWorkloadsInLqREST(kueueMgr)
	lqREST := apirest.NewLqREST()

	visibilityServerResources := map[string]rest.Storage{
		"clusterqueues":                  cqREST,
		"clusterqueues/pendingworkloads": pendingWorkloadsInCqREST,
		"localqueues":                    lqREST,
		"localqueues/pendingworkloads":   pendingWorkloadsInLqREST,
	}
	apiGroupInfo.VersionedResourcesStorageMap[v1alpha1.GroupVersion.Version] = visibilityServerResources
	return server.InstallAPIGroups(&apiGroupInfo)
}
