variable "aws_access_key"  {}
variable "aws_secret_key"  {}
variable "aws_region"  {}
variable "app_slug"  {}
variable "build_id"  {}
variable "build_ami"  {}

provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region = "${var.aws_region}"
}

resource "aws_security_group" "service" {
  name = "${var.app_slug}-${var.build_id}-sg"
  description = "Allow all internal traffic and maintenance"

  ingress { # Maintenance
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress { # Internal
    from_port = 0
    to_port = 65535
    protocol = "-1"
    self = true
  }
}

resource "aws_elb" "service" {
  name = "${var.app_slug}-${var.build_id}-elb"
  availability_zones = ["${aws_instance.service.*.availability_zone}"]

  listener {
    instance_port = 8000
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }

  instances = ["${aws_instance.service.*.id}"]
}

resource "aws_instance" "service" {
  instance_type = "t2.micro"
  ami = "${var.build_ami}"
  security_groups = ["${aws_security_group.service.name}"]

  count = 2
  lifecycle = {
    create_before_destroy = true
  }
}
