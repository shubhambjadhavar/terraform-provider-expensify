package expensify

import (
	"os"
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccExpensifyPolicyDataSource_basic(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccExpensifyProviderPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccExpensifyPolicyDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.expensify_policy.policy", "policy_name", "shubhamj"),
					resource.TestCheckResourceAttr("data.expensify_policy.policy", "policy_id", "56B042862350ADD2"),
					resource.TestCheckResourceAttr("data.expensify_policy.policy", "plan", "corporate"),
					resource.TestCheckResourceAttr("data.expensify_policy.policy", "owner", "shubhamj@clevertapdemo.ml"),
					resource.TestCheckResourceAttr("data.expensify_policy.policy", "output_currency", "INR"),
				),
			},
		},
	})
}

func testAccExpensifyPolicyDataSourceConfig() string {
	return fmt.Sprintf(`
	data "expensify_policy" "policy" {
		policy_id = "56B042862350ADD2"
	}`)
}
