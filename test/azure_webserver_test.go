package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "56d67ba7-7129-4796-810e-a6159bb4b456"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "sing1249",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Confirm NIC exists and is connected to the VM
	nics := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)
	assert.NotEmpty(t, nics, "VM should have at least one NIC connected")

	// Confirm the VM is running the correct Ubuntu version
	expectedUbuntuSKU := "22_04-lts-gen2" 
	vmImage := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)
	assert.Equal(t, "Canonical", vmImage.Publisher, "The VM should be from Canonical")
	assert.Equal(t, "UbuntuServer", vmImage.Offer, "The VM should be an Ubuntu Server")
	assert.Equal(t, expectedUbuntuSKU, vmImage.SKU, "The VM should be running the expected Ubuntu version")
}
