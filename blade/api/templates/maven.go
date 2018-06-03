package templates

import (
	"fmt"
	"os"

	"github.com/biezhi/blade-cli/blade/utils"
	"github.com/mkideal/cli"
)

var _ = register("Maven", Maven)

// Maven create maven application
func Maven(ctx *cli.Context, cfg *BaseConfig) error {
	appDir := cfg.Name

	if cfg.RenderType == "Web Application" {
		cfg.TplDependency = getMavenDependency()
	} else {
		cfg.TplDependency = ""
	}

	// create dir
	if err := os.Mkdir(appDir, 0755); err != nil {
		return err
	}

	// create pom.xml
	pomPath := appDir + "/pom.xml"
	if flag, _ := utils.Exists(pomPath); !flag {
		utils.WriteTemplate("tpl_pom", pomPath, TplPom, cfg)
		fmt.Println("\n\ncreate file success:", pomPath)
	}

	packageXML := appDir + "/package.xml"
	utils.WriteFile(packageXML, TplPackageXML)
	PrintLine(packageXML)

	WriteCommon(cfg)
	return nil
}

func getMavenDependency() string {
	return `<dependency>
			<groupId>com.bladejava</groupId>
			<artifactId>blade-template-jetbrick</artifactId>
			<version>` + GetRepoLatestVersion("blade-template-jetbrick", "0.1.3") + `</version>
		</dependency>`
}
