package plugin

import (
	"fmt"
	"reflect"
)

// you should define your plugin struct name of plugin file name.
type Stdout struct {
	config interface{}
	prefix string
}

// you should implement Plugin interface.
func (so *Stdout) Notice(s string) {
	fmt.Printf("%s%s\n", so.prefix, s)
}

func (so *Stdout) SetPluginConfig(cfg interface{}) {
	so.config = cfg
	if cfg, ok := cfg.(map[string]interface{}); ok {
		if prefix, ok := cfg["prefix"].(string); ok {
			so.prefix = prefix
		}
	}
}

// you should register your plugin struct in init function.
func init() {
	registerPluginType(reflect.TypeOf(Stdout{}))
}
