package main

// cmd_knowledge_graph.go — alpha --knowledge-graph: manage docker services + graph update/init (CLI).

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// runGraphifyUpdate runs Python graphify update + understand --update and returns combined output.
// Used by CLI (alpha --update) and MCP update tool.
func runGraphifyUpdate(alphaDir, projectRoot, uaBin string) string {
	return runGraphifyUpdateTarget(alphaDir, projectRoot, uaBin, "--update")
}

// runGraphifyUpdateTarget runs Python graphify update against a specific target path.
// understandFlag is "--update" for incremental or "--start" for initial full scan.
// rtk err filters output to errors/warnings only — silent on success (Unix convention).
func runGraphifyUpdateTarget(alphaDir, target, uaBin, understandFlag string, extraGfyArgs ...string) string {
	var sb strings.Builder

	gfyArgs := append([]string{"update", target}, extraGfyArgs...)
	var gfyCmd *exec.Cmd
	// Use graphify-core (Go binary) — available both in Docker (/usr/local/bin/graphify-core)
	// and natively (via binPath). On host, wrap with rtk err for token-efficient output.
	gfyCoreBin := "graphify-core"
	if !inDocker() {
		// On host: look for Go binary alongside alpha
		if exe, err := os.Executable(); err == nil {
			candidate := filepath.Join(filepath.Dir(exe), "graphify-core")
			if _, err := os.Stat(candidate); err == nil {
				gfyCoreBin = candidate
			}
		}
		gfyCmd = exec.Command("rtk", append([]string{"err", gfyCoreBin}, gfyArgs...)...)
	} else {
		gfyCmd = exec.Command(gfyCoreBin, gfyArgs...)
	}
	gfyCmd.Dir = target
	gfyCmd.Env = append(os.Environ(), "PROJECT_ROOT="+alphaDir, "ALPHA_ROOT="+target)
	out1, _ := gfyCmd.CombinedOutput()
	if len(strings.TrimSpace(string(out1))) > 0 {
		sb.Write(out1)
	}

	var uaCmd *exec.Cmd
	if inDocker() {
		uaCmd = exec.Command(uaBin, understandFlag)
	} else {
		uaCmd = exec.Command("rtk", "err", uaBin, understandFlag)
	}
	uaCmd.Dir = alphaDir
	uaCmd.Env = append(os.Environ(), "PROJECT_ROOT="+alphaDir, "ALPHA_ROOT="+target)
	out2, _ := uaCmd.CombinedOutput()
	if len(strings.TrimSpace(string(out2))) > 0 {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.Write(out2)
	}

	if sb.Len() == 0 {
		return "OK"
	}
	return sb.String()
}

func handleKnowledgeGraph(alphaDir, projectRoot, uaBin string) {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: alpha --knowledge-graph [start|stop|restart|status|logs|update|init] [args...]")
		os.Exit(1)
	}
	subCmd := os.Args[2]

	globalMode := os.Getenv("ALPHA_GLOBAL") == "1"
	var composeFile string
	var composeEnv []string
	if globalMode {
		composeFile = filepath.Join(alphaDir, "docker-compose.global.yml")
		projectID := alphaProjectID(projectRoot)
		composeEnv = append(os.Environ(),
			"ALPHA_HOME="+alphaDir,
			"HOST_PROJECT_ROOT="+projectRoot,
			"ALPHA_PROJECT_ID="+projectID,
			"ALPHA_GLOBAL=1",
		)
	} else {
		composeFile = filepath.Join(alphaDir, "docker-compose.yml")
		composeEnv = append(os.Environ(), "HOST_PROJECT_ROOT="+projectRoot)
	}

	runCompose := func(extraArgs ...string) {
		args := append([]string{"docker", "compose", "-f", composeFile}, extraArgs...)
		cmd := exec.Command("rtk", args...)
		cmd.Env = composeEnv
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		cmd.Run()
	}

	switch subCmd {
	case "start":
		runCompose("--profile", "dashboard", "up", "-d")

	case "stop":
		runCompose("--profile", "dashboard", "down")

	case "restart":
		runCompose("--profile", "dashboard", "down")
		runCompose("--profile", "dashboard", "up", "-d")

	case "status":
		runCompose("ps")

	case "logs":
		logsArgs := []string{"--profile", "dashboard", "logs"}
		grepPattern := ""
		follow := false
		for i, a := range os.Args[3:] {
			switch {
			case a == "-f" || a == "--follow":
				follow = true
			case a == "--grep" && i+1 < len(os.Args[3:]):
				grepPattern = os.Args[3:][i+1]
			case strings.HasPrefix(a, "--grep="):
				grepPattern = strings.TrimPrefix(a, "--grep=")
			}
		}
		if follow {
			logsArgs = append(logsArgs, "-f")
		}
		if grepPattern != "" {
			// raw docker for pipe to grep
			logCmd := exec.Command("docker", append([]string{"compose", "-f", composeFile}, logsArgs...)...)
			logCmd.Env = composeEnv
			logCmd.Stderr = os.Stderr
			grepCmd := exec.Command("grep", "--line-buffered", grepPattern)
			grepCmd.Stdin, _ = logCmd.StdoutPipe()
			grepCmd.Stdout, grepCmd.Stderr = os.Stdout, os.Stderr
			grepCmd.Start()
			logCmd.Run()
			grepCmd.Wait()
		} else {
			// rtk log: dedup + filter repetitive lines
			args := append([]string{"log", "docker", "compose", "-f", composeFile}, logsArgs...)
			cmd := exec.Command("rtk", args...)
			cmd.Env = composeEnv
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			cmd.Run()
		}

	case "update":
		target := projectRoot
		if len(os.Args) > 3 {
			target = os.Args[3]
		}
		out := runGraphifyUpdateTarget(alphaDir, target, uaBin, "--update")
		fmt.Print(out)

	case "init":
		force := false
		target := projectRoot
		for _, a := range os.Args[3:] {
			if a == "--force" {
				force = true
			} else {
				target = a
			}
		}
		graphPath := filepath.Join(alphaProjectDataDir(alphaDir, projectRoot), "graph.json")
		if !force {
			if _, err := os.Stat(graphPath); err == nil {
				fmt.Println("Knowledge graph already initialized. Use 'alpha --knowledge-graph update' for incremental updates.")
				os.Exit(0)
			}
		}
		var extraArgs []string
		if force {
			extraArgs = append(extraArgs, "--force")
		}
		out := runGraphifyUpdateTarget(alphaDir, target, uaBin, "--start", extraArgs...)
		fmt.Print(out)

	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", subCmd)
		fmt.Fprintln(os.Stderr, "Usage: alpha --knowledge-graph [start|stop|restart|status|logs|update|init]")
		os.Exit(1)
	}
}
