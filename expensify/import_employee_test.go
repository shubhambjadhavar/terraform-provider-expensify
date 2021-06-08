package expensify

import (
	"os"
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExpensifyEmployeeResource_import_basic(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccExpensifyProviderPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckExpensifyEmployeeResourceImporterBasic(),
			},
			{
				ResourceName: "expensify_employee.employee",
				ImportState: true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckExpensifyEmployeeResourceImporterBasic() string {
	return fmt.Sprintf(`
	resource "expensify_employee" "employee" {
		employee_email = "abhishiek@clevertapdemo.ml"
    	        manager_email = "shubham@clevertapdemo.ml"
    	        policy_id = "E95AFCD33ABE2BB8"
    	        employee_id = "1003"
	}`)
}
