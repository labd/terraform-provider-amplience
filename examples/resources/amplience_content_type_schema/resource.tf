resource "amplience_content_type_schema" "tag" {
  body             = <<EOL
{
  "$id": "https://tf-amplience-provider.com/tag",
  "$schema": "http://json-schema.org/draft-07/schema#",
  "allOf": [
    {
      "$ref": "http://bigcontent.io/cms/schema/v1/core#/definitions/content"
    }
  ],
  "title": "Tag",
  "description": "Tag",
  "type": "object",
  "properties": {
    "title": {
      "title": "Title",
      "allOf": [
        {
          "$ref": "http://bigcontent.io/cms/schema/v1/localization#/definitions/localized-string"
        }
      ]
    },
    "color": {
      "title": "Color",
      "type": "string",
      "format": "color"
    }
  },
  "propertyOrder": [
    "title",
    "color"
  ],
  "required": [
    "title",
    "color"
  ]
}
  EOL
  schema_id        = "https://tf-amplience-provider.com/tag"
  validation_level = "CONTENT_TYPE"
}
