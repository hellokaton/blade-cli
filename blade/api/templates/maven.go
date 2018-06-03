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

	param := make(map[string]string)
	param["BladeVersion"] = GetRepoLatestVersion("blade-mvc", "2.0.8-R1")
	param["AppName"] = cfg.Name
	param["PackageName"] = cfg.PackageName
	param["Version"] = cfg.Version
	param["BuildTool"] = "maven"
	param["RenderType"] = cfg.RenderType

	if cfg.RenderType == "Web Application" {
		param["TplDependency"] = getMavenDependency()
	} else {
		param["TplDependency"] = ""
	}

	// create dir
	if err := os.Mkdir(appDir, 0755); err != nil {
		return err
	}

	// create pom.xml
	pomPath := appDir + "/pom.xml"
	if flag, _ := utils.Exists(pomPath); !flag {
		utils.WriteTemplate("tpl_pom", pomPath, TplPom, param)
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
