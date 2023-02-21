resource "prowlarr_download_client_torrent_blackhole" "example" {
  enable                = true
  priority              = 1
  name                  = "Example"
  magnet_file_extension = ".magnet"
  torrent_folder        = "/torrent/"
}