/* **********************************************************
 * Copyright 2018 VMware, Inc.  All rights reserved. -- VMware Confidential
 * **********************************************************/

package v1alpha1

import (
	"github.com/golang/glog"
	"gitlab.eng.vmware.com/iaas-platform/vm-operator/pkg/apis"
	"gitlab.eng.vmware.com/iaas-platform/vm-operator/pkg/apis/vmoperator/v1alpha1"
	"gitlab.eng.vmware.com/iaas-platform/vm-operator/pkg/client/clientset_generated/clientset"
	"gitlab.eng.vmware.com/iaas-platform/vm-operator/pkg/openapi"
	"gitlab.eng.vmware.com/iaas-platform/vm-operator/pkg/vmprovider"
	"gitlab.eng.vmware.com/iaas-platform/vm-operator/pkg/vmprovider/providers/vsphere"
	"gitlab.eng.vmware.com/iaas-platform/vm-operator/test/integration"
	"testing"

	"github.com/kubernetes-incubator/apiserver-builder-alpha/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"

	vmrest "gitlab.eng.vmware.com/iaas-platform/vm-operator/pkg/apis/vmoperator/rest"
)

var testenv *test.TestEnvironment
var config *rest.Config
var cs *clientset.Clientset
var vcsim *integration.VcSimInstance

func TestV1alpha1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "v1 Suite", []Reporter{test.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	vcsim = integration.NewVcSimInstance()
	address, port := vcsim.Start()

	provider, err := vsphere.NewVSphereVmProviderFromConfig(integration.NewIntegrationVmOperatorConfig(address, port))
	if err != nil {
		glog.Fatalf("Failed to create vSphere provider: %v", err)
	}

	vmprovider.RegisterVmProvider(provider)

	if err := v1alpha1.RegisterRestProvider(vmrest.NewVirtualMachineImagesREST(provider)); err != nil {
		glog.Fatalf("Failed to register REST provider: %s", err)
	}

	testenv = test.NewTestEnvironment()
	config = testenv.Start(apis.GetAllApiBuilders(), openapi.GetOpenAPIDefinitions)
	cs = clientset.NewForConfigOrDie(config)
})

var _ = AfterSuite(func() {
	testenv.Stop()
	vcsim.Stop()
})
