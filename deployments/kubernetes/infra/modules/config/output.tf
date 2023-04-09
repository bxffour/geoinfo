output "config_toml" {
  value = data.template_file.config_renderer.rendered
}

output "secret_toml" {
  value = data.template_file.secret_renderer.rendered
}
