# Blade Cli

## Install

**For MacOSX**

```bash
brew tap biezhi/tap && brew install blade
```

upgrade 

```bash
brew upgrade blade
```

**For Windows**

Download [blade] binary, copy `blade` to `c:/windows/system32`

upgrade cli

```bash
blade upgrade
```

## Usage

```bash
» blade

    __, _,   _, __, __,
    |_) |   /_\ | \ |_
    |_) | , | | |_/ |
    ~   ~~~ ~ ~ ~   ~~~
    :: Blade Cli :: (v0.0.1)

    Inspired by https://lets-blade.com

Options:

  -h, --help      display help information
  -v, --version   display blade cli version

Commands:

  help    display help information
  new     create blade application by template
  serve   start blade application
  build   build application as jar or dir
```

## Example

**Create Application**

```bash
» blade new hello
? please input package name (e.g: com.bladejava.example): com.bladedemo.hello

? choose a build tool:  [Use arrows to move, type to filter]
❯ Maven
  Gradle

⠳ creating project, please wait...

create file success: hello/pom.xml
create file success: hello/.blade
create file success: hello/package.xml
create file success: hello/src/main/resources/app.properties
create file success: hello/src/main/resources/templates/index.html
create file success: hello/src/main/java/com/bladejava/example/Application.java
create file success: hello/src/main/java/com/bladejava/example/controller/IndexController.java

application hello create successful!

    $ cd hello
    $ blade serve
```

**Run Application**

cd app dir

```bash
» blade serve
[INFO] Scanning for projects...
[INFO]
[INFO] ------------------------------------------------------------------------
[INFO] Building hello 0.0.1
[INFO] ------------------------------------------------------------------------
[INFO]
[INFO] --- maven-resources-plugin:2.6:resources (default-resources) @ hello ---
[INFO] Using 'UTF-8' encoding to copy filtered resources.
[INFO] Copying 2 resources
[INFO]
[INFO] --- maven-compiler-plugin:2.5.1:compile (default-compile) @ hello ---
[INFO] Compiling 2 source files to /Users/biezhi/workspace/golang/src/github.com/biezhi/blade-cli/hello/target/classes
[INFO]
[INFO] --- exec-maven-plugin:1.6.0:java (default-cli) @ hello ---
2018/06/02 18:15:42  INFO [          _(:3」∠)_ ]             c.b.s.n.NettyServer : environment.jdk.version    » 1.8.0_101
2018/06/02 18:15:42  INFO [          _(:3」∠)_ ]             c.b.s.n.NettyServer : environment.user.dir       » /Users/biezhi/workspace/golang/src/github.com/biezhi/blade-cli/hello
2018/06/02 18:15:42  INFO [          _(:3」∠)_ ]             c.b.s.n.NettyServer : environment.java.io.tmpdir » /var/folders/y7/fdpr6jzx1rs6x0jmty2h6lvw0000gn/T/
2018/06/02 18:15:42  INFO [          _(:3」∠)_ ]             c.b.s.n.NettyServer : environment.user.timezone  » Asia/Shanghai
2018/06/02 18:15:42  INFO [          _(:3」∠)_ ]             c.b.s.n.NettyServer : environment.file.encoding  » UTF-8
2018/06/02 18:15:42  INFO [          _(:3」∠)_ ]             c.b.s.n.NettyServer : environment.classpath      » /Users/biezhi/workspace/golang/src/github.com/biezhi/blade-cli/hello/target/classes/

                                         __, _,   _, __, __,
                                         |_) |   /_\ | \ |_
                                         |_) | , | | |_/ |
                                         ~   ~~~ ~ ~ ~   ~~~
                                     :: Blade :: (v2.0.8-BETA3)

2018/06/02 18:15:42  INFO [          _(:3」∠)_ ]            c.b.m.r.RouteMatcher : » Add route  GET     /
2018/06/02 18:15:42  INFO [          _(:3」∠)_ ]             c.b.s.n.NettyServer : » Register bean: [com.blade.Environment@137a568d, com.bladejava.example.controller.IndexController@62c7c070]
2018/06/02 18:15:42  INFO [          _(:3」∠)_ ]             c.b.s.n.NettyServer : » Watched environment: true
2018/06/02 18:15:42  INFO [          _(:3」∠)_ ]             c.b.s.n.NettyServer : » Use NioEventLoopGroup
2018/06/02 18:15:42  INFO [          _(:3」∠)_ ]             c.b.s.n.NettyServer : » Blade initialize successfully, Time elapsed: 275 ms
2018/06/02 18:15:42  INFO [          _(:3」∠)_ ]             c.b.s.n.NettyServer : » Blade start with  0.0.0.0:9000
2018/06/02 18:15:42  INFO [          _(:3」∠)_ ]             c.b.s.n.NettyServer : » Open browser access http://127.0.0.1:9000 ⚡
```

## License

[MIT](LICENSE)