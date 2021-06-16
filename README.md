# Terraform Provider Expensify

This Expensify provider allows Terraform to add users to the Expensify policies, read users of the Expensify policies, update users in the Expensify policies, remove users from the Expensify policies, create a new policy, read a policy details.<br>


## Requirements

* [Go](https://golang.org/doc/install) >= 1.16 (To build the provider plugin)<br>
* [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x <br/>
* Account: [Expensify](https://www.expensify.com/)(Required a Verified Domain and an Expensify Policy with Control/Collect Plan)<br>
* [API Documentation](https://integrations.expensify.com/Integration-Server/doc/)<br>


## API Authentication

*Generate credentials using policy's admin accounts*
1. To authenticate API, we need a pair of credentials: partnerUserID and partnerUserSecret.<br>
2. For this, go to https://www.expensify.com/tools/integrations/ and generate the credentials.<br>
3. A pair of credentials: partnerUserID and partnerUserSecret will be generated and shown on the page.<br>


## Installing the Provider

*For Windows using Command Prompt*
1. Clone the repository, add all the dependencies, create a vendor directory that contains all dependencies and generate a binary. For this, run the following commands: <br>
```
git clone https://github.com/shubhambjadhavar/terraform-provider-expensify.git
cd terraform-provider-expensify
go mod init terraform-provider-expensify
go mod tidy
go mod vendor
go build -o terraform-provider-expensify.exe
```
2. Move the generated binary to `%APPDATA%/terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${OS_ARCH}`. For this, run the following commands: <br>  
```
mkdir -p %APPDATA%/terraform.d/plugins/expensify.com/employee/expensify/1.0.0/windows_amd64
move terraform-provider-expensify.exe %APPDATA%\terraform.d\plugins\expensify.com\employee\expensify\1.0.0\windows_amd64
```

*For Windows Manually*
1. Download the latest compiled binary from [GitHub releases](https://github.com/shubhambjadhavar/terraform-provider-expensify/releases) and unzip/untar the archive.<br>
2. Move the binary to `%APPDATA%/terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${OS_ARCH}`.<br>

Create your Terraform configurations as shown in [example usage](#example-usage) and run `terraform init`. This will find plugin locally.<br>


## Note

* Update the fields `manager_email`, `approves_to`, `over_limit_approver`, and `approval_limit` only if Approval Mode for policy is Advanced Approval.<br>
* API not allow overwriting manually set values for `first_name` and `last_name` in their Expensify account.<br>
* The fields `first_name` and `last_name` are set at account level.<br>
* Once the value of any attribute is set, it cannot be set back to null through provider. But, you can set it to null via UI.<br> 


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