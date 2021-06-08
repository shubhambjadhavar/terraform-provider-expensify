package expensify

import(
	"os"
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExpensifyPolicyResource_basic(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccExpensifyProviderPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckExpensifyPolicyResourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("expensify_policy.policy", "policy_name", "testing"),
					resource.TestCheckResourceAttr("expensify_policy.policy", "plan", "corporate"),
				),
			},
		},
	})
}

func testAccCheckExpensifyPolicyResourceBasic() string {
	return fmt.Sprintf(`
	resource "expensify_policy" "policy" {
		policy_name = "testing"
    	        plan = "corporate"
	}`)
}
