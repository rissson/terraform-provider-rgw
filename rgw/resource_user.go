package rgw

import (
	"context"

	rgwadmin "github.com/ceph/go-ceph/rgw/admin"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thoas/go-funk"
)

func schemaUser() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Properties
		"user_id": &schema.Schema{
			Description: "The ID the user is referred by.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"display_name": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"email": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"suspended": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
			Optional: true,
		},
		"max_buckets": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
			Optional: true,
		},
		// Only for creation and modification
		"generate_key": &schema.Schema{
			Description: "Only used for creation and modification. If true, a new key will be generated for the user. Default: true for creation, false for modification.",
			Type:        schema.TypeBool,
			Optional:    true,
		},
		"key_type": &schema.Schema{
			Description: "Only use for creation and modification when `generate_key` is true.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"user_caps": &schema.Schema{
			Description: "Only used to set user capabilities. To get user capabilities, use `caps` read-only attribute instead.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		// Only for deletion
		"purge_data": &schema.Schema{
			Description: "Only used when deleting the user. Check Ceph RGW Admin Ops API documentation for details.",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		// Computed
		"subusers": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"permissions": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"keys": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"user": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"access_key": &schema.Schema{
						Type:      schema.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"secret_key": &schema.Schema{
						Type:      schema.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},
		"swift_keys": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"user": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"secret_key": &schema.Schema{
						Type:      schema.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},
		"caps": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"perm": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"op_mask": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"default_placement": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"default_storage_class": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"placement_tags": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"bucket_quota": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: schemaQuota(),
			},
		},
		"user_quota": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: schemaQuota(),
			},
		},
		"type": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "This resource can be used to manage rgw users.",
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema:        schemaUser(),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func rgwUserFromSchemaUser(d *schema.ResourceData) rgwadmin.User {
	user := rgwadmin.User{
		ID: d.Get("user_id").(string),
	}

	if displayName, ok := d.GetOk("display_name"); ok {
		user.DisplayName = displayName.(string)
	}

	if email, ok := d.GetOk("email"); ok {
		user.Email = email.(string)
	}

	if suspended, ok := d.GetOk("suspended"); ok {
		suspended := suspended.(int)
		user.Suspended = &suspended
	}

	if maxBuckets, ok := d.GetOk("max_buckets"); ok {
		maxBuckets := maxBuckets.(int)
		user.MaxBuckets = &maxBuckets
	}

	return user
}

func flattenRgwKey(key rgwadmin.UserKeySpec) interface{} {
	return map[string]interface{}{
		"user":       key.User,
		"access_key": key.AccessKey,
		"secret_key": key.SecretKey,
	}
}

func flattenRgwUserCap(userCap rgwadmin.UserCapSpec) interface{} {
	return map[string]interface{}{
		"type": userCap.Type,
		"perm": userCap.Perm,
	}
}

func flattenRgwUser(user rgwadmin.User) interface{} {
	return map[string]interface{}{
		"user_id":               user.ID,
		"display_name":          user.DisplayName,
		"email":                 user.Email,
		"suspended":             user.Suspended,
		"max_buckets":           user.MaxBuckets,
		"subusers":              user.Subusers,
		"keys":                  funk.Map(user.Keys, flattenRgwKey),
		"swift_keys":            user.SwiftKeys,
		"caps":                  funk.Map(user.Caps, flattenRgwUserCap),
		"op_mask":               user.OpMask,
		"default_placement":     user.DefaultPlacement,
		"default_storage_class": user.DefaultStorageClass,
		"placement_tags":        user.PlacementTags,
		"bucket_quota":          []interface{}{flattenRgwQuota(user.BucketQuota, user.ID)},
		"user_quota":            []interface{}{flattenRgwQuota(user.UserQuota, user.ID)},
		"type":                  user.Type,
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*rgwadmin.API)
	var diags diag.Diagnostics

	user := rgwUserFromSchemaUser(d)

	if generateKey, ok := d.GetOk("generate_key"); ok {
		generateKey := generateKey.(bool)
		user.GenerateKey = &generateKey
	}

	if keyType, ok := d.GetOk("key_type"); ok {
		user.KeyType = keyType.(string)
	}

	if userCaps, ok := d.GetOk("user_caps"); ok {
		user.UserCaps = userCaps.(string)
	}

	user, err := api.CreateUser(ctx, user)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.ID)

	diags = resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*rgwadmin.API)
	var diags diag.Diagnostics

	userID, ok := d.GetOk("user_id")
	if ok {
		d.SetId(userID.(string))
	}

	user, err := api.GetUser(ctx, rgwadmin.User{ID: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}

	for key, value := range flattenRgwUser(user).(map[string]interface{}) {
		err := d.Set(key, value)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*rgwadmin.API)
	var diags diag.Diagnostics

	user := rgwUserFromSchemaUser(d)

	user, err := api.ModifyUser(ctx, user)
	if err != nil {
		return diag.FromErr(err)
	}

	diags = resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := m.(*rgwadmin.API)
	var diags diag.Diagnostics

	user := rgwadmin.User{
		ID: d.Id(),
	}

	if purgeData, ok := d.GetOk("purge_data"); ok {
		purgeData := purgeData.(int)
		user.PurgeData = &purgeData
	}

	err := api.RemoveUser(ctx, user)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
