/* **********************************************************
 * Copyright 2018 VMware, Inc.  All rights reserved. -- VMware Confidential
 * **********************************************************/

// This file was automatically generated by informer-gen

package internalversion

import (
	internalinterfaces "vmware.com/kubevsphere/pkg/client/informers_generated/internalversion/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// VirtualMachines returns a VirtualMachineInformer.
	VirtualMachines() VirtualMachineInformer
	// VirtualMachineImages returns a VirtualMachineImageInformer.
	VirtualMachineImages() VirtualMachineImageInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// VirtualMachines returns a VirtualMachineInformer.
func (v *version) VirtualMachines() VirtualMachineInformer {
	return &virtualMachineInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// VirtualMachineImages returns a VirtualMachineImageInformer.
func (v *version) VirtualMachineImages() VirtualMachineImageInformer {
	return &virtualMachineImageInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
