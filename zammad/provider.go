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
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/o11ydev/terraform-provider-zammad/internal/client"
)

func New() tfprovider.Provider {
	return &provider{}
}

type provider struct{}

// GetSchema
func (p *provider) GetSchema(_ context.Context) (schema.Schema, diag.Diagnostics) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Required: true,
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}, nil
}

func (p *provider) Metadata(ctx context.Context, req tfprovider.MetadataRequest, resp *tfprovider.MetadataResponse) {
	resp.TypeName = "zammad"
}

// Provider schema struct
type providerData struct {
	Token types.String `tfsdk:"token"`
	Host  types.String `tfsdk:"host"`
}

func (p *provider) Configure(ctx context.Context, req tfprovider.ConfigureRequest, resp *tfprovider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide a token to the provider
	var token string
	if config.Token.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as token",
		)
		return
	}

	if config.Token.IsNull() {
		token = os.Getenv("ZAMMAD_TOKEN")
	} else {
		token = config.Token.ValueString()
	}

	if token == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find token",
			"Token cannot be an empty string",
		)
		return
	}

	// User must specify a host
	var host string
	if config.Host.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as host",
		)
		return
	}

	if config.Host.IsNull() {
		host = os.Getenv("ZAMMAD_HOST")
	} else {
		host = config.Host.ValueString()
	}

	if host == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find host",
			"Host cannot be an empty string",
		)
		return
	}

	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	transport := logging.NewSubsystemLoggingHTTPTransport("Zammad", t)

	// Create a new Zammad client and set it to the provider client
	c, err := client.New(host, token, transport)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create zammad client:\n\n"+err.Error(),
		)
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

// Resources - Defines provider resources
func (p *provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewZammadTicketPriority,
		NewZammadOrganization,
	}
}

func (p *provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
