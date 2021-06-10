This Terraform provider enables create, read, update, delete, and import operations for Expensify policy users.


## Requirements

* [Go](https://golang.org/doc/install) >= 1.16 (To build the provider plugin)<br>
* [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x <br/>
* Application: [Expensify](https://www.expensify.com/) (API is supported in collect and control policy plans)
* [API Documentation](https://integrations.expensify.com/Integration-Server/doc/employeeUpdater.html)


## Application Account

### Setup<a id="setup"></a>
1. Create an expensify account at https://www.expensify.com/<br>
2. To create a policy, go to `Settings -> Policies -> Group -> click on New Policy`.<br>
3. After creating the policy, for policy ID, go to `Settings -> Policies -> Group -> Select the appropriate policy` and note the policy ID from the URL.<br>
   For example, in Policy url - ```"https://www.expensify.com/policy?param={policyID:22E95AFCD33ABE2BB8}", "22E95AFCD33ABE2BB8" is Policy ID```

### API Authentication
1. To authenticate API, we need a pair of credentials: partnerUserID and partnerUserSecret.<br>
2. For this, go to https://www.expensify.com/tools/integrations/ and generate the credentials.<br>
3. A pair of credentials: partnerUserID and partnerUserSecret will be generated and shown on the page.<br>


## Building The Provider
1. Clone the repository, add all the dependencies and create a vendor directory that contains all dependencies. For this, run the following commands: <br>
```
cd terraform-provider-expensify
go mod init terraform-provider-expensify
go mod tidy
go mod vendor
```

## Managing terraform plugins
*For Windows:*
1. Run the following command to create a vendor sub-directory (`%APPDATA%/terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${OS_ARCH}`) which will consist of all terraform plugins. <br> 
Command: 
```bash
mkdir -p %APPDATA%/terraform.d/plugins/expensify.com/employee/expensify/1.0.0/windows_amd64
```
2. Run `go build -o terraform-provider-expensify.exe` to generate the binary in present working directory. <br>
3. Run this command to move this binary file to the appropriate location.
 ```
 move terraform-provider-expensify.exe %APPDATA%\terraform.d\plugins\expensify.com\employee\expensify\1.0.0\windows_amd64
 ``` 
 <p align="center">
 [OR]
 </p>
 
3. Manually move the file from current directory to destination directory (`%APPDATA%\terraform.d\plugins\expensify.com\employee\expensify\1.0.0\windows_amd64`).<br>


## Working with terraform

### Application Credential Integration in terraform
1. Add `terraform` block and `provider` block as shown in [example usage](#example-usage).
2. Get a pair of credentials: partnerUserID and partnerUserSecret. For this, visit https://www.expensify.com/tools/integrations/.
3. Assign the above credentials to the respective field in the `provider` block.

### Basic Terraform Commands
1. `terraform init` - To initialize a working directory containing Terraform configuration files.
2. `terraform plan` - To create an execution plan. Displays the changes to be done.
3. `terraform apply` - To execute the actions proposed in a Terraform plan. Apply the changes.

### Create User
1. Add the `employee_email`, `manager_email`, `policy_id`, `first_name`, `last_name` in the respective field in `resource` block as shown in [example usage](#example-usage).
2. Refer to [setup](#setup) for the policy ID.
3. Run the basic terraform commands.<br>
4. On successful execution, sends an account setup mail to user.<br>

### Update the user
1. Update the data of the user in the `resource` block as show in [example usage](#example-usage) and run the basic terraform commands to update user. 
   User is not allowed to update `employee_email` and `policy_id`.
2. To remove a user from a policy set `is_terminated` field in `resource` block to `true` and vice versa to re-add.

### Read the User Data
Add `data` and `output` blocks as shown in the [example usage](#example-usage) and run the basic terraform commands.

### Delete the user
Delete the `resource` block of the user and run `terraform apply`.

### Import a User Data
1. Write manually a `resource` configuration block for the user as shown in [example usage](#example-usage). Imported user will be mapped to this block.
2. Run the command `terraform import expensify_employee.employee [POLICY_ID]:[EMAIL_ID]` to import user.
3. Refer to [setup](#setup) for the policy ID.
4. Run `terraform plan`, if output shows `0 to addd, 0 to change and 0 to destroy` user import is successful, otherwise recheck the employee data in `resource` block with employee data in the policy in Expensify website. 


## Example Usage<a id="example-usage"></a>

```
terraform{
    required_providers {
        expensify = {
            version = "1.0.0"
            source = "expensify.com/employee/expensify"
        }
    }
}

provider "expensify" {
    partner_user_id = "_REPLACE_PARTNER_USER_ID_"
    partner_user_secret = "_REPLACE_PARTNER_USER_SECRET_" 
}

resource "expensify_employee" "employee"{
    employee_email = "employee@domain.com"
    manager_email = "manager@domain.com"
    policy_id = "22E95AFCD33ABE2BB8"
    employee_id = "101"
    first_name = "Dummy"
    last_name = "Employee"
    approves_to = "approver@domain.com"
    approval_limit = 5
    over_limit_approver = "overlimitapprover@domain.com"
    is_terminated = false
}

output "resource_employee"{
    value = expensify_employee.employee
}

data "expensify_employee" "employee" {
    policy_id = "22E95AFCD33ABE2BB8"
    employee_email = "employee@domain.com" 
}

output "datasouce_employee"{
    value = data.expensify_employee.employee
}
```


## Argument Reference

* `partner_user_id` (Optional, String) - The Expensify Partner User ID. This may also be set via the `"PARTNER_USER_ID"` environment variable.
* `partner_user_secret` (Optional, String) - The Expensify Partner User Secret. This may also be set via the `"PARTNER_USER_SECRET"` environment variable.
* `employee_email` (Required, String) - The email address of the employee.
* `manager_email` (Required, String) - Manager email address.
* `policy_id` (Required, String) - The ID of policy for which employee is to be added.
* `first_name` (Optional, String) - First name of the employee in Expensify. 
* `last_name` (Optional, String) - Last name of the employee in Expensify. 
* `is_terminated` (Optional, Boolean) - If set to true, the employee will be removed from the policy.
* `employee_id` (Optional, String) - Unique ID of the Employee.
* `over_limit_approver` (Optional, String) - over limit approver email address. Required if an `approval_limit` is specified.
* `approval_limit` (Optional, Float) - Specifies limit of report total.
* `approves_to` (Optional, String) - approver email address.


## Exceptions

* Updating of the fields `manager_email`, `approves_to`, `over_limit_approver`, and `approval_limit` is meaningful only if Approval Mode for policy is Advanced Approval.
* Updating `first_name` and `last_name` in any one policy will automatically update them in other policies.
* Not allowed overwriting `first_name` and `last_name` values manually set by the employee in their Expensify account.
* To add an employee to multiple policies, write multiple `resource` block with different policy ID.
