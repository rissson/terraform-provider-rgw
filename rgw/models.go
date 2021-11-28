package rgw

import "github.com/hashicorp/terraform-plugin-framework/types"

// User -
type User struct {
	// Properties
	ID          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
	Email       types.String `tfsdk:"email"`
	Suspended   types.Number `tfsdk:"suspended"`
	MaxBuckets  types.Number `tfsdk:"max_buckets"`
	// Only for creation and modification
	GenerateKey types.Bool   `tfsdk:"generate_key"`
	KeyType     types.String `tfsdk:"key_type"`
	UserCaps    types.String `tfsdk:"user_caps"`
	// Only for deletion
	PurgeData types.Bool `tfsdk:"purge_data"`
	// Computed
	Subusers            []Subuser      `tfsdk:"subusers"`
	Keys                []Key          `tfsdk:"keys"`
	SwiftKeys           []SwiftKey     `tfsdk:"swift_keys"`
	Caps                []Cap          `tfsdk:"caps"`
	OpMask              types.String   `tfsdk:"op_mask"`
	DefaultPlacement    types.String   `tfsdk:"default_placement"`
	DefaultStorageClass types.String   `tfsdk:"default_storage_class"`
	PlacementTags       []types.String `tfsdk:"placement_tags"`
	BucketQuota         Quota          `tfsdk:"bucket_quota"`
	UserQuota           Quota          `tfsdk:"user_quota"`
	TempURLKeys         []interface{}  `tfsdk:"temp_url_keys"`
	Type                types.String   `tfsdk:"type"`
	MfaIds              []interface{}  `tfsdk:"mfa_ids"`
}

// Subuser -
type Subuser struct {
	ID          types.String `tfsdk:"id"`
	Permissions types.String `tfsdk:"permissions"`
}

// Key -
type Key struct {
	AccessKey types.String `tfsdk:"access_key"`
	SecretKey types.String `tfsdk:"secret_key"`
}

// SwiftKey -
type SwiftKey struct {
	User      types.String `tfsdk:"user"`
	SecretKey types.String `tfsdk:"secret_key"`
}

// Cap -
type Cap struct {
	Type types.String `tfsdk:"type"`
	Perm types.String `tfsdk:"perm"`
}

// Quota -
type Quota struct {
	UID        types.String `tfsdk:"user_id"`
	QuotaType  types.String `tfsdk:"quota_type"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	MaxSize    types.Int64  `tfsdk:"max_size"`
	MaxSizeKb  types.Number `tfsdk:"max_size_kb"`
	MaxObjects types.Int64  `tfsdk:"max_objects"`
	// Computed
	CheckOnRaw types.Bool `tfsdk:"check_on_raw"`
}
