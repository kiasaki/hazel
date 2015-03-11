package main

import (
	"encoding/json"
	"fmt"
)

type App struct {
	Slug  string `json:"slug"`
	Stack string `json:"stack"`
}

func (a *App) Get() error {
	bucket := awsS3.Bucket(cfg.Bucket)
	path := fmt.Sprintf("/apps/%s/config.json", a.Slug)

	bytes, err := bucket.Get(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, a)
	if err != nil {
		return err
	}

	return nil
}
