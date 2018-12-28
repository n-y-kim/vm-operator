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

// FakeVirtualMachineImages implements VirtualMachineImageInterface
type FakeVirtualMachineImages struct {
	Fake *FakeVmoperatorV1beta1
	ns   string
}

var virtualmachineimagesResource = schema.GroupVersionResource{Group: "vmoperator.vmware.com", Version: "v1beta1", Resource: "virtualmachineimages"}

var virtualmachineimagesKind = schema.GroupVersionKind{Group: "vmoperator.vmware.com", Version: "v1beta1", Kind: "VirtualMachineImage"}

// Get takes name of the virtualMachineImage, and returns the corresponding virtualMachineImage object, and an error if there is any.
func (c *FakeVirtualMachineImages) Get(name string, options v1.GetOptions) (result *v1beta1.VirtualMachineImage, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(virtualmachineimagesResource, c.ns, name), &v1beta1.VirtualMachineImage{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VirtualMachineImage), err
}

// List takes label and field selectors, and returns the list of VirtualMachineImages that match those selectors.
func (c *FakeVirtualMachineImages) List(opts v1.ListOptions) (result *v1beta1.VirtualMachineImageList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(virtualmachineimagesResource, virtualmachineimagesKind, c.ns, opts), &v1beta1.VirtualMachineImageList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.VirtualMachineImageList{}
	for _, item := range obj.(*v1beta1.VirtualMachineImageList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested virtualMachineImages.
func (c *FakeVirtualMachineImages) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(virtualmachineimagesResource, c.ns, opts))

}

// Create takes the representation of a virtualMachineImage and creates it.  Returns the server's representation of the virtualMachineImage, and an error, if there is any.
func (c *FakeVirtualMachineImages) Create(virtualMachineImage *v1beta1.VirtualMachineImage) (result *v1beta1.VirtualMachineImage, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(virtualmachineimagesResource, c.ns, virtualMachineImage), &v1beta1.VirtualMachineImage{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VirtualMachineImage), err
}

// Update takes the representation of a virtualMachineImage and updates it. Returns the server's representation of the virtualMachineImage, and an error, if there is any.
func (c *FakeVirtualMachineImages) Update(virtualMachineImage *v1beta1.VirtualMachineImage) (result *v1beta1.VirtualMachineImage, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(virtualmachineimagesResource, c.ns, virtualMachineImage), &v1beta1.VirtualMachineImage{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VirtualMachineImage), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeVirtualMachineImages) UpdateStatus(virtualMachineImage *v1beta1.VirtualMachineImage) (*v1beta1.VirtualMachineImage, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(virtualmachineimagesResource, "status", c.ns, virtualMachineImage), &v1beta1.VirtualMachineImage{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VirtualMachineImage), err
}

// Delete takes name of the virtualMachineImage and deletes it. Returns an error if one occurs.
func (c *FakeVirtualMachineImages) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(virtualmachineimagesResource, c.ns, name), &v1beta1.VirtualMachineImage{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVirtualMachineImages) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(virtualmachineimagesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1beta1.VirtualMachineImageList{})
	return err
}

// Patch applies the patch and returns the patched virtualMachineImage.
func (c *FakeVirtualMachineImages) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.VirtualMachineImage, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(virtualmachineimagesResource, c.ns, name, data, subresources...), &v1beta1.VirtualMachineImage{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VirtualMachineImage), err
}
