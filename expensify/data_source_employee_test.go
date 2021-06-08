package expensify

import (
	"os"
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExpensifyEmployeeDataSource_basic(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccExpensifyProviderPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccExpensifyEmployeeDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.expensify_employee.employee", "employee_email", "shubham@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("data.expensify_employee.employee", "policy_id", "E95AFCD33ABE2BB8"),
					resource.TestCheckResourceAttr("data.expensify_employee.employee", "manager_email", "shubham@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("data.expensify_employee.employee", "role", "admin"),
				),
			},
		},
	})
}

func testAccExpensifyEmployeeDataSourceConfig() string {
	return fmt.Sprintf(`
	data "expensify_employee" "employee" {
		policy_id = "E95AFCD33ABE2BB8"
    	        employee_email = "shubham@clevertapdemo.ml"
	}`)
}
