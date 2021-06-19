package expensify

import (
	"fmt"
	"time"
	"strings"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		DeleteContext: resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyImporter,
		},
		Schema: map[string]*schema.Schema{
			"policy_name":  &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_id": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"plan": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Default: "team",
				ForceNew: true,
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
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		policyId, err := apiClient.NewPolicy(policyName, plan)
		if err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.SetId(policyId)
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	return diags
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	policyId := d.Id()
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
		return nil
	})
	if retryErr!=nil {
		if strings.Contains(retryErr.Error(), "not exist")==true {
			d.SetId("")
			return diags
		}
		return diag.FromErr(retryErr)
	}	
	return diags
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}

func resourcePolicyImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	apiClient := m.(*client.Client)
	policyId := d.Id()
	policy, err := apiClient.GetPolicy(policyId)
	if err != nil {
		return nil, err
	}
	d.Set("owner", policy.Owner)
	d.Set("policy_id", policy.PolicyId)
	d.Set("policy_name", policy.PolicyName)
	d.Set("plan", policy.Plan)
	d.Set("output_currency", policy.OutputCurrency)
	return []*schema.ResourceData{d}, nil
}