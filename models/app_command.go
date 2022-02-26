package models

type AppCommand struct {
	Type              uint8              `json:"type"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	DefaultPermission bool               `json:"default_permission"`
	Options           []AppCommandOption `json:"options,omitempty"`
}

type AppCommandOption struct {
	Type        uint8              `json:"type"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Required    bool               `json:"required"`
	Choices     []AppCommandChoice `json:"choices,omitempty"`
	Options     []AppCommandOption `json:"options,omitempty"`
	MinValue    *float64           `json:"min_value,omitempty"`
	MaxValue    *float64           `json:"max_value,omitempty"`
}

type AppCommandChoice struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
