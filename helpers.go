package celeritas

import "os"

func (c *Celertias) CreateDirIfNotExist(path string) error {
	const mode = 0775
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, mode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Celertias) CreateFileIfNotExist(path string) error {
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}

		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}

	return nil
}
