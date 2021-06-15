package expensify

import (
	"os"
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExpensifyPolicyResource_import_basic(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccExpensifyProviderPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckExpensifyPolicyResourceImporterBasic(),
			},
			{
				ResourceName: "expensify_policy.policy",
				ImportState: true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckExpensifyPolicyResourceImporterBasic() string {
	return fmt.Sprintf(`
	resource "expensify_policy" "policy" {
		policy_name = "test"
	}`)
}