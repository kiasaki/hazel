package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"github.com/jessevdk/go-flags"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
)

var cfg struct {
	Port         string `short:"p" long:"port" env:"PORT" description:"Port to bind to" default:"8000"`
	AwsAccessKey string `long:"aws-access-key" env:"AWS_ACCESS_KEY" description:"Aws access key" required:"true"`
	AwsSecretKey string `long:"aws-secret-key" env:"AWS_SECRET_KEY" description:"Aws secret key" required:"true"`
	Bucket       string `short:"b" long:"bucket" env:"BUCKET" description:"AWS bucket with config" required:"true"`
}

var awsAuth *aws.Auth
var awsS3 *s3.S3

func main() {
	initConfig()
	initAws()
	initServer()
}

func initConfig() {
	_, err := flags.Parse(&cfg)
	if err != nil {
		os.Exit(1)
	}
}

func initAws() {
	awsAuth = aws.NewAuth(cfg.AwsAccessKey, cfg.AwsSecretKey, "", time.Time{})
	awsS3 = s3.New(*awsAuth, aws.USEast)
}

func initServer() {
	registerHandlers()

	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatal(err)
	}

	goji.DefaultMux.Compile()

	http.Handle("/", goji.DefaultMux)

	log.Println("Starting Goji on", listener.Addr())
	graceful.HandleSignals()
	graceful.PreHook(func() { log.Printf("Goji received signal, gracefully stopping") })
	graceful.PostHook(func() { log.Printf("Goji stopped") })

	err = graceful.Serve(listener, http.DefaultServeMux)
	if err != nil {
		log.Fatal(err)
	}

	graceful.Wait()
}
