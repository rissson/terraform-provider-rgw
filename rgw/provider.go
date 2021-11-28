package rgw

import (
	"context"
	"os"

	rgwadmin "github.com/ceph/go-ceph/rgw/admin"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var stderr = os.Stderr

// New provider
func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	api        *rgwadmin.API
}

func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"endpoint": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"access_key": {
				Type:      types.StringType,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"secret_key": {
				Type:      types.StringType,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
		},
	}, nil
}

type providerData struct {
	Endpoint  types.String `tfsdk:"endpoint"`
	AccessKey types.String `tfsdk:"access_key"`
	SecretKey types.String `tfsdk:"secret_key"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var config providerData
	var err error

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var endpoint string
	if config.Endpoint.Unknown {
		resp.Diagnostics.AddWarning(
			"Unable to create API client",
			"Cannot use unknown value as endpoint",
		)
		return
	}

	if config.Endpoint.Null {
		endpoint = os.Getenv("RGW_ENDPOINT")
	} else {
		endpoint = config.Endpoint.Value
	}

	if endpoint == "" {
		resp.Diagnostics.AddError(
			"Unable to find endpoint",
			"Endpoint cannot be an empty string",
		)
		return
	}

	var accessKey string
	if config.AccessKey.Unknown {
		resp.Diagnostics.AddWarning(
			"Unable to create API client",
			"Cannot use unknown value as access_key",
		)
		return
	}

	if config.AccessKey.Null {
		accessKey = os.Getenv("RGW_ACCESS_KEY")
	} else {
		accessKey = config.AccessKey.Value
	}

	if accessKey == "" {
		resp.Diagnostics.AddError(
			"Unable to find access_key",
			"Access key cannot be an empty string",
		)
		return
	}

	var secretKey string
	if config.SecretKey.Unknown {
		resp.Diagnostics.AddWarning(
			"Unable to create API client",
			"Cannot use unknown value as secret_key",
		)
		return
	}

	if config.SecretKey.Null {
		secretKey = os.Getenv("RGW_secret_KEY")
	} else {
		secretKey = config.SecretKey.Value
	}

	if secretKey == "" {
		resp.Diagnostics.AddError(
			"Unable to find secret_key",
			"Secret key cannot be an empty string",
		)
		return
	}

	// TODO: add support for configuring HTTPClient
	api, err := rgwadmin.New(endpoint, accessKey, secretKey, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create API client",
			err.Error(),
		)
		return
	}

	p.api = api
	p.configured = true
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{}, nil
}

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}
