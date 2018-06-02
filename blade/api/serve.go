package api

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/biezhi/blade-cli/blade/api/templates"
	"github.com/biezhi/blade-cli/blade/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/mkideal/cli"
)

// Serve run blade application
func Serve() *cli.Command {
	return &cli.Command{
		Name:        "serve",
		Desc:        "start blade application",
		CanSubRoute: true,
		Fn: func(ctx *cli.Context) error {
			e, err := NewEngine()
			if err != nil {
				return err
			}
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigs
				e.Stop()
			}()
			e.Run()
			return nil
		},
	}
}

// Engine blade serve engine
type Engine struct {
	config  *templates.BladeConf
	watcher *fsnotify.Watcher

	eventCh        chan string
	watcherStopCh  chan bool
	buildRunCh     chan bool
	buildRunStopCh chan bool
	binStopCh      chan bool
	exitCh         chan bool

	mu         sync.RWMutex
	binRunning bool
	watchers   uint

	ll sync.Mutex // lock for logger
}

// NewEngine ...
func NewEngine() (*Engine, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	var conf templates.BladeConf
	bytes := utils.ReadFileAsByte(".blade")
	if err := json.Unmarshal(bytes, &conf); err != nil {
		return nil, err
	}

	return &Engine{
		config:         &conf,
		watcher:        watcher,
		eventCh:        make(chan string, 1000),
		watcherStopCh:  make(chan bool, 10),
		buildRunCh:     make(chan bool, 1),
		buildRunStopCh: make(chan bool, 1),
		binStopCh:      make(chan bool),
		exitCh:         make(chan bool),
		binRunning:     false,
		watchers:       0,
	}, nil
}

// Run run run
func (e *Engine) Run() {
	var err error
	if err = e.watching("src/main"); err != nil {
		os.Exit(1)
	}
	e.start()
	e.cleanup()
}

func (e *Engine) watching(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			err = e.watcher.Add(path)
			if err != nil {
				return err
			}
			log.Printf("watching %s", path)

			go func() {
				e.withLock(func() {
					e.watchers++
				})
				defer func() {
					e.withLock(func() {
						e.watchers--
					})
				}()

				for {
					select {
					case <-e.watcherStopCh:
						return
					case ev := <-e.watcher.Events:
						log.Printf("event: %+v", ev)
						if !e.isIncludeExt(ev.Name) {
							break
						}
						log.Printf("%s has changed", ev.Name)
						e.eventCh <- ev.Name
					case err := <-e.watcher.Errors:
						log.Printf("error: %s", err.Error())
					}
				}
			}()
		}
		return nil
	})
}

// Endless loop and never return
func (e *Engine) start() {
	firstRunCh := make(chan bool, 1)
	firstRunCh <- true

	for {
		var filename string

		select {
		case <-e.exitCh:
			return
		case filename = <-e.eventCh:
			// time.Sleep(e.config.d())
			e.flushEvents()
			if !e.isIncludeExt(filename) {
				continue
			}
			log.Printf("%s has changed", filename)
		case <-firstRunCh:
			// go down
			break
		}

		select {
		case <-e.buildRunCh:
			e.buildRunStopCh <- true
		default:
		}
		e.withLock(func() {
			if e.binRunning {
				e.binStopCh <- true
			}
		})
		go e.buildRun()
	}
}

func (e *Engine) buildRun() {
	e.buildRunCh <- true

	select {
	case <-e.buildRunStopCh:
		return
	default:
	}

	err := e.runBin()
	if err != nil {
		log.Printf("failed to run, error: %s", err.Error())
	}

	<-e.buildRunCh
}

func (e *Engine) flushEvents() {
	for {
		select {
		case <-e.eventCh:
			log.Printf("flushing events")
		default:
			return
		}
	}
}

func (e *Engine) runBin() error {
	var err error
	log.Println("running...")
	cmd, stdout, stderr, err := utils.StartCmd(`mvn compile exec:java -Dexec.mainClass="` + e.config.MainClass + `"`)
	if err != nil {
		return err
	}
	e.withLock(func() {
		e.binRunning = true
	})

	go io.Copy(os.Stderr, stderr)
	go io.Copy(os.Stdout, stdout)

	go func(cmd *exec.Cmd) {
		<-e.binStopCh
		log.Printf("trying to kill cmd %+v", cmd.Args)
		pid, err := utils.KillCmd(cmd)
		if err != nil {
			log.Printf("failed to kill PID %d, error: %s", pid, err.Error())
			if cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
				os.Exit(1)
			}
		} else {
			log.Printf("cmd killed, pid: %d", pid)
		}
		e.withLock(func() {
			e.binRunning = false
		})
	}(cmd)
	return nil
}

func (e *Engine) cleanup() {
	log.Println("cleaning...")
	defer log.Println("see you again~")

	e.withLock(func() {
		if e.binRunning {
			e.binStopCh <- true
		}
	})

	e.withLock(func() {
		for i := 0; i < int(e.watchers); i++ {
			e.watcherStopCh <- true
		}
	})

	var err error
	if err = e.watcher.Close(); err != nil {
		log.Printf("failed to close watcher, error: %s", err.Error())
	}
}

// Stop the air
func (e *Engine) Stop() {
	e.exitCh <- true
}

func (e *Engine) withLock(f func()) {
	e.mu.Lock()
	f()
	e.mu.Unlock()
}

func (e *Engine) isIncludeExt(path string) bool {
	ext := filepath.Ext(path)
	for _, v := range e.config.IncludeExt {
		if ext == "."+strings.TrimSpace(v) {
			return true
		}
	}
	return false
}

func cmdPath(path string) string {
	return strings.Split(path, " ")[0]
}
