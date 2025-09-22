// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure DevOpsProvider satisfies various provider interfaces.
var _ provider.Provider = &DevOpsProvider{}
var _ provider.ProviderWithFunctions = &DevOpsProvider{}
var _ provider.ProviderWithEphemeralResources = &DevOpsProvider{}

// DevOpsProvider defines the provider implementation.
type DevOpsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// DevOpsProviderModel describes the provider data model.
type DevOpsProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *DevOpsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "devops-bootcamp"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
// Note: DevOps API does not require authentication
func (p *DevOpsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "URI for DevOps API. May also be provided via DEVOPS_ENDPOINT environment variable.",
			},
		},
	}
}

func (p *DevOpsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config DevOpsProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for the endpoint,
	// it must be a known value.

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown DevOps API Endpoint",
			"The provider cannot create the DevOps API client as there is an unknown configuration value for the DevOps API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the DEVOPS_ENDPOINT environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	endpoint := os.Getenv("DEVOPS_ENDPOINT")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing DevOps API Endpoint",
			"The provider cannot create the DevOps API client as there is a missing or empty value for the DevOps API endpoint. "+
				"Set the endpoint value in the configuration or use the DEVOPS_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new DevOps client using the configuration values
	client, err := NewClient(endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create DevOps API Client",
			"An unexpected error occurred when creating the DevOps API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"DevOps Client Error: "+err.Error(),
		)
		return
	}

	// Make the DevOps client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *DevOpsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewEngineerResource,
	}
}

func (p *DevOpsProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *DevOpsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEngineersDataSource,
	}
}

func (p *DevOpsProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DevOpsProvider{
			version: version,
		}
	}
}
