package plugin

import (
	"reflect"
)

type Plugin interface {
	Notice(s string)
}

var pluginTypeMap = map[string]reflect.Type{}

func registerPluginType(t reflect.Type) {
	pluginTypeMap[t.String()] = t
}

func CreatePlugin(typeName string) (Plugin, bool) {
	i := reflect.New(pluginTypeMap[typeName]).Elem().Interface()
	// Type Assertion
	p, ok := i.(Plugin)
	return p, ok
}
