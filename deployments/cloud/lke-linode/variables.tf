variable "token" {
  description = "The linode API personal access token. (required)"
}

variable "k8s_version" {
  description = "Kubernetes version to use for this cluster. (required)"
  default = "1.23"
}

variable "label" {
  description = "The unique label to assign to this cluster. (required)"
  default = "default_lke_cluster"
}

variable "region" {
  description = "The region where the cluster will be located at. (required)"
  default = "eu-central"
}

variable "tags" {
  description = "Tags to apply to the cluster for organisational purposes. (optional)"
  type = list(string)
  default = [ "testing", "k8s_cluster" ]
}

variable "pools" {
  description = "The node pool specifications for the K8s cluster. (required)"
  type = list(object({
      type = string
      count = number
  }))

  default = [ {
    count = 2
    type = "g6-standard-2"
  } ]
}