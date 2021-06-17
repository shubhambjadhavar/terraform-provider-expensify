package expensify

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"terraform-provider-expensify/client"
)

func dataSourceUser() *schema.Resource{
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
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

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	employee := client.Employee{
		EmployeeEmail: d.Get("employee_email").(string),
		PolicyId: d.Get("policy_id").(string),
	}
	body, err := apiClient.GetEmployee(&employee)
	if err!=nil{
		return diag.FromErr(err)
	}
	d.Set("role", body.Role)
	d.Set("manager_email", body.ManagerEmail)
	d.Set("employee_id", body.EmployeeId)
	d.Set("approves_to", body.ApprovesTo)
	d.Set("over_limit_approver", body.OverLimitApprover)
	d.Set("approval_limit", body.ApprovalLimit)
	d.SetId(body.PolicyId + ":" + body.EmployeeEmail)
	return diags
}