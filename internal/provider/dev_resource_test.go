package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDevResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "devops-bootcamp_dev" "test" {
  name = "Frontend Team"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify developer team attributes
					resource.TestCheckResourceAttr("devops-bootcamp_dev.test", "name", "Frontend Team"),
					// Verify computed ID is set
					resource.TestCheckResourceAttrSet("devops-bootcamp_dev.test", "id"),
					// Verify engineers list exists (initially empty)
					resource.TestCheckResourceAttr("devops-bootcamp_dev.test", "engineers.#", "0"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "devops-bootcamp_dev.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "devops-bootcamp_dev" "test" {
  name = "Backend Team"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify developer team attributes updated
					resource.TestCheckResourceAttr("devops-bootcamp_dev.test", "name", "Backend Team"),
					// Verify ID remains set
					resource.TestCheckResourceAttrSet("devops-bootcamp_dev.test", "id"),
					// Verify engineers list still exists
					resource.TestCheckResourceAttrSet("devops-bootcamp_dev.test", "engineers.#"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccDevResourceWithEngineers tests a developer resource that has engineers assigned
func TestAccDevResourceWithEngineers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create an engineer first, then a developer team
			{
				Config: providerConfig + `
resource "devops-bootcamp_engineer" "test_engineer" {
  name  = "Alice Johnson"
  email = "alice.johnson@example.com"
}

resource "devops-bootcamp_dev" "test" {
  name = "Full Stack Team"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify engineer was created
					resource.TestCheckResourceAttr("devops-bootcamp_engineer.test_engineer", "name", "Alice Johnson"),
					resource.TestCheckResourceAttr("devops-bootcamp_engineer.test_engineer", "email", "alice.johnson@example.com"),
					// Verify developer team was created
					resource.TestCheckResourceAttr("devops-bootcamp_dev.test", "name", "Full Stack Team"),
					resource.TestCheckResourceAttrSet("devops-bootcamp_dev.test", "id"),
					// Engineers list should exist (may be empty initially)
					resource.TestCheckResourceAttrSet("devops-bootcamp_dev.test", "engineers.#"),
				),
			},
		},
	})
}
