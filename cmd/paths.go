package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func has(s []string, val string) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}

	return false
}

func getEnvFilePaths(options *Options) ([]string, error) {
	paths := make([]string, 0, 5)
	paths = append(paths, fmt.Sprintf("%s/%s", options.Dir, options.Source))

	for _, target := range options.Target {
		paths = append(paths, fmt.Sprintf("%s/%s", options.Dir, target))
	}

	if !options.All {
		if options.Env != "" {
			path := fmt.Sprintf("%s/.env.%s", options.Dir, options.Env)
			if fileExists(path) {
				paths = append(paths, path)
			}
		}

		if options.Local {
			if options.Env != "" {
				path := fmt.Sprintf("%s/.env.%s.local", options.Dir, options.Env)
				if fileExists(path) {
					paths = append(paths, path)
				}
			}

			path := fmt.Sprintf("%s/.env.local", options.Dir)
			if fileExists(path) {
				paths = append(paths, path)
			}
		}
	} else {
		files, err := os.ReadDir(options.Dir)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			qualifiedPath := fmt.Sprintf("%s/%s", options.Dir, file.Name())
			if strings.HasPrefix(file.Name(), ".env") && !has(paths, qualifiedPath) {
				paths = append(paths, qualifiedPath)
			}
		}
	}

	return paths, nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func validateFilePaths(paths []string) error {
	for _, path := range paths {
		if !fileExists(path) {
			return fmt.Errorf("file \"%s\" does not exist", path)
		}
	}

	return nil
}
