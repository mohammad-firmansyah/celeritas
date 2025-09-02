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

package yourpkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Use a consistent, correct receiver name.
type Celeritas struct{}

// EnsureFile guarantees that `path` is a regular file.
// - Creates the parent directory if needed.
// - If `path` exists as a directory or symlink, it is handled safely:
//   * a directory is renamed to a timestamped backup before creating the file
//   * a symlink is rejected (safer) unless it points to a regular file
func (c *Celertias) EnsureFile(path string) error {
	// 1) Ensure parent directory exists
	parent := filepath.Dir(path)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return fmt.Errorf("ensure parent dir: %w", err)
	}

	// 2) Inspect existing path using Lstat to see symlinks/directories as-is
	info, err := os.Lstat(path)
	switch {
	case err == nil:
		switch mode := info.Mode(); {
		case mode.IsRegular():
			// Already a regular file — OK
			return nil
		case mode.IsDir():
			// Path is a directory — rename it away, then create a file.
			backup := path + ".backup-" + time.Now().Format("20060102-150405")
			if rerr := os.Rename(path, backup); rerr != nil {
				return fmt.Errorf("'.env' is a directory; failed to rename to %q: %w", backup, rerr)
			}
			// Now fall through to create the file fresh below.
		case mode&os.ModeSymlink != 0:
			// Resolve symlink and only accept if it points to a regular file.
			target, rerr := os.Readlink(path)
			if rerr != nil {
				return fmt.Errorf("readlink %s: %w", path, rerr)
			}
			resolved := target
			if !filepath.IsAbs(resolved) {
				resolved = filepath.Join(parent, target)
			}
			tinfo, terr := os.Stat(resolved)
			if terr != nil || !tinfo.Mode().IsRegular() {
				return fmt.Errorf("%s is a symlink to a non-regular file (%q); please fix", path, resolved)
			}
			// Points to a regular file — accept.
			return nil
		default:
			return fmt.Errorf("%s exists but is not a regular file; mode=%v", path, mode)
		}
	case errors.Is(err, os.ErrNotExist):
		// Does not exist — proceed to create below.
	default:
		return fmt.Errorf("stat %s: %w", path, err)
	}

	// 3) Create the file atomically; if another process raced us and created it,
	//    treat that as success if it is indeed a regular file.
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			// Re-check: if it is a regular file now, accept; otherwise, error.
			if info2, e2 := os.Lstat(path); e2 == nil && info2.Mode().IsRegular() {
				return nil
			}
		}
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()

	// Optionally, write a newline to ensure the file is not zero-length.
	// _, _ = f.WriteString("\n")

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
