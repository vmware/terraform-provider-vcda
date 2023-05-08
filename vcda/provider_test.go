/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcda

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider

var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"vcda": testAccProvider,
	}
}

func (at *AccTests) TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func (at *AccTests) TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testProviders() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"vcda": func() (*schema.Provider, error) { return Provider(), nil },
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(VcdaIP); v == "" {
		t.Fatal(VcdaIP + " must be set for acceptance tests")
	}
	if v := os.Getenv(LocalUser); v == "" {
		t.Fatal(LocalUser + " must be set for acceptance tests")
	}
	if v := os.Getenv(LocalPassword); v == "" {
		t.Fatal(LocalPassword + " must be set for acceptance tests")
	}
	if v := os.Getenv(VsphereUser); v == "" {
		t.Fatal(VsphereUser + " must be set for acceptance tests")
	}
	if v := os.Getenv(VspherePassword); v == "" {
		t.Fatal(VspherePassword + " must be set for acceptance tests")
	}
	if v := os.Getenv(VsphereServer); v == "" {
		t.Fatal(VsphereServer + " must be set for acceptance tests")
	}
	if v := os.Getenv(DatacenterID); v == "" {
		t.Fatal(DatacenterID + " must be set for acceptance tests")
	}
}

type AccTests struct{ Test *testing.T }

func TestRunner(t *testing.T) {
	t.Run("provider", func(t *testing.T) {
		test := AccTests{Test: t}
		test.TestProvider(t)
		test.TestProvider_impl(t)
	})

	t.Run("cloud", func(t *testing.T) {
		test := AccTests{Test: t}
		test.TestAccVcdaAppliancePassword_basic(t)
		test.TestAccVcdaCloudDirectorReplicationManager_basic(t)
		test.TestAccVcdaTunnel_basic(t)
	})

	t.Run("manager", func(t *testing.T) {
		test := AccTests{Test: t}
		test.TestAccVcdaVcenterReplicationManager_basic(t)
		test.TestAccVcdaReplicator_basic(t)
	})

	t.Run("datasource", func(t *testing.T) {
		test := AccTests{Test: t}
		test.TestAccVcdaDataSourceRemoteServicesThumbprint_basic(t)
		test.TestAccVcdaDataSourceServiceCert_basic(t)
	})
}
