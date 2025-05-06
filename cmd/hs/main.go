package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/hurricanerix/http-helper/build"
	"github.com/hurricanerix/http-helper/config"
	"github.com/hurricanerix/http-helper/middleware"
	"github.com/hurricanerix/http-helper/platforms/python"
	"github.com/hurricanerix/http-helper/platforms/s3"
	"github.com/joho/godotenv"
)

func init() {
	loadEnv()
}

const defaultServerPipeline = "logger, error, request_id, bandwidth, ttfb, cors, mime, etag"
const defaultServerHandler = "python"

const defaultServerIdleTimeout = 5 * time.Second
const defaultServerReadTimeout = 5 * time.Second
const defaultServerWriteTimeout = 5 * time.Second

func main() {
	address := flag.String("bind", "127.0.0.1", "Bind to this address.")
	port := flag.Int("port", 8000, "Bind to this port.")
	directory := flag.String("d", ".", "Serve this directory.")
	showDiff := flag.Bool("diff", false, "Display the changes made at compile time, suitable for patching.")
	showVersion := flag.Bool("version", false, "Display the version and exit.")

	flag.Usage = func() {
		fmt.Printf("Usage: %s [FLAGS]\n", os.Args[0])
		fmt.Println("")
		fmt.Println("Build Info:")
		fmt.Println("  Built with:", build.GoVersion())
		fmt.Printf("  Version: %s", build.CommitHash())
		if build.SourceModified() {
			fmt.Printf(" (modified)")
		}
		fmt.Println("")
		fmt.Println("  Commit Date:", build.CommitDate())

		fmt.Println("")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showDiff {
		fmt.Println(build.SourceDiff())
		return
	}

	if *showVersion {
		fmt.Println(build.CommitHash())
		return
	}

	directoryAbsolutePath, err := filepath.Abs(*directory)
	if err != nil {
		panic(err)
	}

	p := getPipeline(config.StringEnv("HH_SERVER_PIPELINE", defaultServerPipeline))
	h := getHandler(config.StringEnv("HH_SERVER_HANDLER", defaultServerHandler), directoryAbsolutePath)
	s := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", *address, *port),
		Handler:      wrap(h, p),
		IdleTimeout:  config.DurationEnv("HH_SERVER_IDLE_TIMEOUT", defaultServerIdleTimeout),
		ReadTimeout:  config.DurationEnv("HH_SERVER_READ_TIMEOUT", defaultServerReadTimeout),
		WriteTimeout: config.DurationEnv("HH_SERVER_WRITE_TIMEOUT", defaultServerWriteTimeout),
	}

	fmt.Printf("Serving HTTP on %s port %d (http://%s/)\n", *address, *port, s.Addr)
	log.Fatal(s.ListenAndServe())
}

type pipeline []stage
type stage func(h http.Handler) http.Handler

func wrap(h http.Handler, p pipeline) http.Handler {
	for i := len(p) - 1; i != -1; i-- {
		h = p[i](h)
	}
	return h
}

func getHandler(name string, dir string) http.Handler {
	switch name {
	case "s3":
		return s3.Handler{
			Directory: dir,
		}
	case "python":
		fallthrough
	default:
		return python.Handler{
			Directory: dir,
		}
	}
}

func getPipeline(rawPipeline string) pipeline {
	stageNames := strings.Split(strings.ReplaceAll(rawPipeline, " ", ""), ",")

	p := make([]stage, len(stageNames))
	for i := range stageNames {
		p[i] = getStage(stageNames[i])
	}

	return p
}

func getStage(stageName string) stage {
	switch stageName {
	case "logger":
		return middleware.Logger
	case "error":
		return middleware.Error
	case "request_id":
		return middleware.RequestID
	case "bandwidth":
		return middleware.Bandwidth
	case "ttfb":
		return middleware.TTFB
	case "cors":
		return middleware.CORS
	case "mime":
		return middleware.Mime
	case "etag":
		return middleware.ETag
	case "python.logger":
		return python.Logger
	}

	return middleware.NOP
}

func loadEnv() {
	if err := godotenv.Load(".env"); err == nil {
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	globalConfig := path.Join(homeDir, ".config/http_helper.env")
	godotenv.Load(globalConfig)
}
