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

variable "db_user" {
  description = "Postgresql database user"
}

variable "db_port" {
  description = "Postgresql database port"
}

variable "db_password" {
  description = "Postgresql database password"
}

variable "db_host" {
  description = "Postgresql database host"
}

variable "db_sslrootcert" {
  description = "Postgresql database ssl rootcert"
}

variable "db_user" {

}

variable "db_user" {

}

variable "db_user" {

}

variable "db_user" {

}
