package templates

import (
	"fmt"
	"os"

	"github.com/biezhi/blade-cli/blade/utils"
	"github.com/mkideal/cli"
)

var _ = register("Gradle", Gradle)

// Gradle create gradle application
func Gradle(ctx *cli.Context, cfg *BaseConfig) error {
	appDir := cfg.Name
	if cfg.RenderType == "Web Application" {
		cfg.TplDependency = gradleTplDependency()
	} else {
		cfg.TplDependency = ""
	}

	// create dir
	if err := os.Mkdir(appDir, 0755); err != nil {
		return err
	}

	// create build.gradle
	buildPath := appDir + "/build.gradle"
	if flag, _ := utils.Exists(buildPath); !flag {
		utils.WriteTemplate("tpl_build_gradle", buildPath, TplGradleBuild, cfg)
		fmt.Println("\n\ncreate file success:", buildPath)
	}

	settingPath := appDir + "/setting.gradle"
	if flag, _ := utils.Exists(settingPath); !flag {
		utils.WriteTemplate("tpl_setting_gradle", settingPath, TplGradleSetting, cfg)
		PrintLine(settingPath)
	}

	WriteCommon(cfg)
	return nil
}

func gradleTplDependency() string {
	return `compile 'com.bladejava:blade-template-jetbrick:` + GetRepoLatestVersion("blade-template-jetbrick", "0.1.3") + `'
	`
}
