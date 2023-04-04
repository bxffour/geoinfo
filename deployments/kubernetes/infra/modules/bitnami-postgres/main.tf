# resource "kubernetes_persistent_volume_claim" "postgresql-pvc" {
#   metadata = {
#     name       = "postgres-pvc"
#     finalizers = null
#   }

#   spec {
#     access_modes = ["ReadWriteOnce"]
#     resources {
#       requests = {
#         storage = "2Gi"
#       }
#     }

#     storage_class_name = "standard"
#   }
# }

resource "helm_release" "bitmami-postgresql" {
  name       = var.name
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "postgresql"
  version    = "12.2.2"

  set_sensitive {
    name  = "auth.postgresPassword"
    value = var.postgresPassword
  }

  set {
    name  = "auth.database"
    value = var.database
  }

  set {
    name  = "tls.enabled"
    value = var.tls.enabled
  }

  set {
    name  = "tls.certificatesSecret"
    value = var.tls.secret
  }

  set {
    name  = "tls.certFilename"
    value = var.tls.cert
  }

  set {
    name  = "tls.certKeyFilename"
    value = var.tls.key
  }

  set {
    name  = "tls.certCAFilename"
    value = var.tls.ca
  }

  set {
    name  = "primary.persistence.storageClass"
    value = var.storageClass
  }

  set {
    name  = "primary.persistence.size"
    value = var.storageSize
  }

  set {
    name  = "image.tag"
    value = "15.2.0-debian-11-r5"
  }
}
