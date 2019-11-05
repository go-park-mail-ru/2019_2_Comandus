package model

import "github.com/microcosm-cc/bluemonday"

type Role struct {
	Role string `json:"role"`
	Label string `json:"label"`
	Avatar string `json:"avatar"`
}

func (role *Role) Sanitize (sanitizer *bluemonday.Policy)  {
	sanitizer.Sanitize(role.Role)
	sanitizer.Sanitize(role.Label)
	sanitizer.Sanitize(role.Avatar)
}