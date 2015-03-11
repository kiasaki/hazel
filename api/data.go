package main

import (
	"fmt"
)

type App struct {
	Slug string
}

func (a App) Exists() (bool, error) {
	bucket := awsS3.Bucket(cfg.Bucket)
	path := fmt.Sprintf("/apps/%s/config.json", a.Slug)
	return bucket.Exists(path)
}
