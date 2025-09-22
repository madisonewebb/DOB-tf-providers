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
	_ datasource.DataSource              = &operationsDataSource{}
	_ datasource.DataSourceWithConfigure = &operationsDataSource{}
)

// NewOperationsDataSource is a helper function to simplify the provider implementation.
func NewOperationsDataSource() datasource.DataSource {
	return &operationsDataSource{}
}

// operationsDataSource is the data source implementation.
type operationsDataSource struct {
	client *Client
}

// Metadata returns the data source type name.
func (d *operationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_operations"
}

// Schema defines the schema for the data source.
func (d *operationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of operations teams.",
		Attributes: map[string]schema.Attribute{
			"operations": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of operations teams",
				NestedObject: schema.NestedAttributeObject{
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
							Description: "List of engineers in this operations team",
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
func (d *operationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get operations from the API
	operations, err := d.client.GetOperations()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read DevOps Operations",
			"An unexpected error occurred when reading the DevOps operations. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"DevOps Client Error: "+err.Error(),
		)
		return
	}

	// Map response body to model
	var state operationsDataSourceModel
	for _, operation := range operations {
		// Map engineers for this operations team
		var engineersState []operationEngineerModel
		for _, engineer := range operation.Engineers {
			engineerState := operationEngineerModel{
				ID:    types.StringValue(engineer.ID),
				Name:  types.StringValue(engineer.Name),
				Email: types.StringValue(engineer.Email),
			}
			engineersState = append(engineersState, engineerState)
		}

		operationState := operationsModel{
			ID:        types.StringValue(operation.ID),
			Name:      types.StringValue(operation.Name),
			Engineers: engineersState,
		}

		state.Operations = append(state.Operations, operationState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *operationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// operationsDataSourceModel maps the data source schema data.
type operationsDataSourceModel struct {
	Operations []operationsModel `tfsdk:"operations"`
}

// operationsModel maps operations schema data.
type operationsModel struct {
	ID        types.String              `tfsdk:"id"`
	Name      types.String              `tfsdk:"name"`
	Engineers []operationEngineerModel  `tfsdk:"engineers"`
}

// operationEngineerModel maps engineer data within operations teams.
type operationEngineerModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}
