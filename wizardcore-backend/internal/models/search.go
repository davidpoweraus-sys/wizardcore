package models

type SearchResults struct {
	Pathways  []Pathway  `json:"pathways"`
	Exercises []Exercise `json:"exercises"`
	Users     []User     `json:"users"`
}