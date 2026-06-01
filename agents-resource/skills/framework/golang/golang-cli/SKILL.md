---
name: golang-cli
description: "Golang CLI application development. Use when building, modifying, or reviewing a Go CLI tool — especially for command structure, flag handling, configuration layering, version embedding, exit codes, I/O patterns, signal handling, shell completion, argument validation, and CLI unit testing. Also triggers when code uses cobra, viper, or urfave/cli."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.3"
  openclaw:
    emoji: "💻"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent AskUserQuestion
---

Go CLI applications use Cobra + Viper as the default stack: Cobra provides the command/subcommand/flag structure, Viper handles configuration layering (flags > env vars > config file > defaults). Always set `SilenceUsage: true` and `SilenceErrors: true`, use `PersistentPreRunE` for config initialization, bind all flags to Viper, and write to `cmd.OutOrStdout()` instead of `os.Stdout` to enable test capture.

## References

| File | Purpose |
|------|---------|
| references/details.md | Project structure, root command setup, flags, args, Viper config, exit codes, I/O patterns, signal handling, testing, common mistakes |
| assets/examples/ | Complete code examples: main.go, root.go, serve.go, flags.go, args.go, config.go, exit_codes.go, output.go, signal.go, completion.go, version.go, cli_test.go |
