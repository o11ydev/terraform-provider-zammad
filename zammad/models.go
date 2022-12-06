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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TicketPriority is a zammad ticket priority.
type TicketPriority struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Note          types.String `tfsdk:"note"`
	UIIcon        types.String `tfsdk:"ui_icon"`
	UIColor       types.String `tfsdk:"ui_color"`
	Active        types.Bool   `tfsdk:"active"`
	DefaultCreate types.Bool   `tfsdk:"default_create"`
	CreatedByID   types.Int64  `tfsdk:"created_by_id"`
	UpdatedByID   types.Int64  `tfsdk:"updated_by_id"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

// Organization is a zammad organization.
type Organization struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Note             types.String `tfsdk:"note"`
	Shared           types.Bool   `tfsdk:"shared"`
	Domain           types.String `tfsdk:"domain"`
	DomainAssignment types.Bool   `tfsdk:"domain_assignment"`
	MemberIDs        types.List   `tfsdk:"member_ids"`
	Active           types.Bool   `tfsdk:"active"`
	CreatedByID      types.Int64  `tfsdk:"created_by_id"`
	UpdatedByID      types.Int64  `tfsdk:"updated_by_id"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
}
