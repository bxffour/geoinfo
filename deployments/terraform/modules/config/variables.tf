variable "config_toml" {
  description = "path to the config.toml.tpl file"
  type        = string
}

variable "secret_toml" {
  description = "path to the secret.toml.tpl file"
  type        = string
}

variable "pgclient_certs" {
  description = "path to postgres client certs to be added as kubernetes a secret"
  type = object({
    ca_file = optional(string, "")
    sslkey  = optional(string, "")
    sslcert = optional(string, "")
  })
}

variable "pgserver_certs" {
  description = "path to postgres server certs to be added as kubernetes a secret"
  type = object({
    ca_file = optional(string, "")
    sslkey  = optional(string, "")
    sslcert = optional(string, "")
  })
}

variable "api_config" {
  description = "Geoinfo configuration parameters"

  type = object({
    port           = optional(number, 8080)
    env            = optional(string, "development")
    max_open_conns = optional(number, 25)
    max_idle_conns = optional(number, 25)
    max_idle_time  = optional(string, "15m")
    rps            = optional(number, 4)
    burst          = optional(number, 8)
    enabled        = optional(bool, false)
  })

  default = {}
}

variable "db_config" {
  description = "Postgresql database config parameter"

  type = object({
    user        = string
    port        = optional(number, 5432)
    password    = string
    host        = string
    dbname      = string
    sslrootcert = optional(string, "")
    sslcert     = optional(string, ""),
    sslkey      = optional(string, ""),
    sslmode     = optional(string, "disable"),
  })
}
