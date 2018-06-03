package templates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/biezhi/blade-cli/blade/utils"
	"github.com/biezhi/moe"
	"github.com/mkideal/cli"
	"github.com/mkideal/pkg/debug"
)

// BaseConfig base config
type BaseConfig struct {
	Name        string `cli:"-"`
	PackageName string
	Version     string
	RenderType  string
	BuildTool   string
}

type maker func(*cli.Context, BaseConfig) error

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
	err := fn(ctx, cfg)
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
	StartDelay     uint     `json:"start_delay"`
	Interval       uint     `json:"interval"`
	LastModifyTime int64    `json:"last_modify_time"`
	IncludeExt     []string `json:"include_ext"`
	ExcludeDir     []string `json:"exclude_dir"`
}

// CreateReloadConf create reload config file
func CreateReloadConf(param map[string]string) {
	confPath := param["AppName"] + "/.blade"

	conf := &BladeConf{
		PackageName:    param["PackageName"],
		AppName:        param["AppName"],
		MainClass:      param["PackageName"] + ".Application",
		StartDelay:     5,
		Interval:       3,
		LastModifyTime: 100,
		IncludeExt:     []string{"java", "properties"},
	}
	content, _ := json.Marshal(conf)

	utils.WriteFile(confPath, string(content))
	PrintLine(confPath)
}
