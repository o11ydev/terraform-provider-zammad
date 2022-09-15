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

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/o11ydev/terraform-provider-zammad/internal/client"
)

type resourceOrganizationType struct{}

// Order Resource schema
func (r resourceOrganizationType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:          types.StringType,
				Computed:      true,
				PlanModifiers: []tfsdk.AttributePlanModifier{tfsdk.UseStateForUnknown()},
			},
			"name": {
				Type:     types.StringType,
				Required: true,
			},
			"shared": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				Description:   "Customers in the organization can see each other's items.",
				PlanModifiers: []tfsdk.AttributePlanModifier{&defaultTrue{}, tfsdk.UseStateForUnknown()},
			},
			"domain": {
				Type:     types.StringType,
				Computed: true,
				Optional: true,
			},
			"domain_assignment": {
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
				Description: "Assign users based on user domain.",
			},
			"active": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: []tfsdk.AttributePlanModifier{&defaultTrue{}, tfsdk.UseStateForUnknown()},
			},
			"note": {
				Type:     types.StringType,
				Optional: true,
			},
			"created_by_id": {
				Type:          types.Int64Type,
				Computed:      true,
				PlanModifiers: []tfsdk.AttributePlanModifier{tfsdk.UseStateForUnknown()},
			},
			"updated_by_id": {
				Type:     types.Int64Type,
				Computed: true,
			},
			"created_at": {
				Type:          types.StringType,
				Computed:      true,
				PlanModifiers: []tfsdk.AttributePlanModifier{tfsdk.UseStateForUnknown()},
			},
			"updated_at": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}

// New resource instance
func (r resourceOrganizationType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceOrganization{
		p: *(p.(*provider)),
	}, nil
}

type resourceOrganization struct {
	p provider
}

// Create a new resource
func (r resourceOrganization) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan Organization
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgreq := &client.Organization{
		Name:             plan.Name.Value,
		Active:           plan.Active.Value,
		Note:             plan.Note.Value,
		Domain:           plan.Domain.Value,
		DomainAssignment: plan.DomainAssignment.Value,
		Shared:           plan.Shared.Value,
	}

	org, err := r.p.client.CreateOrganization(orgreq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating organization",
			"Could not create organization, unexpected error: "+err.Error(),
		)
		return
	}

	result := Organization{
		ID:               types.String{Value: strconv.Itoa(org.ID)},
		Name:             types.String{Value: org.Name},
		Note:             types.String{Value: org.Note},
		Shared:           types.Bool{Value: org.Shared},
		Domain:           types.String{Value: org.Domain},
		DomainAssignment: types.Bool{Value: org.DomainAssignment},
		Active:           types.Bool{Value: org.Active},
		CreatedByID:      types.Int64{Value: int64(org.CreatedByID)},
		UpdatedByID:      types.Int64{Value: int64(org.UpdatedByID)},
		CreatedAt:        types.String{Value: org.CreatedAt},
		UpdatedAt:        types.String{Value: org.UpdatedAt},
	}
	if plan.Note.Null && org.Note == "" {
		result.Note = types.String{Null: true}
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceOrganization) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state Organization
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, err := strconv.Atoi(state.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ID",
			"Could convert id "+state.ID.Value+": "+err.Error(),
		)
		return

	}

	neworg, err := r.p.client.GetOrganization(orgID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading organization",
			"Could not read organization "+state.ID.Value+": "+err.Error(),
		)
		return
	}

	state.Name = types.String{Value: neworg.Name}
	state.Active = types.Bool{Value: neworg.Active}
	state.Shared = types.Bool{Value: neworg.Shared}
	state.Domain = types.String{Value: neworg.Domain}
	state.DomainAssignment = types.Bool{Value: neworg.DomainAssignment}
	state.UpdatedAt = types.String{Value: neworg.UpdatedAt}
	state.UpdatedByID = types.Int64{Value: int64(neworg.UpdatedByID)}
	state.CreatedAt = types.String{Value: neworg.CreatedAt}
	state.CreatedByID = types.Int64{Value: int64(neworg.CreatedByID)}
	if state.Note.Null && neworg.Note == "" {
		state.Note = types.String{Null: true}
	} else {
		state.Note = types.String{Value: neworg.Note}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceOrganization) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
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

	orgID, err := strconv.Atoi(state.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ID",
			"Could convert id "+state.ID.Value+": "+err.Error(),
		)
		return
	}

	updatedOrg := &client.Organization{
		ID:               orgID,
		Name:             plan.Name.Value,
		Note:             plan.Note.Value,
		Domain:           plan.Domain.Value,
		DomainAssignment: plan.DomainAssignment.Value,
		Active:           plan.Active.Value,
		Shared:           plan.Shared.Value,
	}

	org, err := r.p.client.UpdateOrganization(updatedOrg)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error update organization",
			"Could not update organization "+state.ID.Value+": "+err.Error(),
		)
		return
	}

	result := Organization{
		ID:               types.String{Value: strconv.Itoa(org.ID)},
		Name:             types.String{Value: org.Name},
		Note:             types.String{Value: org.Note},
		Domain:           types.String{Value: org.Domain},
		DomainAssignment: types.Bool{Value: org.DomainAssignment},
		Active:           types.Bool{Value: org.Active},
		Shared:           types.Bool{Value: org.Shared},
		CreatedByID:      types.Int64{Value: int64(org.CreatedByID)},
		UpdatedByID:      types.Int64{Value: int64(org.UpdatedByID)},
		CreatedAt:        types.String{Value: org.CreatedAt},
		UpdatedAt:        types.String{Value: org.UpdatedAt},
	}
	if plan.Note.Null && org.Note == "" {
		result.Note = types.String{Null: true}
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceOrganization) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state Organization
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, err := strconv.Atoi(state.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ID",
			"Could convert id "+state.ID.Value+": "+err.Error(),
		)
		return
	}

	err = r.p.client.DeleteOrganization(&client.Organization{ID: orgID})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting organization",
			"Could not delete organization "+state.ID.Value+": "+err.Error(),
		)
		return
	}
}

// Import resource
func (r resourceOrganization) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
