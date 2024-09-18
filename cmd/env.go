package cmd

import (
	"bufio"
	"os"
	"strings"
)

type EnvVar struct {
	Key       string
	Value     string
	Duplicate bool
}

func readEnvFiles(paths []string) ([]map[string]*EnvVar, error) {
	envs := make([]map[string]*EnvVar, 0, len(paths))

	for _, path := range paths {
		env, err := readEnvFile(path)
		if err != nil {
			return nil, err
		}
		envs = append(envs, env)
	}

	return envs, nil
}

func readEnvFile(path string) (map[string]*EnvVar, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	env := make(map[string]*EnvVar)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) == 0 {
			continue
		}

		_, duplicate := env[parts[0]]

		envVar := &EnvVar{
			Key:       parts[0],
			Duplicate: duplicate,
		}

		if len(parts) > 1 {
			envVar.Value = trim(parts[1])
		}

		env[envVar.Key] = envVar
	}

	return env, nil
}

func trim(s string) string {
	s = strings.TrimSpace(s)

	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		return s[1 : len(s)-1]
	}

	if strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'") {
		return s[1 : len(s)-1]
	}

	return s
}

func makeEnvVarIndex(envs []map[string]*EnvVar) map[string]int {
	keysIndex := make(map[string]int, len(envs[0]))

	for _, env := range envs {
		for k := range env {
			keysIndex[k] = 0
		}
	}

	return keysIndex
}

func getKeys[T comparable, U any](m map[T]U) []T {
	s := make([]T, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

func addSystemEnvs(envs []map[string]*EnvVar, keysIndex map[string]int) []map[string]*EnvVar {
	sys := make(map[string]*EnvVar, len(keysIndex))
	for k := range keysIndex {
		sys[k] = &EnvVar{
			Key:   k,
			Value: os.Getenv(k),
		}
	}

	newEnvs := make([]map[string]*EnvVar, 0, 1+len(envs))
	newEnvs = append(newEnvs, envs[0])
	newEnvs = append(newEnvs, sys)
	newEnvs = append(newEnvs, envs[1:]...)

	return newEnvs
}
