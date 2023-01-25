package github

import (
	"encoding/json"

	"github.com/razonyang/api/internal/helper"
)

type Dependents struct {
	Repositories int `json:"repositories"`
	Packages     int `json:"packages"`
}

func (d *Dependents) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"repositories":          d.Repositories,
		"repositories_humanize": helper.FormatInt(d.Repositories),
		"packages":              d.Packages,
		"packages_humanize":     helper.FormatInt(d.Packages),
	})
}
