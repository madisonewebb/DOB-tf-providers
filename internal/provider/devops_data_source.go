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
	_ datasource.DataSource              = &devopsDataSource{}
	_ datasource.DataSourceWithConfigure = &devopsDataSource{}
)

// NewDevOpsDataSource is a helper function to simplify the provider implementation.
func NewDevOpsDataSource() datasource.DataSource {
	return &devopsDataSource{}
}

// devopsDataSource is the data source implementation.
type devopsDataSource struct {
	client *Client
}

// Metadata returns the data source type name.
func (d *devopsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devops"
}

// Schema defines the schema for the data source.
func (d *devopsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of DevOps teams (combination of developer and operations teams).",
		Attributes: map[string]schema.Attribute{
			"devops": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of DevOps teams",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for the DevOps team",
						},
						"dev": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Developer team within this DevOps team",
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
									Description: "List of engineers in the developer team",
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
						"ops": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Operations team within this DevOps team",
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:    true,
									Description: "Unique identifier for the operations team",
								},
								"name": schema.StringAttribute{
									Computed:    true,
									Description: "Name of the operations team",
								},
								"engineers": schema.ListNestedAttribute{
									Computed:    true,
									Description: "List of engineers in the operations team",
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
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *devopsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get devops from the API
	devopsTeams, err := d.client.GetDevOps()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read DevOps Teams",
			"An unexpected error occurred when reading the DevOps teams. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"DevOps Client Error: "+err.Error(),
		)
		return
	}

	// Map response body to model
	var state devopsDataSourceModel
	for _, devopsTeam := range devopsTeams {
		// Map developer engineers
		var devEngineersState []devopsEngineerModel
		for _, engineer := range devopsTeam.Dev.Engineers {
			engineerState := devopsEngineerModel{
				ID:    types.StringValue(engineer.ID),
				Name:  types.StringValue(engineer.Name),
				Email: types.StringValue(engineer.Email),
			}
			devEngineersState = append(devEngineersState, engineerState)
		}

		// Map operations engineers
		var opsEngineersState []devopsEngineerModel
		for _, engineer := range devopsTeam.Ops.Engineers {
			engineerState := devopsEngineerModel{
				ID:    types.StringValue(engineer.ID),
				Name:  types.StringValue(engineer.Name),
				Email: types.StringValue(engineer.Email),
			}
			opsEngineersState = append(opsEngineersState, engineerState)
		}

		devopsState := devopsModel{
			ID: types.StringValue(devopsTeam.ID),
			Dev: devopsDevModel{
				ID:        types.StringValue(devopsTeam.Dev.ID),
				Name:      types.StringValue(devopsTeam.Dev.Name),
				Engineers: devEngineersState,
			},
			Ops: devopsOpsModel{
				ID:        types.StringValue(devopsTeam.Ops.ID),
				Name:      types.StringValue(devopsTeam.Ops.Name),
				Engineers: opsEngineersState,
			},
		}

		state.DevOps = append(state.DevOps, devopsState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *devopsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// devopsDataSourceModel maps the data source schema data.
type devopsDataSourceModel struct {
	DevOps []devopsModel `tfsdk:"devops"`
}

// devopsModel maps devops schema data.
type devopsModel struct {
	ID  types.String    `tfsdk:"id"`
	Dev devopsDevModel  `tfsdk:"dev"`
	Ops devopsOpsModel  `tfsdk:"ops"`
}

// devopsDevModel maps developer team data within DevOps teams.
type devopsDevModel struct {
	ID        types.String            `tfsdk:"id"`
	Name      types.String            `tfsdk:"name"`
	Engineers []devopsEngineerModel   `tfsdk:"engineers"`
}

// devopsOpsModel maps operations team data within DevOps teams.
type devopsOpsModel struct {
	ID        types.String            `tfsdk:"id"`
	Name      types.String            `tfsdk:"name"`
	Engineers []devopsEngineerModel   `tfsdk:"engineers"`
}

// devopsEngineerModel maps engineer data within DevOps teams.
type devopsEngineerModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}
