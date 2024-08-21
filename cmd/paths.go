package cmd

import (
	"errors"
	"fmt"
	"os"
)

func getEnvFilePaths(options *Options) []string {
	envFilePaths := make([]string, 0, 5)
	envFilePaths = append(envFilePaths, fmt.Sprintf("%s/%s", options.Dir, options.Source))

	for _, target := range options.Target {
		envFilePaths = append(envFilePaths, fmt.Sprintf("%s/%s", options.Dir, target))
	}

	if options.Env != "" {
		path := fmt.Sprintf("%s/.env.%s", options.Dir, options.Env)
		if fileExists(path) {
			envFilePaths = append(envFilePaths, path)
		}
	}

	if options.Local {
		if options.Env != "" {
			path := fmt.Sprintf("%s/.env.%s.local", options.Dir, options.Env)
			if fileExists(path) {
				envFilePaths = append(envFilePaths, path)
			}
		}

		path := fmt.Sprintf("%s/.env.local", options.Dir)
		if fileExists(path) {
			envFilePaths = append(envFilePaths, path)
		}
	}

	return envFilePaths
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
