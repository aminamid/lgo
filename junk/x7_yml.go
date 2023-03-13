package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed x7.yml
	yamlval string
	//go:embed x7.json
	jsonval string

	err error
)

func main() {
	ret1 := make(map[string]interface{})
	err = json.Unmarshal([]byte(jsonval), &ret1)
	if err != nil{ fmt.Printf("%#v\n", err)}
	fmt.Printf("ret1 ---\n%#v\n---\n", ret1)

	ret2 := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(yamlval), &ret2)
	if err != nil {fmt.Printf("%#v\n", err)}

	fmt.Printf("ret2 ---\n%#v\n---\n", ret2)
	d, err := yaml.Marshal(ret2)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	fmt.Printf("%s\n", d)

}
