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
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBasicOrganizationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOrganizationResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zammad_organization.test", "name", "one"),
					resource.TestCheckResourceAttr("zammad_organization.test", "active", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "zammad_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccOrganizationResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zammad_organization.test", "name", "two"),
					resource.TestCheckResourceAttr("zammad_organization.test", "active", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAdvancedOrganizationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAdvancedOrganizationResourceConfig("one", "false", "One Priority", "example.com", "true", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zammad_organization.test", "name", "one"),
					resource.TestCheckResourceAttr("zammad_organization.test", "active", "false"),
					resource.TestCheckResourceAttr("zammad_organization.test", "note", "One Priority"),
					resource.TestCheckResourceAttr("zammad_organization.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("zammad_organization.test", "domain_assignment", "true"),
					resource.TestCheckResourceAttr("zammad_organization.test", "shared", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "zammad_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAdvancedOrganizationResourceConfig("one", "true", "Second Priority", "example.example.com", "false", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zammad_organization.test", "name", "one"),
					resource.TestCheckResourceAttr("zammad_organization.test", "active", "true"),
					resource.TestCheckResourceAttr("zammad_organization.test", "note", "Second Priority"),
					resource.TestCheckResourceAttr("zammad_organization.test", "domain", "example.example.com"),
					resource.TestCheckResourceAttr("zammad_organization.test", "domain_assignment", "false"),
					resource.TestCheckResourceAttr("zammad_organization.test", "shared", "true"),
				),
			},
			// Back to original
			{
				Config: testAccAdvancedOrganizationResourceConfig("one", "false", "One Priority", "example.com", "true", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zammad_organization.test", "name", "one"),
					resource.TestCheckResourceAttr("zammad_organization.test", "active", "false"),
					resource.TestCheckResourceAttr("zammad_organization.test", "note", "One Priority"),
					resource.TestCheckResourceAttr("zammad_organization.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("zammad_organization.test", "domain_assignment", "true"),
					resource.TestCheckResourceAttr("zammad_organization.test", "shared", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccOrganizationResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "zammad_organization" "test" {
	name = "%s"
}
`, name)
}

func testAccAdvancedOrganizationResourceConfig(name, active, note, domain, domainAssignment, shared string) string {
	return fmt.Sprintf(`
resource "zammad_organization" "test" {
	name = "%s"
	active = "%s"
	note = "%s"
	domain = "%s"
	domain_assignment = "%s"
	shared = "%s"
}
`, name, active, note, domain, domainAssignment, shared)
}
