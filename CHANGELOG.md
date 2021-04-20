0.2.0 (2021-04-20)
==================
- Refactored the codebase to use a separate Go SDK for Amplience
- Add the following new resources:
  - Content repository
  - Content type schema
  - Content type assignment
  - Content type
- Add the following new data sources
  - Hub
  - Content repository

0.1.1 (2021-04-01)
==================
- Improved Error messages for more user friendly error responses and propagation of Amplience API error response
- Resource Webhook: Fixed a bug where creating a Webhook without a `custom_payload` would return an error

0.1.0 (2021-03-17)
==================
- **New resource:** `amplience_content_repository`
- **New resource:** `amplience_content_webhook`
