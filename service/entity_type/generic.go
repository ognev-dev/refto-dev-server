package entity

// Generic type is data without specific type
type Generic struct {
	// Title (name) of entity
	Title string `json:"title"`

	// Description of entity
	Description string `json:"description"`

	// URL to the preview picture
	PreviewURL string `json:"preview_url"`

	// List of helpful links associated with data
	Links []Link `json:"links"`

	Topics []Topic `json:"topics"`
}

type Link struct {
	Icon  string `json:"icon"`
	Label string `json:"label"`
	Addr  string `json:"addr"`
}

type Topic string
