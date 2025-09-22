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

# Test data source - read all engineers
data "devops-bootcamp_engineers" "all" {}

# Test resource - create a new engineer
resource "devops-bootcamp_engineer" "terraform_test" {
  name  = "Terraform Test Engineer"
  email = "terraform@example.com"
}

# Outputs
output "all_engineers" {
  description = "All engineers from the API"
  value       = data.devops-bootcamp_engineers.all.engineers
}

output "new_engineer_id" {
  description = "ID of the newly created engineer"
  value       = devops-bootcamp_engineer.terraform_test.id
}
