package client

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"github.com/clarketm/json"
)

var (
    Errors = make(map[int]string)
)

func init() {
	Errors[500] = "{\"responseMessage\":Bad Request,\"responseCode\":500}"
	Errors[404] = "{\"responseMessage\":User Does Not Exist ,\"responseCode\":404}"
	Errors[409] = "{\"responseMessage\":User Already Exist,\"responseCode\":409}"
	Errors[401] = "{\"responseMessage\":Unautharized Access,\"responseCode\":401}"
}

type Client struct {
	partnerUserId string
	partnerUserSecret string
	httpClient *http.Client
}

type Credentials struct{
	PartnerUserId string `json:"partnerUserID"`
	PartnerUserSecret string `json:"partnerUserSecret"`
}

type InputSettings struct{
	Type string `json:"type"`
	Entity string `json:"entity,omitempty"`
	Fields []string `json:"fields,omitempty"`
	PolicyIdList []string `json:"policyIDList,omitempty"`
	PolicyName string `json:"policyName,omitempty"`
	Plan string `json:"plan,omitempty"`
}

type Categories struct{
	Action string `json:"action"`
	DataArray []CategoriesData `json:"data"`
}

type CategoriesData struct{
	Name string `json:"name"`
	Type string `json:"type"`
	Enabled bool `json:"enabled"`
	GlCode string `json:"glCode"`
	PayrollCode string `json:"payrollCode"`
	AreCommentsRequired bool `json:"areCommentsRequired"`
	CommentHint string `json:"commentHint"`
	MaxExpenseAmount int `json:"maxExpenseAmount"`
}

type ReportFields struct{
	Action string `json:"action"`
	DataArray []ReportFieldsData `json:"data"`
}

type ReportFieldsData struct {
	Name string `json:"name"`
	Type string `json:"type"`
	DefaultValue string `json:"defaultValue"`
	ValuesArray []Values `json:"values"`
}

type Values struct{
	Value string `json:"value"`
	Enabled bool `json:"enabed"`
	ExternalId string `json:"externalID"`
}

type Tags struct{
	Source string `json:"source"`
	DataArray []TagsData `json:"data"`
}

type TagsData struct{
	Name string `json:"name"`
	SetRequired bool `json:"setRequired"`
	TagArray []Tag `json:"tags"`	
}

type Tag struct{
	Name string `json:"name"`
	Enabled bool `json:"enabled"`
	GlCode string `json:"glCode"`
}

type RequestJobDescription struct{
	Type string `json:"type"`
	DryRun bool `json:"dry-run,omitempty"`
	Credential Credentials `json:"credentials"`
	DataSource string `json:"dataSource,omitempty"`
	InputSetting InputSettings `json:"inputSettings"`
	Category Categories `json:"categories,omitempty"`
	ReportField ReportFields `json:"reportField,omitempty"`
	TagArray Tags `json:"tags,omitempty"`
}

type Response struct{
	ResponseCode int `json:"responseCode"`
	PolicyInfo map[string]interface{} `json:"policyInfo,omitempty"`
	UpdatedEmployeesCount int `json:"updatedEmployeesCount,omitempty"`
	SkippedEmployees []Message `json:"skippedEmployees,omitempty"`
	PolicyId string `json:"policyID,omitempty"`
	PolicyList []Policy `json:"policyList"`
}

type Message struct{
	Reason string `json:"reason"`
}

type EmployeesList struct{
	Employees []Employee `json:"Employees"`
}

type Employee struct{
	EmployeeEmail string `json:"employeeEmail"`
	ManagerEmail string `json:"managerEmail"`
	PolicyId string `json:"policyID"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	EmployeeId string `json:"employeeID"`
	ApprovalLimit float64 `json:"approvalLimit"`
	OverLimitApprover string `json:"overLimitApprover"`
	IsTerminated bool `json:"isTerminated"`
	ApprovesTo string `json:"approvesTo"`
	Role string `json:"role"`
}

type Policy struct{
	PolicyId string `json:"id"`
	Owner string `json:"owner"`
	PolicyName string `json:"name"`
	OutputCurrency string `json:"outputCurrency"`
	Plan string `json:"type"`
}

func NewClient(partnerUserId, partnerUserSecret string) *Client {
	return &Client{
		partnerUserId: partnerUserId,
		partnerUserSecret: partnerUserSecret,
		httpClient: &http.Client{},
	}
}

func (c *Client) NewEmployee(employeesList *EmployeesList) error {
	credentials := Credentials{
		PartnerUserId: c.partnerUserId,
		PartnerUserSecret: c.partnerUserSecret,
	}
	inputSettings := InputSettings{
		Type: "employees",
		Entity: "generic",
	}
	requestJobDescription := RequestJobDescription{
		Type: "update",
		DryRun: false,
		Credential: credentials,
		DataSource: "request",
		InputSetting: inputSettings,
	}
	employee := Employee{
		PolicyId: employeesList.Employees[0].PolicyId,
		EmployeeEmail: employeesList.Employees[0].EmployeeEmail,
	}
	_, err := c.GetEmployee(&employee)
	if err == nil {
		log.Println("[CREATE ERROR]:", Errors[409])
		return fmt.Errorf(Errors[409])
	}
	employeesListMarshal, err := json.Marshal(employeesList)
	if err != nil {
		log.Println("[CREATE ERROR]:", err)
		return err
	}
	requestJobDescriptionMarshal, err := json.Marshal(requestJobDescription)
	if err != nil {
		log.Println("[CREATE ERROR]:", err)
		return err
	}
	parms := url.Values{}
	parms.Add("requestJobDescription", string(requestJobDescriptionMarshal))
	parms.Add("data", string(employeesListMarshal))
	body := strings.NewReader(parms.Encode())
	res, err := c.httpRequest("POST", body)
	if err != nil {
		log.Println("[CREATE ERROR]:", err)
		return err
	}
	resp := Response{}
	err = json.Unmarshal(res, &resp)
	if err != nil {
		log.Println("[CREATE ERROR]:", err)
		return err
	}
	if resp.ResponseCode!=200{
		log.Println("[CREATE ERROR]:", fmt.Errorf(string(res)))
		return fmt.Errorf(string(res))
	}
	if resp.UpdatedEmployeesCount!=1{
		err := ""
		for _,v := range(resp.SkippedEmployees){
			err = err + v.Reason
		}
		log.Println("[CREATE ERROR]:", err)
		return fmt.Errorf(err)
	}
	return nil
}

func (c *Client) GetEmployee(employee *Employee) (*Employee, error) {
	credentials := Credentials{
		PartnerUserId: c.partnerUserId,
		PartnerUserSecret: c.partnerUserSecret,
	}
	inputSettings := InputSettings{
		Type: "policy",
		Fields: []string{"employees",},
        PolicyIdList: []string{employee.PolicyId,},
	}
	requestJobDescription := RequestJobDescription{
		Type: "get",
		Credential: credentials,
		InputSetting: inputSettings,
	}
	requestJobDescriptionMarshal, err := json.Marshal(requestJobDescription)
	if err != nil {
		log.Println("[READ ERROR]:", err)
		return nil, err
	}
	parms := url.Values{}
	parms.Add("requestJobDescription", string(requestJobDescriptionMarshal))
	body := strings.NewReader(parms.Encode())
	res, err := c.httpRequest("POST", body)
	if err != nil {
		log.Println("[READ ERROR]:", err)
		return nil, err
	}
	resp := Response{}
	err = json.Unmarshal(res, &resp)
	if err != nil {
		log.Println("[READ ERROR]:", err)
		return nil, err
	}
	if resp.ResponseCode!=200{
		log.Println("[READ ERROR]:", fmt.Errorf(string(res)))
		return nil, fmt.Errorf(string(res))
	}
	employees := resp.PolicyInfo[employee.PolicyId].(map[string]interface{})
	employeeList := employees["employees"].([]interface{})
	userExist := 0
	for _, v := range(employeeList){
		temp := v.(map[string]interface{})
		if temp["email"].(string)==employee.EmployeeEmail {
			userExist = 1
			if temp["role"]!=nil {
				employee.Role = temp["role"].(string)
			}
			if temp["submitsTo"]!=nil {
				employee.ManagerEmail = temp["submitsTo"].(string)
			}
			if temp["employeeID"]!=nil {
				employee.EmployeeId = temp["employeeID"].(string)
			}
			if temp["forwardsTo"]!=nil {
				employee.ApprovesTo = temp["forwardsTo"].(string)
			}
			if temp["overLimitForwardsTo"]!=nil {
				employee.OverLimitApprover = temp["overLimitForwardsTo"].(string)
			}
			if temp["approvalLimit"]!=nil {
				employee.ApprovalLimit = temp["approvalLimit"].(float64)
			}
			break
		}
	}
	if userExist!=1 {
		log.Println("[READ ERROR]:", Errors[404])
		return nil, fmt.Errorf(Errors[404])
	}
	return employee, nil
}

func (c *Client) UpdateEmployee(employeesList *EmployeesList) error {
	credentials := Credentials{
		PartnerUserId: c.partnerUserId,
		PartnerUserSecret: c.partnerUserSecret,
	}
	inputSettings := InputSettings{
		Type: "employees",
		Entity: "generic",
	}
	requestJobDescription := RequestJobDescription{
		Type: "update",
		DryRun: false,
		Credential: credentials,
		DataSource: "request",
		InputSetting: inputSettings,
	}
	employee := Employee{
		PolicyId: employeesList.Employees[0].PolicyId,
		EmployeeEmail: employeesList.Employees[0].EmployeeEmail,
	}
	_, err := c.GetEmployee(&employee)
	if err != nil && (strings.Contains(fmt.Sprintf("%v", err),"\"responseCode\":404")) {
		log.Println("[UPDATE ERROR]:", err)
		return err
	}
	employeesListMarshal, err := json.Marshal(employeesList)
	if err != nil {
		log.Println("[UPDATE ERROR]:", err)
		return err
	}
	requestJobDescriptionMarshal, err := json.Marshal(requestJobDescription)
	if err != nil {
		log.Println("[UPDATE ERROR]:", err)
		return err
	}
	parms := url.Values{}
	parms.Add("requestJobDescription", string(requestJobDescriptionMarshal))
	parms.Add("data", string(employeesListMarshal))
	body := strings.NewReader(parms.Encode())
	res, err := c.httpRequest("POST", body)
	if err != nil {
		log.Println("[UPDATE ERROR]:", err)
		return err
	}
	resp := Response{}
	err = json.Unmarshal(res, &resp)
	if err != nil {
		log.Println("[UPDATE ERROR]:", err)
		return err
	}
	if resp.ResponseCode!=200{
		log.Println("[UPDATE ERROR]:", fmt.Errorf(string(res)))
		return fmt.Errorf(string(res))
	}
	if resp.UpdatedEmployeesCount!=1{
		err := ""
		for _,v := range(resp.SkippedEmployees){
			err = err + v.Reason
		}
		log.Println("[UPDATE ERROR]:", err)
		return fmt.Errorf(err)
	}
	return nil
}

func (c *Client) ActivateEmployee(employeesList *EmployeesList) error {
	credentials := Credentials{
		PartnerUserId: c.partnerUserId,
		PartnerUserSecret: c.partnerUserSecret,
	}
	inputSettings := InputSettings{
		Type: "employees",
		Entity: "generic",
	}
	requestJobDescription := RequestJobDescription{
		Type: "update",
		DryRun: false,
		Credential: credentials,
		DataSource: "request",
		InputSetting: inputSettings,
	}
	employeesListMarshal, err := json.Marshal(employeesList)
	if err != nil {
		log.Println("[ACTIVATE ERROR]:", err)
		return err
	}
	requestJobDescriptionMarshal, err := json.Marshal(requestJobDescription)
	if err != nil {
		log.Println("[ACTIVATE ERROR]:", err)
		return err
	}
	parms := url.Values{}
	parms.Add("requestJobDescription", string(requestJobDescriptionMarshal))
	parms.Add("data", string(employeesListMarshal))
	body := strings.NewReader(parms.Encode())
	res, err := c.httpRequest("POST", body)
	if err != nil {
		log.Println("[ACTIVATE ERROR]:", err)
		return err
	}
	resp := Response{}
	err = json.Unmarshal(res, &resp)
	if err != nil {
		log.Println("[ACTIVATE ERROR]:", err)
		return err
	}
	if resp.ResponseCode!=200{
		log.Println("[ACTIVATE ERROR]:", fmt.Errorf(string(res)))
		return fmt.Errorf(string(res))
	}
	if resp.UpdatedEmployeesCount!=1{
		err := ""
		for _,v := range(resp.SkippedEmployees){
			err = err + v.Reason
		}
		log.Println("[ACTIVATE ERROR]:", err)
		return fmt.Errorf(err)
	}
	return nil
}

func (c *Client) DeleteEmployee(employeesList *EmployeesList) error {
	credentials := Credentials{
		PartnerUserId: c.partnerUserId,
		PartnerUserSecret: c.partnerUserSecret,
	}
	inputSettings := InputSettings{
		Type: "employees",
		Entity: "generic",
	}
	requestJobDescription := RequestJobDescription{
		Type: "update",
		DryRun: false,
		Credential: credentials,
		DataSource: "request",
		InputSetting: inputSettings,
	}
	employee := Employee{
		PolicyId: employeesList.Employees[0].PolicyId,
		EmployeeEmail: employeesList.Employees[0].EmployeeEmail,
	}
	_, err := c.GetEmployee(&employee)
	if err != nil && (strings.Contains(fmt.Sprintf("%v", err),"\"responseCode\":404")) {
		log.Println("[DELETE ERROR]:", err)
		return err
	}
	employeesList.Employees[0].IsTerminated = true
	employeesListMarshal, err := json.Marshal(employeesList)
	if err != nil {
		log.Println("[DELETE ERROR]:", err)
		return err
	}
	requestJobDescriptionMarshal, err := json.Marshal(requestJobDescription)
	if err != nil {
		log.Println("[DELETE ERROR]:", err)
		return err
	}
	parms := url.Values{}
	parms.Add("requestJobDescription", string(requestJobDescriptionMarshal))
	parms.Add("data", string(employeesListMarshal))
	body := strings.NewReader(parms.Encode())
	res, err := c.httpRequest("POST", body)
	if err != nil {
		log.Println("[DELETE ERROR]:", err)
		return err
	}
	resp := Response{}
	err = json.Unmarshal(res, &resp)
	if err != nil {
		log.Println("[DELETE ERROR]:", err)
		return err
	}
	if resp.ResponseCode!=200{
		log.Println("[DELETE ERROR]:", fmt.Errorf(string(res)))
		return fmt.Errorf(string(res))
	}
	if resp.UpdatedEmployeesCount!=1{
		err := ""
		for _,v := range(resp.SkippedEmployees){
			err = err + v.Reason
		}
		log.Println("[DELETE ERROR]:", err)
		return fmt.Errorf(err)
	}
	return nil
}

func (c *Client) NewPolicy(policyName string, plan string) (string, error){
	credentials := Credentials{
		PartnerUserId: c.partnerUserId,
		PartnerUserSecret: c.partnerUserSecret,
	}
	inputSettings := InputSettings{
		Type: "policy",
		PolicyName: policyName,
		Plan: plan,
	}
	requestJobDescription := RequestJobDescription{
		Type: "create",
		Credential: credentials,
		InputSetting: inputSettings,
	}
	requestJobDescriptionMarshal, err := json.Marshal(requestJobDescription)
	if err != nil {
		log.Println("[CREATE ERROR]:", err)
		return "", err
	}
	parms := url.Values{}
	parms.Add("requestJobDescription", string(requestJobDescriptionMarshal))
	body := strings.NewReader(parms.Encode())
	res, err := c.httpRequest("POST", body)
	if err != nil {
		log.Println("[CREATE ERROR]:", err)
		return "", err
	}
	resp := Response{}
	err = json.Unmarshal(res, &resp)
	if err != nil {
		log.Println("[CREATE ERROR]:", err)
		return "", err
	}
	if resp.ResponseCode!=200{
		log.Println("[CREATE ERROR]:", fmt.Errorf(string(res)))
		return "", fmt.Errorf(string(res))
	}
	return resp.PolicyId, nil
}

func (c *Client) GetPolicy(policyId string) (*Policy, error){
	credentials := Credentials{
		PartnerUserId: c.partnerUserId,
		PartnerUserSecret: c.partnerUserSecret,
	}
	inputSettings := InputSettings{
		Type: "policyList",
	}
	requestJobDescription := RequestJobDescription{
		Type: "get",
		Credential: credentials,
		InputSetting: inputSettings,
	}
	requestJobDescriptionMarshal, err := json.Marshal(requestJobDescription)
	if err != nil {
		log.Println("[READ ERROR]:", err)
		return nil, err
	}
	parms := url.Values{}
	parms.Add("requestJobDescription", string(requestJobDescriptionMarshal))
	body := strings.NewReader(parms.Encode())
	res, err := c.httpRequest("POST", body)
	if err != nil {
		log.Println("[READ ERROR]:", err)
		return nil, err
	}
	resp := Response{}
	err = json.Unmarshal(res, &resp)
	if err != nil {
		log.Println("[RAED ERROR]:", err)
		return nil, err
	}
	if resp.ResponseCode!=200{
		log.Println("[READ ERROR]:", fmt.Errorf(string(res)))
		return nil, fmt.Errorf(string(res))
	}
	for _, v := range resp.PolicyList {
		if v.PolicyId == policyId {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("policy does not exist")
}

func (c *Client) UpdatePolicy(policyId string, categories Categories, reportField ReportFields, tags Tags) (error){
	credentials := Credentials{
		PartnerUserId: c.partnerUserId,
		PartnerUserSecret: c.partnerUserSecret,
	}
	inputSettings := InputSettings{
		Type: "policy",
		PolicyIdList: []string{policyId,},
	}
	requestJobDescription := RequestJobDescription{
		Type: "update",
		Credential: credentials,
		InputSetting: inputSettings,
		TagArray: tags,
		Category: categories,
		ReportField: reportField,
	}
	requestJobDescriptionMarshal, err := json.Marshal(requestJobDescription)
	if err != nil {
		log.Println("[UPDATE ERROR]:", err)
		return err
	}
	parms := url.Values{}
	parms.Add("requestJobDescription", string(requestJobDescriptionMarshal))
	body := strings.NewReader(parms.Encode())
	res, err := c.httpRequest("POST", body)
	if err != nil {
		log.Println("[UPDATE ERROR]:", err)
		return err
	}
	resp := Response{}
	err = json.Unmarshal(res, &resp)
	if err != nil {
		log.Println("[UPDATE ERROR]:", err)
		return err
	}
	if resp.ResponseCode!=200{
		log.Println("[UPDATE ERROR]:", fmt.Errorf(string(res)))
		return fmt.Errorf(string(res))
	}
	return nil
}

func (c *Client) httpRequest(method string, body *strings.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, c.requestPath(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func (c *Client) requestPath() string {
	return fmt.Sprintf("https://integrations.expensify.com/Integration-Server/ExpensifyIntegrations")
}