/* **********************************************************
 * Copyright 2018 VMware, Inc.  All rights reserved. -- VMware Confidential
 * **********************************************************/
package fake

import (
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
	v1beta1 "vmware.com/kubevsphere/pkg/client/clientset_generated/clientset/typed/vmoperator/v1beta1"
)

type FakeVmoperatorV1beta1 struct {
	*testing.Fake
}

func (c *FakeVmoperatorV1beta1) VirtualMachines(namespace string) v1beta1.VirtualMachineInterface {
	return &FakeVirtualMachines{c, namespace}
}

func (c *FakeVmoperatorV1beta1) VirtualMachineImages(namespace string) v1beta1.VirtualMachineImageInterface {
	return &FakeVirtualMachineImages{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeVmoperatorV1beta1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
