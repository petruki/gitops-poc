package model

type Domain struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Activated   bool    `json:"activated"`
	Group       []Group `json:"group"`
}

type Group struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Activated   bool     `json:"activated"`
	Config      []Config `json:"config"`
}

type Config struct {
	Key         string     `json:"key"`
	Description string     `json:"description"`
	Activated   bool       `json:"activated"`
	Strategies  []Strategy `json:"strategies"`
	Components  []string   `json:"components"`
}

type Strategy struct {
	Strategy  string   `json:"strategy"`
	Activated bool     `json:"activated"`
	Operation string   `json:"operation"`
	Values    []string `json:"values"`
}

type Snapshot struct {
	Domain Domain `json:"domain"`
}
