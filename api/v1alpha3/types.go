/*
Copyright 2020 The Kubernetes Authors.

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

package v1alpha3

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	// ControlPlane machine label
	ControlPlane string = "control-plane"
	// Node machine label
	Node string = "node"
)

// Network encapsulates the state of Azure networking resources.
type Network struct {
	// SecurityGroups is a map from the role/kind of the security group to its unique name, if any.
	SecurityGroups map[SecurityGroupRole]SecurityGroup `json:"securityGroups,omitempty"`

	// APIServerLB is the Kubernetes API server load balancer.
	APIServerLB LoadBalancer `json:"apiServerLb,omitempty"`

	// APIServerIP is the Kubernetes API server public IP address.
	APIServerIP PublicIP `json:"apiServerIp,omitempty"`
}

// NetworkSpec specifies what the Azure networking resources should look like.
type NetworkSpec struct {
	// Vnet is the configuration for the Azure virtual network.
	// +optional
	Vnet VnetSpec `json:"vnet,omitempty"`

	// Subnets is the configuration for the control-plane subnet and the node subnet.
	// +optional
	Subnets Subnets `json:"subnets,omitempty"`
}

// VnetSpec configures an Azure virtual network.
type VnetSpec struct {
	// ResourceGroup is the name of the resource group of the existing virtual network
	// or the resource group where a managed virtual network should be created.
	ResourceGroup string `json:"resourceGroup,omitempty"`

	// ID is the identifier of the virtual network this provider should use to create resources.
	ID string `json:"id,omitempty"`

	// Name defines a name for the virtual network resource.
	Name string `json:"name"`

	// CidrBlock is the CIDR block to be used when the provider creates a managed virtual network.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// Tags is a collection of tags describing the resource.
	Tags Tags `json:"tags,omitempty"`
}

// IsManaged returns true if the vnet is managed.
func (v *VnetSpec) IsManaged(clusterName string) bool {
	return v.ID == "" || v.Tags.HasOwned(clusterName)
}

// Subnets is a slice of Subnet.
type Subnets []*SubnetSpec

// SecurityGroupRole defines the unique role of a security group.
type SecurityGroupRole string

const (
	// SecurityGroupNode defines a Kubernetes workload node role
	SecurityGroupNode = SecurityGroupRole(Node)

	// SecurityGroupControlPlane defines a Kubernetes control plane node role
	SecurityGroupControlPlane = SecurityGroupRole(ControlPlane)
)

// SecurityGroup defines an Azure security group.
type SecurityGroup struct {
	ID           string       `json:"id,omitempty"`
	Name         string       `json:"name,omitempty"`
	IngressRules IngressRules `json:"ingressRule,omitempty"`
	Tags         Tags         `json:"tags,omitempty"`
}

// SecurityGroupProtocol defines the protocol type for a security group rule.
type SecurityGroupProtocol string

const (
	// SecurityGroupProtocolAll is a wildcard for all IP protocols
	SecurityGroupProtocolAll = SecurityGroupProtocol("*")

	// SecurityGroupProtocolTCP represents the TCP protocol in ingress rules
	SecurityGroupProtocolTCP = SecurityGroupProtocol("Tcp")

	// SecurityGroupProtocolUDP represents the UDP protocol in ingress rules
	SecurityGroupProtocolUDP = SecurityGroupProtocol("Udp")
)

// IngressRule defines an Azure ingress rule for security groups.
type IngressRule struct {
	Description string                `json:"description"`
	Protocol    SecurityGroupProtocol `json:"protocol"`

	// SourcePorts - The source port or range. Integer or range between 0 and 65535. Asterix '*' can also be used to match all ports.
	SourcePorts *string `json:"sourcePorts,omitempty"`

	// DestinationPorts - The destination port or range. Integer or range between 0 and 65535. Asterix '*' can also be used to match all ports.
	DestinationPorts *string `json:"destinationPorts,omitempty"`

	// Source - The CIDR or source IP range. Asterix '*' can also be used to match all source IPs. Default tags such as 'VirtualNetwork', 'AzureLoadBalancer' and 'Internet' can also be used. If this is an ingress rule, specifies where network traffic originates from.
	Source *string `json:"source,omitempty"`

	// Destination - The destination address prefix. CIDR or destination IP range. Asterix '*' can also be used to match all source IPs. Default tags such as 'VirtualNetwork', 'AzureLoadBalancer' and 'Internet' can also be used.
	Destination *string `json:"destination,omitempty"`
}

// IngressRules is a slice of Azure ingress rules for security groups.
type IngressRules []*IngressRule

// PublicIP defines an Azure public IP address.
type PublicIP struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
	DNSName   string `json:"dnsName,omitempty"`
}

// LoadBalancer defines an Azure load balancer.
type LoadBalancer struct {
	ID               string           `json:"id,omitempty"`
	Name             string           `json:"name,omitempty"`
	SKU              SKU              `json:"sku,omitempty"`
	FrontendIPConfig FrontendIPConfig `json:"frontendIpConfig,omitempty"`
	BackendPool      BackendPool      `json:"backendPool,omitempty"`
	Tags             Tags             `json:"tags,omitempty"`
}

// FrontendIPConfig - DO NOT USE
// this empty struct is here to preserve backwards compatibility and should be removed in v1alpha4
type FrontendIPConfig struct{}

// SKU defines an Azure load balancer SKU.
type SKU string

const (
	// SKUBasic is the value for the Azure load balancer Basic SKU
	SKUBasic = SKU("Basic")
	// SKUStandard is the value for the Azure load balancer Standard SKU
	SKUStandard = SKU("Standard")
)

// BackendPool defines a load balancer backend pool
type BackendPool struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

// VMState describes the state of an Azure virtual machine.
type VMState string

const (
	// VMStateCreating ...
	VMStateCreating VMState = "Creating"
	// VMStateDeleting ...
	VMStateDeleting VMState = "Deleting"
	// VMStateFailed ...
	VMStateFailed VMState = "Failed"
	// VMStateMigrating ...
	VMStateMigrating VMState = "Migrating"
	// VMStateSucceeded ...
	VMStateSucceeded VMState = "Succeeded"
	// VMStateUpdating ...
	VMStateUpdating VMState = "Updating"
)

// VM describes an Azure virtual machine.
type VM struct {
	ID               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
	// Hardware profile
	VMSize string `json:"vmSize,omitempty"`
	// Storage profile
	Image         Image  `json:"image,omitempty"`
	OSDisk        OSDisk `json:"osDisk,omitempty"`
	StartupScript string `json:"startupScript,omitempty"`
	// State - The provisioning state, which only appears in the response.
	State    VMState    `json:"vmState,omitempty"`
	Identity VMIdentity `json:"identity,omitempty"`
	Tags     Tags       `json:"tags,omitempty"`

	// Addresses contains the addresses associated with the Azure VM.
	Addresses []corev1.NodeAddress `json:"addresses,omitempty"`
}

// Image defines information about the image to use for VM creation.
// There are three ways to specify an image: by ID, Marketplace Image or SharedImageGallery
// One of ID, SharedImage or Marketplace should be set.
type Image struct {
	// ID specifies an image to use by ID
	// +optional
	ID *string `json:"id,omitempty"`

	// SharedGallery specifies an image to use from an Azure Shared Image Gallery
	// +optional
	SharedGallery *AzureSharedGalleryImage `json:"sharedGallery,omitempty"`

	// Marketplace specifies an image to use from the Azure Marketplace
	// +optional
	Marketplace *AzureMarketplaceImage `json:"marketplace,omitempty"`
}

// AzureMarketplaceImage defines an image in the Azure Marketplace to use for VM creation
type AzureMarketplaceImage struct {
	// Publisher is the name of the organization that created the image
	// +kubebuilder:validation:MinLength=1
	Publisher string `json:"publisher"`
	// Offer specifies the name of a group of related images created by the publisher.
	// For example, UbuntuServer, WindowsServer
	// +kubebuilder:validation:MinLength=1
	Offer string `json:"offer"`
	// SKU specifies an instance of an offer, such as a major release of a distribution.
	// For example, 18.04-LTS, 2019-Datacenter
	// +kubebuilder:validation:MinLength=1
	SKU string `json:"sku"`
	// Version specifies the version of an image sku. The allowed formats
	// are Major.Minor.Build or 'latest'. Major, Minor, and Build are decimal numbers.
	// Specify 'latest' to use the latest version of an image available at deploy time.
	// Even if you use 'latest', the VM image will not automatically update after deploy
	// time even if a new version becomes available.
	// +kubebuilder:validation:MinLength=1
	Version string `json:"version"`
}

// AzureSharedGalleryImage defines an image in a Shared Image Gallery to use for VM creation
type AzureSharedGalleryImage struct {
	// SubscriptionID is the identifier of the subscription that contains the shared image gallery
	// +kubebuilder:validation:MinLength=1
	SubscriptionID string `json:"subscriptionID"`
	// ResourceGroup specifies the resource group containing the shared image gallery
	// +kubebuilder:validation:MinLength=1
	ResourceGroup string `json:"resourceGroup"`
	// Gallery specifies the name of the shared image gallery that contains the image
	// +kubebuilder:validation:MinLength=1
	Gallery string `json:"gallery"`
	// Name is the name of the image
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// Version specifies the version of the marketplace image. The allowed formats
	// are Major.Minor.Build or 'latest'. Major, Minor, and Build are decimal numbers.
	// Specify 'latest' to use the latest version of an image available at deploy time.
	// Even if you use 'latest', the VM image will not automatically update after deploy
	// time even if a new version becomes available.
	// +kubebuilder:validation:MinLength=1
	Version string `json:"version"`
}

// AvailabilityZone specifies an Azure Availability Zone
//
// Deprecated: Use FailureDomain instead
type AvailabilityZone struct {
	ID      *string `json:"id,omitempty"`
	Enabled *bool   `json:"enabled,omitempty"`
}

// VMIdentity defines the identity of the virtual machine, if configured.
// +kubebuilder:validation:Enum=None;SystemAssigned;UserAssigned
type VMIdentity string

const (
	// VMIdentityNone ...
	VMIdentityNone VMIdentity = "None"
	// VMIdentitySystemAssigned ...
	VMIdentitySystemAssigned VMIdentity = "SystemAssigned"
	// VMIdentityUserAssigned ...
	VMIdentityUserAssigned VMIdentity = "UserAssigned"
)

// UserAssignedIdentity defines the user-assigned identities provided
// by the user to be assigned to Azure resources.
type UserAssignedIdentity struct {
	// ProviderID is the identification ID of the user-assigned Identity, the format of an identity is:
	// 'azure:///subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{identityName}'
	ProviderID string `json:"providerID"`
}

// OSDisk defines the operating system disk for a VM.
type OSDisk struct {
	OSType      string      `json:"osType"`
	DiskSizeGB  int32       `json:"diskSizeGB"`
	ManagedDisk ManagedDisk `json:"managedDisk"`
}

// ManagedDisk defines the managed disk options for a VM.
type ManagedDisk struct {
	StorageAccountType string `json:"storageAccountType"`
}

// SubnetRole defines the unique role of a subnet.
type SubnetRole string

const (
	// SubnetNode defines a Kubernetes workload node role
	SubnetNode = SubnetRole(Node)

	// SubnetControlPlane defines a Kubernetes control plane node role
	SubnetControlPlane = SubnetRole(ControlPlane)
)

// SubnetSpec configures an Azure subnet.
type SubnetSpec struct {
	// Role defines the subnet role (eg. Node, ControlPlane)
	Role SubnetRole `json:"role,omitempty"`

	// ID defines a unique identifier to reference this resource.
	ID string `json:"id,omitempty"`

	// Name defines a name for the subnet resource.
	Name string `json:"name"`

	// CidrBlock is the CIDR block to be used when the provider creates a managed Vnet.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// InternalLBIPAddress is the IP address that will be used as the internal LB private IP.
	// For the control plane subnet only.
	InternalLBIPAddress string `json:"internalLBIPAddress,omitempty"`

	// SecurityGroup defines the NSG (network security group) that should be attached to this subnet.
	SecurityGroup SecurityGroup `json:"securityGroup,omitempty"`
}
