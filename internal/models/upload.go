package models

import "mime/multipart"

type File struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Photo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type WriteFile struct {
	UserID int                   `json:"-"`
	Name   string                `json:"name"`
	File   *multipart.FileHeader `json:"-"`
}

type WritePhoto struct {
	UserID int    `json:"-"`
	Name   string `json:"name"`
	File   *multipart.FileHeader `json:"-"`
}
