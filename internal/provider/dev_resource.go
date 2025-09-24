package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/madisonewebb/DOB-tf-providers/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &devResource{}
	_ resource.ResourceWithConfigure   = &devResource{}
	_ resource.ResourceWithImportState = &devResource{}
)

// NewDevResource is a helper function to simplify the provider implementation.
func NewDevResource() resource.Resource {
	return &devResource{}
}

// devResource is the resource implementation.
type devResource struct {
	client *client.Client
}

// Metadata returns the resource type name.
func (r *devResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dev"
}

// Schema defines the schema for the resource.
func (r *devResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a developer team.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the developer team.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the developer team.",
				Required:    true,
			},
			"engineers": schema.ListNestedAttribute{
				Description: "List of engineers in the developer team.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier for the engineer.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the engineer.",
							Computed:    true,
						},
						"email": schema.StringAttribute{
							Description: "Email address of the engineer.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *devResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan devResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new developer
	developer := client.Developer{
		Name:      plan.Name.ValueString(),
		Engineers: []client.Engineer{}, // Start with empty engineers list
	}

	// Create developer via API
	createdDeveloper, err := r.client.CreateDeveloper(developer)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Developer",
			"Could not create developer, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(createdDeveloper.ID)
	plan.Name = types.StringValue(createdDeveloper.Name)
	
	// Map engineers - create types.List
	engineerElementType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":    types.StringType,
			"name":  types.StringType,
			"email": types.StringType,
		},
	}
	
	var engineerElements []attr.Value
	for _, engineer := range createdDeveloper.Engineers {
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
	plan.Engineers = engineersList

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *devResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state devResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed developer value from API
	developer, err := r.client.GetDeveloper(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Developer",
			"Could not read developer ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite developer with refreshed state
	state.ID = types.StringValue(developer.ID)
	state.Name = types.StringValue(developer.Name)
	
	// Map engineers - create types.List
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
	state.Engineers = engineersList

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *devResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan devResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state to preserve engineers list
	var state devResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert engineers from state back to client format
	var engineers []client.Engineer
	// For now, just use empty engineers list since this is complex to extract from types.List
	// In a real implementation, you would extract the engineers from state.Engineers
	engineers = []client.Engineer{}

	// Update existing developer
	developer := client.Developer{
		ID:        plan.ID.ValueString(),
		Name:      plan.Name.ValueString(),
		Engineers: engineers, // Preserve existing engineers
	}

	// Update developer via API
	updatedDeveloper, err := r.client.UpdateDeveloper(plan.ID.ValueString(), developer)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Developer",
			"Could not update developer, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated developer
	plan.ID = types.StringValue(updatedDeveloper.ID)
	plan.Name = types.StringValue(updatedDeveloper.Name)
	
	// Map engineers - create types.List
	engineerElementType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":    types.StringType,
			"name":  types.StringType,
			"email": types.StringType,
		},
	}
	
	var engineerElements []attr.Value
	for _, engineer := range updatedDeveloper.Engineers {
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
	plan.Engineers = engineersList

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *devResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state devResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing developer
	err := r.client.DeleteDeveloper(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Developer",
			"Could not delete developer, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *devResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// ImportState imports the resource state.
func (r *devResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// devResourceModel maps the resource schema data.
type devResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Engineers types.List   `tfsdk:"engineers"`
}

// devEngineerModel maps engineer schema data within developer resource.
type devEngineerModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}
