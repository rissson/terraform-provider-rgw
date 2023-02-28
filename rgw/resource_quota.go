package rgw

import (
	"context"
	"fmt"

	rgwadmin "github.com/ceph/go-ceph/rgw/admin"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func schemaQuota() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Properties
		"user_id": {
			Description: "The ID of the user to set the quota for.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"type": {
			Description: "`user` or `bucket`",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"check_on_raw": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"max_size": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"max_size_kb": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"max_objects": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
	}
}

func resourceQuota() *schema.Resource {
	return &schema.Resource{
		Description:   "This resource can be used to set the quota for a rgw user. Refer to the Ceph RGW Admin Ops API documentation for values documentation. Upon deletion, quota is disabled.",
		CreateContext: resourceQuotaCreate,
		ReadContext:   resourceQuotaRead,
		UpdateContext: resourceQuotaUpdate,
		DeleteContext: resourceQuotaDelete,
		Schema:        schemaQuota(),
	}
}

func rgwQuotaFromSchemaQuota(d *schema.ResourceData) rgwadmin.QuotaSpec {
	enabled := d.Get("enabled").(bool)
	quota := rgwadmin.QuotaSpec{
		UID:        d.Get("user_id").(string),
		QuotaType:  d.Get("type").(string),
		Enabled:    &enabled,
		CheckOnRaw: d.Get("check_on_raw").(bool),
	}

	if maxSize, ok := d.GetOk("max_size"); ok {
		maxSize := int64(maxSize.(int))
		quota.MaxSize = &maxSize
	}

	if maxSizeKb, ok := d.GetOk("max_size_kb"); ok {
		maxSizeKb := maxSizeKb.(int)
		quota.MaxSizeKb = &maxSizeKb
	}

	if maxObjects, ok := d.GetOk("max_objects"); ok {
		maxObjects := int64(maxObjects.(int))
		quota.MaxObjects = &maxObjects
	}

	return quota
}

func flattenRgwQuota(quota rgwadmin.QuotaSpec, userID string) interface{} {
	q := map[string]interface{}{
		"type":         quota.QuotaType,
		"enabled":      quota.Enabled,
		"check_on_raw": quota.CheckOnRaw,
		"max_size":     quota.MaxSize,
		"max_size_kb":  quota.MaxSizeKb,
		"max_objects":  quota.MaxObjects,
	}
	if userID != "" {
		q["user_id"] = userID
	}
	return q
}

func resourceQuotaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*rgwadmin.API)
	var diags diag.Diagnostics

	quota := rgwQuotaFromSchemaQuota(d)

	err := api.SetUserQuota(ctx, quota)
	if err != nil {
		return diag.FromErr(err)
	}

	diags = resourceQuotaRead(ctx, d, m)

	return diags
}

func resourceQuotaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*rgwadmin.API)
	var diags diag.Diagnostics

	userID := d.Get("user_id").(string)
	user, err := api.GetUser(ctx, rgwadmin.User{ID: userID})
	if err != nil {
		return diag.FromErr(err)
	}

	quotaType := d.Get("type")

	var quota rgwadmin.QuotaSpec
	if quotaType == "user" {
		quota = user.UserQuota
	} else {
		quota = user.BucketQuota
	}

	id := fmt.Sprintf("%s_%s", quotaType, userID)
	d.SetId(id)

	for key, value := range flattenRgwQuota(quota, userID).(map[string]interface{}) {
		err := d.Set(key, value)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.Set("type", quotaType)

	return diags
}

func resourceQuotaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*rgwadmin.API)
	var diags diag.Diagnostics

	quota := rgwQuotaFromSchemaQuota(d)

	err := api.SetUserQuota(ctx, quota)
	if err != nil {
		return diag.FromErr(err)
	}

	diags = resourceQuotaRead(ctx, d, m)

	return diags
}

func resourceQuotaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*rgwadmin.API)
	var diags diag.Diagnostics

	quota := rgwQuotaFromSchemaQuota(d)
	f := false
	quota.Enabled = &f

	err := api.SetUserQuota(ctx, quota)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
