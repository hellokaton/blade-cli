package api

import (
	"github.com/mkideal/cli"
)

// Generator generator model, service, controller
func Generator() *cli.Command {
	return &cli.Command{
		Name:        "generator",
		Desc:        "generator model, controller, templates",
		Text:        `    blade generator`,
		CanSubRoute: true,
		Fn: func(ctx *cli.Context) error {
			// argv := ctx.Argv().(*newT)
			return nil
		},
	}
}
