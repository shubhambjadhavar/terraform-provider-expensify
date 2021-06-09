This Terraform provider enables create, read, update, delete, and import operations for Expensify policy users.


## Requirements

* [Go](https://golang.org/doc/install) >= 1.16 (To build the provider plugin)<br>
* [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x <br/>
* Application: [Expensify](https://www.expensify.com/) (API support in policies with collect or control plans)


## Application Account

### Setup
1. Create an expensify account at https://www.expensify.com/<br>
2. Sign in to the expensify account.<br>
3. To create a policy, go to `Settings -> Policies -> Group -> click on New Policy`.<br>
4. After creating the policy, for policy ID, go to `Settings -> Policies -> Group -> Select the appropriate policy` and note the policy ID from the link.<br>
   Example: Link of a Polciy - https://www.expensify.com/policy?param={%22policyID%22:%22E95AFCD33ABE2BB8%22}

### API Authentication
1. To authenticate API we need a pair of credentials: partnerUserID and partnerUserSecret.<br>
2. For this go to https://www.expensify.com/tools/integrations/ and genrate the credentials.<br>
3. A pair of credentials: partnerUserID and partnerUserSecret will be generated and shown on the page.<br>


## Building The Provider
1. Clone the repository and add the dependencies. For this run the following commands: <br>
```git clone https://github.com/shubhambjadhavar/terraform-provider-expensify.git
cd terraform-provider-expensify
go mod init terraform-provider-expensify
go mod tidy
```
2. Run `go mod vendor` to create a vendor directory that contains all the provider's dependencies. <br>


## Managing terraform plugins
*For Windows:*
1. Run the following command to create a vendor subdirectory (`%APPDATA%/terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${OS_ARCH}`) which will comprise of all terraform plugins. <br> 
Command: 
```bash
mkdir -p %APPDATA%/terraform.d/plugins/expensify.com/employee/expensify/1.0.0/windows_amd64
```
2. Run `go build -o terraform-provider-expensify.exe` to generate the binary in present working directory. <br>
3. Run this command to move this binary file to appropriate location.
 ```
 move terraform-provider-expensify.exe %APPDATA%\terraform.d\plugins\expensify.com\employee\expensify\1.0.0\windows_amd64
 ``` 
 <p align="center">
 [OR]
 </p>
 
3. Manually move the file from current directory to destination directory (`%APPDATA%\terraform.d\plugins\expensify.com\employee\expensify\1.0.0\windows_amd64`).<br>


## Working with terraform

### Application Credential Integration in terraform
1. Add terraform block and provider block as shown in Example Usage below.
2. Get a pair of credentials: partnerUserID and partnerUserSecret. For this visit https://www.expensify.com/tools/integrations/.
3. Assign the above credentials to the repective field in the provider block.

### Basic Terraform Commands
1. `terraform init` - To initialize a working directory containing Terraform configuration files.
2. `terraform plan` - To create an execution plan. Displays the changes to be done.
3. `terraform apply` - To execute the actions proposed in a Terraform plan. Apply the chages.

#### Create User
1. Add the employee emial, manager emial, policy id, first name, last name in the respective field in resource block as shown in Example Usage below.
2. For policy ID, go to `Settings -> Policies -> Group -> Select the appropriate policy` and note the policy ID from the link.<br>
   Example: Link of a Polciy - https://www.expensify.com/policy?param={%22policyID%22:%22E95AFCD33ABE2BB8%22}
3. Run the basic terraform commands.
4. On successful execution, creates a user and sends an account setup mail to the user.
5. Setup the account using the link provided in the mail.

#### Update the user
1. Update the data of the user in the resource block as show in Example Usage below and run the basic terraform commands to update user. User is allowed to update `employee_id`, `first_name`, `last_name`, `manager_email`, `approves to`, `over limit approver`, and `approval limit`. Updating `first name` and `last name` in any one policy will atomatically update them in other policies.
2. To remove a user from a policy update value of `is_terminated` field in resource block to `true`.
3. To readd the removed user the policy update value of `is_terminated` field in resource block back to `false`.

#### Read the User Data
Add data and output blocks as shown in the Example Usage below and run the basic terraform commands.

#### Delete the user
Delete the resource block of the particular user and run `terraform apply`.

#### Import a User Data
1. Write manually a resource configuration block for the User as shown in Example Usage below. Imported user will be mapped to this block.
2. Run the command `terraform import expensify_employee.employee [POLICY_ID]:[EMAIL_ID]` to import user.
3. For policy ID, go to `Settings -> Policies -> Group -> Select the appropriate policy` and note the policy ID from the link.<br>
   Example: Link of a Polciy - https://www.expensify.com/policy?param={%22policyID%22:%22E95AFCD33ABE2BB8%22}
4. Run `terraform plan`, if output show `0 to addd, 0 to change and 0 to destroy` user import is successful, otherwise recheck the employee data in resource block with employee data in the policy in Expensify Website. 


## Example Usage
```terraform
terraform{
    required_providers {
        expensify = {
            version = "1.0.0"
            source = "expensify.com/employee/expensify"
        }
    }
}

provider "expensify" {
    partner_user_id = ""
    partner_user_secret = "" 
}

resource "expensify_employee" "employee"{
    employee_email = "[EMPLOYEE_EMAIL]"
    manager_email = "[MANAGER_EMAIL]"
    policy_id = "[POLICY_ID]"
    employee_id = "[EMPLOYEE_ID]"
    first_name = "[FIRST_NAME]"
    last_name = "[LAST_NAME]"
    approves_to = "[APPROVES_TO]"
    approval_limit = "[APPROVAL_LIMIT]"
    over_limit_approver = "[OVER_LIMIT_APPROVER]"
    is_terminated = false
}

output "policy1"{
    value = expensify_employee.employee
}

data "expensify_employee" "employee" {
    policy_id = "[POLICY_ID]"
    employee_email = "[EMPLOYEE_EMAIL]" 
}

output "datasouce_employee"{
    value = data.expensify_employee.employee
}
```


## Argument Reference

* `partner_user_id`      - (Required, String)  - The Expensify Partner User ID
* `partner_user_secret`  - (Required, String)  - The Expensify Partner User Secret
* `employee_email`       - (Required, String)  - The email address of the employee.
* `manager_email`        - (Required, String)  - Who the employee should submit reports to.
* `policy_id`            - (Required, String)  - The id of policy for which employee is to be added.
* `first_name`           - (Optional, String)  - First name of the employee in Expensify. Not allowed overwrite values manually set by the employee in their Expensify account.
* `last_name`            - (Optional, String)  - Last name of the employee in Expensify. Not allowed overwrite values manually set by the employee in their Expensify account.
* `is_terminated`        - (Optional, Boolean) - If set to true, the employee will be removed from the policy.
* `employee_id`          - (Optional, String)  - Unique ID of the Employee.
* `over_limit_approver`  - (Optional, String)  - Who the manager should forward reports to if a report is over approval_limit. Required if an approval_limit is specified.
* `approver_limit`       - (Optional, Float)   - Specifies limit of report total.
* `approves_to`          - (Optional, String)  - Who the employee should forward the report to.


## Exceptions

* Updating of the fields `manager emial`, `approves to`, `over limit approver`, and `approval limit` is meaningful only if Approval Mode for policy is Advanced Approval.
* Updating `first name` and `last name` in any one policy will atomatically update them in other policies.
* To add an employee to multiple policies write multiple resource block with different policy ID.
