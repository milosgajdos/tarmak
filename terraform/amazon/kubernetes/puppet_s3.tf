resource "aws_s3_bucket_object" "puppet-tar-gz" {
  key          = "${data.template_file.stack_name.rendered}/puppet.tar.gz"
  bucket       = "${data.terraform_remote_state.hub_state.secrets_bucket.0}"
  content_type = "application/tar+gzip"
  source       = "puppet.tar.gz"
  etag         = "${md5(file("puppet.tar.gz"))}"
}
