package hugo

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Config struct {
	Module Module `yaml:"module"`
}

type Module struct {
	HugoVersion HugoVersion `json:"hugo_version" yaml:"hugoVersion"`
}

type HugoVersion struct {
	Extended bool   `json:"extended"`
	Min      string `json:"min"`
	Max      string `json:"max"`
}

func (m Module) requirements() string {
	r := []string{}
	if m.HugoVersion.Min != "" {
		r = append(r, fmt.Sprintf(">=%s", m.HugoVersion.Min))
	}
	if m.HugoVersion.Max != "" {
		r = append(r, fmt.Sprintf("<=%s", m.HugoVersion.Max))
	}
	if m.HugoVersion.Extended {
		r = append(r, "extended")
	}
	if len(r) > 0 {
		return strings.Join(r, " ")
	}

	return "*"
}

func (m Module) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"hugo_version": m.HugoVersion,
		"requirements": m.requirements(),
	})
}
