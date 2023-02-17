package celeritas

import "os"

func (c *Celertias) CreateDirIfNotExist(path string) error {
	const mode = 775
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, mode)
		if err != nil {
			return err
		}
	}

	return nil
}
