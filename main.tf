terraform {
  required_providers {
    devops-bootcamp = {
      source = "liatr.io/terraform/devops-bootcamp"
    }
  }
}

provider "devops-bootcamp" {
  endpoint = "http://localhost:8080"
}

resource "devops-bootcamp_engineer" "gob" {
  name  = "Gobs"
  email = "gobs@goblins.com"
}

resource "devops-bootcamp_dev" "goblin" {
  name = "Goblin"
}

resource "devops-bootcamp_dev" "backend_team" {
  name = "Backend Development Team"
}

data "devops-bootcamp_engineers" "all" {}
data "devops-bootcamp_developers" "all" {}

output "engineers" {
  value = data.devops-bootcamp_engineers.all.engineers
}
output "devs" {
  value = data.devops-bootcamp_developers.all.developers
}