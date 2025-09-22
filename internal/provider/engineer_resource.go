package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &engineerResource{}
	_ resource.ResourceWithConfigure   = &engineerResource{}
	_ resource.ResourceWithImportState = &engineerResource{}
)

// NewEngineerResource is a helper function to simplify the provider implementation.
func NewEngineerResource() resource.Resource {
	return &engineerResource{}
}

// engineerResource is the resource implementation.
type engineerResource struct {
	client *Client
}

// Metadata returns the resource type name.
func (r *engineerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer"
}

// Schema defines the schema for the resource.
func (r *engineerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an engineer.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the engineer.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the engineer.",
				Required:    true,
			},
			"email": schema.StringAttribute{
				Description: "Email address of the engineer.",
				Required:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *engineerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan engineerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new engineer
	engineer := Engineer{
		Name:  plan.Name.ValueString(),
		Email: plan.Email.ValueString(),
	}

	// Create engineer via API
	createdEngineer, err := r.client.CreateEngineer(engineer)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Engineer",
			"Could not create engineer, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(createdEngineer.ID)
	plan.Name = types.StringValue(createdEngineer.Name)
	plan.Email = types.StringValue(createdEngineer.Email)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *engineerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state engineerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed engineer value from API
	engineer, err := r.client.GetEngineer(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Engineer",
			"Could not read engineer ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite engineer with refreshed state
	state.ID = types.StringValue(engineer.ID)
	state.Name = types.StringValue(engineer.Name)
	state.Email = types.StringValue(engineer.Email)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *engineerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan engineerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update existing engineer
	engineer := Engineer{
		ID:    plan.ID.ValueString(),
		Name:  plan.Name.ValueString(),
		Email: plan.Email.ValueString(),
	}

	// Update engineer via API
	updatedEngineer, err := r.client.UpdateEngineer(plan.ID.ValueString(), engineer)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Engineer",
			"Could not update engineer, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated engineer
	plan.ID = types.StringValue(updatedEngineer.ID)
	plan.Name = types.StringValue(updatedEngineer.Name)
	plan.Email = types.StringValue(updatedEngineer.Email)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *engineerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state engineerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing engineer
	err := r.client.DeleteEngineer(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Engineer",
			"Could not delete engineer, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *engineerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// ImportState imports the resource state.
func (r *engineerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// engineerResourceModel maps the resource schema data.
type engineerResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}
