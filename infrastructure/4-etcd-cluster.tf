variable "etcd_ips" {
  default = {
    "0" = "10.10.1.100"
    "1" = "10.10.1.101"
    "2" = "10.10.1.102"
  }
}

resource "aws_instance" "etcd_cluster" {
  count = "3"
  ami = "${lookup(var.coreos_ami, var.aws_region)}"
  instance_type = "t2.small"
  key_name = "${aws_key_pair.key.key_name}"
  availability_zone = "${var.aws_az}"
  subnet_id = "${aws_subnet.subnet.id}"
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
    Name = "${var.prefix}-${var.aws_az}-etcd-${count.index}"
    Application = "${var.prefix}-etcd"
    Terraformed = "true"
  }

  user_data = <<EOF
#cloud-config

hostname: ${"etcd-cluster-${var.aws_az}-${count.index}"}
manage_etc_hosts: localhost

ssh_authorized_keys:
  - ${file(var.aws_key_public)}

coreos:
  fleet:
    metadata: app=etcd,region=${var.aws_region},az=${var.aws_az}
  etcd:
    name: ${"etcd-cluster-${count.index}"}
    discovery: ${var.etcd_discovery_url}
    addr: ${lookup(var.etcd_ips, count.index)}:4001
    peer-addr: ${lookup(var.etcd_ips, count.index)}:7001
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
    "1" = "10.10.1.201"
  }
}
resource "aws_instance" "compute_cluster" {
  count = "2"
  ami = "${lookup(var.coreos_ami, var.aws_region)}"
  instance_type = "m3.medium"
  key_name = "${aws_key_pair.key.key_name}"
  availability_zone = "${var.aws_az}"
  subnet_id = "${aws_subnet.subnet.id}"
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
    Name = "${var.prefix}-${var.aws_az}-compute-${count.index}"
    Application = "${var.prefix}-compute"
    Terraformed = "true"
  }

  user_data = <<EOF
#cloud-config

hostname: ${"compute-cluster-${var.aws_az}-${count.index}"}
manage_etc_hosts: localhost

ssh_authorized_keys:
  - ${file(var.aws_key_public)}

coreos:
  fleet:
    metadata: app=compute,region=${var.aws_region},az=${var.aws_az},web=true,worker=true
  etcd:
    proxy: on
    discovery: ${var.etcd_discovery_url}
  units:
    - name: etcd.service
      command: start
    - name: fleet.service
      command: start
EOF
}
