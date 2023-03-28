provider "kubernetes" {
  config_path = "~/.kube/config"
  # config_context = "kind-kind"
}

data "template_file" "config_renderer" {
  template = file("${path.module}/config.toml.tpl")

  vars = {
    port           = var.port
    environment    = var.env
    max_open_conns = var.max_open_conns
    max_idle_conns = var.max_idle_conns
    max_idle_time  = var.max_idle_time
    rps            = var.rps
    burst          = var.burst
    enabled        = var.enabled
  }
}

data "template_file" "secret_renderer" {
  template = file("${path.module}/secrets.toml.tpl")

  vars = {
    db_user        = ""
    db_port        = 5432
    db_password    = ""
    db_dbname      = ""
    db_host        = "localhost"
    db_sslrootcert = ""
    db_sslcert     = ""
    db_sslkey      = ""
    db_sslmode     = ""
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
