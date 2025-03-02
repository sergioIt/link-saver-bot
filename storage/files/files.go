package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"link-saver-bot/lib/e"
	"link-saver-bot/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0777

func New(basePath string) Storage {

	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {

	defer func() { err = e.Wrap("can't save page", err) }()

	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}

	fileName, err := fileName(page)

	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, fileName)

	file, err := os.Create(filepath.Join(filePath, page.URL))
	if err != nil {
		return err
	}
	defer file.Close()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {

	defer func() { err = e.Wrap("can't pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano()) //@todo: fix deprecation

	n := rand.Intn(len(files))

	file := files[n]

	// open selected file and decode it
	return s.DecodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(page *storage.Page) error {

	fileName, err := fileName(page)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	filePath := filepath.Join(s.basePath, page.UserName, fileName)

	if err := os.Remove(filePath); err != nil {
		return e.Wrap(fmt.Sprintf("can't remove file %s", filePath), err)
	}

	return nil
}

func (s Storage) Exists(page *storage.Page) (bool, error) {

	fileName, err := fileName(page)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}

	filePath := filepath.Join(s.basePath, page.UserName, fileName)

	// _, err = os.Stat(filePath)

	switch _, err = os.Stat(filePath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:

		msg := fmt.Sprintf("can't check if file %s exists", filePath)

		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) DecodePage(filepath string) (*storage.Page, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer file.Close()

	var page storage.Page

	if err := gob.NewDecoder(file).Decode(&page); err != nil {
		return nil, err
	}

	return &page, nil
}

func fileName(page *storage.Page) (string, error) {
	return page.Hash()
}
