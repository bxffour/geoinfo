data "template_file" "config_renderer" {
  template = var.config_toml

  vars = {
    port           = var.api_config.port
    environment    = var.api_config.env
    max_open_conns = var.api_config.max_open_conns
    max_idle_conns = var.api_config.max_idle_conns
    max_idle_time  = var.api_config.max_idle_time
    rps            = var.api_config.rps
    burst          = var.api_config.burst
    enabled        = var.api_config.enabled
  }
}

data "template_file" "secret_renderer" {
  template = var.secret_toml

  vars = {
    db_user        = var.db_config.user
    db_port        = var.db_config.port
    db_password    = var.db_config.password
    db_dbname      = var.db_config.host
    db_host        = var.db_config.dbname
    db_sslrootcert = var.db_config.sslrootcert
    db_sslcert     = var.db_config.sslcert
    db_sslkey      = var.db_config.sslkey
    db_sslmode     = var.db_config.sslmode
  }
}

resource "kubernetes_secret" "config_file" {
  metadata {
    name = "geoinfo-config"
  }

  data = {
    "config.toml" = data.template_file.config_renderer.rendered
  }
}

resource "kubernetes_secret" "secret_file" {
  metadata {
    name = "geoinfo-secret"
  }

  data = {
    "secrets.toml" = data.template_file.secret_renderer.rendered
  }
}

resource "kubernetes_secret" "pgclient_certs" {
  count = length(compact([
    for value in values(var.pgclient_certs) :
    value
  ])) > 0 ? 1 : 0


  metadata {
    name = "pgclient-ssl"
  }

  type = "kubernetes.io/tls"

  data = {
    "ca.pem"  = var.pgclient_certs.ca_file
    "tls.key" = var.pgclient_certs.sslkey
    "tls.crt" = var.pgclient_certs.sslcert
  }
}

resource "kubernetes_secret" "pgserver_certs" {
  count = length(compact([
    for value in values(var.pgserver_certs) :
    value
  ])) > 0 ? 1 : 0

  metadata {
    name = "pgserver-ssl"
  }

  type = "kubernetes.io/tls"

  data = {
    "ca.pem"  = var.pgserver_certs.ca_file
    "tls.key" = var.pgserver_certs.sslkey
    "tls.crt" = var.pgserver_certs.sslcert
  }
}

