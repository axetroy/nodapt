package util

import "os"

func EnsureDir(dirPath string) error {
	// Check if the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Directory does not exist, create it
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	} else if err != nil {
		// Return any other error encountered
		return err
	}
	return nil
}
