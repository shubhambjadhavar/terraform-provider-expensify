# expensify_policy

Provides a resource to create and read a Expensify policies.<br>


## Example Usage

```
resource "expensify_policy" "policy"{
    policy_name = "demo"
    plan = "corporate"
}
```


## Argument Reference

The following arguments are supported:

* `policy_name` (Required, String) - Name of the policy.
* `plan` (Optional, String) - Defines the plan for the policy. Supported values are `team` (Collect) and `corporate` (Control). Default value is `team`.


## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `policy_id` (String) - The Policy ID.
* `owner` (String) - Email address of policy owner.
* `output_currency` (String) - Currency for the policy.