package model

import (
	"fmt"
	"strings"
)

type Unicorn struct {
	Name         string   `json:"name"`
	Capabilities []string `json:"capabilities"`
}

func (u Unicorn) String() string {
	return fmt.Sprintf("Unicorn{Name: %s, Capabilities: %s}", u.Name, strings.Join(u.Capabilities, ","))
}
