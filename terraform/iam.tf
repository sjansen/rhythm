data "aws_iam_policy_document" "media" {
  statement {
    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.media.arn}/*"]
    principals {
      type        = "AWS"
      identifiers = [aws_cloudfront_origin_access_identity.cdn.iam_arn]
    }
  }
}


resource "aws_iam_policy" "fn-logs" {
  name = "${var.fn}-logs"
  path = "/"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    }
  ]
}
EOF
}


resource "aws_iam_policy" "fn-ssm" {
  name = "${var.fn}-ssm"
  path = "/"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ssm:GetParameters"
      ],
      "Resource": [
        "arn:aws:ssm:*:*:parameter/${var.secret_name}"
      ]
    }
  ]
}
EOF
}


resource "aws_iam_policy" "fn-xray" {
    name = "${var.fn}-xray"
    policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": {
        "Effect": "Allow",
        "Action": [
            "xray:PutTraceSegments",
            "xray:PutTelemetryRecords"
        ],
        "Resource": [
            "*"
        ]
    }
}
EOF
}


resource "aws_iam_role" "fn" {
  name = var.fn

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      }
    }
  ]
}
EOF
}


resource "aws_iam_role_policy_attachment" "fn-logs" {
  policy_arn = aws_iam_policy.fn-logs.arn
  role = aws_iam_role.fn.name
}


resource "aws_iam_role_policy_attachment" "fn-ssm" {
  policy_arn = aws_iam_policy.fn-ssm.arn
  role = aws_iam_role.fn.name
}


resource "aws_iam_role_policy_attachment" "fn-xray" {
  policy_arn = aws_iam_policy.fn-xray.arn
  role = aws_iam_role.fn.name
}
