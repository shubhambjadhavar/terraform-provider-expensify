package client

import (
	"os"
	"log"
	"io/ioutil"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/clarketm/json"
)

func init(){
	file, err := os.Open("../credentials.json")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(file)
	if err!=nil {
		log.Fatal(err)
	}
	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err!=nil {
		log.Fatal(err)
	}
	os.Setenv("EXPENSIFY_USER_ID", res["EXPENSIFY_USER_ID"].(string))
	os.Setenv("EXPENSIFY_USER_SECRET", res["EXPENSIFY_USER_SECRET"].(string))
}

func TestExpensifyClient_NewEmployee(t *testing.T) {
	testCases := []struct{
		testName string
		employeesList *EmployeesList
		employee *Employee
		expectErr bool
	}{
		{
			testName: "employee successfully created",
			employeesList: &EmployeesList{
				Employees: []Employee{
					{
						PolicyId: "56B042862350ADD2",
						EmployeeEmail: "abhishiek@clevertapdemo.ml",
						ManagerEmail: "shubham@clevertapdemo.ml",
						EmployeeId: "1003",
						FirstName: "Abhishiek",
						LastName: "Singh",
					},
				},
			},
			employee: &Employee{
				PolicyId: "56B042862350ADD2",
				EmployeeEmail: "abhishiek@clevertapdemo.ml",
				ManagerEmail: "shubham@clevertapdemo.ml",
				EmployeeId: "1003",
				Role: "user",
			},
			expectErr: false,
		},
		{
			testName: "employee already exists",
			employeesList: &EmployeesList{
				Employees: []Employee{
					{
						PolicyId: "56B042862350ADD2",
						EmployeeEmail: "abhishiek@clevertapdemo.ml",
						ManagerEmail: "shubham@clevertapdemo.ml",
						EmployeeId: "1003",
						FirstName: "Abhishiek",
						LastName: "Singh",
					},
				},
			},
			employee: nil,
			expectErr: true,
		},
		{
			testName: "invalid policy id",
			employeesList: &EmployeesList{
				Employees: []Employee{
					{
						PolicyId: "E95AFCD33ABE2BB7",
						EmployeeEmail: "ashutosh@clevertapdemo.ml",
						ManagerEmail: "shubham@clevertapdemo.ml",
						EmployeeId: "1002",
						FirstName: "Ashutosh",
						LastName: "Verma",
					},
				},
			},
			employee: nil,
			expectErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			apiClient := NewClient(os.Getenv("EXPENSIFY_USER_ID"), os.Getenv("EXPENSIFY_USER_SECRET"))
			err := apiClient.NewEmployee(tc.employeesList)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			employee := Employee{
				PolicyId: tc.employeesList.Employees[0].PolicyId,
				EmployeeEmail: tc.employeesList.Employees[0].EmployeeEmail,
			}
			body, err := apiClient.GetEmployee(&employee)
			assert.NoError(t, err)
			assert.Equal(t, tc.employee, body)
		})
	}
}

func TestExpensifyClient_GetEmployee(t *testing.T) {
	testCases := []struct{
		testName string
		employee *Employee
		expectErr bool
		expectedResp *Employee
	}{
		{
			testName: "employee exists",
			employee: &Employee{
				PolicyId: "56B042862350ADD2",
				EmployeeEmail: "abhishiek@clevertapdemo.ml",
			},
			expectErr: false,
			expectedResp: &Employee{
				PolicyId: "56B042862350ADD2",
				EmployeeEmail: "abhishiek@clevertapdemo.ml",
				ManagerEmail: "shubham@clevertapdemo.ml",
				EmployeeId: "1003",
				Role: "user",
			},
		},
		{
			testName: "employee does not exist",
			employee: &Employee{
				PolicyId: "56B042862350ADD2",
				EmployeeEmail: "ashutosh@clevertapdemo.ml",
			},
			expectErr: true,
			expectedResp: nil,
		},
		{
			testName: "Invalid Policy ID",
			employee: &Employee{
				PolicyId: "E95AFCD33ABE2BB7",
				EmployeeEmail: "abhishiek@clevertapdemo.ml",
			},
			expectErr: true,
			expectedResp: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			apiClient := NewClient(os.Getenv("EXPENSIFY_USER_ID"), os.Getenv("EXPENSIFY_USER_SECRET"))
			body, err := apiClient.GetEmployee(tc.employee)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResp, body)
		})
	}
}

func TestExpensifyClient_UpdateEmployee(t *testing.T) {
	testCases := []struct{
		testName string
		employeesList *EmployeesList
		employee *Employee
		expectErr bool
	}{
		{
			testName: "employee exists",
			employeesList: &EmployeesList{
				Employees: []Employee{
					{
						PolicyId: "56B042862350ADD2",
						EmployeeEmail: "abhishiek@clevertapdemo.ml",
						ManagerEmail: "shubham@clevertapdemo.ml",
						EmployeeId: "1003",
						FirstName: "Abhishiek",
						LastName: "Singh",
					},
				},
			},
			employee: &Employee{
				PolicyId: "56B042862350ADD2",
				EmployeeEmail: "abhishiek@clevertapdemo.ml",
				ManagerEmail: "shubham@clevertapdemo.ml",
				EmployeeId: "1003",
				Role: "user",
			},
			expectErr: false,
		},
		{
			testName: "employee does not exist",
			employeesList: &EmployeesList{
				Employees: []Employee{
					{
						PolicyId: "56B042862350ADD2",
						EmployeeEmail: "ashutosh@clevertapdemo.ml",
						ManagerEmail: "shubham@clevertapdemo.ml",
						EmployeeId: "1002",
						FirstName: "Ashutosh",
						LastName: "Verma",
					},
				},
			},
			employee: nil,
			expectErr: true,
		},
		{
			testName: "invalid policy id",
			employeesList: &EmployeesList{
				Employees: []Employee{
					{
						PolicyId: "E95AFCD33ABE2BB7",
						EmployeeEmail: "abhishiek@clevertapdemo.ml",
						ManagerEmail: "shubham@clevertapdemo.ml",
						EmployeeId: "1003",
						FirstName: "Abhishiek",
						LastName: "Singh",
					},
				},
			},
			employee: nil,
			expectErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			apiClient := NewClient(os.Getenv("EXPENSIFY_USER_ID"), os.Getenv("EXPENSIFY_USER_SECRET"))
			err := apiClient.UpdateEmployee(tc.employeesList)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			employee := Employee{
				PolicyId: tc.employeesList.Employees[0].PolicyId,
				EmployeeEmail: tc.employeesList.Employees[0].EmployeeEmail,
			}
			body, err := apiClient.GetEmployee(&employee)
			assert.NoError(t, err)
			assert.Equal(t, tc.employee, body)
		})
	}
}

func TestExpensifyClient_DeleteEmployee(t *testing.T) {
	testCases := []struct{
		testName string
		employeesList *EmployeesList
		expectErr bool
	}{
		{
			testName: "employee exists",
			employeesList: &EmployeesList{
				Employees: []Employee{
					{
						PolicyId: "56B042862350ADD2",
						EmployeeEmail: "abhishiek@clevertapdemo.ml",
						ManagerEmail: "shubham@clevertapdemo.ml",
						EmployeeId: "1003",
						FirstName: "Abhishiek",
						LastName: "Singh",
					},
				},
			},
			expectErr: false,
		},
		{
			testName: "employee does not exist",
			employeesList: &EmployeesList{
				Employees: []Employee{
					{
						PolicyId: "56B042862350ADD2",
						EmployeeEmail: "abhishiek@clevertapdemo.ml",
						ManagerEmail: "shubham@clevertapdemo.ml",
						EmployeeId: "1003",
						FirstName: "Abhishiek",
						LastName: "Singh",
					},
				},
			},
			expectErr: true,
		},
		{
			testName: "invalid policy id",
			employeesList: &EmployeesList{
				Employees: []Employee{
					{
						PolicyId: "E95AFCD33ABE2BB7",
						EmployeeEmail: "abhishiek@clevertapdemo.ml",
						ManagerEmail: "shubham@clevertapdemo.ml",
						EmployeeId: "1003",
						FirstName: "Abhishiek",
						LastName: "Singh",
					},
				},
			},
			expectErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			apiClient := NewClient(os.Getenv("EXPENSIFY_USER_ID"), os.Getenv("EXPENSIFY_USER_SECRET"))
			err := apiClient.DeleteEmployee(tc.employeesList)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			employee := Employee{
				PolicyId: tc.employeesList.Employees[0].PolicyId,
				EmployeeEmail: tc.employeesList.Employees[0].EmployeeEmail,
			}
			_, err = apiClient.GetEmployee(&employee)
			assert.Error(t, err)
		})
	}
}

func TestExpensifyClient_NewPolicy(t *testing.T) {
	testCases := []struct{
		testName string
		policyName string
		plan string
		expectErr bool
		expectedResp *Policy
	}{
		{
			testName: "policy successfully created",
			policyName: "test create",
			plan: "corporate",
			expectErr: false,
			expectedResp: &Policy{
				Owner: "shubhamj@clevertapdemo.ml",
				PolicyName: "test create",
				OutputCurrency: "USD",
				Plan: "corporate",
			},
		},
		{
			testName: "invalid plan type",
			policyName: "test create",
			plan: "corporat",
			expectErr: true,
			expectedResp: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			apiClient := NewClient(os.Getenv("EXPENSIFY_USER_ID"), os.Getenv("EXPENSIFY_USER_SECRET"))
			policyId, err := apiClient.NewPolicy(tc.policyName, tc.plan)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			policy, err := apiClient.GetPolicy(policyId)
			assert.NoError(t, err)
			tc.expectedResp.PolicyId = policyId
			assert.Equal(t, tc.expectedResp, policy)
		})
	}
}

func TestExpensifyClient_GetPolicy(t *testing.T) {
	testCases := []struct{
		testName string
		policyId string
		expectErr bool
		expectedResp *Policy
	}{
		{
			testName: "policy does not exist",
			policyId: "E95AFCD33ABE2BB",
			expectErr: true,
			expectedResp: nil,
		},
		{
			testName: "policy exist",
			policyId: "56B042862350ADD2",
			expectErr: false,
			expectedResp: &Policy{
				PolicyId: "56B042862350ADD2",
				Owner: "shubhamj@clevertapdemo.ml",
				PolicyName: "shubhamj",
				OutputCurrency: "INR",
				Plan: "corporate",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			apiClient := NewClient(os.Getenv("EXPENSIFY_USER_ID"), os.Getenv("EXPENSIFY_USER_SECRET"))
			policy, err := apiClient.GetPolicy(tc.policyId)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResp, policy)
		})
	}
}