package expensify

import(
	"time"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"terraform-provider-expensify/client"
)

func dataSourcePolicy() *schema.Resource{
	return &schema.Resource{
		ReadContext:   dataSourcePolicyRead,
		Schema: map[string]*schema.Schema{
			"policy_name":  &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"policy_id": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
			"plan": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"owner": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"output_currency": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	policyId := d.Get("policy_id").(string)
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		policy, err := apiClient.GetPolicy(policyId)
		if err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.Set("owner", policy.Owner)
		d.Set("policy_id", policy.PolicyId)
		d.Set("policy_name", policy.PolicyName)
		d.Set("plan", policy.Plan)
		d.Set("output_currency", policy.OutputCurrency)
		d.SetId(policyId)
		return nil
	})
	if retryErr!=nil {
		return diag.FromErr(retryErr)
	}
	return diags
}