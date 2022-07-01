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

func TestAccBasicTicketPriorityResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTicketPriorityResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "name", "one"),
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "active", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "zammad_ticket_priority.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccTicketPriorityResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "name", "two"),
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "active", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAdvancedTicketPriorityResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAdvancedTicketPriorityResourceConfig("one", "false", "One Priority", "red", "fa-ok"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "name", "one"),
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "active", "false"),
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "note", "One Priority"),
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "ui_color", "red"),
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "ui_icon", "fa-ok"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "zammad_ticket_priority.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAdvancedTicketPriorityResourceConfig("one", "false", "Updated prio", "green", "fa-nok"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "name", "one"),
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "active", "false"),
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "note", "Updated prio"),
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "ui_color", "green"),
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "ui_icon", "fa-nok"),
				),
			},
			// Back to original
			{
				Config: testAccAdvancedTicketPriorityResourceConfig("one", "false", "One Priority", "red", "fa-ok"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "name", "one"),
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "active", "false"),
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "note", "One Priority"),
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "ui_color", "red"),
					resource.TestCheckResourceAttr("zammad_ticket_priority.test", "ui_icon", "fa-ok"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTicketPriorityResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "zammad_ticket_priority" "test" {
	name = "%s"
}
`, name)
}

func testAccAdvancedTicketPriorityResourceConfig(name, active, note, uiColor, uiIcon string) string {
	return fmt.Sprintf(`
resource "zammad_ticket_priority" "test" {
	name = "%s"
	active = "%s"
	note = "%s"
	ui_color = "%s"
	ui_icon = "%s"
}
`, name, active, note, uiColor, uiIcon)
}
