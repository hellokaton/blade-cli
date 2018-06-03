package templates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/biezhi/blade-cli/blade/utils"
	"github.com/biezhi/moe"
	"github.com/mkideal/cli"
	"github.com/mkideal/pkg/debug"
)

// BaseConfig base config
type BaseConfig struct {
	Name          string `cli:"-"`
	PackageName   string
	Version       string
	RenderType    string
	BuildTool     string
	BladeVersion  string
	TplDependency string
}

type maker func(*cli.Context, *BaseConfig) error

var templatesMap = make(map[string]maker)

func register(tool string, fn maker) bool {
	if _, ok := templatesMap[tool]; ok {
		debug.Panicf("repeat register template %s", tool)
	}
	templatesMap[tool] = fn
	return true
}

// New new application
func New(ctx *cli.Context, cfg BaseConfig) error {
	clr := ctx.Color()
	moe := moe.New(clr.Bold("creating project, please wait...")).Spinner("dots3").Color(moe.Green).Start()

	fn, ok := templatesMap[cfg.BuildTool]
	if !ok {
		return fmt.Errorf("unsupported template type %s", clr.Yellow(cfg.BuildTool))
	}
	err := fn(ctx, &cfg)
	moe.Stop()
	if err == nil {
		fmt.Printf("application %s create successful!\n", cfg.Name)
		fmt.Printf("\n    $ cd %s", cfg.Name)
		fmt.Printf("\n    $ blade serve\n")
	}
	return err
}

// PrintLine print create file success
func PrintLine(message string) {
	fmt.Println("create file success:", message)
}

// GetRepoLatestVersion get repo latest version
func GetRepoLatestVersion(artifactID, defaultVersion string) string {
	resp, err := http.Get("http://search.maven.org/solrsearch/select?q=g:%20com.bladejava%20+AND+a:%20" + artifactID + "%20&rows=1&wt=json")
	if err != nil {
		fmt.Println(err)
		return defaultVersion
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("body:", string(body))
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println(err)
		return defaultVersion
	}
	response := result["response"].(map[string]interface{})
	docs := response["docs"].([]interface{})
	doc := docs[0].(map[string]interface{})
	return doc["latestVersion"].(string)
}

// BladeConf blade config
type BladeConf struct {
	PackageName    string   `json:"package"`
	AppName        string   `json:"app_name"`
	MainClass      string   `json:"main_class"`
	BuildTool      string   `json:"build_tool"`
	StartDelay     uint     `json:"start_delay"`
	Interval       uint     `json:"interval"`
	LastModifyTime int64    `json:"last_modify_time"`
	IncludeExt     []string `json:"include_ext"`
	ExcludeDir     []string `json:"exclude_dir"`
}

// CreateReloadConf create reload config file
func CreateReloadConf(cfg *BaseConfig) {
	confPath := cfg.Name + "/.blade"

	conf := &BladeConf{
		PackageName:    cfg.PackageName,
		AppName:        cfg.Name,
		MainClass:      cfg.PackageName + ".Application",
		BuildTool:      cfg.BuildTool,
		StartDelay:     5,
		Interval:       3,
		LastModifyTime: 100,
		IncludeExt:     []string{"java", "properties"},
	}
	content, _ := json.Marshal(conf)

	utils.WriteFile(confPath, string(content))
	PrintLine(confPath)
}

func WriteCommon(cfg *BaseConfig) {
	appDir := cfg.Name

	CreateReloadConf(cfg)

	gitignorePath := appDir + "/.gitignore"
	utils.WriteFile(gitignorePath, TplGitignore)
	PrintLine(gitignorePath)

	// create java、resources dir
	packagePath := appDir + "/src/main/java/" + strings.Replace(cfg.PackageName, ".", "/", -1)
	configPath := packagePath + "/config"
	controllerPath := packagePath + "/controller"

	applicationPath := packagePath + "/Application.java"
	bootstrapPath := configPath + "/Bootstrap.java"
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
		utils.WriteTemplate("tpl_application", applicationPath, TplApplication, cfg)
		PrintLine(applicationPath)
	}

	// create Bootstrap
	if flag, _ := utils.Exists(bootstrapPath); !flag {
		utils.WriteTemplate("tpl_bootstrap", bootstrapPath, TplBootstrap, cfg)
		PrintLine(applicationPath)
	}

	// create controller
	if flag, _ := utils.Exists(indexController); !flag {
		utils.WriteTemplate("tpl_controller", indexController, TplController, cfg)
		PrintLine(indexController)
	}
	fmt.Println("")
}
