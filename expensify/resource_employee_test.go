package expensify

import(
	"os"
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExpensifyEmployeeResource_basic(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccExpensifyProviderPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckExpensifyEmployeeResourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("expensify_employee.employee", "employee_email", "abhishiek@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "manager_email", "shubham@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "policy_id", "56B042862350ADD2"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "employee_id", "1003"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "first_name", "Abhishiek"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "last_name", "Singh"),
				),
			},
		},
	})
}

func testAccCheckExpensifyEmployeeResourceBasic() string {
	return fmt.Sprintf(`
	resource "expensify_employee" "employee" {
		employee_email = "abhishiek@clevertapdemo.ml"
    		manager_email = "shubham@clevertapdemo.ml"
    		policy_id = "56B042862350ADD2"
    		employee_id = "1003"
    		first_name = "Abhishiek"
    		last_name = "Singh"
	}`)
}


func TestAccExpensifyEmployeeResource_update(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccExpensifyProviderPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckExpensifyEmployeeResourceUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("expensify_employee.employee", "employee_email", "abhishiek@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "manager_email", "shubham@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "policy_id", "56B042862350ADD2"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "employee_id", "1003"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "first_name", "Abhishiek"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "last_name", "Singh"),	
				),
			},
			{
				Config: testAccCheckExpensifyEmployeeResourceUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("expensify_employee.employee", "employee_email", "abhishiek@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "manager_email", "ashutosh@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "policy_id", "56B042862350ADD2"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "employee_id", "1003"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "first_name", "Abhishiek"),
					resource.TestCheckResourceAttr("expensify_employee.employee", "last_name", "Singh Delhi"),
				),
			},
		},
	})
}

func testAccCheckExpensifyEmployeeResourceUpdatePre() string {
	return fmt.Sprintf(`
	resource "expensify_employee" "employee" {
		employee_email = "abhishiek@clevertapdemo.ml"
    		manager_email = "shubham@clevertapdemo.ml"
    		policy_id = "56B042862350ADD2"
    		employee_id = "1003"
    		first_name = "Abhishiek"
    		last_name = "Singh"
	}`)
}

func testAccCheckExpensifyEmployeeResourceUpdatePost() string {
	return fmt.Sprintf(`
	resource "expensify_employee" "employee" {
		employee_email = "abhishiek@clevertapdemo.ml"
    		manager_email = "ashutosh@clevertapdemo.ml"
    		policy_id = "56B042862350ADD2"
    		employee_id = "1003"
    		first_name = "Abhishiek"
    		last_name = "Singh Delhi"
	}`)
}