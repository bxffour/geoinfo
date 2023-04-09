variable "name" {
  description = "helm release name"
  default     = "postgresql"
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

