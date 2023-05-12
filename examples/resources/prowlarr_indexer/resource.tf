resource "prowlarr_indexer" "example" {
  enable          = true
  name            = "HDBits"
  implementation  = "HDBits"
  config_contract = "HDBitsSettings"
  protocol        = "torrent"
  tags            = [1, 2, 5]

  fields = [
    {
      name       = "username"
      text_value = "test"
    },
    {
      name       = "apiKey"
      text_value = "test"
    },
    {
      name      = "codecs"
      set_value = [1, 5]
    },
    {
      name      = "mediums"
      set_value = [1, 3]
    },
    {
      name         = "torrentBaseSettings.seedRatio"
      number_value = 0.5
    },
    {
      name         = "torrentBaseSettings.seedTime"
      number_value = 5
    },
  ]
}