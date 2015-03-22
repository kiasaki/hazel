package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/kiasaki/hazel/api/client"
)

const APIURLEnvName = "HAZEL_API_URL"

func ApiUrlFlag(f *flag.FlagSet) *string {
	defaultApiUrl := os.Getenv(APIURLEnvName)
	if defaultApiUrl == "" {
		defaultApiUrl = "http://localhost:6201"
	}
	return f.String("api-url", defaultApiUrl, "Hazel API Url")
}

func NewApiClient(url string) (client.Client, error) {
	contents, err := ioutil.ReadFile(tokenFilename)
	return client.NewClient(url, string(contents)), err
}
