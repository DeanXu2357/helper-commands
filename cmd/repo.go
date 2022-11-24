package cmd

import (
	"fmt"
	"github.com/DeanXu2357/helper-commands/tpl"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type repoGenerator struct {
	Path           string
	FileName       string
	EntityImport   string
	EntityName     string
	CollectionName string
}

func NewGenerator(importFrom, entityName, filePath string) *repoGenerator {
	fileName := toSnakeCase(entityName) + "_repo"

	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	p := root + "/" + fileName
	if filePath != "" {
		p = root + "/" + fileName
	}

	return &repoGenerator{
		Path:           p,
		FileName:       fileName,
		EntityImport:   importFrom,
		EntityName:     toCapitalFirst(entityName),
		CollectionName: toLowerFirst(entityName),
	}
}

func (r *repoGenerator) Create() error {
	// check if AbsolutePath exists
	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(r.Path, 0754); err != nil {
			return err
		}
	}

	cmdFile, err := os.Create(fmt.Sprintf("%s/%s.go", r.Path, r.FileName))
	if err != nil {
		return err
	}
	defer cmdFile.Close()

	commandTemplate := template.Must(template.New("sub").Parse(string(tpl.RepoTemplate())))
	err = commandTemplate.Execute(cmdFile, r)
	if err != nil {
		return err
	}
	return nil
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(s string) string {
	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func toCapitalFirst(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}

func toLowerFirst(s string) string {
	return strings.ToLower(string(s[0])) + s[1:]
}
