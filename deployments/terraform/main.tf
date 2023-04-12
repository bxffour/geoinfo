provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
  # config_context = "kind-kind"
}

module "secrets" {
  source = "./modules/config"

  config_toml = fileexists(var.templates.config) ? file(var.templates.config) : ""
  secret_toml = fileexists(var.templates.creds) ? file(var.templates.creds) : ""

  pgclient_certs = {
    ca_file = fileexists(var.pgclient_certs.ca_file) ? file(var.pgclient_certs.ca_file) : ""
    sslkey  = fileexists(var.pgclient_certs.sslkey) ? file(var.pgclient_certs.sslkey) : ""
    sslcert = fileexists(var.pgclient_certs.sslcert) ? file(var.pgclient_certs.sslcert) : ""
  }

  pgserver_certs = {
    ca_file = fileexists(var.pgserver_certs.ca_file) ? file(var.pgserver_certs.ca_file) : ""
    sslkey  = fileexists(var.pgserver_certs.sslkey) ? file(var.pgserver_certs.sslkey) : ""
    sslcert = fileexists(var.pgserver_certs.sslcert) ? file(var.pgserver_certs.sslcert) : ""
  }

  db_config = {
    user     = "crest"
    port     = 5432
    password = "nana"
    host     = "localhost"
    dbname   = "crest_countries"
  }
}

module "bitnami_postgresql" {
  source = "./modules/bitnami-postgres"

  depends_on = [module.secrets]

  postgresPassword = var.postgresPassword
}

resource "helm_release" "geoinfo" {
  name             = "geoinfo-api"
  chart            = "../geoinfo"
  namespace        = "default"
  create_namespace = true
  cleanup_on_fail  = true
  lint             = true

  values = [
    "${file("../geoinfo/values.yaml")}"
  ]
}
