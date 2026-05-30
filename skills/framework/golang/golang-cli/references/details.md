# Go CLI Best Practices

Use Cobra + Viper as the default stack for Go CLI applications. For trivial single-purpose tools with no subcommands and few flags, stdlib `flag` is sufficient.

## Quick Reference

| Concern             | Package / Tool                       |
| ------------------- | ------------------------------------ |
| Commands & flags    | `github.com/spf13/cobra`             |
| Configuration       | `github.com/spf13/viper`             |
| Colored output      | `github.com/fatih/color`             |
| Table output        | `github.com/olekukonko/tablewriter`  |
| Interactive prompts | `github.com/charmbracelet/bubbletea` |
| Version injection   | `go build -ldflags`                  |
| Distribution        | `goreleaser`                         |

## Project Structure

```
myapp/
├── cmd/
│   └── myapp/
│       ├── main.go              # package main, only calls Execute()
│       ├── root.go              # Root command + Viper init
│       ├── serve.go             # "serve" subcommand
│       └── version.go           # "version" subcommand
├── go.mod
└── go.sum
```

See `assets/examples/main.go` for minimal main.go.

## Root Command

- `SilenceUsage: true` MUST be set — prevents printing full usage on every error
- `SilenceErrors: true` MUST be set — lets you control error output format yourself
- `PersistentPreRunE` runs before every subcommand (config always initialized)
- Logs go to stderr, output goes to stdout

See `assets/examples/root.go`.

## Flags

- **Persistent** flags inherited by all subcommands (`--config`)
- **Local** flags only for the defining command (`--port`)
- Use `MarkFlagRequired`, `MarkFlagsMutuallyExclusive`, `MarkFlagsOneRequired` for constraints
- Always bind flags to Viper: `viper.BindPFlag` so env vars and config files work

See `assets/examples/flags.go`.

## Argument Validation

| Validator | Description |
| --- | --- |
| `cobra.NoArgs` | Fails if any args provided |
| `cobra.ExactArgs(n)` | Requires exactly n args |
| `cobra.MinimumNArgs(n)` | Requires at least n args |
| `cobra.MaximumNArgs(n)` | Allows at most n args |
| `cobra.RangeArgs(min, max)` | Between min and max |

See `assets/examples/args.go`.

## Viper Configuration Precedence

1. CLI flags (highest)
2. Environment variables
3. Config file
4. Defaults (lowest)

Use `viper.SetEnvPrefix` to namespace env vars (e.g., `MYAPP_PORT`). See `assets/examples/config.go`.

## Exit Codes

| Code | Meaning |
| --- | --- |
| 0 | Success |
| 1 | General error |
| 2 | Usage error |
| 126 | Cannot execute |
| 127 | Command not found |
| 128+N | Signal N |

See `assets/examples/exit_codes.go`.

## I/O Patterns

- **stdout vs stderr**: NEVER write diagnostic output to stdout — use stderr for logs/errors
- **Detecting pipe vs terminal**: check `os.ModeCharDevice` on stdout
- **Machine-readable output**: support `--output` flag for table/json/plain formats
- Use `cmd.OutOrStdout()` and `cmd.ErrOrStderr()` in commands (not `os.Stdout`)

See `assets/examples/output.go`.

## Signal Handling

Use `signal.NotifyContext` to propagate cancellation through context. See `assets/examples/signal.go`.

## Shell Completions

Cobra generates bash, zsh, fish, PowerShell completions automatically. See `assets/examples/completion.go`.

## Version Embedding

Inject at compile time via `go build -ldflags "-X main.version=v1.0.0"`. See `assets/examples/version.go`.

## Testing CLI Commands

Test by executing programmatically and capturing output. Use `cmd.OutOrStdout()` so tests can redirect to a buffer. See `assets/examples/cli_test.go`.

## Common Mistakes

| Mistake | Fix |
| --- | --- |
| Writing to `os.Stdout` directly | Use `cmd.OutOrStdout()` |
| Calling `os.Exit()` inside `RunE` | Return an error |
| Not binding flags to Viper | Call `viper.BindPFlag` |
| Missing `viper.SetEnvPrefix` | `PORT` collides; use `MYAPP_PORT` |
| Logging to stdout | Logs go to stderr |
| Printing usage on every error | Set `SilenceUsage: true` |
| Config file required | Ignore `viper.ConfigFileNotFoundError` |
| Not using `PersistentPreRunE` | Config init must run before any subcommand |
| Hardcoded version string | Inject via `ldflags` at build time |
| Not supporting `--output` format | Add JSON/table/plain for machine consumption |
