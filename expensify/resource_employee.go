package expensify

import (
	"fmt"
	"regexp"
	"context"
	"strings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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

func resourceEmployee() *schema.Resource{
	return &schema.Resource{
		CreateContext: resourceEmployeeCreate,
		ReadContext:   resourceEmployeeRead,
		UpdateContext: resourceEmployeeUpdate,
		DeleteContext: resourceEmployeeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Optional: true,
				Default: false,
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

func resourceEmployeeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
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
	employees[0].IsTerminated = d.Get("is_terminated").(bool)
	employees[0].ApprovesTo = d.Get("approves_to").(string)
	employeesList := client.EmployeesList{
		Employees: employees,
	}
	err := apiClient.NewEmployee(&employeesList)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(employeesList.Employees[0].PolicyId + ":" + employeesList.Employees[0].EmployeeEmail)
	return diags
}

func resourceEmployeeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	parts := resourceEmployeeParseId(d.Id())
	employee := client.Employee{
		EmployeeEmail: parts[1],
		PolicyId: parts[0],
	}
	if isTerminated := d.Get("is_terminated").(bool); isTerminated==true {
		employees := make([]client.Employee, 1)
		employees[0].EmployeeEmail = parts[1]
		employees[0].ManagerEmail = d.Get("manager_email").(string)
		employees[0].PolicyId = parts[0]
		employees[0].FirstName = d.Get("first_name").(string)
		employees[0].LastName = d.Get("last_name").(string)
		employees[0].EmployeeId = d.Get("employee_id").(string)
		employees[0].ApprovalLimit = d.Get("approval_limit").(float64)
		employees[0].OverLimitApprover = d.Get("over_limit_approver").(string)
		employees[0].IsTerminated = false
		employees[0].ApprovesTo = d.Get("approves_to").(string)
		employeesList := client.EmployeesList{
			Employees: employees,
		}
		err := apiClient.ActivateEmployee(&employeesList)
		if err!=nil{
			return diag.FromErr(err)
		}
	}
	body, err := apiClient.GetEmployee(&employee)
	if err!=nil{
		if strings.Contains(err.Error(), "\"responseCode\":404")==true {
			d.SetId("")
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Employee does not exist. Creating a new employee with given details.",
			})
			return diags
		}
		return diag.FromErr(err)
	}
	d.Set("employee_email", body.EmployeeEmail)
	d.Set("manager_email", body.ManagerEmail)
	d.Set("employee_id", body.EmployeeId)
	d.Set("policy_id", body.PolicyId)
	d.Set("approves_to", body.ApprovesTo)
	d.Set("over_limit_approver", body.OverLimitApprover)
	d.Set("approval_limit", body.ApprovalLimit)
	d.Set("is_terminated", false)
	return diags
}

func resourceEmployeeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
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
	parts := resourceEmployeeParseId(d.Id())
	employees := make([]client.Employee, 1)
	employees[0].EmployeeEmail = parts[1]
	employees[0].ManagerEmail = d.Get("manager_email").(string)
	employees[0].PolicyId = parts[0]
	employees[0].FirstName = d.Get("first_name").(string)
	employees[0].LastName = d.Get("last_name").(string)
	employees[0].EmployeeId = d.Get("employee_id").(string)
	employees[0].ApprovalLimit = d.Get("approval_limit").(float64)
	employees[0].OverLimitApprover = d.Get("over_limit_approver").(string)
	employees[0].IsTerminated = d.Get("is_terminated").(bool)
	employees[0].ApprovesTo = d.Get("approves_to").(string)
	employeesList := client.EmployeesList{
		Employees: employees,
	}
	err := apiClient.UpdateEmployee(&employeesList)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceEmployeeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics{
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	parts := resourceEmployeeParseId(d.Id())
	employees := make([]client.Employee, 1)
	employees[0].EmployeeEmail = parts[1]
	employees[0].ManagerEmail = d.Get("manager_email").(string)
	employees[0].PolicyId = parts[0]
	employees[0].EmployeeId = d.Get("employee_id").(string)
	employees[0].FirstName = d.Get("first_name").(string)
	employees[0].LastName = d.Get("last_name").(string)
	employees[0].IsTerminated = true
	employees[0].ApprovalLimit = d.Get("approval_limit").(float64)
	employees[0].OverLimitApprover = d.Get("over_limit_approver").(string)
	employees[0].ApprovesTo = d.Get("approves_to").(string)
	employeesList := client.EmployeesList{
		Employees: employees,
	}
	err := apiClient.DeleteEmployee(&employeesList)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceEmployeeParseId(id string) ([]string) {
	parts := strings.Split(id, ":")
  	return parts
}