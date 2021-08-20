package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"unicode"
)

// Cache is a cache generator.
type Cache struct {
	Name    string
	Path    string
	Service string
	Package string
}

// Generate generate a cache template.
func (c *Cache) Generate() error {
	body, err := c.execute()
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	to := path.Join(wd, c.Path)
	if _, err := os.Stat(to); os.IsNotExist(err) {
		if err := os.MkdirAll(to, 0700); err != nil {
			return err
		}
	}
	name := path.Join(to, Lcfirst(c.Name))
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", c.Name)
	}
	return ioutil.WriteFile(name, body, 0644)
}

// 首字母小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
