package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/madisonewebb/DOB-tf-providers/internal/client"
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
	client *client.Client
}

// Metadata returns the data source type name.
func (d *developersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_developers"
}

// Schema defines the schema for the data source.
func (d *developersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"developers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"engineers": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
									},
									"name": schema.StringAttribute{
										Computed: true,
									},
									"email": schema.StringAttribute{
										Computed: true,
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
		// Map engineers for this developer - create types.List
		engineerElementType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":    types.StringType,
				"name":  types.StringType,
				"email": types.StringType,
			},
		}
		
		var engineerElements []attr.Value
		for _, engineer := range developer.Engineers {
			engineerObj, diags := types.ObjectValue(
				engineerElementType.AttrTypes,
				map[string]attr.Value{
					"id":    types.StringValue(engineer.ID),
					"name":  types.StringValue(engineer.Name),
					"email": types.StringValue(engineer.Email),
				},
			)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			engineerElements = append(engineerElements, engineerObj)
		}
		
		engineersList, diags := types.ListValue(engineerElementType, engineerElements)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		developerState := developerDataModel{
			ID:        types.StringValue(developer.ID),
			Name:      types.StringValue(developer.Name),
			Engineers: engineersList,
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

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// developersDataSourceModel maps the data source schema data.
type developersDataSourceModel struct {
	Developers []developerDataModel `tfsdk:"developers"`
}

// developerDataModel maps developer schema data.
type developerDataModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Engineers types.List   `tfsdk:"engineers"`
}

// devEngineerDataModel maps engineer schema data within developer data source.
type devEngineerDataModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}
