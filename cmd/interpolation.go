package cmd

import (
	"regexp"
)

func interpolate(envs []map[string]*EnvVar, keys []string, interpolated map[string]bool, maxIters int) error {
	re := regexp.MustCompile(`\${(.*?)}`)

	if interpolated == nil {
		interpolated = make(map[string]bool)
	}

	keysIndex := make(map[string]int)
	for _, k := range keys {
		keysIndex[k] = 0
	}

	interpolates := make(map[string]*EnvVar)

	for _, k := range keys {
		for _, env := range envs {
			envVar, ok := env[k]

			if ok && envVar != nil && re.MatchString(envVar.Value) {
				interpolates[k] = envVar
			}
		}
	}

	for _, envVar := range interpolates {
		value := re.ReplaceAllStringFunc(envVar.Value, func(s string) string {
			submatches := re.FindAllStringSubmatch(s, -1)
			if len(submatches) > 0 {
				for _, v := range submatches {
					if len(v) > 1 {
						interpolatedValue := findValue(envs, v[1])

						if !re.MatchString(interpolatedValue) {
							interpolated[v[1]] = true
							// interpolatedValue = fmt.Sprintf("<fg=cyan>%s</>", interpolatedValue)
						}

						return interpolatedValue
					}
				}
			}

			return s
		})

		envVar.Value = value
	}

	for k, v := range interpolated {
		if v {
			delete(keysIndex, k)
			delete(interpolated, k)
		}
	}

	if maxIters > 0 && len(keysIndex) > 0 {
		maxIters--
		return interpolate(envs, getKeys(keysIndex), interpolated, maxIters)
	}

	return nil
}

func findValue(envs []map[string]*EnvVar, key string) string {
	for i := len(envs) - 1; i >= 0; i-- {
		envVar, ok := envs[i][key]
		if ok && envVar.Value != "" {
			return envVar.Value
		}
	}

	return ""
}
