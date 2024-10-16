resource "prowlarr_host" "test" {
  launch_browser  = true
  port            = 9696
  url_base        = ""
  bind_address    = "*"
  application_url = ""
  instance_name   = "Prowlarr"
  proxy = {
    enabled = false
  }
  ssl = {
    enabled                = false
    certificate_validation = "enabled"
  }
  logging = {
    log_level      = "info"
    log_size_limit = 1
  }
  backup = {
    folder    = "/backup"
    interval  = 5
    retention = 10
  }
  authentication = {
    method   = "none"
    required = "disabledForLocalAddresses"
  }
  update = {
    mechanism = "docker"
    branch    = "develop"
  }
}
