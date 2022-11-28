# v0.3.14 (2022-10-21)

- Add check in delete for archiving

# v0.3.13 (2022-10-19)

- Check for version before deleting content type
- Improve error handling

# v0.3.12 (2022-10-17)

- Update content type (schema) queries to allow for querying archived schemas
- Add resilience to creating and updating content types and content type schemas

# v0.3.8 (2022-07-15)

- upgraded to go 1.17 to trigger arm64 builds

# v0.3.7 (2022-03-14)

- Add icon and visualization when importing/reading content type
- add unarchive option when provider receives 409 on create (#15)
- Update documentation

# v0.3.6 (2021-12-10)

- update amplience sdk
- recreate index on suffix change

# v0.3.5 (2021-09-28)

- update provider and handle create corruption

# v0.3.4 (2021-09-24)

- **New resource:** `amplience_search_index`
- Update amplience sdk version
- Fix archiving content types

# 0.3.3 (2021-05-21)

- Fix resource_content_type deletion (archiving)

# 0.3.2 (2021-05-18)

- Fix the resource_content_type_assignemt resource

# 0.2.0 (2021-04-20)

- Refactored the codebase to use a separate Go SDK for Amplience
- Add the following new resources:
  - Content repository
  - Content type schema
  - Content type assignment
  - Content type
- Add the following new data sources

  - Hub
  - Content repository

  # 0.1.1 (2021-04-01)

- Improved Error messages for more user friendly error responses and propagation of Amplience API error response
- Resource Webhook: Fixed a bug where creating a Webhook without a `custom_payload` would return an error

  # 0.1.0 (2021-03-17)

- **New resource:** `amplience_content_repository`
- **New resource:** `amplience_content_webhook`
