package cmd

import (
	"path/filepath"

	"github.com/michielnijenhuis/cli"
)

func renderEnvsInTable(c *cli.Ctx, envs []map[string]*EnvVar, keys []string, paths []string, options *Options) {
	headers := makeTableHeaders(envs, paths, options)
	rows := makeTableRows(envs, keys, options)

	table := cli.NewTable(c.Output)
	table.SetRows(rows)
	table.SetHeaders(headers)
	table.SetStyleByName("box")

	c.NewLine(1)
	table.Render()
}

func renderLegend(c *cli.Ctx) {
	c.NewLine(1)
	c.Writeln(cli.Dim("Legend:\n") + "<fg=red>x</> Missing   <fg=magenta>?</> Undefined   <fg=yellow>!</> Divergent")
}

func makeTableHeaders(envs []map[string]*EnvVar, paths []string, options *Options) []string {
	headers := make([]string, 0, 1+len(envs))
	headers = append(headers, "")

	for i, path := range paths {
		base := filepath.Base(path)
		headers = append(headers, base)

		if i == 0 && options.System {
			headers = append(headers, "os")
		}
	}

	if options.Result {
		headers = append(headers, "results")
	}

	return headers
}

func makeTableRows(envs []map[string]*EnvVar, keys []string, options *Options) [][]*cli.TableCell {
	rows := make([][]*cli.TableCell, 0, len(keys))

	for _, k := range keys {
		row := make([]*cli.TableCell, 0, 1+len(envs))

		var result string

		if options.Result {
			result = findValue(envs, k)
			sourceVar, exists := envs[0][k]

			missing := result == ""
			undefined := !exists
			divergent := false

			if exists && sourceVar.Value != "" && result != sourceVar.Value {
				divergent = true
			}

			headerKey := "<options=bold>" + k + "</>"
			if missing {
				headerKey = "<fg=red>x</> " + headerKey
			} else if undefined {
				headerKey = "<fg=magenta>?</> " + headerKey
			} else if divergent {
				headerKey = "<fg=yellow>!</> " + headerKey
			} else {
				headerKey = "  " + headerKey
			}

			row = append(row, cli.NewTableCell(headerKey))
		} else {
			row = append(row, cli.NewTableCell("<options=bold>"+k+"</>"))
		}

		for _, env := range envs {

			envVar, ok := env[k]
			if ok {
				row = append(row, cli.NewTableCell(getEnvVarCellValue(envVar.Value, options)))
			} else {
				row = append(row, cli.NewTableCell(""))
			}
		}

		if options.Result {
			row = append(row, cli.NewTableCell(getEnvVarCellValue(result, options)))
		}

		rows = append(rows, row)
	}

	return rows
}

func getEnvVarCellValue(value string, options *Options) string {
	if options.Truncate > 3 && len(value) > options.Truncate {
		value = cli.Truncate(value, options.Truncate, "...")
	}

	return value
}
