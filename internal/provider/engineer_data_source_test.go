package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEngineersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "devops-bootcamp_engineers" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify that engineers are returned
					resource.TestCheckResourceAttrSet("data.devops-bootcamp_engineers.test", "engineers.#"),
					// Verify the first engineer to ensure all attributes are set
					resource.TestCheckResourceAttrSet("data.devops-bootcamp_engineers.test", "engineers.0.id"),
					resource.TestCheckResourceAttrSet("data.devops-bootcamp_engineers.test", "engineers.0.name"),
					resource.TestCheckResourceAttrSet("data.devops-bootcamp_engineers.test", "engineers.0.email"),
				),
			},
		},
	})
}
