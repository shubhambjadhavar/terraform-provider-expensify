package expensify

import (
	"fmt"
	"strings"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"terraform-provider-expensify/client"
)

func validatePlan(v interface{}, k string) (warns []string, errs []error) {
	value := v.(string)
	plans := []string{"team","corporate"} 
	for _, t := range plans{
		if t==value{
			return
		}
	}
	errs = append(errs, fmt.Errorf("%q must be either \"team\" or \"corporate\"", k))
	return
}

func resourcePolicy() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_name":  &schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
			"policy_id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"plan": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Default: "team",
				ValidateFunc: validatePlan,
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

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	policyName := d.Get("policy_name").(string)
	plan := d.Get("plan").(string)
	policyId, err := apiClient.NewPolicy(policyName, plan)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(policyId)
	return diags
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	policyId := d.Id()
	policy, err := apiClient.GetPolicy(policyId)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			d.SetId("")
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Policy does not exist. Creating a new policy with given details.",
			})
			return diags
		}
		return diag.FromErr(err)
	}
	d.Set("owner", policy.Owner)
	d.Set("policy_id", policy.PolicyId)
	d.Set("policy_name", policy.PolicyName)
	d.Set("plan", policy.Plan)
	d.Set("output_currency", policy.OutputCurrency)
	return diags
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
	var diags diag.Diagnostics
	return diags
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}