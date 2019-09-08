resource "aws_s3_bucket" "logs" {
  bucket        = var.logs
  acl           = "log-delivery-write"
  force_destroy = false
  lifecycle_rule {
    id                                     = "cleanup"
    enabled                                = true
    abort_incomplete_multipart_upload_days = 3
    expiration {
      days = 90
    }
    noncurrent_version_expiration {
      days = 30
    }
  }
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
  versioning {
    enabled = var.protect_logs
  }
}


resource "aws_s3_bucket" "media" {
  bucket        = var.media
  acl           = "private"
  force_destroy = false
  lifecycle_rule {
    id                                     = "cleanup"
    enabled                                = true
    abort_incomplete_multipart_upload_days = 3
    expiration {
      expired_object_delete_marker = true
    }
    noncurrent_version_expiration {
      days = 30
    }
  }
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
  versioning {
    enabled = var.protect_media
  }
}


resource "aws_s3_bucket_policy" "media" {
  bucket = aws_s3_bucket.media.id
  policy = data.aws_iam_policy_document.media.json
}
