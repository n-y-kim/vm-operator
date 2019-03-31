/*
Copyright 2017 The Kubernetes Authors.

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

// This file was autogenerated by apiregister-gen. Do not edit it manually!

package controller

import (
	"github.com/kubernetes-incubator/apiserver-builder-alpha/pkg/controller"
	"gitlab.eng.vmware.com/iaas-platform/vm-operator/pkg/controller/sharedinformers"
	"gitlab.eng.vmware.com/iaas-platform/vm-operator/pkg/controller/virtualmachine"
	"gitlab.eng.vmware.com/iaas-platform/vm-operator/pkg/controller/virtualmachineclass"
	"gitlab.eng.vmware.com/iaas-platform/vm-operator/pkg/controller/virtualmachineimage"
	"gitlab.eng.vmware.com/iaas-platform/vm-operator/pkg/controller/virtualmachineservice"
	"k8s.io/client-go/rest"
)

func GetAllControllers(config *rest.Config) ([]controller.Controller, chan struct{}) {
	shutdown := make(chan struct{})
	si := sharedinformers.NewSharedInformers(config, shutdown)
	return []controller.Controller{
		virtualmachine.NewVirtualMachineController(config, si),
		virtualmachineclass.NewVirtualMachineClassController(config, si),
		virtualmachineimage.NewVirtualMachineImageController(config, si),
		virtualmachineservice.NewVirtualMachineServiceController(config, si),
	}, shutdown
}
