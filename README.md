# Terraform provider for VMware Cloud Director Availability

The official Terraform provider
for [VMware Cloud Director Availability](https://www.vmware.com/products/cloud-director-availability.html).

This provider allows you to perform initial configuration of Cloud Director Replication Manager or vCenter Replication
Manager services
including changing appliance password, adding external replicator/s and setting tunneling service.

Learn more:

* Read the provider [documentation](https://registry.terraform.io/providers/hashicorp/vcda/latest/docs).

# Requirements

- [Terraform](https://www.terraform.io/downloads.html)
- [Go](https://golang.org/doc/install) 1.19 (to build the provider plugin)

# Building the Provider

The instructions outlined below are specific to Mac OS or Linux OS only.

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (please
check
the [requirements](https://github.com/vmware/terraform-provider-for-vmware-cloud-director-availability#requirements)
before proceeding).

First, you will want to clone the repository
to : `$GOPATH/src/github.com/vmware/terraform-provider-for-vmware-cloud-director-availability`

```sh
mkdir -p $GOPATH/src/github.com/vmware
cd $GOPATH/src/github.com/vmware
git clone git@github.com:vmware/terraform-provider-for-vmware-cloud-director-availability.git
```

After the clone is complete, you can enter the provider directory and build the provider.

```sh
cd $GOPATH/src/github.com/vmware/terraform-provider-for-vmware-cloud-director-availability
go get
go build -o terraform-provider-for-vmware-cloud-director-availability
```

After the build is complete, if your terraform running folder does not match your GOPATH environment, you need to copy
the `terraform-provider-for-vmware-cloud-director-availability` executable to your running folder and
re-run `terraform init` to make terraform aware of your local provider executable.

# Using the Provider

To use a released provider in your Terraform environment,
run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the
provider. To specify a particular provider version when installing released providers, see
the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

To instead use a custom-built provider in your Terraform environment (e.g. the provider binary from the build
instructions above), follow the instructions
to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) After placing the
custom-built provider into your plugins directory, run `terraform init` to initialize it.

For either installation method, documentation about the provider specific configuration options can be found on
the [provider's website](https://www.terraform.io/docs/providers/vmware-cloud-director-availability/index.html).

## Controlling the provider version

Note that you can also control the provider version. This requires the use of a
`provider` block in your Terraform configuration if you have not added one
already.

The syntax is as follows:

```sh
provider "vcda" {
  version = "~> 1.0"
  ...
}
```

Version locking uses a pessimistic operator, so this version lock would mean
anything within the 1.x namespace, including or after 1.0.0. [Read
more](https://www.terraform.io/docs/configuration/providers.html#provider-versions) on provider version control.

# Automated Installation (Recommended)

Download and initialization of Terraform providers is with the “terraform init” command. This applies to the VCDA
provider as well. Once the provider block for the VCDA provider is specified in your .tf file, “terraform init” will
detect a need for the provider and download it to your environment.
You can list versions of providers installed in your environment by running “terraform version” command:

```sh
$ terraform version
Terraform v0.12.20
+ provider.vcda (unversioned)
```

# Manual Installation

**NOTE:** Unless you are [developing](#developing-the-provider) or require a
pre-release bugfix or feature, you will want to use the officially released
version of the provider (see [the section above](#using-the-provider)).

**NOTE:** Note that if the provider is manually copied to your running folder (rather than fetched with the “terraform
init” based on provider block), Terraform is not aware of the version of the provider you’re running. It will appear as
“unversioned”:

```sh
$ terraform version
Terraform v0.12.20
+ provider.vcda (unversioned)
```

Since Terraform has no indication of version, it cannot upgrade in a native way, based on the “version” attribute in
provider block.
In addition, this may cause difficulties in housekeeping and issue reporting.

# Developing the Provider

**NOTE:** Before you start work on a feature, please make sure to check the
[issue tracker][gh-issues] and existing [pull requests][gh-prs] to ensure that
work is not being duplicated. For further clarification, you can also ask in a
new issue.

[gh-issues]: https://github.com/vmware/terraform-provider-for-vmware-cloud-director-availability/issues

[gh-prs]: https://github.com/vmware/terraform-provider-for-vmware-cloud-director-availability/pulls

See [the section above](#building-the-provider) for details on building the
provider.

# Testing the Provider

Set required environment variables based as per your infrastructure settings

```sh
$ # VCDA appliance management IP and credentials 
$ export VCDA_IP=xxx
$ export LOCAL_USER=xxx
$ export LOCAL_PASSWORD=xxxx
$ # vSphere server and credentials
$ export VSPHERE_SERVER=xxx
$ export VSPHERE_USER=xxx
$ export VSPHERE_PASSWORD=xxx
$ # The managed object ID of the datacenter in which the appliances are placed 
$ export DC_ID=xxx
```

*Note:* Acceptance tests create real resources and modifies the existing infrastructure.

# License

Copyright 2023 VMware, Inc.

The Terraform provider for VMware Cloud Director Availability is available
under [MPL2.0 license](https://github.com/vmware/terraform-provider-for-vmware-cloud-director-availability/blob/main/LICENSE).
