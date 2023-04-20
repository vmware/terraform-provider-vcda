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

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
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
