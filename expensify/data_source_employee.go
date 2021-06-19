package expensify

import (
	"time"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"terraform-provider-expensify/client"
)

func dataSourceEmployee() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceEmployeeRead,
		Schema: map[string]*schema.Schema{
			"policy_id": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
			"employee_email": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				ValidateFunc: validateEmail,
			},
			"role": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"manager_email": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"employee_id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"approval_limit": &schema.Schema{
				Type: schema.TypeFloat,
				Computed: true,
			},
			"over_limit_approver": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"approves_to": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEmployeeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	employee := client.Employee{
		EmployeeEmail: d.Get("employee_email").(string),
		PolicyId: d.Get("policy_id").(string),
	}
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		body, err := apiClient.GetEmployee(&employee)
		if err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.Set("role", body.Role)
		d.Set("manager_email", body.ManagerEmail)
		d.Set("employee_id", body.EmployeeId)
		d.Set("approves_to", body.ApprovesTo)
		d.Set("over_limit_approver", body.OverLimitApprover)
		d.Set("approval_limit", body.ApprovalLimit)
		d.SetId(body.PolicyId + ":" + body.EmployeeEmail)
		return nil
	})
	if retryErr!=nil {
		return diag.FromErr(retryErr)
	}
	return diags
}