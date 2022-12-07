// Copyright 2022 The Terraform Provider for Zammad Authors
// spdx-license-identifier: apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zammad

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/o11ydev/terraform-provider-zammad/internal/client"
)

// Order Resource schema
func (r resourceOrganization) GetSchema(_ context.Context) (schema.Schema, diag.Diagnostics) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"shared": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Customers in the organization can see each other's items.",
				PlanModifiers: []planmodifier.Bool{&defaultTrue{}, boolplanmodifier.UseStateForUnknown()},
			},
			"member_ids": schema.ListAttribute{
				ElementType: types.Int64Type,
				Computed:    true,
				Optional:    true,
			},
			"domain": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"domain_assignment": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Assign users based on user domain.",
			},
			"active": schema.BoolAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.Bool{&defaultTrue{}, boolplanmodifier.UseStateForUnknown()},
			},
			"note": schema.StringAttribute{
				Optional: true,
			},
			"created_by_id": schema.Int64Attribute{
				Computed:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"updated_by_id": schema.Int64Attribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}, nil
}

func NewZammadOrganization() resource.Resource {
	return &resourceOrganization{}
}

func (r *resourceOrganization) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client, _ = req.ProviderData.(*client.Client)
}

func (r *resourceOrganization) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

type resourceOrganization struct {
	client *client.Client
}

// Create a new resource
func (r resourceOrganization) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Organization
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgreq := &client.Organization{
		Name:             plan.Name.ValueString(),
		Active:           plan.Active.ValueBool(),
		Note:             plan.Note.ValueString(),
		Domain:           plan.Domain.ValueString(),
		DomainAssignment: plan.DomainAssignment.ValueBool(),
		Shared:           plan.Shared.ValueBool(),
	}

	org, err := r.client.CreateOrganization(orgreq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating organization",
			"Could not create organization, unexpected error: "+err.Error(),
		)
		return
	}

	members := make([]attr.Value, len(org.MemberIDs))
	for i := range org.MemberIDs {
		members[i] = types.Int64Value(int64(org.MemberIDs[i]))
	}

	result := Organization{
		ID:               types.StringValue(strconv.Itoa(org.ID)),
		Name:             types.StringValue(org.Name),
		Note:             types.StringValue(org.Note),
		Shared:           types.BoolValue(org.Shared),
		Domain:           types.StringValue(org.Domain),
		DomainAssignment: types.BoolValue(org.DomainAssignment),
		Active:           types.BoolValue(org.Active),
		CreatedByID:      types.Int64Value(int64(org.CreatedByID)),
		UpdatedByID:      types.Int64Value(int64(org.UpdatedByID)),
		CreatedAt:        types.StringValue(org.CreatedAt),
		UpdatedAt:        types.StringValue(org.UpdatedAt),
		MemberIDs:        types.ListValueMust(types.Int64Type, members),
	}
	if plan.Note.IsNull() && org.Note == "" {
		result.Note = types.StringNull()
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceOrganization) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Organization
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ID",
			"Could convert id "+state.ID.ValueString()+": "+err.Error(),
		)
		return

	}

	neworg, err := r.client.GetOrganization(orgID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading organization",
			"Could not read organization "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(neworg.Name)
	state.Active = types.BoolValue(neworg.Active)
	state.Shared = types.BoolValue(neworg.Shared)
	state.Domain = types.StringValue(neworg.Domain)
	state.DomainAssignment = types.BoolValue(neworg.DomainAssignment)
	state.UpdatedAt = types.StringValue(neworg.UpdatedAt)
	state.UpdatedByID = types.Int64Value(int64(neworg.UpdatedByID))
	state.CreatedAt = types.StringValue(neworg.CreatedAt)
	state.CreatedByID = types.Int64Value(int64(neworg.CreatedByID))
	if state.Note.IsNull() && neworg.Note == "" {
		state.Note = types.StringNull()
	} else {
		state.Note = types.StringValue(neworg.Note)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceOrganization) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan Organization
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state Organization
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ID",
			"Could convert id "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	updatedOrg := &client.Organization{
		ID:               orgID,
		Name:             plan.Name.ValueString(),
		Note:             plan.Note.ValueString(),
		Domain:           plan.Domain.ValueString(),
		DomainAssignment: plan.DomainAssignment.ValueBool(),
		Active:           plan.Active.ValueBool(),
		Shared:           plan.Shared.ValueBool(),
	}

	org, err := r.client.UpdateOrganization(updatedOrg)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error update organization",
			"Could not update organization "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	result := Organization{
		ID:               types.StringValue(strconv.Itoa(org.ID)),
		Name:             types.StringValue(org.Name),
		Note:             types.StringValue(org.Note),
		Domain:           types.StringValue(org.Domain),
		DomainAssignment: types.BoolValue(org.DomainAssignment),
		Active:           types.BoolValue(org.Active),
		Shared:           types.BoolValue(org.Shared),
		CreatedByID:      types.Int64Value(int64(org.CreatedByID)),
		UpdatedByID:      types.Int64Value(int64(org.UpdatedByID)),
		CreatedAt:        types.StringValue(org.CreatedAt),
		UpdatedAt:        types.StringValue(org.UpdatedAt),
	}
	if plan.Note.IsNull() && org.Note == "" {
		result.Note = types.StringNull()
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceOrganization) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Organization
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ID",
			"Could convert id "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	err = r.client.DeleteOrganization(&client.Organization{ID: orgID})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting organization",
			"Could not delete organization "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}

// Import resource
func (r resourceOrganization) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Save the import identifier in the id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
