---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "looker_user_roles Resource - terraform-provider-looker"
subcategory: ""
description: |-
  
---

# looker_user_roles (Resource)



## Example Usage

```terraform
resource "looker_user_roles" "user_roles" {
  user_id  = looker_user.user.id
  role_ids = [looker_role.role.id]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **role_ids** (Set of Number)
- **user_id** (String)

### Optional

- **id** (String) The ID of this resource.

