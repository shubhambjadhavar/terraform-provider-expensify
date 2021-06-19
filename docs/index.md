# Expensify Provider

This Expensify provider allows Terraform to add users to the Expensify policies, read users of the Expensify policies, update users in the Expensify policies, remove users from the Expensify policies, create a new policy, read a policy details.<br>


## Relevant Expensify API Documentation

* https://integrations.expensify.com/Integration-Server/doc/#authentication
* https://integrations.expensify.com/Integration-Server/doc/#policy-creator
* https://integrations.expensify.com/Integration-Server/doc/#policy-list-getter
* https://integrations.expensify.com/Integration-Server/doc/employeeUpdater.html
* https://integrations.expensify.com/Integration-Server/doc/#policy-getter


## API Authentication

*Generate credentials from account which is domain and policy admin*
1. To authenticate API, we need a pair of credentials: partnerUserID and partnerUserSecret.<br>
2. For this, go to https://www.expensify.com/tools/integrations/ and generate the credentials.<br>
3. A pair of credentials: partnerUserID and partnerUserSecret will be generated and shown on the page.<br>


## Provider Arguments

The provider configuration block accepts the following arguments. In most cases it is recommended to set them via the indicated environment variables in order to keep credential information out of the configuration.<br>

* `partner_user_id` (Required, String) - The Expensify Partner User ID. This can also be set via the `"EXPENSIFY_PARTNER_USER_ID"` environment variable.
* `partner_user_secret` (Required, String) - The Expensify Partner User Secret. This can also be set via the `"EXPENSIFY_PARTNER_USER_SECRET"` environment variable.


## Example Usage

```
provider "expensify" {
    partner_user_id = "_REPLACE_PARTNER_USER_ID_"
    partner_user_secret = "_REPLACE_PARTNER_USER_SECRET_" 
}

resource "expensify_policy" "policy"{
    policy_name = "demo"
    plan = "corporate"
}

data "expensify_policy" "policy" {
    policy_id = "22E95AFCD33ABE2BB8"
}

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

data "expensify_employee" "employee" {
    policy_id = "22E95AFCD33ABE2BB8"
    employee_email = "employee@domain.com" 
}
```