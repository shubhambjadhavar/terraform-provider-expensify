This terraform provider allows to perform Create, Read, Update, Remove, Import Expensify Employee(s) to a Policy.


## Requirements

* [Go](https://golang.org/doc/install) >= 1.16 <br>
* [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x <br/>
* [Expensify](https://www.expensify.com/) Collect/Control policy plan


## Setup Expensify Account
 :heavy_exclamation_mark:  [IMPORTANT] : This provider can only be successfully tested on a Expensify account with Control or Collect plan. <br>

1. Create a expensify account. (https://www.expensify.com/)<br>
2. Sign in to the expensify account.<br>
3. Create a policy. For this go to traverse to `Settings -> Policies -> Group` and click on `New Policy`.<br>
4. After creating policy we will get Policy ID in which we are going to manage users.
5. Go to https://www.expensify.com/tools/integrations/.<br>
A pair of credentials: partnerUserID and partnerUserSecret will be generated and shown on the page.<br>


## Initialise Expensify Provider in local machine 
1. Clone the repository  to $GOPATH/src/github.com/expensify/terraform-provider-expensify <br>
2. Add the partnerUserID and partnerUserSecret generated to respective fields in `main.tf` <br>
3. Run the following command :
```golang
go mod init terraform-provider-expensify
go mod tidy
```
4. Run `go mod vendor` to create a vendor directory that contains all the provider's dependencies. <br>

## Installation
1. Run the following command to create a vendor subdirectory which will comprise of  all provider dependencies. <br>
```
%APPDATA%/terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}
``` 
Command: 
```bash
mkdir -p %APPDATA%/terraform.d/plugins/expensify.com/employee/expensify/1.0.0/[OS_ARCH]
```
For eg. `mkdir -p %APPDATA%/terraform.d/plugins/expensify.com/employee/expensify/1.0.0/windows_amd64`<br>

2. Run `go build -o terraform-provider-expensify.exe`. This will save the binary (`.exe`) file in the main/root directory. <br>
3. Run this command to move this binary file to appropriate location.
 ```
 move terraform-provider-expensify.exe %APPDATA%\terraform.d\plugins\expensify.com\employee\expensify\1.0.0\[OS_ARCH]
 ``` 
Otherwise you can manually move the file from current directory to destination directory.<br>


[OR]

1. Download required binaries <br>
2. move binary `%APPDATA%/terraform.d/plugins/[architecture name]/`


## Run the Terraform provider

#### Create User
1. Add the employee emial, manager emial, policy id, first name, last name, employee id in the respective field in `main.tf`
2. To add employee in multiple policies create multiple resource blocks with different Policy ID value.
3. Initialize the terraform provider `terraform init`
4. Check the changes applicable using `terraform plan` and apply using `terraform apply`
5. You will see that a user has been successfully created and an account setup mail has been sent to the user.
6. Setup the account using the link provided in the mail.

#### Update the user
Update the data of the user in the `main.tf` file and apply using `terraform apply`. Updating of the fields `manager emial`, `approves to`, `over limit approver`, and `approval limit` is useful only if Approval Mode for policy is Advanced Approval. Updating `first name` and `last name` in any one policy will atomatically update them in other policies.

#### Read the User Data
Add data and output blocks in the `main.tf` file and run `terraform plan` to read user data.

#### Activate/Deactivate the user
Change the is terminated field of User from `false` to `true` or viceversa and run `terraform apply`.

#### Delete the user
Delete the resource block of the particular user from `main.tf` file and run `terraform apply`.

#### Import a User Data
1. Write manually a resource configuration block for the User in `main.tf`, to which the imported object will be mapped.
2. Run the command `terraform import expensify_employee.policy1 [POLICY_ID]:[EMAIL_ID]`
4. Check for the attributes in the `.tfstate` file and fill them accordingly in resource block.


### Testing the Provider
1. Navigate to the test file directory.
2. Run command `go test` . This command will give combined test result for the execution or errors if any failure occur.
3. If you want to see test result of each test function individually while running test in a single go, run command `go test -v`
4. To check test cover run `go test -cover`


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

resource "expensify_employee" "policy1"{
    employee_email = "[EMPLOYEE_EMAIL]"
    manager_email = "[MANAGER_EMAIL]"
    policy_id = "[POLICY_ID]"
    employee_id = "[EMPLOYEE_ID]"
    first_name = "[FIRST_NAME]"
    last_name = "[LAST_NAME]"
    is_terminated = false
}

output "policy1"{
    value = expensify_employee.policy1
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

* `partner_user_id`      - The Expensify Partner User ID
* `partner_user_secret`  - The Expensify Partner User Secret
* `employee_email`       - The email id associated with the user account.
* `manager_email`        - The email id associated with the manager account.
* `policy_id`            - The id of policy for which employee is to be added
* `first_name`           - First name of the User.
* `last_name`            - Last Name / Family Name / Surname of the User.
* `is_terminated`        - User account present in policy or not.
* `employee_id`          - Unique ID of the User.