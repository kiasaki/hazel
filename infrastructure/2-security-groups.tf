resource "aws_security_group" "internal" {
  name = "${var.prefix}-${var.aws_region}-internal"
  description = "VPC internal traffic"
  vpc_id = "${aws_vpc.vpc.id}"

  ingress {
    from_port = 0
    to_port = 65535
    protocol = "tcp"
    cidr_blocks = ["10.10.0.0/16"]
  }
  ingress {
    from_port = 0
    to_port = 65535
    protocol = "udp"
    cidr_blocks = ["10.10.0.0/16"]
  }
}

resource "aws_security_group" "maintenance" {
  name = "${var.prefix}-${var.aws_region}-maintenance"
  description = "Cluster maintenance (SSH)"
  vpc_id = "${aws_vpc.vpc.id}"

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
