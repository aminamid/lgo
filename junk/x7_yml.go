package main

import (
    "fmt"
    "gopkg.in/yaml.v3"
)

var (
    yamlval = `
key-1: val-1
key-2: 2
key-3:
  key-3-1: val-3-1
4: 4
`
)

func main() {
    ret2 := make(map[string]interface{})
    yaml.Unmarshal([]byte(yamlval), &ret2)
    fmt.Printf("%#v\n", ret2)
    d, err := yaml.Marshal(ret2)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	fmt.Printf("%s\n", d)
}
