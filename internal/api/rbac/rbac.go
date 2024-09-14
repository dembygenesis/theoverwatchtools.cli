package rbac

import "embed"

type Rbac struct {
	File embed.FS
	Name string
}

//go:embed model.conf
var modelConf embed.FS

//go:embed policies.conf
var policesConf embed.FS

var ModelConf = &Rbac{
	File: modelConf,
	Name: "model.conf",
}

var PoliciesConf = &Rbac{
	File: policesConf,
	Name: "policies.conf",
}
