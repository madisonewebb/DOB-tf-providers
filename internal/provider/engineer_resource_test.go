package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEngineerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "devops-bootcamp_engineer" "test" {
  name  = "John Doe"
  email = "john.doe@example.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify engineer attributes
					resource.TestCheckResourceAttr("devops-bootcamp_engineer.test", "name", "John Doe"),
					resource.TestCheckResourceAttr("devops-bootcamp_engineer.test", "email", "john.doe@example.com"),
					// Verify computed ID is set
					resource.TestCheckResourceAttrSet("devops-bootcamp_engineer.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "devops-bootcamp_engineer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "devops-bootcamp_engineer" "test" {
  name  = "Jane Smith"
  email = "jane.smith@example.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify engineer attributes updated
					resource.TestCheckResourceAttr("devops-bootcamp_engineer.test", "name", "Jane Smith"),
					resource.TestCheckResourceAttr("devops-bootcamp_engineer.test", "email", "jane.smith@example.com"),
					// Verify ID remains set
					resource.TestCheckResourceAttrSet("devops-bootcamp_engineer.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
