variable "etcd_ips" {
  default = {
    "0" = "10.10.1.100"
    "1" = "10.10.1.101"
    "2" = "10.10.1.102"
    "3" = "10.10.2.100"
    "4" = "10.10.2.101"
    "5" = "10.10.2.102"
  }
}
variable "etcd_subnets" {
  default = {
    "0" = "${aws_subnet.subnet1.id}"
    "1" = "${aws_subnet.subnet1.id}"
    "2" = "${aws_subnet.subnet1.id}"
    "3" = "${aws_subnet.subnet2.id}"
    "4" = "${aws_subnet.subnet2.id}"
    "5" = "${aws_subnet.subnet2.id}"
  }
}
variable "etcd_azs" {
  default = {
    "0" = "d"
    "1" = "d"
    "2" = "d"
    "3" = "e"
    "4" = "e"
    "5" = "e"
  }
}

resource "aws_instance" "etcd_cluster" {
  count = "6"
  ami = "${lookup(var.coreos_ami, var.aws_region)}"
  instance_type = "t2.small"
  key_name = "${aws_key_pair.key.key_name}"
  availability_zone = "${var.aws_region}${lookup(var.etcd_azs, count.index)}"
  subnet_id = "${lookup(var.etcd_subnets, count.index)}"
  private_ip = "${lookup(var.etcd_ips, count.index)}"
  associate_public_ip_address = true
  security_groups = [
    "${aws_security_group.maintenance.id}",
    "${aws_security_group.internal.id}"
  ]

  connection {
    user = "core"
    key_file = "${var.key_path}"
  }

  tags {
    Name = "${var.prefix}-${var.aws_region}-etcd-${count.index}"
    Application = "${var.prefix}-etcd"
    Terraformed = "true"
  }

  user_data = <<EOF
#cloud-config

hostname: ${"etcd-cluster-${count.index}"}
manage_etc_hosts: localhost

ssh_authorized_keys:
  - ${file(var.aws_key_public)}

coreos:
  fleet:
    metadata: app=etcd,region=${var.aws_region}
  etcd:
    name: ${"etcd-cluster-${count.index}"}
    discovery: ${var.etcd_discovery_url}
    addr: ${lookup(var.etcd_ips, count.index)}:3201
    peer-addr: ${lookup(var.etcd_ips, count.index)}:3301
  units:
    - name: etcd.service
      command: start
    - name: fleet.service
      command: start
EOF
}

variable "compute_ips" {
  default = {
    "0" = "10.10.1.200"
    "1" = "10.10.2.200"
  }
}
variable "compute_subnets" {
  default = {
    "0" = "${aws_subnet.subnet1.id}"
    "1" = "${aws_subnet.subnet2.id}"
  }
}
variable "compute_azs" {
  default = {
    "0" = "${var.aws_region}d"
    "1" = "${var.aws_region}e"
  }
}
resource "aws_instance" "compute_cluster" {
  count = "2"
  ami = "${lookup(var.coreos_ami, var.aws_region)}"
  instance_type = "m3.medium"
  key_name = "${aws_key_pair.key.key_name}"
  availability_zone = "${lookup(var.compute_azs, count.index)}"
  subnet_id = "${lookup(var.compute_subnets, count.index)}"
  private_ip = "${lookup(var.compute_ips, count.index)}"
  associate_public_ip_address = true
  security_groups = [
    "${aws_security_group.maintenance.id}",
    "${aws_security_group.internal.id}"
  ]

  connection {
    user = "core"
    key_file = "${var.key_path}"
  }

  tags {
    Name = "${var.prefix}-${var.aws_region}-compute-${count.index}"
    Application = "${var.prefix}-compute"
    Terraformed = "true"
  }

  user_data = <<EOF
#cloud-config

hostname: ${"compute-cluster-${count.index}"}
manage_etc_hosts: localhost

ssh_authorized_keys:
  - ${file(var.aws_key_public)}

coreos:
  fleet:
    metadata: app=compute,region=${var.aws_region},web=true,worker=true
  etcd:
    proxy: on
    listen-client-urls: ${lookup(var.etcd_ips, count.index)}:3201
    discovery: ${var.etcd_discovery_url}
  units:
    - name: etcd.service
      command: start
    - name: fleet.service
      command: start
EOF
}
