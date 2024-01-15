package model

type Domain struct {
	Group []Group `json:"group,omitempty"`
}

type Group struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Activated   *bool    `json:"activated,omitempty"`
	Config      []Config `json:"config,omitempty"`
}

type Config struct {
	Key         string     `json:"key,omitempty"`
	Description string     `json:"description,omitempty"`
	Activated   *bool      `json:"activated,omitempty"`
	Strategies  []Strategy `json:"strategies,omitempty"`
	Components  []string   `json:"components,omitempty"`
}

type Strategy struct {
	Strategy  string   `json:"strategy,omitempty"`
	Activated *bool    `json:"activated,omitempty"`
	Operation string   `json:"operation,omitempty"`
	Values    []string `json:"values,omitempty"`
}

type Snapshot struct {
	Domain Domain `json:"domain,omitempty"`
}
