# Terraform provider for VMware Cloud Director Availability

The official Terraform provider
for [VMware Cloud Director Availability](https://www.vmware.com/products/cloud-director-availability.html).

The VMware Cloud Director Availability provider can be used to perform initial configuration of either of the two
appliance roles:

* Cloud Director Replication Management Appliance
* vCenter Replication Management Appliance

The initial configuration includes changing the initial password of the **root** user of the appliance, adding external
Replicator Service instances and configuring the Tunnel Service.

* For information about the Terraform provider, see the [provider documentation](https://registry.terraform.io/providers/hashicorp/vcda/latest/docs).
* For information about VMware Cloud Director Availability,
  see the [ Product Page](https://www.vmware.com/products/cloud-director-availability.html).
* For more information,
  see the [VMware Cloud Director Availability Documentation](https://docs.vmware.com/en/VMware-Cloud-Director-Availability/index.html).

# Requirements

- Verify that you download and install [Terraform](https://www.terraform.io/downloads.html). For information about Terraform, see [What is Terraform?](https://developer.hashicorp.com/terraform/intro)
- Verify that to build the provider plugin you download and install [Go](https://golang.org/doc/install) 1.19.

# Building the Provider

Note: The following instructions apply only to Mac OS or Linux OS.

To work on the provider, you must first download and install [Go](http://www.golang.org) on your local computer. For more information, see the [requirements](#requirements) before proceeding.

First, clone the repository
to the `$GOPATH/src/github.com/vmware/terraform-provider-for-vmware-cloud-director-availability`
path or to a directory of your
choice:

```sh
mkdir -p $GOPATH/src/github.com/vmware
cd $GOPATH/src/github.com/vmware
git clone git@github.com:vmware/terraform-provider-for-vmware-cloud-director-availability.git
```

After the clone completes, navigate into the provider directory and build the provider by using either of the following two commands:

- By using the Makefile command:

  This command places the executable in the `$GOPATH/bin` directory.
    ```sh
    make build
    ```
  To find where your `$GOPATH` directory is located, run:
    ```sh
    go env GOPATH
    ```

- Alternatively, perform a manual build:

  This command places the executable in your current directory.
    ```sh
    cd $GOPATH/src/github.com/vmware/terraform-provider-for-vmware-cloud-director-availability
    go get
    go build -o terraform-provider-vcda
    ```

  After the build completes, if your terraform running directory does not match your `$GOPATH` environment, first you must copy
  the `terraform-provider-vcda` executable into your running directory. Then
  re-run `terraform init` to make the terraform aware of your local provider executable.

# Using the Provider

To use a released provider in your Terraform environment,
run the [`terraform init`](https://www.terraform.io/docs/commands/init.html) command and Terraform automatically installs the
provider. For information about specifying a specific provider version when installing released providers, see
the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

Alternatively, to use a custom-built provider in your Terraform environment, for example the provider binary from the build
instructions above, follow the instructions
to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-plugins). After placing the
custom-built provider into your plugins directory, to initialize it run `terraform init`.

For information about either installation method and for documentation about the provider-specific configuration options, see the [VMware Cloud Director Availability provider's website](https://www.terraform.io/docs/providers/vmware-cloud-director-availability/index.html).

# Automated Installation (Recommended)

To download and initialize the Terraform providers, including the VMware Cloud Director Availability provider, use the `terraform init` command.
In your `.tf` file, once you specify the provider block for the VMware Cloud Director Availability provider,
`terraform init` detects the need for the provider and downloads it to your environment.
To list the versions of the installed providers in your environment, run the `terraform version` command.

# Manual Installation

**NOTE:** Unless you require a pre-release version or you are [Developing the provider](#developing-the-provider), you must use the officially released
version of the provider, as per the above [Using the Provider](#using-the-provider) section.

After building the provider as per the [Building the Provider](#building-the-provider) section, perform the following:

Move the `terraform-provider-vcda` binary into the OS-specific directory on your system by using the following commands.
If the directory does not exist, create it under the `.terraform.d/plugins` directory -
this is where Terraform searches for the executable file.

- Linux:
    ```sh 
    $ mkdir -p ~/.terraform.d/plugins/terraform.example.com/vmware/vcda/1.0.0/linux_amd64
    $ mv terraform-provider-vcda ~/.terraform.d/plugins/terraform.example.com/vmware/vcda/1.0.0/linux_amd64
    ```

- Mac OS (Intel x86-64):
    ```sh 
    $ mkdir -p ~/.terraform.d/plugins/terraform.example.com/vmware/vcda/1.0.0/darwin_amd64
    $ mv terraform-provider-vcda ~/.terraform.d/plugins/terraform.example.com/vmware/vcda/1.0.0/darwin_amd64
    ```

- Mac OS (M1/M2 ARM64):
    ```sh 
    $ mkdir -p ~/.terraform.d/plugins/terraform.example.com/vmware/vcda/1.0.0/darwin_arm64
    $ mv terraform-provider-vcda ~/.terraform.d/plugins/terraform.example.com/vmware/vcda/1.0.0/darwin_arm64
    ```

Then navigate to a directory where the terraform scripts/files are located.

Ensure that the `vcda` provider path matches the above plugin path in your terraform modules:

```terraform
terraform {
  required_providers {
    vcda = {
      source  = "terraform.example.com/vmware/vcda"
      version = ">=1.0"
    }
  }
}
```

After placing the custom-built provider into your plugins directory, to initialize it run `terraform init`.
For information about the custom provider plugins,
see [Install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-plugins).

# Developing the Provider

**NOTE:**  To ensure that no
work is being duplicated, before you start working on a feature, check the
[Issue Tracker][gh-issues] and the existing [Pull Requests][gh-prs]. For further clarification, you can also ask in a
new issue.

[gh-issues]: https://github.com/vmware/terraform-provider-for-vmware-cloud-director-availability/issues

[gh-prs]: https://github.com/vmware/terraform-provider-for-vmware-cloud-director-availability/pulls

For more information, see the [Building the Provider](#building-the-provider) section and for instructions
about manually loading the provider for development, see the [Manual Installation](#manual-installation) section.

For information about how to execute acceptance tests, see the [Testing the Provider](#testing-the-provider) section.
Depending on your changeset, add new acceptance tests or modify the existing ones.
Ensure that all Acceptance tests pass.

# Testing the Provider

Set the required environment variables in `scripts/env_variables.sh` based on your infrastructure settings.

Before running acceptance tests, load the environment variables:

```sh
source ./scripts/env_variables.sh
```

**Running Acceptance tests:**

All Acceptance tests

```sh
make testacc TESTS="''"
```

Provider tests

```sh
make testacc TESTS=/provider
```

Cloud tests

```sh
make testacc TESTS=/cloud
```

Manager tests

```sh
make testacc TESTS=/manager
```

Data sources tests

```sh
make testacc TESTS=/datasource
```

**Note:** Acceptance tests create real resources and modify the existing infrastructure.

# License

Copyright 2023 VMware, Inc.

The Terraform provider for VMware Cloud Director Availability is available
under [MPL2.0 license](https://github.com/vmware/terraform-provider-for-vmware-cloud-director-availability/blob/main/LICENSE).
