variable "templates" {
  description = "paths to template files"
  type = object({
    config = optional(string, "templates/config.yaml.tpl")
    creds  = optional(string, "templates/secrets.yaml.tpl")
  })

  default = {}
}

variable "port" {
  type        = number
  description = "application port"
  default     = 8080
}

variable "env" {
  description = "application environment"
  default     = "development"
}

variable "max_open_conns" {
  type        = number
  description = "postgresql maximum database connections"
  default     = 25
}

variable "max_idle_conns" {
  type        = number
  description = "application port"
  default     = 25
}

variable "max_idle_time" {
  description = "application port"
  default     = "15m"
}

variable "rps" {
  type        = number
  description = "application port"
  default     = 4
}

variable "burst" {
  type        = number
  description = "application port"
  default     = 8
}

variable "enabled" {
  type        = bool
  description = "application port"
  default     = false
}

variable "pgclient_certs" {
  description = "path to postgres client certs to be added as kubernetes a secret"
  type = object({
    ca_file = optional(string, "certs/ca.pem")
    sslkey  = optional(string, "certs/pg-client-key.pem")
    sslcert = optional(string, "certs/pg-client.pem")
  })

  default = {}
}

variable "pgserver_certs" {
  description = "path to postgres server certs to be added as kubernetes a secret"
  type = object({
    ca_file = optional(string, "certs/ca.pem")
    sslkey  = optional(string, "certs/pg-server-key.pem")
    sslcert = optional(string, "certs/pg-server.pem")
  })

  default = {}
}

variable "postgresPassword" {
  description = "Password for the postgre admin user."
  sensitive   = true
}

variable "database" {
  description = "Name for a custom database to create"
  default     = "geoinfo"
}

variable "storageClass" {
  description = "PVC Storage Class for PostgreSQL Primary data volume"
  default     = "standard"
}

variable "storageSize" {
  description = "PVC Storage Request for PostgreSQL volume"
  default     = "2Gi"
}

variable "tls" {
  description = "Postgres tls config parameters"
  type = object({
    enabled = optional(bool, false)
    secret  = optional(string, "pgserver-ssl")
    cert    = optional(string, "tls.crt")
    key     = optional(string, "tls.key")
    ca      = optional(string, "ca.pem")
  })

  default = {}
}

