# expensify_employee

Provides a resource to create and manage a Expensify policy users.<br>


**Note** the following behaviors regarding employee arguments: <br>

* Updating the fields `manager_email`, `approves_to`, `over_limit_approver`, and `approval_limit` is meaningful only if Approval Mode for policy is Advanced Approval.<br>
* API not allow overwriting manually set values for `first_name` and `last_name` in their Expensify account.<br>
* The fields `first_name` and `last_name` are set at account level.<br>
* Once the value of any attribute is set, it cannot be set back to null through provider. But, you can set it to null via UI.<br>


## Example Usage

```
resource "expensify_employee" "employee"{
    employee_email = "employee@domain.com"
    manager_email = "manager@domain.com"
    policy_id = "22E95AFCD33ABE2BB8"
    employee_id = "1"
    first_name = "Dummy"
    last_name = "Employee"
    approves_to = "approver@domain.com"
    approval_limit = 5
    over_limit_approver = "overlimitapprover@domain.com"
}
```


## Argument Reference

The following arguments are supported:

* `employee_email` (Required, String) - The email address of the employee.
* `manager_email` (Required, String) - Manager email address.
* `policy_id` (Required, String) - The ID of policy for which employee is to be added.
* `first_name` (Optional, String) - First name of the employee in Expensify. 
* `last_name` (Optional, String) - Last name of the employee in Expensify.
* `employee_id` (Optional, String) - Unique ID of the Employee.
* `over_limit_approver` (Optional, String) - over limit approver email address. Required if `approval_limit` is specified.
* `approval_limit` (Optional, Float) - Specifies limit of report total.
* `approves_to` (Optional, String) - approver email address.


## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `is_terminated` (Bool) - Specifies whether employee is terminated or not
