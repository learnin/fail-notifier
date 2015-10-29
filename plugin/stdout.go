package plugin

import (
	"fmt"
	"reflect"
)

// you should define your plugin struct name of plugin file name.
type Stdout struct {
}

// you should implement Plugin interface.
func (a Stdout) Notice(s string) {
	fmt.Printf("%s", s)
}

// you should register your plugin struct in init function.
func init() {
	registerPluginType(reflect.TypeOf(Stdout{}))
}
