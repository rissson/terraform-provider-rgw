package main

import (
	"context"
	"terraform-provider-rgw/rgw"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func main() {
	tfsdk.Serve(context.Background(), rgw.New, tfsdk.ServeOpts{
		Name: "rgw",
	})
}
