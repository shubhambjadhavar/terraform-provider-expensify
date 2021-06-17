package expensify

import(
	"os"
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExpensifyUserResource_basic(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccExpensifyProviderPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckExpensifyUserResourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("expensify_user.employee", "employee_email", "abhishiek@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("expensify_user.employee", "manager_email", "shubham@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("expensify_user.employee", "policy_id", "56B042862350ADD2"),
					resource.TestCheckResourceAttr("expensify_user.employee", "employee_id", "1003"),
					resource.TestCheckResourceAttr("expensify_user.employee", "first_name", "Abhishiek"),
					resource.TestCheckResourceAttr("expensify_user.employee", "last_name", "Singh"),
				),
			},
		},
	})
}

func testAccCheckExpensifyUserResourceBasic() string {
	return fmt.Sprintf(`
	resource "expensify_user" "employee" {
		employee_email = "abhishiek@clevertapdemo.ml"
    		manager_email = "shubham@clevertapdemo.ml"
    		policy_id = "56B042862350ADD2"
    		employee_id = "1003"
    		first_name = "Abhishiek"
    		last_name = "Singh"
	}`)
}


func TestAccExpensifyUserResource_update(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccExpensifyProviderPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckExpensifyUserResourceUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("expensify_user.employee", "employee_email", "abhishiek@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("expensify_user.employee", "manager_email", "shubham@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("expensify_user.employee", "policy_id", "56B042862350ADD2"),
					resource.TestCheckResourceAttr("expensify_user.employee", "employee_id", "1003"),
					resource.TestCheckResourceAttr("expensify_user.employee", "first_name", "Abhishiek"),
					resource.TestCheckResourceAttr("expensify_user.employee", "last_name", "Singh"),	
				),
			},
			{
				Config: testAccCheckExpensifyUserResourceUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("expensify_user.employee", "employee_email", "abhishiek@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("expensify_user.employee", "manager_email", "ashutosh@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("expensify_user.employee", "policy_id", "56B042862350ADD2"),
					resource.TestCheckResourceAttr("expensify_user.employee", "employee_id", "1003"),
					resource.TestCheckResourceAttr("expensify_user.employee", "first_name", "Abhishiek"),
					resource.TestCheckResourceAttr("expensify_user.employee", "last_name", "Singh Delhi"),
				),
			},
		},
	})
}

func testAccCheckExpensifyUserResourceUpdatePre() string {
	return fmt.Sprintf(`
	resource "expensify_user" "employee" {
		employee_email = "abhishiek@clevertapdemo.ml"
    		manager_email = "shubham@clevertapdemo.ml"
    		policy_id = "56B042862350ADD2"
    		employee_id = "1003"
    		first_name = "Abhishiek"
    		last_name = "Singh"
	}`)
}

func testAccCheckExpensifyUserResourceUpdatePost() string {
	return fmt.Sprintf(`
	resource "expensify_user" "employee" {
		employee_email = "abhishiek@clevertapdemo.ml"
    		manager_email = "ashutosh@clevertapdemo.ml"
    		policy_id = "56B042862350ADD2"
    		employee_id = "1003"
    		first_name = "Abhishiek"
    		last_name = "Singh Delhi"
	}`)
}