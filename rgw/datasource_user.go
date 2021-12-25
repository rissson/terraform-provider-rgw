package rgw

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceUserRead,
		Schema:      schemaUser(),
	}
}
