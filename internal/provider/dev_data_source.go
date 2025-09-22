package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &developersDataSource{}
	_ datasource.DataSourceWithConfigure = &developersDataSource{}
)

// NewDevelopersDataSource is a helper function to simplify the provider implementation.
func NewDevelopersDataSource() datasource.DataSource {
	return &developersDataSource{}
}

// developersDataSource is the data source implementation.
type developersDataSource struct {
	client *Client
}

// Metadata returns the data source type name.
func (d *developersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_developers"
}

// Schema defines the schema for the data source.
func (d *developersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of developer teams.",
		Attributes: map[string]schema.Attribute{
			"developers": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of developer teams",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for the developer team",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the developer team",
						},
						"engineers": schema.ListNestedAttribute{
							Computed:    true,
							Description: "List of engineers in this developer team",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:    true,
										Description: "Unique identifier for the engineer",
									},
									"name": schema.StringAttribute{
										Computed:    true,
										Description: "Name of the engineer",
									},
									"email": schema.StringAttribute{
										Computed:    true,
										Description: "Email address of the engineer",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *developersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get developers from the API
	developers, err := d.client.GetDevelopers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read DevOps Developers",
			"An unexpected error occurred when reading the DevOps developers. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"DevOps Client Error: "+err.Error(),
		)
		return
	}

	// Map response body to model
	var state developersDataSourceModel
	for _, developer := range developers {
		// Map engineers for this developer team
		var engineersState []developerEngineerModel
		for _, engineer := range developer.Engineers {
			engineerState := developerEngineerModel{
				ID:    types.StringValue(engineer.ID),
				Name:  types.StringValue(engineer.Name),
				Email: types.StringValue(engineer.Email),
			}
			engineersState = append(engineersState, engineerState)
		}

		developerState := developersModel{
			ID:        types.StringValue(developer.ID),
			Name:      types.StringValue(developer.Name),
			Engineers: engineersState,
		}

		state.Developers = append(state.Developers, developerState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *developersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// developersDataSourceModel maps the data source schema data.
type developersDataSourceModel struct {
	Developers []developersModel `tfsdk:"developers"`
}

// developersModel maps developers schema data.
type developersModel struct {
	ID        types.String               `tfsdk:"id"`
	Name      types.String               `tfsdk:"name"`
	Engineers []developerEngineerModel   `tfsdk:"engineers"`
}

// developerEngineerModel maps engineer data within developer teams.
type developerEngineerModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}