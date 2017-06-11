package cmder

import (
	"os"
	"path/filepath"
)

//Program returns program name.
func Program() string {
	return os.Args[0]
}

//MoveTo moves program to given dest.
func MoveTo(dir string) error {
	cur, err := os.Getwd()
	if err != nil {
		return err
	}

	oldPath := filepath.Join(cur, Program())
	newPath := filepath.Join(dir, Program())

	err = os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}

	return nil
}
