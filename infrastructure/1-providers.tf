provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region     = "${var.aws_region}"
}

resource "aws_key_pair" "key" {
  key_name = "${var.prefix}-${var.aws_region}-key"
  public_key = "${file(var.aws_key_public)}"
}
