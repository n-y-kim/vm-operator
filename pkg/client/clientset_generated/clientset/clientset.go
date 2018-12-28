/* **********************************************************
 * Copyright 2018 VMware, Inc.  All rights reserved. -- VMware Confidential
 * **********************************************************/
package clientset

import (
	glog "github.com/golang/glog"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
	vmoperatorv1beta1 "vmware.com/kubevsphere/pkg/client/clientset_generated/clientset/typed/vmoperator/v1beta1"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	VmoperatorV1beta1() vmoperatorv1beta1.VmoperatorV1beta1Interface
	// Deprecated: please explicitly pick a version if possible.
	Vmoperator() vmoperatorv1beta1.VmoperatorV1beta1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	vmoperatorV1beta1 *vmoperatorv1beta1.VmoperatorV1beta1Client
}

// VmoperatorV1beta1 retrieves the VmoperatorV1beta1Client
func (c *Clientset) VmoperatorV1beta1() vmoperatorv1beta1.VmoperatorV1beta1Interface {
	return c.vmoperatorV1beta1
}

// Deprecated: Vmoperator retrieves the default version of VmoperatorClient.
// Please explicitly pick a version.
func (c *Clientset) Vmoperator() vmoperatorv1beta1.VmoperatorV1beta1Interface {
	return c.vmoperatorV1beta1
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.vmoperatorV1beta1, err = vmoperatorv1beta1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		glog.Errorf("failed to create the DiscoveryClient: %v", err)
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.vmoperatorV1beta1 = vmoperatorv1beta1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.vmoperatorV1beta1 = vmoperatorv1beta1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
