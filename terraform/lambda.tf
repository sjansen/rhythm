resource "aws_lambda_function" "fn" {
  function_name    = var.fn
  filename         = "../server/build/lambda-handler.zip"
  handler          = "lambda-handler"
  source_code_hash = filebase64sha256("../server/build/lambda-handler.zip")
  role             = aws_iam_role.fn.arn

  runtime     = "go1.x"
  memory_size = 128
  timeout     = 15

  environment {
    variables = {
      RHYTHM_SECRET_NAME = "/${var.secret_name}"
    }
  }

  tracing_config {
    mode = "Active"
  }
}


resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowExecutionFromAPIGW"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.fn.arn
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${join("", aws_api_gateway_deployment.default.*.execution_arn)}*/*/*"
}
