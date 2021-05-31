package project

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/1024casts/snake/cmd/snake/internal/base"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

// Project is a project template.
type Project struct {
	Name string
}

// New new a project from remote repo.
func (p *Project) New(ctx context.Context, dir string, layout string) error {
	to := path.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s already exists\n", p.Name)
		override := false
		prompt := &survey.Confirm{
			Message: "📂 Do you want to override the folder ?",
			Help:    "Delete the existing folder and create the project.",
		}
		survey.AskOne(prompt, &override)
		if !override {
			return err
		}
		os.RemoveAll(to)
	}

	fmt.Printf("🚀 Creating service %s, layout repo is %s, please wait a moment.\n\n", p.Name, layout)
	repo := base.NewRepo(layout)
	if err := repo.CopyTo(ctx, to, p.Name, []string{".git", ".github"}); err != nil {
		return err
	}
	os.Rename(
		path.Join(to, "cmd", "server"),
		path.Join(to, "cmd", p.Name),
	)
	base.Tree(to, dir)

	fmt.Printf("\n🍺 Project creation succeeded %s\n", color.GreenString(p.Name))
	fmt.Print("💻 Use the following command to start the project 👇:\n\n")

	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	fmt.Println(color.WhiteString("$ go build"))
	fmt.Println(color.WhiteString("$ ./%s\n", p.Name))
	fmt.Println("🤝 Thanks for using Snake")
	return nil
}
