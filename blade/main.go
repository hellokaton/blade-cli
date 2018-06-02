package main

import (
	"fmt"
	"os"

	"github.com/biezhi/blade-cli/blade/api"
	"github.com/mkideal/cli"
)

const BladeCliVersion = "v0.0.1"

// root command
type rootT struct {
	cli.Helper
	Version bool `cli:"v,version" usage:"display blade cli version"`
}

// Banner blade banner
const Banner = `
    __, _,   _, __, __,
    |_) |   /_\ | \ |_
    |_) | , | | |_/ |
    ~   ~~~ ~ ~ ~   ~~~
    :: Blade Cli :: (` + BladeCliVersion + `) 
	
    Inspired by https://lets-blade.com`

var root = &cli.Command{
	Desc: Banner,
	Argv: func() interface{} { return new(rootT) },
	Fn: func(ctx *cli.Context) error {
		if len(ctx.FormValues()) == 0 {
			ctx.WriteUsage()
			return nil
		}
		argv := ctx.Argv().(*rootT)
		if argv.Version {
			ctx.String("\n" + BladeCliVersion + "\n")
		}
		return nil
	},
}

func main() {
	if err := cli.Root(root,
		cli.Tree(cli.HelpCommand("display help information")),
		cli.Tree(api.New()),
		cli.Tree(api.Serve()),
		cli.Tree(api.Build()),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
