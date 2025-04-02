package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"link-saver-bot/lib/e"
	"link-saver-bot/storage"
	"log/slog"
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
	slog.Info("Initializing file storage", "path", basePath)
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	// Create user directory
	dirPath := filepath.Join(s.basePath, page.UserName)
	if err := os.MkdirAll(dirPath, defaultPerm); err != nil {
		slog.Error("Failed to create directory", "path", dirPath, "error", err)
		return err
	}

	// Get hashed filename
	fileName, err := fileName(page)
	if err != nil {
		slog.Error("Failed to generate filename", "error", err)
		return err
	}

	// Create the full file path
	filePath := filepath.Join(dirPath, fileName)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		slog.Error("Failed to create file", "path", filePath, "error", err)
		return err
	}
	defer file.Close()

	// Encode and save the page
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		slog.Error("Failed to encode page", "path", filePath, "error", err)
		return err
	}

	slog.Info("Page saved successfully", "user", page.UserName, "file", fileName)
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		slog.Error("Failed to read directory", "path", path, "error", err)
		return nil, err
	}

	if len(files) == 0 {
		slog.Info("No saved pages found", "user", userName)
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano()) //@todo: fix deprecation

	n := rand.Intn(len(files))
	file := files[n]

	slog.Info("Random page selected", "user", userName, "file", file.Name())
	return s.DecodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(page *storage.Page) error {
	fileName, err := fileName(page)
	if err != nil {
		slog.Error("Failed to get filename for removal", "error", err)
		return e.Wrap("can't remove file", err)
	}

	filePath := filepath.Join(s.basePath, page.UserName, fileName)

	if err := os.Remove(filePath); err != nil {
		slog.Error("Failed to remove file", "path", filePath, "error", err)
		return e.Wrap(fmt.Sprintf("can't remove file %s", filePath), err)
	}

	slog.Info("Page removed successfully", "user", page.UserName, "file", fileName)
	return nil
}

func (s Storage) Exists(page *storage.Page) (bool, error) {
	fileName, err := fileName(page)
	if err != nil {
		slog.Error("Failed to get filename for existence check", "error", err)
		return false, e.Wrap("can't check if file exists", err)
	}

	filePath := filepath.Join(s.basePath, page.UserName, fileName)

	switch _, err = os.Stat(filePath); {
	case errors.Is(err, os.ErrNotExist):
		slog.Error("Page does not exist", "user", page.UserName, "file", fileName)
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", filePath)
		slog.Error("Failed to check file existence", "path", filePath, "error", err)
		return false, e.Wrap(msg, err)
	}

	slog.Info("Page exists", "user", page.UserName, "file", fileName)
	return true, nil
}

func (s Storage) DecodePage(filepath string) (*storage.Page, error) {
	file, err := os.Open(filepath)
	if err != nil {
		slog.Error("Failed to open file for decoding", "path", filepath, "error", err)
		return nil, e.Wrap("can't decode page", err)
	}
	defer file.Close()

	var page storage.Page

	if err := gob.NewDecoder(file).Decode(&page); err != nil {
		slog.Error("Failed to decode page", "path", filepath, "error", err)
		return nil, err
	}

	slog.Info("Page decoded successfully", "user", page.UserName)
	return &page, nil
}

func fileName(page *storage.Page) (string, error) {
	return page.Hash()
}
