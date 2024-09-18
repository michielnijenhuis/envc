package cmd

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/michielnijenhuis/cli"
)

var Command = &cli.Command{
	Name:        "envc",
	Description: "Prints all variables in the available .env file(s), sorted alphabetically",
	Arguments: []cli.Arg{
		&cli.StringArg{
			Name:        "dir",
			Description: "The directory to check for .env files",
			Value:       "./",
		},
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "source",
			Description: "The base .env file to compare with",
			Value:       ".env.example",
		},
		&cli.ArrayFlag{
			Name:        "target",
			Description: "The .env file(s) to compare against",
			Value:       []string{".env"},
		},
		&cli.StringFlag{
			Name:        "env",
			Shortcuts:   []string{"e"},
			Description: "The environment .env to include in the comparison (e.g. `.env.dev`)",
		},
		&cli.BoolFlag{
			Name:        "local",
			Shortcuts:   []string{"l"},
			Description: "Whether to include local .env files (e.g. `.env.local` or `.env.dev.local`)",
		},
		&cli.StringFlag{
			Name:        "skip",
			Shortcuts:   []string{"s"},
			Description: "Comma separated list of variable name patterns to skip",
		},
		&cli.StringFlag{
			Name:        "pattern",
			Shortcuts:   []string{"p"},
			Description: "Comma separated list of variable name patterns to focus on",
		},
		&cli.BoolFlag{
			Name:        "interpolate",
			Shortcuts:   []string{"i"},
			Description: "Interpolate env var values that refer to other env vars",
		},
		&cli.BoolFlag{
			Name:        "result",
			Shortcuts:   []string{"r"},
			Description: "Include conjuncted .env file result",
		},
		&cli.StringFlag{
			Name:        "truncate",
			Shortcuts:   []string{"t"},
			Description: "Whether to truncate long values or not",
			Value:       "40",
		},
		&cli.BoolFlag{
			Name:        "system",
			Description: "Include os values",
		},
		&cli.BoolFlag{
			Name:        "all",
			Shortcuts:   []string{"a"},
			Description: "Include all .env files that can be found",
		},
	},
	RunE: handle,
}

type Options struct {
	Dir         string
	Source      string
	Target      []string
	Env         string
	Local       bool
	Truncate    int
	Skip        string
	Pattern     string
	Interpolate bool
	Result      bool
	System      bool
	All         bool
}

func handle(c *cli.Ctx) error {
	options, err := getOptions(c)
	if err != nil {
		return err
	}

	paths, err := getEnvFilePaths(options)
	if err != nil {
		return err
	}

	err = validateFilePaths(paths)
	if err != nil {
		return err
	}

	envs, err := readEnvFiles(paths)
	if err != nil {
		return err
	}

	keysIndex := makeEnvVarIndex(envs)

	if options.System {
		envs = addSystemEnvs(envs, keysIndex)
	}

	if options.Skip != "" {
		keysIndex = applySkipPattern(keysIndex, options.Skip)
	} else if options.Pattern != "" {
		keysIndex = applyTargetPattern(keysIndex, options.Pattern)
	}

	keys := getKeys(keysIndex)
	sort.Strings(keys)

	if options.Interpolate {
		err := interpolate(envs, keys, nil, 3)
		if err != nil {
			return err
		}
	}

	renderEnvsInTable(c, envs, keys, paths, options)

	if options.Result {
		renderLegend(c)
	}

	return nil
}

func getOptions(c *cli.Ctx) (*Options, error) {
	options := &Options{
		Dir:         strings.TrimSuffix(c.String("dir"), "/"),
		Source:      c.String("source"),
		Target:      c.Array("target"),
		Env:         c.String("env"),
		Local:       c.Bool("local"),
		Skip:        c.String("skip"),
		Pattern:     c.String("pattern"),
		Interpolate: c.Bool("interpolate"),
		Result:      c.Bool("result"),
		System:      c.Bool("system"),
		All:         c.Bool("all"),
	}

	if options.Skip != "" && options.Pattern != "" {
		return nil, errors.New("can't have a value for both the 'skip' and 'pattern' options")
	}

	truncate := c.String("truncate")
	integer, err := strconv.Atoi(truncate)
	if err != nil {
		return options, err
	}
	options.Truncate = integer

	return options, nil
}
