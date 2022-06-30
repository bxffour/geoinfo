terraform {
  required_providers {
      linode = {
          source = "linode/linode"
          version = "1.27.2"
      }
  }
}

# Using the Linode Provider
provider "linode" {
    token = var.token
}

# Creating linode_lke_cluster resource to provision 
# Kubernetes cluster
resource "linode_lke_cluster" "k8s_Cluster" {
  k8s_version = var.k8s_version
  label = var.label
  region = var.region
  tags = var.tags

  dynamic "pool" {
      for_each = var.pools
      content {
          type = pool.value["type"]
          count = pool.value["count"]
      }
  }
}

# exporting the cluster's attributes
output "kubeconfig" {
  value = linode_lke_cluster.k8s_Cluster.kubeconfig
  sensitive = true
}

output "api_endpoints" {
  value = linode_lke_cluster.k8s_Cluster.api_endpoints
}

output "status" {
  value = linode_lke_cluster.k8s_Cluster.status
}

output "id" {
  value = linode_lke_cluster.k8s_Cluster.id
}

output "pool" {
  value = linode_lke_cluster.k8s_Cluster.pool
}
