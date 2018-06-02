package templates

import (
	"fmt"
	"os"
	"strings"

	"github.com/biezhi/blade-cli/blade/utils"
	"github.com/mkideal/cli"
)

var _ = register("Maven", Maven)

// Maven create maven application
func Maven(ctx *cli.Context, cfg BaseConfig) error {
	appDir := cfg.Name

	param := make(map[string]string)
	param["BladeVersion"] = GetRepoLatestVersion("blade-mvc", "2.0.8-R1")
	param["JetbrickVersion"] = GetRepoLatestVersion("blade-template-jetbrick", "0.1.3")
	param["AppName"] = cfg.Name
	param["PackageName"] = cfg.PackageName
	param["Version"] = cfg.Version

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

	CreateReloadConf(param)

	gitignorePath := appDir + "/.gitignore"
	utils.WriteFile(gitignorePath, TplGitignore)
	PrintLine(gitignorePath)

	// create java、resources dir
	packagePath := appDir + "/src/main/java/" + strings.Replace(cfg.PackageName, ".", "/", -1)
	controllerPath := packagePath + "/controller"
	templatePath := appDir + "/src/main/resources/templates"

	applicationPath := packagePath + "/Application.java"
	indexController := controllerPath + "/IndexController.java"
	indexHTML := templatePath + "/index.html"
	appProperties := appDir + "/src/main/resources/app.properties"

	os.MkdirAll(packagePath, os.ModePerm)
	os.Mkdir(controllerPath, os.ModePerm)
	os.MkdirAll(appDir+"/src/test/java", os.ModePerm)
	os.MkdirAll(templatePath, os.ModePerm)
	os.Mkdir(appDir+"/src/main/resources/static", os.ModePerm)

	// app.properties
	if flag, _ := utils.Exists(appProperties); !flag {
		utils.WriteFile(appProperties, TplAppProperties)
		PrintLine(appProperties)
	}

	// 创建模板文件
	if flag, _ := utils.Exists(indexHTML); !flag {
		utils.WriteFile(indexHTML, TplIndexHTML)
		PrintLine(indexHTML)
	}

	// 创建运行类
	if flag, _ := utils.Exists(applicationPath); !flag {
		utils.WriteTemplate("tpl_application", applicationPath, TplApplication, param)
		PrintLine(applicationPath)
	}

	// 创建控制器
	if flag, _ := utils.Exists(indexController); !flag {
		utils.WriteTemplate("tpl_controller", indexController, TplController, param)
		PrintLine(indexController)
	}

	fmt.Println("")
	return nil
}
