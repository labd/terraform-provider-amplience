# Resources
## Content Type resource
```terraform
variable "amplience_content_repository" {
  id = "12345"
}

resource "amplience_content_type" "standard" {
  content_repository_id = local.amplience_content_repository.id
  label = "test_content_repo_label"
  name = "test_content_repo_name"
}
```
