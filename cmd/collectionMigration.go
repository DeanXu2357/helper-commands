/*
Copyright © 2022 Dean_Hsu <jasugun0000@gmail.com>

*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/DeanXu2357/helper-commands/tpl"
	"github.com/spf13/cobra"
	"log"
	"os"
	"text/template"
	"time"
)

// collectionMigrationCmd represents the collectionMigration command
var collectionMigrationCmd = &cobra.Command{
	Use:   "collectionMigration [collection name]",
	Short: "撰寫 migration 新增 collection",
	Long:  `撰寫 migration 新增 collection`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("collectionMigration called")

		g := NewCreateCollectionGenerator(args[0])

		cobra.CheckErr(g.Create())
	},
}

func init() {
	rootCmd.AddCommand(collectionMigrationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// collectionMigrationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// collectionMigrationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type createColGenerator struct {
	path           string
	version        string
	collectionName string
}

func NewCreateCollectionGenerator(colName string) *createColGenerator {
	version := time.Now().Format("20060102")

	fileName := fmt.Sprintf("%s_Create_%s_Collection.js", version, colName)

	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &createColGenerator{
		path:           fmt.Sprintf("%s/migrations/%s", root, fileName),
		version:        version,
		collectionName: colName,
	}
}

func (g *createColGenerator) Create() error {
	if _, err := os.Stat(g.path); os.IsExist(err) {
		return errors.New("file already exist")
	}

	cmdFile, err := os.Create(g.path)
	if err != nil {
		return err
	}
	defer cmdFile.Close()

	commandTemplate := template.Must(
		template.New("sub").Parse(string(tpl.MigrationCreateTemplate())),
	)
	err = commandTemplate.Execute(cmdFile, g)
	if err != nil {
		return err
	}
	return nil
}
