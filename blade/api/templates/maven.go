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
		cfg.TplDependency = getMavenTplDependency()
	} else {
		cfg.TplDependency = ""
	}

	if cfg.DBType == "MySQL" {
		cfg.MySQLDependency = getMavenMySQLDependency()
	} else {
		cfg.MySQLDependency = ""
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

func getMavenTplDependency() string {
	return `<dependency>
			<groupId>com.bladejava</groupId>
			<artifactId>blade-template-jetbrick</artifactId>
			<version>` + GetRepoLatestVersion("com.bladejava", "blade-template-jetbrick", "0.1.3") + `</version>
		</dependency>`
}

func getMavenMySQLDependency() string {
	return `<dependency>
			<groupId>mysql</groupId>
			<artifactId>mysql-connector-java</artifactId>
			<version>5.1.46</version>
		</dependency>
		<dependency>
            <groupId>io.github.biezhi</groupId>
            <artifactId>anima</artifactId>
            <version>` + GetRepoLatestVersion("io.github.biezhi", "anima", "0.2.2") + `</version>
		</dependency>
		<dependency>
            <groupId>com.alibaba</groupId>
            <artifactId>druid</artifactId>
            <version>1.1.10</version>
        </dependency>
		`
}
