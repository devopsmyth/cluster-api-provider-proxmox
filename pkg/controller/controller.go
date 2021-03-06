/*
 * Copyright 2019 Rafael Fernández López <ereslibre@ereslibre.es>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controller

import (
	"k8s.io/klog"

	"github.com/ereslibre/cluster-api-provider-proxmox/pkg/cloud/proxmox"
	"k8s.io/client-go/kubernetes"
	clusterapi "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(manager.Manager) error

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager) error {
	for _, f := range AddToManagerFuncs {
		if err := f(m); err != nil {
			return err
		}
	}

	return nil
}

func getActuatorParams(mgr manager.Manager) proxmox.ActuatorParams {
	config := mgr.GetConfig()

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Could not create kubernetes client to talk to the apiserver: %v", err)
	}

	clusterClient, err := clusterapi.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Could not create cluster-api client to talk to the apiserver: %v", err)
	}

	return proxmox.ActuatorParams{
		KubeClient:    kubeClient,
		Client:        mgr.GetClient(),
		ClusterClient: clusterClient,
		Scheme:        mgr.GetScheme(),
		EventRecorder: mgr.GetRecorder("proxmox-controller"),
	}

}
