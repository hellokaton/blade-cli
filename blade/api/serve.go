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
			s, err := NewServe()
			if err != nil {
				return err
			}
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigs
				s.Stop()
			}()
			s.Run()
			return nil
		},
	}
}

// BladeServe blade serve engine
type BladeServe struct {
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

// NewServe ...
func NewServe() (*BladeServe, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	var conf templates.BladeConf
	bytes := utils.ReadFileAsByte(".blade")
	if err := json.Unmarshal(bytes, &conf); err != nil {
		return nil, err
	}

	return &BladeServe{
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

// Run run
func (s *BladeServe) Run() {
	var err error
	if err = s.watching("src/main"); err != nil {
		os.Exit(1)
	}
	s.start()
	s.cleanup()
}

func (s *BladeServe) watching(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			err = s.watcher.Add(path)
			if err != nil {
				return err
			}
			// log.Printf("watching %s", path)

			go func() {
				s.withLock(func() {
					s.watchers++
				})
				defer func() {
					s.withLock(func() {
						s.watchers--
					})
				}()

				for {
					select {
					case <-s.watcherStopCh:
						return
					case ev := <-s.watcher.Events:
						if !s.isIncludeExt(ev.Name) {
							break
						}
						log.Printf("%s has changed", ev.Name)
						s.eventCh <- ev.Name
					case err := <-s.watcher.Errors:
						log.Printf("error: %s", err.Error())
					}
				}
			}()
		}
		return nil
	})
}

// Endless loop and never return
func (s *BladeServe) start() {
	firstRunCh := make(chan bool, 1)
	firstRunCh <- true

	for {
		var filename string

		select {
		case <-s.exitCh:
			return
		case filename = <-s.eventCh:
			// time.Sleep(e.config.d())
			s.flushEvents()
			if !s.isIncludeExt(filename) {
				continue
			}
			log.Printf("%s has changed", filename)
		case <-firstRunCh:
			// go down
			break
		}

		select {
		case <-s.buildRunCh:
			s.buildRunStopCh <- true
		default:
		}
		s.withLock(func() {
			if s.binRunning {
				s.binStopCh <- true
			}
		})
		go s.buildRun()
	}
}

func (s *BladeServe) buildRun() {
	s.buildRunCh <- true

	select {
	case <-s.buildRunStopCh:
		return
	default:
	}

	err := s.runBin()
	if err != nil {
		log.Printf("failed to run, error: %s", err.Error())
	}

	<-s.buildRunCh
}

func (s *BladeServe) flushEvents() {
	for {
		select {
		case <-s.eventCh:
			log.Printf("flushing events")
		default:
			return
		}
	}
}

func (s *BladeServe) runBin() error {
	var err error
	log.Println("running...")

	shell := `mvn compile exec:java -Dexec.mainClass="` + s.config.MainClass + `"`
	if s.config.BuildTool == "gradle" {
		shell = "gradle -q run"
	}

	cmd, stdout, stderr, err := utils.StartCmd(shell)
	if err != nil {
		return err
	}
	s.withLock(func() {
		s.binRunning = true
	})

	go io.Copy(os.Stderr, stderr)
	go io.Copy(os.Stdout, stdout)

	go func(cmd *exec.Cmd) {
		<-s.binStopCh
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
		s.withLock(func() {
			s.binRunning = false
		})
	}(cmd)
	return nil
}

func (s *BladeServe) cleanup() {
	log.Println("cleaning...")
	defer log.Println("see you again~")

	s.withLock(func() {
		if s.binRunning {
			s.binStopCh <- true
		}
	})

	s.withLock(func() {
		for i := 0; i < int(s.watchers); i++ {
			s.watcherStopCh <- true
		}
	})

	var err error
	if err = s.watcher.Close(); err != nil {
		log.Printf("failed to close watcher, error: %s", err.Error())
	}
}

// Stop the air
func (s *BladeServe) Stop() {
	s.exitCh <- true
}

func (s *BladeServe) withLock(f func()) {
	s.mu.Lock()
	f()
	s.mu.Unlock()
}

func (s *BladeServe) isIncludeExt(path string) bool {
	ext := filepath.Ext(path)
	for _, v := range s.config.IncludeExt {
		if ext == "."+strings.TrimSpace(v) {
			return true
		}
	}
	return false
}

func cmdPath(path string) string {
	return strings.Split(path, " ")[0]
}
