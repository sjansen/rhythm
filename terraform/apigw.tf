resource "aws_api_gateway_rest_api" "gw" {
  name = var.gw
  endpoint_configuration {
    types = ["REGIONAL"]
  }
}


resource "aws_api_gateway_deployment" "default" {
  depends_on = [
    "aws_api_gateway_integration.public",
  ]

  rest_api_id = join("", aws_api_gateway_rest_api.gw.*.id)
  stage_name  = "default"
}


resource "aws_api_gateway_method" "public" {
  authorization = "NONE"
  rest_api_id   = join("", aws_api_gateway_rest_api.gw.*.id)
  resource_id   = join("", aws_api_gateway_resource.proxy.*.id)
  http_method   = "ANY"
}


resource "aws_api_gateway_integration" "public" {
  rest_api_id = join("", aws_api_gateway_rest_api.gw.*.id)
  resource_id = join("", aws_api_gateway_method.public.*.resource_id)
  http_method = join("", aws_api_gateway_method.public.*.http_method)

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.fn.invoke_arn
}


resource "aws_api_gateway_resource" "proxy" {
  rest_api_id = join("", aws_api_gateway_rest_api.gw.*.id)
  parent_id   = join("", aws_api_gateway_rest_api.gw.*.root_resource_id)
  path_part   = "{proxy+}"
}
