---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "materialize_database Data Source - terraform-provider-materialize"
subcategory: ""
description: |-
  
---

# materialize_database (Data Source)



## Example Usage

```terraform
data "materialize_database" "all" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `databases` (List of Object) The databases in the account (see [below for nested schema](#nestedatt--databases))
- `id` (String) The ID of this resource.

<a id="nestedatt--databases"></a>
### Nested Schema for `databases`

Read-Only:

- `id` (String)
- `name` (String)
