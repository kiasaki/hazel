variable "coreos_ami" {
  description = "CoreOS ami's across regions"
  default = {
    eu-central-1   =  "ami-8ec1f293"
    ap-northeast-1 =  "ami-e85c46e9"
    sa-east-1      =  "ami-2de95630"
    ap-southeast-2 =  "ami-4dd3a777"
    ap-southeast-1 =  "ami-72dcf620"
    us-east-1      =  "ami-8297d4ea"
    us-west-2      =  "ami-f1702bc1"
    us-west-1      =  "ami-24b5ad61"
    eu-west-1      =  "ami-5d911f2a"
  }
}

variable "prefix" {
  default = "c" # For cluster
  description = "Ressources prefix to avoid collisions"
}

variable "aws_region" {
  default = "us-east-1"
  description = "The region of AWS, for AMI lookups."
}
variable "aws_access_key" {
  description = "AWS Access key."
}
variable "aws_secret_key" {
  description = "AWS Secret key."
}
variable "aws_key_public" {
  description = "Path to public key."
}
variable "aws_zone_id" {
  description = "ZoneID of TLD"
}
variable "aws_zone_tld" {
  description = "Zone TLD"
}

variable "etcd_discovery_url" {
  description = "etcd discovery url for the CoreOS cluster"
}
