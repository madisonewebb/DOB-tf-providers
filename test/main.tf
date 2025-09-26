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

# Test the developers data source
data "devops-bootcamp_developers" "all" {}

# Test the developer resource
resource "devops-bootcamp_dev" "test_team" {
  name = "Updated Development Team"
}

# Outputs
output "all_developers" {
  description = "All developers from the API"
  value       = data.devops-bootcamp_developers.all.developers
}
output "all_engineers" {
  description = "All engineers from the API"
  value       = data.devops-bootcamp_engineers.all.engineers
}
