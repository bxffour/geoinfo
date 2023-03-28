provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

module "secrets" {
  source = "./config"

  port = 7077
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
