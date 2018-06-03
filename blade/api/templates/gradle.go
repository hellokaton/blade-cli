package templates

import (
	"fmt"
	"os"
	"strings"

	"github.com/biezhi/blade-cli/blade/utils"
	"github.com/mkideal/cli"
)

var _ = register("Gradle", Gradle)

// Gradle create gradle application
func Gradle(ctx *cli.Context, cfg BaseConfig) error {
	appDir := cfg.Name

	param := make(map[string]string)
	param["BladeVersion"] = GetRepoLatestVersion("blade-mvc", "2.0.8-R1")
	param["AppName"] = cfg.Name
	param["PackageName"] = cfg.PackageName
	param["Version"] = cfg.Version
	param["BuildTool"] = "gradle"

	if cfg.RenderType == "Web Application" {
		param["TplDependency"] = gradleTplDependency()
	} else {
		param["TplDependency"] = ""
	}

	// create dir
	if err := os.Mkdir(appDir, 0755); err != nil {
		return err
	}

	// create build.gradle
	buildPath := appDir + "/build.gradle"
	if flag, _ := utils.Exists(buildPath); !flag {
		utils.WriteTemplate("tpl_build_gradle", buildPath, TplGradleBuild, param)
		fmt.Println("\n\ncreate file success:", buildPath)
	}

	settingPath := appDir + "/setting.gradle"
	if flag, _ := utils.Exists(settingPath); !flag {
		utils.WriteTemplate("tpl_setting_gradle", settingPath, TplGradleSetting, param)
		PrintLine(settingPath)
	}

	CreateReloadConf(param)

	gitignorePath := appDir + "/.gitignore"
	utils.WriteFile(gitignorePath, TplGitignore)
	PrintLine(gitignorePath)

	// create java、resources dir
	packagePath := appDir + "/src/main/java/" + strings.Replace(cfg.PackageName, ".", "/", -1)
	controllerPath := packagePath + "/controller"

	applicationPath := packagePath + "/Application.java"
	indexController := controllerPath + "/IndexController.java"
	appProperties := appDir + "/src/main/resources/app.properties"

	os.MkdirAll(packagePath, os.ModePerm)
	os.MkdirAll(controllerPath, os.ModePerm)
	os.MkdirAll(appDir+"/src/test/java", os.ModePerm)
	os.MkdirAll(appDir+"/src/main/resources/static", os.ModePerm)

	// app.properties
	if flag, _ := utils.Exists(appProperties); !flag {
		utils.WriteFile(appProperties, TplAppProperties)
		PrintLine(appProperties)
	}

	if cfg.RenderType == "Web Application" {
		templatePath := appDir + "/src/main/resources/templates"
		indexHTML := templatePath + "/index.html"
		os.MkdirAll(templatePath, os.ModePerm)

		// create template file
		if flag, _ := utils.Exists(indexHTML); !flag {
			utils.WriteFile(indexHTML, TplIndexHTML)
			PrintLine(indexHTML)
		}
	}

	// create Application
	if flag, _ := utils.Exists(applicationPath); !flag {
		utils.WriteTemplate("tpl_application", applicationPath, TplApplication, param)
		PrintLine(applicationPath)
	}

	// create controller
	if flag, _ := utils.Exists(indexController); !flag {
		utils.WriteTemplate("tpl_controller", indexController, TplController, param)
		PrintLine(indexController)
	}

	fmt.Println("")
	return nil
}

func gradleTplDependency() string {
	return `compile 'com.bladejava:blade-template-jetbrick:` + GetRepoLatestVersion("blade-template-jetbrick", "0.1.3") + `'
	`
}
