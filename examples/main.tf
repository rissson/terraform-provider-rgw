terraform {
  required_providers {
    rgw = {
      source = "risson/rgw"
    }
  }
}

provider "rgw" {
  endpoint = "https://s3.example.org"
}

resource "rgw_user" "my_user" {
  uid          = "my_user"
  display_name = "My User"
}
