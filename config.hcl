app_name = "hermes"

app {

  env = {
    PRIMARY_REGION = "lhr"
    DATABASE_URL = "hermes.db"
    ENVIRONMENT = "staging"
    PORT = "8080"
    HERMES_SENDER_IDENTITY="hermes"
    HONEYCOMB_DSN="api.honeycomb.io"
  }

  port = 8080

  compute {
    cpu      = 1
    memory   = 256
    cpu_kind = "shared"
  }

  process {
    name = "hermes"
  }

  storage {
    name        = "data"
    destination = "/data"
  }
}