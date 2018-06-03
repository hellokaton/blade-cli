package api

import (
	"fmt"

	"github.com/AlecAivazis/survey"
	"github.com/biezhi/blade-cli/blade/api/templates"
	"github.com/mkideal/cli"
)

// New new blade application
func New() *cli.Command {
	return &cli.Command{
		Name:        "new",
		Desc:        "create blade application by template",
		Text:        `    blade new <name>`,
		Argv:        func() interface{} { return new(newT) },
		CanSubRoute: true,
		Fn: func(ctx *cli.Context) error {
			argv := ctx.Argv().(*newT)
			argv.Version = "0.0.1"

			prompt := &survey.Input{
				Message: "please input package name (e.g: com.bladejava.example):",
			}
			survey.AskOne(prompt, &argv.PackageName, nil)

			if argv.PackageName == "" {
				argv.PackageName = "com.bladejava.example"
			}
			fmt.Println("")

			buildPrompt := &survey.Select{
				Message: "select the build tool:",
				Options: []string{"Maven", "Gradle"},
			}
			survey.AskOne(buildPrompt, &argv.BuildTool, nil)
			fmt.Println("")

			renderPrompt := &survey.Select{
				Message: "select the type of site you built:",
				Options: []string{"Web Application", "Restful API"},
			}
			survey.AskOne(renderPrompt, &argv.RenderType, nil)
			fmt.Println("")

			dbPrompt := &survey.Select{
				Message: "used database?",
				Options: []string{"No database", "MySQL"},
			}
			survey.AskOne(dbPrompt, &argv.DBType, nil)
			fmt.Println("")

			argv.BladeVersion = templates.GetRepoLatestVersion("com.bladejava", "blade-mvc", "2.0.8-R1")

			return templates.New(ctx, argv.BaseConfig)
		},
	}
}

type newT struct {
	cli.Helper
	templates.BaseConfig
}

func (t *newT) Validate(ctx *cli.Context) error {
	clr := ctx.Color()
	b := clr.Bold
	if len(ctx.Args()) == 0 || ctx.Args()[0] == "" {
		return fmt.Errorf("application %s is empty", b("name"))
	}
	if len(ctx.Args()) > 1 {
		return fmt.Errorf("too many args for %s", b("name"))
	}
	t.Name = ctx.Args()[0]
	return nil
}
