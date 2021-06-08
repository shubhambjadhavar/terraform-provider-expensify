package expensify

import(
	"os"
	"log"
	"io/ioutil"
	"testing"
	"github.com/clarketm/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
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
	os.Setenv("PARTNER_USER_ID", res["PARTNER_USER_ID"].(string))
	os.Setenv("PARTNER_USER_SECRET", res["PARTNER_USER_SECRET"].(string))
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"expensify": testAccProvider,
	}
}

func TestExpensifyProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestExpensifyProvider_impl(t *testing.T)  {
	var _ *schema.Provider = Provider()
}

func testAccExpensifyProviderPreCheck(t *testing.T) {
	if v := os.Getenv("PARTNER_USER_ID"); v == "" {
		t.Fatal("PARTNER_USER_ID must be set for acceptance tests")
	}
	if v := os.Getenv("PARTNER_USER_SECRET"); v == "" {
		t.Fatal("PARTNER_USER_SECRET must be set for acceptance tests")
	}
}