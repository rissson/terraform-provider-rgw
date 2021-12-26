package rgw

import (
	"context"

	rgwadmin "github.com/ceph/go-ceph/rgw/admin"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RGW_ENDPOINT", nil),
			},
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("RGW_ACCESS_KEY", nil),
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("RGW_SECRET_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"rgw_user":  resourceUser(),
			"rgw_quota": resourceQuota(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"rgw_user": datasourceUser(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	endpoint := d.Get("endpoint").(string)
	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)

	var diags diag.Diagnostics

	if endpoint == "" || accessKey == "" || secretKey == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to find endpoint, access_key or secret_key",
			Detail:   "Those values must be set",
		})
		return nil, diags
	}

	// TODO: add support for configuring HTTPClient
	api, err := rgwadmin.New(endpoint, accessKey, secretKey, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create API client",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return api, diags
}
