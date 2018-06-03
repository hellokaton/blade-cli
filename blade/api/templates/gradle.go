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
		cfg.TplDependency = getGradleTplDependency()
	} else {
		cfg.TplDependency = ""
	}

	if cfg.DBType == "MySQL" {
		cfg.MySQLDependency = getGradleMySQLDependency()
	} else {
		cfg.MySQLDependency = ""
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

func getGradleTplDependency() string {
	return `compile 'com.bladejava:blade-template-jetbrick:` + GetRepoLatestVersion("com.bladejava", "blade-template-jetbrick", "0.1.3") + `'
	`
}

func getGradleMySQLDependency() string {
	return `compile 'mysql:mysql-connector-java:5.1.46'
	 		compile 'io.github.biezhi:anima:` + GetRepoLatestVersion("io.github.biezhi", "anima", "0.2.2") + `'
	 		compile 'com.alibaba:druid:1.1.10'`
}
