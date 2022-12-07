---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zammad_organization Resource - terraform-provider-zammad"
subcategory: ""
description: |-
  
---

# zammad_organization (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)

### Optional

- `active` (Boolean)
- `domain` (String)
- `domain_assignment` (Boolean) Assign users based on user domain.
- `note` (String)
- `shared` (Boolean) Customers in the organization can see each other's items.

### Read-Only

- `created_at` (String)
- `created_by_id` (Number)
- `id` (String) The ID of this resource.
- `updated_at` (String)
- `updated_by_id` (Number)

