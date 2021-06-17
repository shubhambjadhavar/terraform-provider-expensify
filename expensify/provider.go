package expensify

import(
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"terraform-provider-expensify/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"expensify_user_id": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("EXPENSIFY_USER_ID", ""),
			},
			"expensify_user_secret": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("EXPENSIFY_USER_SECRET", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"expensify_user": resourceUser(),
			"expensify_policy": resourcePolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"expensify_user": dataSourceUser(),
			"expensify_policy": dataSourcePolicy(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	expensifyUserId := d.Get("expensify_user_id").(string)
	expensifyUserSecret := d.Get("expensify_user_secret").(string)
	var diags diag.Diagnostics
	return client.NewClient(expensifyUserId, expensifyUserSecret), diags
}