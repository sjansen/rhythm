resource "aws_ssm_parameter" "secret" {
  name        = "/${var.secret_name}"
  description = "TODO"
  type        = "SecureString"
  value       = "invalid"
  overwrite   = false

  lifecycle {
    ignore_changes = [value]
  }
}
