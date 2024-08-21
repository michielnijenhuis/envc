# envc

## Description:
  Prints all variables in the available .env file(s), sorted alphabetically

## Usage:
```sh
$ print [options] [--] [dir]
```

## Arguments:
```
dir                     -- The directory to check for .env files [default: ./]
```

## Options:
```
--source=SOURCE         -- The base .env file to compare with [default: .env.example]
```

```
--target=TARGET         -- The .env file(s) to compare against [default: [.env]] (multiple values allowed)
```

```
-e, --env=ENV           -- The environment .env to include in the comparison (e.g. `.env.dev`)
```

```
-l, --local             -- Whether to include local .env files (e.g. `.env.local` or `.env.dev.local`)
```

```
-r, --result            -- Include conjuncted .env file result
```

```
-t, --truncate=TRUNCATE -- Whether to truncate long values or not [default: 40]
```

```
-s, --skip=SKIP         -- Comma separated list of variable name patterns to skip
```

```
-p, --pattern=PATTERN   -- Comma separated list of variable name patterns to focus on
```

```
-i, --interpolate       -- Interpolate env var values that refer to other env vars
```

```
-h, --help              -- Display help for the print command
```

```
-q, --quiet             -- Do not output any message
```

```
-V, --version           -- Display this application version
```

```
--ansi|--no-ansi        -- Force (or disable --no-ansi) ANSI output
```

```
-n, --no-interaction    -- Do not ask any interactive question
```

```
-v|vv|vvv, --verbose    -- Increase the verbosity of messages: normal (1), verbose (2) or debug (3)
```

## Example:
```sh
$ envc --interpolate --result --env=dev --local
```

This will compare the following .env files:
- .env.example
- .env
- .env.local (if available)
- .env.dev (if available)
- .env.dev.local (if available)

Environment variables that refer to other environment variables will interpolate those values.

The comparison table will also include the final conjuncted .env file result of all included files.