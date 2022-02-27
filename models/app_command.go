package models

type AppCommand struct {
	Type              uint8              `json:"type" mapstructure:"type"`
	Name              string             `json:"name" mapstructure:"name"`
	Description       string             `json:"description" mapstructure:"description"`
	DefaultPermission bool               `json:"default_permission" mapstructure:"default_permission"`
	Options           []AppCommandOption `json:"options,omitempty" mapstructure:"options"`
}

type AppCommandOption struct {
	Type        uint8              `json:"type" mapstructure:"type"`
	Name        string             `json:"name" mapstructure:"name"`
	Description string             `json:"description" mapstructure:"description"`
	Required    bool               `json:"required" mapstructure:"required"`
	Choices     []AppCommandChoice `json:"choices,omitempty" mapstructure:"choices"`
	Options     []AppCommandOption `json:"options,omitempty" mapstructure:"options"`
	MinValue    *float64           `json:"min_value,omitempty" mapstructure:"min_value"`
	MaxValue    *float64           `json:"max_value,omitempty" mapstructure:"max_value"`
}

type AppCommandChoice struct {
	Name  string `json:"name" mapstructure:"name"`
	Value string `json:"value" mapstructure:"value"`
}
