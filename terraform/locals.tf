locals {
  apigw_url_parts = split("/", aws_api_gateway_deployment.default.invoke_url)
}
