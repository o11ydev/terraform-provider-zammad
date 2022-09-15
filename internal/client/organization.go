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

package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

type Organization struct {
	ID               int    `json:"id,omitempty"`
	Name             string `json:"name"`
	Note             string `json:"note"`
	Shared           bool   `json:"shared"`
	Domain           string `json:"domain"`
	DomainAssignment bool   `json:"domain_assignment"`
	Active           bool   `json:"active"`
	MemberIDs        []int  `json:"member_ids"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
	CreatedByID      int    `json:"created_by_id,omitempty"`
	UpdatedByID      int    `json:"updated_by_id,omitempty"`
}

func (c *Client) CreateOrganization(org *Organization) (*Organization, error) {
	rb, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.host+"/api/v1/organizations", bytes.NewReader(rb))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	neworg := &Organization{}
	err = json.Unmarshal(body, neworg)
	if err != nil {
		return nil, err
	}
	return neworg, nil
}

func (c *Client) GetOrganization(id int) (*Organization, error) {
	req, err := http.NewRequest("GET", c.host+"/api/v1/organizations/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	neworg := &Organization{}
	err = json.Unmarshal(body, neworg)
	if err != nil {
		return nil, err
	}
	return neworg, nil
}

func (c *Client) UpdateOrganization(org *Organization) (*Organization, error) {
	rb, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", c.host+"/api/v1/organizations/"+strconv.Itoa(org.ID), bytes.NewReader(rb))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	neworg := &Organization{}
	err = json.Unmarshal(body, neworg)
	if err != nil {
		return nil, err
	}
	return neworg, nil
}

func (c *Client) DeleteOrganization(org *Organization) error {
	req, err := http.NewRequest("DELETE", c.host+"/api/v1/organizations/"+strconv.Itoa(org.ID), nil)
	if err != nil {
		return err
	}
	_, err = c.doRequest(req)
	return err
}
