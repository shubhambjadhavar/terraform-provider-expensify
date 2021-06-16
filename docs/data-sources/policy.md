# expensify_employee

Reads attributes of a Expensify policy.


## Example Usage

```
data "expensify_policy" "policy" {
    policy_id = "22E95AFCD33ABE2BB8"
}
```


## Argument Reference

The following arguments are supported:

* `policy_id` (Required, String) - The Policy ID.


## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `policy_name` (String) - Name of the policy.
* `plan` (String) - Defines the plan for the policy.
* `owner` (String) - Email address of policy owner.
* `output_currency` (String) - Currency for the policy.