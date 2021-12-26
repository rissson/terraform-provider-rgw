resource "rgw_user" "my_user" {
  user_id      = "my_user"
  display_name = "My User"
}

resource "rgw_quota" "user_my_user" {
  user_id     = rgw_user.my_user.user_id
  enabled     = true
  type        = "user"
  max_objects = 10000
  max_size    = -1
}
