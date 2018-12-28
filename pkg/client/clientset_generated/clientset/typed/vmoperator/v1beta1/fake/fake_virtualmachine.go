/* **********************************************************
 * Copyright 2018 VMware, Inc.  All rights reserved. -- VMware Confidential
 * **********************************************************/
package fake

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1beta1 "vmware.com/kubevsphere/pkg/apis/vmoperator/v1beta1"
)

// FakeVirtualMachines implements VirtualMachineInterface
type FakeVirtualMachines struct {
	Fake *FakeVmoperatorV1beta1
	ns   string
}

var virtualmachinesResource = schema.GroupVersionResource{Group: "vmoperator.vmware.com", Version: "v1beta1", Resource: "virtualmachines"}

var virtualmachinesKind = schema.GroupVersionKind{Group: "vmoperator.vmware.com", Version: "v1beta1", Kind: "VirtualMachine"}

// Get takes name of the virtualMachine, and returns the corresponding virtualMachine object, and an error if there is any.
func (c *FakeVirtualMachines) Get(name string, options v1.GetOptions) (result *v1beta1.VirtualMachine, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(virtualmachinesResource, c.ns, name), &v1beta1.VirtualMachine{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VirtualMachine), err
}

// List takes label and field selectors, and returns the list of VirtualMachines that match those selectors.
func (c *FakeVirtualMachines) List(opts v1.ListOptions) (result *v1beta1.VirtualMachineList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(virtualmachinesResource, virtualmachinesKind, c.ns, opts), &v1beta1.VirtualMachineList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.VirtualMachineList{}
	for _, item := range obj.(*v1beta1.VirtualMachineList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested virtualMachines.
func (c *FakeVirtualMachines) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(virtualmachinesResource, c.ns, opts))

}

// Create takes the representation of a virtualMachine and creates it.  Returns the server's representation of the virtualMachine, and an error, if there is any.
func (c *FakeVirtualMachines) Create(virtualMachine *v1beta1.VirtualMachine) (result *v1beta1.VirtualMachine, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(virtualmachinesResource, c.ns, virtualMachine), &v1beta1.VirtualMachine{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VirtualMachine), err
}

// Update takes the representation of a virtualMachine and updates it. Returns the server's representation of the virtualMachine, and an error, if there is any.
func (c *FakeVirtualMachines) Update(virtualMachine *v1beta1.VirtualMachine) (result *v1beta1.VirtualMachine, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(virtualmachinesResource, c.ns, virtualMachine), &v1beta1.VirtualMachine{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VirtualMachine), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeVirtualMachines) UpdateStatus(virtualMachine *v1beta1.VirtualMachine) (*v1beta1.VirtualMachine, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(virtualmachinesResource, "status", c.ns, virtualMachine), &v1beta1.VirtualMachine{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VirtualMachine), err
}

// Delete takes name of the virtualMachine and deletes it. Returns an error if one occurs.
func (c *FakeVirtualMachines) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(virtualmachinesResource, c.ns, name), &v1beta1.VirtualMachine{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVirtualMachines) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(virtualmachinesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1beta1.VirtualMachineList{})
	return err
}

// Patch applies the patch and returns the patched virtualMachine.
func (c *FakeVirtualMachines) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.VirtualMachine, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(virtualmachinesResource, c.ns, name, data, subresources...), &v1beta1.VirtualMachine{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VirtualMachine), err
}
