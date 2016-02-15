package hue

import (
	"encoding/json"
)

// Group - encapsulates the controls for multiple Hue lights in unison
type Group struct {
	ID     string
	Name   string
	Type   string
	Class  string
	Lights []string

	Bridge *Bridge
}


// GetAllGroups - retrieves all groups added to Hue except #0
func (b *Bridge) GetAllGroups() ([]*Group, error) {
	// fetch all the groups
	response, err := b.get("/groups")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// deconstruct the json results
	var results map[string]Group
	err = json.NewDecoder(response.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	// and convert them into groups
	var groups []*Group
	for id, group := range results {
	  group.ID = id
	  group.Bridge = b
		groups = append(groups, &group)
	}

	return groups, nil
}
