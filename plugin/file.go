package plugin

import (
	"bufio"
	"os"
	"reflect"
)

// you should define your plugin struct name of plugin file name.
type File struct {
	config interface{}
	path   string
}

// you should implement Plugin interface.
func (f *File) Notice(s string) {
	file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		println("can't open or create " + f.path + ": " + err.Error())
		return
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	if _, err := w.WriteString(s); err != nil {
		println("can't write to " + f.path + ": " + err.Error())
		return
	}
	w.Flush()
}

func (f *File) SetPluginConfig(cfg interface{}) {
	f.config = cfg
	if cfg, ok := cfg.(map[string]interface{}); ok {
		if path, ok := cfg["path"].(string); ok {
			f.path = path
		}
	}
}

// you should register your plugin struct in init function.
func init() {
	registerPluginType(reflect.TypeOf(File{}))
}
