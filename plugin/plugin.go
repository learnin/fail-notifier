package plugin

import (
	"reflect"
)

type Plugin interface {
	Notice(s string)
	SetPluginConfig(cfg interface{})
}

var pluginTypeMap = map[string]reflect.Type{}

func registerPluginType(t reflect.Type) {
	pluginTypeMap[t.String()] = t
}

func CreatePlugin(typeName string, cfg interface{}) (*Plugin, bool) {
	i := reflect.New(pluginTypeMap[typeName]).Interface()
	// Type Assertion
	p, ok := i.(Plugin)
	if !ok {
		return &p, false
	}
	p.SetPluginConfig(cfg)
	return &p, true
}
