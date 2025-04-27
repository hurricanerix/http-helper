package main

import (
	"flag"
	"fmt"
	"github.com/hurricanerix/http-helper/build"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/hurricanerix/http-helper/config"
	"github.com/hurricanerix/http-helper/handler"
	"github.com/hurricanerix/http-helper/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func init() {
	godotenv.Load()
}

const defaultServerIdleTimeout = 5 * time.Second
const defaultServerReadTimeout = 5 * time.Second
const defaultServerWriteTimeout = 5 * time.Second

const defaultCORSAllowedOrigins = "*"
const defaultCORSAllowedMethods = "HEAD,GET"
const defaultCORSAllowCredentials = true

var base64SourceDiff string

func main() {
	address := flag.String("bind", "127.0.0.1", "bind to this address")
	port := flag.Int("port", 8000, "bind to this port")
	directory := flag.String("d", ".", "serve this directory")
	showDiff := flag.Bool("diff", false, "display the changes made at compile time, suitable for patching.")

	flag.Usage = func() {
		fmt.Printf("Usage: %s [FLAGS]\n", os.Args[0])
		fmt.Println("")
		fmt.Println("Build Info:")
		fmt.Println("  Built with:", build.GoVersion())
		fmt.Println("  Commit Hash:", build.CommitHash())
		fmt.Println("  Commit Date:", build.CommitDate())
		if build.SourceModified() {
			fmt.Printf(" (modified)")
		}
		fmt.Println("")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showDiff {
		fmt.Println(build.SourceDiff())
		return
	}

	directoryAbsolutePath, err := filepath.Abs(*directory)
	if err != nil {
		panic(err)
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   config.StringSliceEnv("HH_CORS_ALLOWED_ORIGINS", defaultCORSAllowedOrigins),
		AllowedMethods:   config.StringSliceEnv("HH_CORS_ALLOWED_METHOD", defaultCORSAllowedMethods),
		AllowCredentials: config.BoolEnv("HH_CORS_ALLOW_CREDENTIALS", defaultCORSAllowCredentials),
	})

	p := pipeline{
		middleware.Logger,
		middleware.Error,
		middleware.RequestID,
		middleware.Bandwidth,
		middleware.TTFB,
		corsMiddleware.Handler,
		middleware.Mime,
		middleware.ETag,
	}

	s := &http.Server{
		Addr: fmt.Sprintf("%s:%d", *address, *port),
		Handler: wrap(handler.File{
			Directory: directoryAbsolutePath,
		}, p),
		IdleTimeout:  config.DurationEnv("HH_SERVER_IDLE_TIMEOUT", defaultServerIdleTimeout),
		ReadTimeout:  config.DurationEnv("HH_SERVER_READ_TIMEOUT", defaultServerReadTimeout),
		WriteTimeout: config.DurationEnv("HH_SERVER_WRITE_TIMEOUT", defaultServerWriteTimeout),
	}

	fmt.Printf("Serving HTTP on %s port %d (http://%s/)\n", *address, *port, s.Addr)
	log.Fatal(s.ListenAndServe())
}

type pipeline []func(h http.Handler) http.Handler

func wrap(h http.Handler, p pipeline) http.Handler {
	for i := len(p) - 1; i != -1; i-- {
		h = p[i](h)
	}
	return h
}
