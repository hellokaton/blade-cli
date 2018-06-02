package api

import (
	"fmt"

	"github.com/mkideal/cli"
)

func Build() *cli.Command {
	return build_
}

type buildT struct {
	cli.Helper
	Name string `cli:"-"`
}

func (t *buildT) Validate(ctx *cli.Context) error {
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

var build_ = &cli.Command{
	Name:        "build",
	Desc:        "build application as jar or dir",
	Text:        `    blade build [name]`,
	Argv:        func() interface{} { return new(newT) },
	CanSubRoute: true,

	Fn: func(ctx *cli.Context) error {
		// argv := ctx.Argv().(*newT)

		return nil
		// return templates.New(argv.Type, ctx, nil)
	},
}
