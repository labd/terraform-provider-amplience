package amplience

type RepositoryContentType struct {
	ID             string   `json:"id,omitempty"`
	ContentTypeURI string   `json:"contentTypeUri"`
	Status         string   `json:"status,omitempty"`
	Settings       Settings `json:"settings"`
}

type Settings struct {
	Label          string          `json:"label"`
	Icons          []Icon          `json:"icons"`
	Visualizations []Visualization `json:"visualizations"`
}

type Icon struct {
	Size int    `json:"size"`
	URL  string `json:"url"`
}

type Visualization struct {
	Label        string `json:"label"`
	TemplatedURI string `json:"templatedUri"`
	Default      bool   `json:"default"`
}

// ContentRepositoryStatusEnum is an enum type
type ContentStatusEnum string

// Enum values for ContentRepository status
const (
	ContentStatusActive  ContentStatusEnum = "ACTIVE"
	ContentStatusDeleted ContentStatusEnum = "DELETED"
)
