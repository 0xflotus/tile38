package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"

	"github.com/tidwall/tile38/controller"
	"github.com/tidwall/tile38/controller/log"
	"github.com/tidwall/tile38/core"
)

var (
	dir         string
	port        int
	host        string
	verbose     bool
	veryVerbose bool
	devMode     bool
	quiet       bool
)

func main() {
	flag.IntVar(&port, "p", 9851, "The listening port.")
	flag.StringVar(&host, "h", "127.0.0.1", "The listening host.")
	flag.StringVar(&dir, "d", "data", "The data directory.")
	flag.BoolVar(&verbose, "v", false, "Enable verbose logging.")
	flag.BoolVar(&quiet, "q", false, "Quiet logging. Totally silent.")
	flag.BoolVar(&veryVerbose, "vv", false, "Enable very verbose logging.")
	flag.BoolVar(&devMode, "dev", false, "Activates dev mode. DEV ONLY.")
	flag.Parse()
	var logw io.Writer = os.Stderr
	if quiet {
		logw = ioutil.Discard
	}
	log.Default = log.New(logw, &log.Config{
		HideDebug: !veryVerbose,
		HideWarn:  !(veryVerbose || verbose),
	})
	core.DevMode = devMode
	core.ShowDebugMessages = veryVerbose

	//  _____ _ _     ___ ___
	// |_   _|_| |___|_  | . |
	//   | | | | | -_|_  | . |
	//   |_| |_|_|___|___|___|

	fmt.Fprintf(logw, `
   _______ _______
  |       |       |
  |____   |   _   |   Tile38 %s (%s) %d bit (%s/%s)
  |       |       |   Host: %s, Port: %d, PID: %d
  |____   |   _   |
  |       |       |   tile38.com
  |_______|_______|
`+"\n", core.Version, core.GitSHA, strconv.IntSize, runtime.GOARCH, runtime.GOOS, host, port, os.Getpid())

	if err := controller.ListenAndServe(host, port, dir); err != nil {
		log.Fatal(err)
	}
}
