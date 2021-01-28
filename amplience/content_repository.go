package amplience

type ContentRepository struct {
	ID           string        `json:"id,omitempty"`
	Name         string        `json:"name"`
	Label        string        `json:"label"`
	Status       string        `json:"status,omitempty"`
	Features     []string      `json:"features,omitempty"`
	Type         string        `json:"type,omitempty"`
	ContentTypes []ContentType `json:"contentTypes,omitempty"`
	ItemLocales  []string      `json:"itemLocales,omitempty"`
}

type ContentType struct {
	HubContentTypeID string `json:"hubContentTypeId"`
	ContentTypeURI   string `json:"contentTypeUri"`
}

// ContentRepositoryStatusEnum is an enum type
type ContentRepositoryStatusEnum string

// Enum values for ContentRepository status
const (
	ContentRepositoryStatusActive  ContentRepositoryStatusEnum = "ACTIVE"
	ContentRepositoryStatusDeleted ContentRepositoryStatusEnum = "DELETED"
)

// Below should go somewhere more generic when/if more types are added
type ResourceType string

const (
	TypeContent ResourceType = "CONTENT"
)
