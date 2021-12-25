package rgw

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "This data source can be used to retrieve information about a user.",
		ReadContext: resourceUserRead,
		Schema:      schemaUser(),
	}
}
