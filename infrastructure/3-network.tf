## VPC
## ============

resource "aws_vpc" "vpc" {
  cidr_block = "10.10.0.0/16"
  enable_dns_support = true
  enable_dns_hostnames = true

  tags {
    Name = "${var.prefix}-${var.aws_region}-vpc"
    Terraformed = "true"
  }
}

## Subnet
## ============
resource "aws_subnet" "subnet1" {
  vpc_id = "${aws_vpc.vpc.id}"
  cidr_block = "10.10.1.0/24"
  availability_zone = "${var.aws_region}d"
  map_public_ip_on_launch = true

  tags {
    Name = "${var.prefix}-${var.aws_region}-subnet1"
    Terraformed = "true"
  }
}
resource "aws_subnet" "subnet2" {
  vpc_id = "${aws_vpc.vpc.id}"
  cidr_block = "10.10.2.0/24"
  availability_zone = "${var.aws_region}e"
  map_public_ip_on_launch = true

  tags {
    Name = "${var.prefix}-${var.aws_region}-subnet2"
    Terraformed = "true"
  }
}

resource "aws_internet_gateway" "internet_gateway" {
  vpc_id = "${aws_vpc.vpc.id}"

  tags {
    Name = "${var.prefix}-${var.aws_region}-internet-gateway"
    Terraformed = "true"
  }
}

resource "aws_route_table" "route_table" {
  vpc_id = "${aws_vpc.vpc.id}"
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.internet_gateway.id}"
  }

  tags {
    Name = "${var.prefix}-${var.aws_region}-route-table-public"
    Terraformed = "true"
  }
}

resource "aws_route_table_association" "public_subnet1" {
  subnet_id = "${aws_subnet.subnet1.id}"
  route_table_id = "${aws_route_table.route_table.id}"
}
resource "aws_route_table_association" "public_subnet2" {
  subnet_id = "${aws_subnet.subnet2.id}"
  route_table_id = "${aws_route_table.route_table.id}"
}
