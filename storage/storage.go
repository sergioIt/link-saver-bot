package storage

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"link-saver-bot/lib/e"
)

type Storage interface {
	Save(page *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(page *Page) error
	Exists(page *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

var ErrNoSavedPages = errors.New("no saved pages exists")

func (p *Page) Hash() (string, error) {

	h := sha256.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
