package templates

// TplPom pom.xml
const TplPom = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>{{ .PackageName }}</groupId>
    <artifactId>{{ .Name }}</artifactId>
    <version>{{ .Version }}</version>

    <properties>
        <maven.compiler.source>1.8</maven.compiler.source>
        <maven.compiler.target>1.8</maven.compiler.target>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
        <blade-mvc.version>{{ .BladeVersion }}</blade-mvc.version>
        <junit.version>4.12</junit.version>
        <unirest-java.version>1.4.9</unirest-java.version>
    </properties>
    
    <dependencies>
        <dependency>
            <groupId>com.bladejava</groupId>
            <artifactId>blade-mvc</artifactId>
            <version>${blade-mvc.version}</version>
        </dependency>
        {{ .TplDependency }}
        <dependency>
            <groupId>junit</groupId>
            <artifactId>junit</artifactId>
            <version>${junit.version}</version>
            <scope>test</scope>
        </dependency>

        <dependency>
            <groupId>com.mashape.unirest</groupId>
            <artifactId>unirest-java</artifactId>
            <version>${unirest-java.version}</version>
            <scope>test</scope>
        </dependency>
    </dependencies>

    <profiles>
        <profile>
            <id>prod</id>
            <activation>
                <activeByDefault>false</activeByDefault>
            </activation>
            <build>
                <resources>
                    <resource>
                        <directory>src/main/java</directory>
                        <filtering>false</filtering>
                        <excludes>
                            <exclude>**/*.java</exclude>
                        </excludes>
                    </resource>
                </resources>
            </build>
        </profile>
        <profile>
            <id>dev</id>
            <activation>
                <activeByDefault>true</activeByDefault>
            </activation>
            <build>
                <resources>
                    <resource>
                        <directory>src/main/java</directory>
                        <filtering>false</filtering>
                    </resource>
                    <resource>
                        <directory>src/main/resources</directory>
                        <filtering>false</filtering>
                    </resource>
                    <resource>
                        <directory>src/main/test</directory>
                        <filtering>false</filtering>
                    </resource>
                    <resource>
                        <directory>src/test/resources</directory>
                        <filtering>false</filtering>
                    </resource>
                </resources>
            </build>
        </profile>
    </profiles>

    <build>
        <finalName>{{ .Name }}</finalName>
        <plugins>
            <plugin>
                <artifactId>maven-compiler-plugin</artifactId>
                <version>2.5.1</version>
                <configuration>
                    <source>1.8</source>
                    <target>1.8</target>
                    <encoding>UTF-8</encoding>
                </configuration>
            </plugin>
            <plugin>
                <artifactId>maven-assembly-plugin</artifactId>
                <configuration>
                    <appendAssemblyId>false</appendAssemblyId>
                    <descriptors>
                        <descriptor>package.xml</descriptor>
                    </descriptors>
                    <outputDirectory>${project.build.directory}/dist/</outputDirectory>
                </configuration>
                <executions>
                    <execution>
                        <id>make-assembly</id>
                        <phase>package</phase>
                        <goals>
                            <goal>single</goal>
                        </goals>
                    </execution>
                </executions>
            </plugin>
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-jar-plugin</artifactId>
                <version>2.4</version>
                <configuration>
                    <archive>
                        <manifest>
                            <mainClass>{{ .PackageName }}.Application</mainClass>
                            <classpathPrefix>lib/</classpathPrefix>
                            <addClasspath>true</addClasspath>
                        </manifest>
                        <manifestEntries>
                            <Class-Path>resources/</Class-Path>
                        </manifestEntries>
                    </archive>
                </configuration>
            </plugin>
        </plugins>
    </build>

</project>`

// TplPackageXML package.xml
const TplPackageXML = `<assembly xmlns="http://maven.apache.org/plugins/maven-assembly-plugin/assembly/1.1.2"
          xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
          xsi:schemaLocation="http://maven.apache.org/plugins/maven-assembly-plugin/assembly/1.1.2 http://maven.apache.org/xsd/assembly-1.1.2.xsd">

    <id>customAssembly</id>
    <!-- dir -->
    <formats>
        <format>dir</format>
    </formats>

    <includeBaseDirectory>false</includeBaseDirectory>

    <fileSets>
        <fileSet>
            <directory>src/main/resources/</directory>
            <outputDirectory>/resources</outputDirectory>
        </fileSet>
    </fileSets>

    <dependencySets>
        <dependencySet>
            <outputDirectory>/lib</outputDirectory>
            <scope>runtime</scope>
            <excludes>
                <exclude>${project.groupId}:${project.artifactId}</exclude>
            </excludes>
        </dependencySet>
        <dependencySet>
            <outputDirectory>/</outputDirectory>
            <includes>
                <include>${project.groupId}:${project.artifactId}</include>
            </includes>
        </dependencySet>
    </dependencySets>

</assembly>`

// TplAppProperties app.properties
const TplAppProperties = `app.version=0.0.1

# log
com.blade.logger.logFile=./logs/app.log`

// TplApplication main java file
const TplApplication = `package {{ .PackageName }};

import com.blade.Blade;

public class Application {

    public static void main(String[] args) {
        Blade.me().start(Application.class, args);
    }

}`

// TplBootstrap bootstrap config
const TplBootstrap = `package {{ .PackageName }}.config;

import com.blade.Blade;
import com.blade.event.BeanProcessor;
import com.blade.ioc.annotation.Bean;{{if ne .TplDependency ""}}
import com.blade.mvc.view.template.JetbrickTemplateEngine;{{end}}

@Bean
public class Bootstrap implements BeanProcessor {
    
    @Override
    public void processor(Blade blade) {
        {{if ne .TplDependency ""}}
        blade.templateEngine(new JetbrickTemplateEngine());{{end}}
    }

}`

// TplController controller file
const TplController = `
package {{ .PackageName }}.controller;

import com.blade.mvc.annotation.GetRoute;
import com.blade.mvc.annotation.Path;
import com.blade.mvc.http.Request;
import com.blade.mvc.http.Response;

@Path
public class IndexController {

	@GetRoute("/")
	public String index(){
		return "index.html";
	}

}
`

// TplIndexHTML index.html
const TplIndexHTML = `<!DOCTYPE html>
	<html>
		<head>
			<title>Hello Blade</title>
			<meta http-equiv="Content-Type" content="text/html;charset=utf-8">
		</head>
	<body>
		<center>
			<h1>Hello Boy!</h1>
		</center>
	</body>
	</html>`

// TplGitignore gitignore
const TplGitignore = `### JetBrains template
# Covers JetBrains IDEs: IntelliJ, RubyMine, PhpStorm, AppCode, PyCharm, CLion, Android Studio and WebStorm
# Reference: https://intellij-support.jetbrains.com/hc/en-us/articles/206544839

# User-specific stuff:
.idea/**/workspace.xml
.idea/**/tasks.xml
.idea/dictionaries

# Sensitive or high-churn files:
.idea/**/dataSources/
.idea/**/dataSources.ids
.idea/**/dataSources.local.xml
.idea/**/sqlDataSources.xml
.idea/**/dynamic.xml
.idea/**/uiDesigner.xml

# Gradle:
.idea/**/gradle.xml
.idea/**/libraries

# CMake
cmake-build-debug/
cmake-build-release/

# Mongo Explorer plugin:
.idea/**/mongoSettings.xml

## File-based project format:
*.iws

## Plugin-specific files:

# IntelliJ
out/

# mpeltonen/sbt-idea plugin
.idea_modules/

# JIRA plugin
atlassian-ide-plugin.xml

# Cursive Clojure plugin
.idea/replstate.xml

# Crashlytics plugin (for Android Studio and IntelliJ)
com_crashlytics_export_strings.xml
crashlytics.properties
crashlytics-build.properties
fabric.properties
### Maven template
target/
pom.xml.tag
pom.xml.releaseBackup
pom.xml.versionsBackup
pom.xml.next
release.properties
dependency-reduced-pom.xml
buildNumber.properties
.mvn/timing.properties

# Avoid ignoring Maven wrapper jar file (.jar files are usually ignored)
!/.mvn/wrapper/maven-wrapper.jar
### Java template
# Compiled class file
*.class

# Log file
*.log
logs/

# Package Files #
*.jar
*.war
*.ear
*.zip
*.tar.gz
*.rar

# virtual machine crash logs, see http://www.java.com/en/download/help/error_hotspot.xml
hs_err_pid*
### Eclipse template

.metadata
bin/
tmp/
*.tmp
*.bak
*.swp
*~.nib
local.properties
.settings/
.loadpath
.recommenders

# External tool builders
.externalToolBuilders/

# Locally stored "Eclipse launch configurations"
*.launch

# sbteclipse plugin
.target

# Build results
[Dd]ebug/
[Dd]ebugPublic/
[Rr]elease/
[Rr]eleases/
x64/
x86/
bld/
[Bb]in/
[Oo]bj/
[Ll]og/

# Visual Studio code coverage results
*.coverage
*.coveragexml

# Web workbench (sass)
.sass-cache/

# Paket dependency manager
.paket/paket.exe
paket-files/

# FAKE - F# Make
.fake/

# JetBrains Rider
.idea/
.blade

### VisualStudioCode template
.vscode/*
!.vscode/settings.json
!.vscode/tasks.json
!.vscode/launch.json
!.vscode/extensions.json
`

// TplGradleBuild build.gradle
const TplGradleBuild = `plugins {
    id 'java'
    id 'idea'
    id 'eclipse'
    id 'maven'
    id 'application'
}

idea {
    module {
        outputDir file('build/classes/main')
        testOutputDir file('build/classes/test')
    }
}
if (project.convention.findPlugin(JavaPluginConvention)) {
    // Change the output directory for the main and test source sets back to the old path
    sourceSets.main.output.resourcesDir = new File(buildDir, "classes/main")
    sourceSets.main.java.outputDir = new File(buildDir, "classes/main")
    sourceSets.test.output.resourcesDir = new File(buildDir, "classes/test")
    sourceSets.test.java.outputDir = new File(buildDir, "classes/test")
}

group '{{ .PackageName }}'
version '{{ .Version }}'

mainClassName = '{{ .PackageName }}.Application'
def libPath = 'build/libs/lib'

sourceCompatibility = 1.8
targetCompatibility = 1.8

tasks.withType(JavaCompile){
    options.encoding = "UTF-8"
}

repositories {
    mavenLocal()
    mavenCentral()
}

dependencies {
    compile 'com.bladejava:blade-mvc:{{ .BladeVersion }}'
    {{ .TplDependency }}
    testCompile 'junit:junit:4.12'
}

jar {
    manifest {
        attributes 'Implementation-Title': 'Gradle'
        attributes 'Main-Class': mainClassName
        attributes 'Built-By': 'biezhi'
        attributes 'Class-Path': 'resources/ ' + configurations.compile.collect { 'lib/' + it.getName() }.join(' ')
    }
}

task wrapper(type: Wrapper) {
    gradleVersion = '4.7'
}

task clearJar(type: Delete) {
    delete libPath
}

task copyJar(type: Copy) {
    from configurations.runtime
    into(libPath)
}

task copyResources(type: Copy) {
    from 'src/main/resources'
    into('build/libs/resources')
}

task release(type: Copy, dependsOn: [build, clearJar, copyJar, copyResources])
`

// TplGradleSetting setting.gradle
const TplGradleSetting = `rootProject.name = '{{ .Name }}'`
