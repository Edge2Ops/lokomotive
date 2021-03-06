output "kubeconfig-admin" {
  value = module.bootkube.kubeconfig-admin
}

output "kubeconfig" {
  value = module.bootkube.kubeconfig-kubelet
}

output "machine_domain" {
  value = var.machine_domain
}

output "cluster_name" {
  value = var.cluster_name
}

output "ssh_keys" {
  value = var.ssh_keys
}

output "libvirtpool" {
  value = libvirt_pool.volumetmp.name
}

output "libvirtbaseid" {
  value = libvirt_volume.base.id
}

# values.yaml content for all deployed charts.
output "pod-checkpointer_values" {
  value = module.bootkube.pod-checkpointer_values
}

output "kube-apiserver_values" {
  value = module.bootkube.kube-apiserver_values
}

output "kubernetes_values" {
  value = module.bootkube.kubernetes_values
}

output "kubelet_values" {
  value = module.bootkube.kubelet_values
}

output "calico_values" {
  value = module.bootkube.calico_values
}
