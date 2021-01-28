# Content Repositories

Content repositories are separate spaces within a Hub in which content and slots are stored.
Through the provider a repository can only be

- Created (name & label fields) 
- Updated (name & label fields)
- Read

At present the Amplience Dynamic Content API does not allow deletion of repositories. For adding content to the repository
with Terraform see [Content Item](resource_content_item.md)

Depending on your needs and requirements you may opt to have multiple repositories in a single Hub or use multiple Hubs
(and therefore multiple provider instances).

For more information see the Amplience [user documentation](https://amplience.com/docs/intro/hubsandrepositories.html) 
and/or the [API documentation](https://amplience.com/docs/api/dynamic-content/management/index.html#tag/Content-Repositories)

## Example Usage
```hcl
resource "amplience_content_repository" "a-generic-content-repository" {
  name = "a highly generic name"
  label = "a good identifying label"
}
```

