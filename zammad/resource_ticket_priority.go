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

type resourceTicketPriorityType struct{}

// Order Resource schema
func (r resourceTicketPriorityType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"note": {
				Type:     types.StringType,
				Optional: true,
			},
			"ui_icon": {
				Type:     types.StringType,
				Optional: true,
			},
			"ui_color": {
				Type:     types.StringType,
				Optional: true,
			},
			"default_create": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: []tfsdk.AttributePlanModifier{tfsdk.UseStateForUnknown()},
			},
			"active": {
				Type:          types.BoolType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: []tfsdk.AttributePlanModifier{&defaultTrue{}, tfsdk.UseStateForUnknown()},
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
func (r resourceTicketPriorityType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceTicketPriority{
		p: *(p.(*provider)),
	}, nil
}

type resourceTicketPriority struct {
	p provider
}

// Create a new resource
func (r resourceTicketPriority) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan TicketPriority
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tpreq := &client.TicketPriority{
		Name:          plan.Name.Value,
		Note:          plan.Note.Value,
		UIColor:       plan.UIColor.Value,
		UIIcon:        plan.UIIcon.Value,
		Active:        plan.Active.Value,
		DefaultCreate: plan.DefaultCreate.Value,
	}

	tp, err := r.p.client.CreateTicketPriority(tpreq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ticket_priotiry",
			"Could not create ticket_priority, unexpected error: "+err.Error(),
		)
		return
	}

	result := TicketPriority{
		ID:            types.String{Value: strconv.Itoa(tp.ID)},
		Name:          types.String{Value: tp.Name},
		Note:          types.String{Value: tp.Note},
		UIColor:       types.String{Value: tp.UIColor},
		UIIcon:        types.String{Value: tp.UIIcon},
		Active:        types.Bool{Value: tp.Active},
		DefaultCreate: types.Bool{Value: tp.DefaultCreate},
		CreatedByID:   types.Int64{Value: int64(tp.CreatedByID)},
		UpdatedByID:   types.Int64{Value: int64(tp.UpdatedByID)},
		CreatedAt:     types.String{Value: tp.CreatedAt},
		UpdatedAt:     types.String{Value: tp.UpdatedAt},
	}
	if plan.UIColor.Null && tp.UIColor == "" {
		result.UIColor = types.String{Null: true}
	}
	if plan.UIIcon.Null && tp.UIIcon == "" {
		result.UIIcon = types.String{Null: true}
	}
	if plan.Note.Null && tp.Note == "" {
		result.Note = types.String{Null: true}
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceTicketPriority) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state TicketPriority
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tpID, err := strconv.Atoi(state.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ID",
			"Could convert id "+state.ID.Value+": "+err.Error(),
		)
		return

	}

	newtp, err := r.p.client.GetTicketPriority(tpID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ticket_priority",
			"Could not read ticket_priority "+state.ID.Value+": "+err.Error(),
		)
		return
	}

	state.Name = types.String{Value: newtp.Name}
	state.CreatedAt = types.String{Value: newtp.CreatedAt}
	state.Active = types.Bool{Value: newtp.Active}
	state.DefaultCreate = types.Bool{Value: newtp.DefaultCreate}
	state.UpdatedAt = types.String{Value: newtp.UpdatedAt}
	state.UpdatedByID = types.Int64{Value: int64(newtp.UpdatedByID)}
	state.CreatedByID = types.Int64{Value: int64(newtp.CreatedByID)}
	if state.UIColor.Null && newtp.UIColor == "" {
		state.UIColor = types.String{Null: true}
	} else {
		state.UIColor = types.String{Value: newtp.UIColor}
	}
	if state.UIIcon.Null && newtp.UIIcon == "" {
		state.UIIcon = types.String{Null: true}
	} else {
		state.UIIcon = types.String{Value: newtp.UIIcon}
	}
	if state.Note.Null && newtp.Note == "" {
		state.Note = types.String{Null: true}
	} else {
		state.Note = types.String{Value: newtp.Note}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceTicketPriority) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan TicketPriority
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state TicketPriority
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tpID, err := strconv.Atoi(state.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ID",
			"Could convert id "+state.ID.Value+": "+err.Error(),
		)
		return
	}

	updatedTP := &client.TicketPriority{
		ID:            tpID,
		Name:          plan.Name.Value,
		Note:          plan.Note.Value,
		UIColor:       plan.UIColor.Value,
		UIIcon:        plan.UIIcon.Value,
		Active:        plan.Active.Value,
		DefaultCreate: plan.DefaultCreate.Value,
	}

	tp, err := r.p.client.UpdateTicketPriority(updatedTP)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error update order",
			"Could not update orderID "+state.ID.Value+": "+err.Error(),
		)
		return
	}

	result := TicketPriority{
		ID:            types.String{Value: strconv.Itoa(tp.ID)},
		Name:          types.String{Value: tp.Name},
		Note:          types.String{Value: tp.Note},
		UIColor:       types.String{Value: tp.UIColor},
		UIIcon:        types.String{Value: tp.UIIcon},
		Active:        types.Bool{Value: tp.Active},
		DefaultCreate: types.Bool{Value: tp.DefaultCreate},
		CreatedByID:   types.Int64{Value: int64(tp.CreatedByID)},
		UpdatedByID:   types.Int64{Value: int64(tp.UpdatedByID)},
		CreatedAt:     types.String{Value: tp.CreatedAt},
		UpdatedAt:     types.String{Value: tp.UpdatedAt},
	}
	if plan.UIColor.Null && tp.UIColor == "" {
		result.UIColor = types.String{Null: true}
	}
	if plan.UIIcon.Null && tp.UIIcon == "" {
		result.UIIcon = types.String{Null: true}
	}
	if plan.Note.Null && tp.Note == "" {
		result.Note = types.String{Null: true}
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r resourceTicketPriority) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state TicketPriority
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tpID, err := strconv.Atoi(state.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ID",
			"Could convert id "+state.ID.Value+": "+err.Error(),
		)
		return
	}

	err = r.p.client.DeleteTicketPriority(&client.TicketPriority{ID: tpID})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting ticket_priority",
			"Could not delete ticket_priority "+state.ID.Value+": "+err.Error(),
		)
		return
	}
}

// Import resource
func (r resourceTicketPriority) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

type defaultTrue struct{}

func (m defaultTrue) Description(ctx context.Context) string {
	return "If value is not configured, defaults to true"
}

func (m defaultTrue) MarkdownDescription(ctx context.Context) string {
	return "If value is not configured, defaults to `true`"
}

func (m defaultTrue) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	var str types.Bool
	diags := tfsdk.ValueAs(ctx, req.AttributePlan, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if !str.IsUnknown() {
		return
	}

	resp.AttributePlan = types.Bool{Value: true}
}
