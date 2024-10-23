package pkg

import "stash-go/model/entity"

type NamedPackage struct {
	Name     string            `json:"name"`
	Packages []*entity.Package `json:"packages"`
}
