# expensify_employee

Reads attributes of a User in Expensify policy.


## Example Usage

```
data "expensify_employee" "employee" {
    policy_id = "22E95AFCD33ABE2BB8"
    employee_email = "employee@domain.com" 
}
```

## Argument Reference

The following arguments are supported:

* `employee_email` (Required, String) - The email address of the employee.
* `policy_id` (Required, String) - The ID of policy in which employee is present.


## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `manager_email` (String) - Manager email address.
* `employee_id` (String) - Unique ID of the Employee.
* `over_limit_approver` (String) - over limit approver email address.
* `approval_limit` (Float) - Specifies limit of report total.
* `approves_to` (String) - approver email address.
* `role` (String) - Specifies the role of employee in the policy.