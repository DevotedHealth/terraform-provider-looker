---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "looker_lookml_model Resource - terraform-provider-looker"
subcategory: ""
description: |-
  
---

# looker_lookml_model (Resource)



## Example Usage

```terraform
resource "looker_lookml_model" "lookml_model" {
  name                        = "LookML Model"
  allowed_db_connection_names = ["bigquery-connection"]
  project_name                = "lookml_model"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)
- `project_name` (String)

### Optional

- `allowed_db_connection_names` (Set of String)
- `id` (String) The ID of this resource.


