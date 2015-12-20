package plugin

import (
	"log/syslog"
	"reflect"
)

// you should define your plugin struct name of plugin file name.
type Syslog struct {
	config   interface{}
	facility string
}

// you should implement Plugin interface.
func (log *Syslog) Notice(s string) {
	w, err := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, "fail-notifier")
	if err != nil {
		println("can't connect to the syslog daemon: " + err.Error())
		return
	}
	defer w.Close()
	if err := w.Err(s); err != nil {
		println("can't write to syslog: " + err.Error())
		return
	}
}

func (log *Syslog) SetPluginConfig(cfg interface{}) {
	// log.config = cfg
	// if cfg, ok := cfg.(map[string]interface{}); ok {
	// 	if facility, ok := cfg["facility"].(string); ok {
	// 		log.facility = facility
	// 	}
	// }
}

// you should register your plugin struct in init function.
func init() {
	registerPluginType(reflect.TypeOf(Syslog{}))
}
