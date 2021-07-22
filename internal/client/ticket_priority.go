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

type TicketPriority struct {
	ID            int    `json:"id,omitempty"`
	Name          string `json:"name"`
	Note          string `json:"note"`
	UIColor       string `json:"ui_color"`
	UIIcon        string `json:"ui_icon"`
	Active        bool   `json:"active"`
	DefaultCreate bool   `json:"default_create"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
	CreatedByID   int    `json:"created_by_id,omitempty"`
	UpdatedByID   int    `json:"updated_by_id,omitempty"`
}

func (c *Client) CreateTicketPriority(tp *TicketPriority) (*TicketPriority, error) {
	rb, err := json.Marshal(tp)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.host+"/api/v1/ticket_priorities", bytes.NewReader(rb))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	newtp := &TicketPriority{}
	err = json.Unmarshal(body, newtp)
	if err != nil {
		return nil, err
	}
	return newtp, nil
}

func (c *Client) GetTicketPriority(id int) (*TicketPriority, error) {
	req, err := http.NewRequest("GET", c.host+"/api/v1/ticket_priorities/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	newtp := &TicketPriority{}
	err = json.Unmarshal(body, newtp)
	if err != nil {
		return nil, err
	}
	return newtp, nil
}

func (c *Client) UpdateTicketPriority(tp *TicketPriority) (*TicketPriority, error) {
	rb, err := json.Marshal(tp)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", c.host+"/api/v1/ticket_priorities/"+strconv.Itoa(tp.ID), bytes.NewReader(rb))
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	newtp := &TicketPriority{}
	err = json.Unmarshal(body, newtp)
	if err != nil {
		return nil, err
	}
	return newtp, nil
}

func (c *Client) DeleteTicketPriority(tp *TicketPriority) error {
	req, err := http.NewRequest("DELETE", c.host+"/api/v1/ticket_priorities/"+strconv.Itoa(tp.ID), nil)
	if err != nil {
		return err
	}
	_, err = c.doRequest(req)
	return err
}
