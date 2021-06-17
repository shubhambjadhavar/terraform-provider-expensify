package expensify

import (
	"os"
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExpensifyUserResource_import_basic(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccExpensifyProviderPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckExpensifyUserResourceImporterBasic(),
			},
			{
				ResourceName: "expensify_user.employee",
				ImportState: true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckExpensifyUserResourceImporterBasic() string {
	return fmt.Sprintf(`
	resource "expensify_user" "employee" {
		employee_email = "abhishiek@clevertapdemo.ml"
    		manager_email = "shubham@clevertapdemo.ml"
    		policy_id = "56B042862350ADD2"
    		employee_id = "1003"
	}`)
}
