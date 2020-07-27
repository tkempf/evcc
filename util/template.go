package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kardianos/osext"
	"github.com/spf13/viper"
)

type TemplateManager struct {
	config map[string]interface{}
}

func NewTemplateManager(config map[string]interface{}) *TemplateManager {
	return &TemplateManager{
		config: config,
	}
}

func (t *TemplateManager) Templates(class string) []interface{} {
	conf, ok := t.config[strings.ToLower(class)]
	if !ok {
		panic("unknown class: " + class)
	}
	mapS, ok := conf.([]interface{})
	if !ok {
		panic(fmt.Sprintf("wrong type for class: %s, got %T", class, conf))
	}
	return mapS
}

func DefaultConfig() (*viper.Viper, error) {
	dist := viper.New()

	// configure search path
	if filename, err := osext.Executable(); err != nil {
		fmt.Println(filepath.Dir(filename))
	}
	if home, err := os.UserHomeDir(); err == nil {
		dist.AddConfigPath(home)
	}
	dist.AddConfigPath(".")

	dist.SetConfigName("evcc.dist")
	return dist, dist.ReadInConfig()
}

func init() {
	v, err := DefaultConfig()
	if err != nil {
		panic(err)
	}

	conf := v.AllSettings()
	tm := NewTemplateManager(conf)
	meters := tm.Templates("meters")

	fmt.Printf("%+v\n", meters)

	for _, el := range meters {
		fmt.Printf("%+v\n", el)
		mapS, ok := el.(map[interface{}]interface{})
		if !ok {
			panic(fmt.Sprintf("wrong type got %T", el))
		}

		typ := mapS["type"].(string)
		if typ == "default" {
			continue
		}

		// fmt.Printf("%+v\n", mapS)
		println()
		println(strings.ToUpper(typ))
		for k, v := range mapS {
			fmt.Printf("%s: %v\n", k.(string), v)
		}
	}

	os.Exit(0)
}
