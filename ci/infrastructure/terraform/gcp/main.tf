variable "environment_name" {
  type = "string"
}

variable "gcp_service_account_key" {
	type = "string"
}

variable "projectid" {
  type = "string"
  default = "cf-backup-and-restore"
}

variable "region" {
  type = "string"
  default = "europe-west1"
}

variable "zone1" {
  type = "string"
  default = "europe-west1-b"
}

provider "google" {
  credentials = "${var.gcp_service_account_key}"
  project = "${var.projectid}"
  region = "${var.region}"
}

resource "google_compute_address" "jumpbox" {
  name = "${var.environment_name}-jumpbox"
}

resource "google_compute_network" "director" {
  name = "${var.environment_name}"
}

resource "google_compute_subnetwork" "director" {
  name          = "${var.environment_name}-director"
  ip_cidr_range = "10.0.0.0/24"
  network       = "${google_compute_network.director.self_link}"
}

resource "google_compute_firewall" "external" {
  name    = "${var.environment_name}-external"
  network = "${google_compute_network.director.name}"

  source_ranges = ["0.0.0.0/0"]

  allow {
    ports    = ["22", "6868", "25555"]
    protocol = "tcp"
  }

  target_tags = ["${var.environment_name}-bosh-open"]
}

resource "google_compute_firewall" "bosh-open" {
  name    = "${var.environment_name}-bosh-open"
  network = "${google_compute_network.director.name}"

  source_tags = ["${var.environment_name}-bosh-open"]

  allow {
    ports    = ["22", "6868", "8443", "8844", "25555"]
    protocol = "tcp"
  }

  target_tags = ["${var.environment_name}-director"]
}

resource "google_compute_firewall" "director-to-internal" {
  name    = "${var.environment_name}-bosh-director"
  network = "${google_compute_network.director.name}"

  source_tags = ["${var.environment_name}-director"]

  allow {
    protocol = "tcp"
  }

  target_tags = ["${var.environment_name}-internal"]
}

resource "google_compute_firewall" "internal-to-director" {
  name    = "${var.environment_name}-internal-to-director"
  network = "${google_compute_network.director.name}"

  source_tags = ["${var.environment_name}-internal"]

  allow {
    ports    = ["4222", "25250", "25777"]
    protocol = "tcp"
  }

  target_tags = ["${var.environment_name}-director"]
}

resource "google_compute_firewall" "jumpbox-to-all" {
  name    = "${var.environment_name}-jumpbox-to-all"
  network = "${google_compute_network.director.name}"

  source_tags = ["${var.environment_name}-jumpbox"]

  allow {
    ports    = ["22", "3389"]
    protocol = "tcp"
  }

  target_tags = ["${var.environment_name}-internal", "${var.environment_name}-director"]
}

resource "google_compute_firewall" "internal-to-internal" {
  name    = "${var.environment_name}-internal"
  network = "${google_compute_network.director.name}"

  source_tags = ["${var.environment_name}-internal"]

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
  }

  allow {
    protocol = "udp"
  }

  target_tags = ["${var.environment_name}-internal"]
}

output "jumpbox-ip" {
  value = "${google_compute_address.jumpbox.address}"
}

output "director-network-name" {
value= "${google_compute_network.director.name}"
}

output "director-subnetwork-cidr-range" {
value= "${google_compute_subnetwork.director.ip_cidr_range}"
}

output "director-subnetwork-name" {
value= "${google_compute_subnetwork.director.name}"
}

output "director-tag" {
  value = "${tolist(google_compute_firewall.internal-to-internal.target_tags)[0]}"
}

output "internal-tag" {
  value = "${tolist(google_compute_firewall.internal-to-internal.target_tags)[0]}"
}

output "jumpbox-tag" {
  value = "${tolist(google_compute_firewall.internal-to-internal.target_tags)[0]}"
}

output "bosh-open-tag" {
  value = "${tolist(google_compute_firewall.external.target_tags)[0]}"
}

output "zone1" {
value = "${var.zone1}"
}

output "projectid" {
  value = "${var.projectid}"
}

output "jumpbox-internal-ip" {
  value = "${cidrhost(google_compute_subnetwork.director.ip_cidr_range, 5)}"
}

output "director-internal-ip" {
  value = "${cidrhost(google_compute_subnetwork.director.ip_cidr_range, 6)}"
}

output "internal-gw" {
  value = "${google_compute_subnetwork.director.gateway_address}"
}
