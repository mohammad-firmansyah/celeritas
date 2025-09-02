package celeritas

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

const version = "1.0.0"

type Celertias struct {
	AppName  string
	Debug    bool
	Version  string
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (c *Celertias) New(rootPath string) error {

	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "logs", "middleware", "public", "tmp"},
	}

	err := c.Init(pathConfig)
	if err != nil {
		return err
	}

	err = c.CheckDotEnv(rootPath)
	if err != nil {
		return err
	}

	// read .env
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// create logger
	infoLog, errLog := c.startLoggers()
	c.InfoLog = infoLog
	c.ErrorLog = errLog
	c.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	c.Version = version

	return nil
}

func (c *Celertias) Init(p initPaths) error {
	root := p.rootPath

	for _, path := range p.folderNames {
		// create folder if it doesn't exist
		err := c.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *Celertias) EnsureFile(path string) error {
	// 1) Ensure parent directory exists
	parent := filepath.Dir(path)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return fmt.Errorf("ensure parent dir: %w", err)
	}

	// 2) Inspect existing path
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return fmt.Errorf("path exists but is a directory: %s", path)
		}
		// It is already a file; nothing to do.
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	// 3) Create the file
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()
	return nil
}

func (c *Celertias) CheckDotEnv(root string) error {
	envPath := filepath.Join(root, ".env")
	return c.EnsureFile(envPath)
}
func (c *Celertias) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog

}
