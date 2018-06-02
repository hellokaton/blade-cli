package api

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/biezhi/blade-cli/blade/utils"
	"github.com/mkideal/cli"
)

// Build package
func Build() *cli.Command {
	return &cli.Command{
		Name:        "build",
		Desc:        "build application as jar or dir",
		Text:        `    blade build`,
		CanSubRoute: true,

		Fn: func(ctx *cli.Context) error {
			cmd, stdout, stderr, err := utils.StartCmd("mvn clean package")
			if err != nil {
				return err
			}

			io.Copy(os.Stdout, stdout)
			errMsg, err := ioutil.ReadAll(stderr)
			if err != nil {
				return err
			}
			// wait for building
			err = cmd.Wait()
			if err != nil {
				e := fmt.Sprintf("stderr: %s, cmd err: %s", string(errMsg), err)
				return errors.New(e)
			}
			return nil
		},
	}

}
