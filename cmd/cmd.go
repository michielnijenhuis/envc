package cmd

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/michielnijenhuis/cli"
)

// TODO: add option to include system values for each env var found

var Command *cli.Command

func init() {
	Command = &cli.Command{
		Name:        "envc",
		Description: "Prints all variables in the available .env file(s), sorted alphabetically",
		Handle:      handle,
	}

	Command.AddArgument(&cli.InputArgument{
		Name:         "dir",
		Description:  "The directory to check for .env files",
		Mode:         cli.InputArgumentOptional,
		DefaultValue: "./",
	})

	Command.AddOption(&cli.InputOption{
		Name:         "source",
		Description:  "The base .env file to compare with",
		Mode:         cli.InputOptionRequired,
		DefaultValue: ".env.example",
	})

	Command.AddOption(&cli.InputOption{
		Name:         "target",
		Description:  "The .env file(s) to compare against",
		Mode:         cli.InputOptionIsArray | cli.InputOptionRequired,
		DefaultValue: []string{".env"},
	})

	Command.AddOption(&cli.InputOption{
		Name:        "env",
		Shortcut:    "e",
		Description: "The environment .env to include in the comparison (e.g. `.env.dev`)",
		Mode:        cli.InputOptionRequired,
	})

	Command.AddOption(&cli.InputOption{
		Name:        "local",
		Shortcut:    "l",
		Description: "Whether to include local .env files (e.g. `.env.local` or `.env.dev.local`)",
		Mode:        cli.InputOptionBool,
	})

	Command.AddOption(&cli.InputOption{
		Name:        "result",
		Shortcut:    "r",
		Description: "Include conjuncted .env file result",
		Mode:        cli.InputOptionBool,
	})

	Command.AddOption(&cli.InputOption{
		Name:         "truncate",
		Shortcut:     "t",
		Description:  "Whether to truncate long values or not",
		Mode:         cli.InputOptionRequired,
		DefaultValue: "40",
	})

	Command.AddOption(&cli.InputOption{
		Name:        "skip",
		Shortcut:    "s",
		Description: "Comma separated list of variable name patterns to skip",
		Mode:        cli.InputOptionRequired,
	})

	Command.AddOption(&cli.InputOption{
		Name:        "pattern",
		Shortcut:    "p",
		Description: "Comma separated list of variable name patterns to focus on",
		Mode:        cli.InputOptionRequired,
	})

	Command.AddOption(&cli.InputOption{
		Name:        "interpolate",
		Shortcut:    "i",
		Description: "Interpolate env var values that refer to other env vars",
		Mode:        cli.InputOptionBool,
	})

	Command.AddOption(&cli.InputOption{
		Name:        "system",
		Description: "Include os values",
		Mode:        cli.InputOptionBool,
	})
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
}

func handle(c *cli.Command) (int, error) {
	options, err := getOptions(c)
	if err != nil {
		return 1, err
	}

	paths := getEnvFilePaths(options)
	err = validateFilePaths(paths)
	if err != nil {
		return 1, err
	}

	envs, err := readEnvFiles(paths)
	if err != nil {
		return 1, err
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
			return 1, err
		}
	}

	renderEnvsInTable(c.Output(), envs, keys, paths, options)

	return 0, nil
}

func getOptions(c *cli.Command) (*Options, error) {
	options := &Options{}

	// dir
	dir, err := c.StringArgument("dir")
	if err != nil {
		return nil, err
	}
	options.Dir = strings.TrimSuffix(dir, "/")

	// source
	source, err := c.StringOption("source")
	if err != nil {
		return nil, err
	}
	options.Source = source

	// targets
	targets, err := c.ArrayOption("target")
	if err != nil {
		return nil, err
	}
	// validate?
	options.Target = targets

	// env
	env, err := c.StringOption("env")
	if err != nil {
		return nil, err
	}
	options.Env = env

	// local
	local, err := c.BoolOption("local")
	if err != nil {
		return nil, err
	}
	options.Local = local

	// skip
	skip, err := c.StringOption("skip")
	if err != nil {
		return nil, err
	}
	options.Skip = skip

	// pattern
	pattern, err := c.StringOption("pattern")
	if err != nil {
		return nil, err
	}
	options.Pattern = pattern

	if options.Skip != "" && options.Pattern != "" {
		return nil, errors.New("can't have a value for both the 'skip' and 'pattern' options")
	}

	// interpolate
	interpolate, err := c.BoolOption("interpolate")
	if err != nil {
		return nil, err
	}
	options.Interpolate = interpolate

	// truncate
	truncate, err := c.StringOption("truncate")
	if err != nil {
		return nil, err
	}
	integer, err := strconv.Atoi(truncate)
	if err != nil {
		return nil, err
	}
	options.Truncate = integer

	// result
	includeResult, err := c.BoolOption("result")
	if err != nil {
		return nil, err
	}
	options.Result = includeResult

	// system
	system, err := c.BoolOption("system")
	if err != nil {
		return nil, err
	}
	options.System = system

	return options, nil
}
