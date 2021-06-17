package expensify

import (
	"fmt"
	"time"
	"regexp"
	"context"
	"strings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"terraform-provider-expensify/client"
)

func validateEmail(v interface{}, k string) (warns []string, errs []error) {
	value := v.(string)
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !(emailRegex.MatchString(value)) {
		errs = append(errs, fmt.Errorf("Expected EmailId is not valid %s", k))
		return warns, errs
	}
	return
}

func resourceUser() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceUserImporter,
		},
		Schema: map[string]*schema.Schema{
			"employee_email": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				ValidateFunc: validateEmail,
			},
			"manager_email": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validateEmail,
			},
			"policy_id": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
			"first_name": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_name": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"employee_id": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"approval_limit": &schema.Schema{
				Type: schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"over_limit_approver": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validateEmail,
			},
			"is_terminated": &schema.Schema{
				Type: schema.TypeBool,
				Computed: true,
			},
			"approves_to": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validateEmail,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	employees := make([]client.Employee, 1)
	employees[0].EmployeeEmail = d.Get("employee_email").(string)
	employees[0].ManagerEmail = d.Get("manager_email").(string)
	employees[0].PolicyId = d.Get("policy_id").(string)
	employees[0].FirstName = d.Get("first_name").(string)
	employees[0].LastName = d.Get("last_name").(string)
	employees[0].EmployeeId = d.Get("employee_id").(string)
	employees[0].ApprovalLimit = d.Get("approval_limit").(float64)
	employees[0].OverLimitApprover = d.Get("over_limit_approver").(string)
	employees[0].ApprovesTo = d.Get("approves_to").(string)
	employeesList := client.EmployeesList{
		Employees: employees,
	}
	var err error
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if err = apiClient.NewEmployee(&employeesList); err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(employeesList.Employees[0].PolicyId + ":" + employeesList.Employees[0].EmployeeEmail)
	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	parts := resourceUserParseId(d.Id())
	employee := client.Employee{
		EmployeeEmail: parts[1],
		PolicyId: parts[0],
	}
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		body, err := apiClient.GetEmployee(&employee)
		if err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.Set("employee_email", body.EmployeeEmail)
		d.Set("manager_email", body.ManagerEmail)
		d.Set("employee_id", body.EmployeeId)
		d.Set("policy_id", body.PolicyId)
		d.Set("approves_to", body.ApprovesTo)
		d.Set("over_limit_approver", body.OverLimitApprover)
		d.Set("approval_limit", body.ApprovalLimit)
		d.Set("is_terminated", false)
		return nil
	})
	if retryErr!=nil {
		if strings.Contains(retryErr.Error(), "User Does Not Exist")==true {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Employee" + d.Id() +  "does not exist. Creating a new employee with given details.",
			})
			d.SetId("")
			return diags
		}
		return diag.FromErr(retryErr)
	}
	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	if d.HasChange("employee_email") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Employee not allowed to change employee_email",
			Detail:   "Employee not allowed to change employee_email",
		})
	}
	if d.HasChange("policy_id") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Employee not allowed to change policy_id",
			Detail:   "Employee not allowed to change policy_id",
		})
	}
	if diags.HasError() {
		return diags
	}
	parts := resourceUserParseId(d.Id())
	employees := make([]client.Employee, 1)
	employees[0].EmployeeEmail = parts[1]
	employees[0].ManagerEmail = d.Get("manager_email").(string)
	employees[0].PolicyId = parts[0]
	employees[0].FirstName = d.Get("first_name").(string)
	employees[0].LastName = d.Get("last_name").(string)
	employees[0].EmployeeId = d.Get("employee_id").(string)
	employees[0].ApprovalLimit = d.Get("approval_limit").(float64)
	employees[0].OverLimitApprover = d.Get("over_limit_approver").(string)
	employees[0].ApprovesTo = d.Get("approves_to").(string)
	employeesList := client.EmployeesList{
		Employees: employees,
	}
	var err error
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if err = apiClient.UpdateEmployee(&employeesList); err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	parts := resourceUserParseId(d.Id())
	employees := make([]client.Employee, 1)
	employees[0].EmployeeEmail = parts[1]
	employees[0].PolicyId = parts[0]
	employees[0].IsTerminated = true
	employeesList := client.EmployeesList{
		Employees: employees,
	}
	var err error
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if err = apiClient.DeleteEmployee(&employeesList); err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceUserImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	apiClient := m.(*client.Client)
	parts := resourceUserParseId(d.Id())
	employee := client.Employee{
		EmployeeEmail: parts[1],
		PolicyId: parts[0],
	}
	body, err := apiClient.GetEmployee(&employee)
	if err!=nil{
		return nil, err
	}
	d.Set("employee_email", body.EmployeeEmail)
	d.Set("manager_email", body.ManagerEmail)
	d.Set("employee_id", body.EmployeeId)
	d.Set("policy_id", body.PolicyId)
	d.Set("approves_to", body.ApprovesTo)
	d.Set("over_limit_approver", body.OverLimitApprover)
	d.Set("approval_limit", body.ApprovalLimit)
	d.Set("is_terminated", false)
	return []*schema.ResourceData{d}, nil
}

func resourceUserParseId(id string) ([]string) {
	parts := strings.Split(id, ":")
  	return parts
}