package model

import (
	"fmt"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// CmdNew represents the proto command.
var CmdNew = &cobra.Command{
	Use:   "model",
	Short: "Generate the model file.",
	Long:  "Generate the model file via database table.",
	Run:   run,
}

var (
	filename    string
	host        string
	port        string
	database    string
	tableName   string
	packageName string
	structName  string
	user        string
	password    string
	format      string
	targetDir   string
)

func init() {
	// default value
	host = "localhost"
	port = "3306"
	packageName = "model"
	user = "root"
	password = "123456"
	targetDir = "internal/model"

	CmdNew.Flags().StringVarP(&filename, "filename", "f", filename, "model filename")
	CmdNew.Flags().StringVar(&host, "host", host, "database host addr")
	CmdNew.Flags().StringVar(&port, "port", port, "database port")
	CmdNew.Flags().StringVarP(&database, "database", "d", database, "database name")
	CmdNew.Flags().StringVarP(&tableName, "table", "t", tableName, "table name")
	CmdNew.Flags().StringVarP(&user, "user", "u", user, "database username")
	CmdNew.Flags().StringVarP(&password, "password", "p", password, "password for database")
	CmdNew.Flags().StringVar(&packageName, "package", packageName, "package name")
	CmdNew.Flags().StringVarP(&structName, "struct", "s", structName, "model struct name")
	CmdNew.Flags().StringVar(&format, "format", format, "add json annotations (default)")
	CmdNew.Flags().StringVar(&targetDir, "target-dir", targetDir, "model target dir")
}

func run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// TODO: read database config

	// get file name
	if len(filename) == 0 {
		prompt := &survey.Input{
			Message: "What is file name ?",
			Help:    "Created file name.",
		}
		err = survey.AskOne(prompt, &filename)
		if filename == "" || err != nil {
			return
		}
	}
	if len(host) == 0 {
		prompt := &survey.Input{
			Message: "What is host for database?",
			Help:    "Created host for database.",
			Default: host,
		}
		err = survey.AskOne(prompt, &host)
		if host == "" || err != nil {
			return
		}
	}
	if len(port) == 0 {
		prompt := &survey.Input{
			Message: "What is port for database?",
			Help:    "Created port for database.",
			Default: port,
		}
		err = survey.AskOne(prompt, &port)
		if port == "" || err != nil {
			return
		}
	}
	if len(database) == 0 {
		prompt := &survey.Input{
			Message: "What is database name ?",
			Help:    "Created database name.",
		}
		err = survey.AskOne(prompt, &database)
		if database == "" || err != nil {
			return
		}
	}
	if len(tableName) == 0 {
		prompt := &survey.Input{
			Message: "What is table name ?",
			Help:    "Created table name.",
		}
		err = survey.AskOne(prompt, &tableName)
		if tableName == "" || err != nil {
			return
		}
	}
	if len(user) == 0 {
		prompt := &survey.Input{
			Message: "What is user for database?",
			Help:    "Created user name.",
			Default: user,
		}
		err = survey.AskOne(prompt, &user)
		if user == "" || err != nil {
			return
		}
	}
	if len(password) == 0 {
		prompt := &survey.Input{
			Message: "What is password for database?",
			Help:    "Created password.",
		}
		err = survey.AskOne(prompt, &password)
		if user == "" || err != nil {
			return
		}
	}
	if len(packageName) == 0 {
		prompt := &survey.Input{
			Message: "What is package name for model?",
			Help:    "Created package name.",
			Default: packageName,
		}
		err = survey.AskOne(prompt, &packageName)
		if packageName == "" || err != nil {
			return
		}
	}
	if len(structName) == 0 {
		prompt := &survey.Input{
			Message: "What is struct name for model?",
			Help:    "Created struct name.",
		}
		err = survey.AskOne(prompt, &structName)
		if filename == "" || err != nil {
			return
		}
	}

	m := &Model{
		Filename:    filename,
		Host:        host,
		MysqlPort:   port,
		Database:    database,
		TableName:   tableName,
		PackageName: packageName,
		StructName:  structName,
		User:        user,
		Password:    password,
		Format:      format,
		TargetDir:   path.Join(wd, targetDir),
	}
	if err := m.Generate(); err != nil {
		fmt.Println(err)
		return
	}
}
