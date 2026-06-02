package main

// cmd_knowledge_graph.go — alpha --knowledge-graph: manage docker services + graph update/init (CLI).

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// runGraphifyUpdate runs graphify update + understand --update and returns combined output.
// Used by MCP update tool.
func runGraphifyUpdate(alphaDir, projectRoot, uaBin string) string {
	return runGraphifyUpdateTarget(alphaDir, projectRoot, uaBin, "--update")
}

// runGraphifyUpdateTarget runs graphify update + understand for MCP callers that need combined output.
func runGraphifyUpdateTarget(alphaDir, target, uaBin, understandFlag string, extraGfyArgs ...string) string {
	var sb strings.Builder
	if out := runGraphifyStep(alphaDir, target, extraGfyArgs...); out != "" {
		sb.WriteString(out)
	}
	if out := runUnderstandStep(alphaDir, target, uaBin, understandFlag); out != "" {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(out)
	}
	if sb.Len() == 0 {
		return "OK"
	}
	return sb.String()
}

// runGraphifyStep runs the graphify graph build and returns error output only.
func runGraphifyStep(alphaDir, target string, extraArgs ...string) string {
	gfyArgs := append([]string{"update", target}, extraArgs...)
	gfyBin := "graphify"
	if !inDocker() {
		if exe, err := os.Executable(); err == nil {
			candidate := filepath.Join(filepath.Dir(exe), "graphify")
			if _, err := os.Stat(candidate); err == nil {
				gfyBin = candidate
			}
		}
	}
	var cmd *exec.Cmd
	if inDocker() {
		cmd = exec.Command(gfyBin, gfyArgs...)
	} else {
		cmd = exec.Command("rtk", append([]string{"err", gfyBin}, gfyArgs...)...)
	}
	cmd.Dir = target
	cmd.Env = append(os.Environ(), "PROJECT_ROOT="+alphaDir, "ALPHA_ROOT="+target)
	return captureErrors(cmd)
}

// runUnderstandStep runs the understand scan and returns error output only.
func runUnderstandStep(alphaDir, target, uaBin, flag string) string {
	var cmd *exec.Cmd
	if inDocker() {
		cmd = exec.Command(uaBin, flag)
	} else {
		cmd = exec.Command("rtk", "err", uaBin, flag)
	}
	cmd.Dir = alphaDir
	cmd.Env = append(os.Environ(), "PROJECT_ROOT="+alphaDir, "ALPHA_ROOT="+target)
	return captureErrors(cmd)
}

// captureErrors runs a command and returns real error output, filtering rtk's "[ok]..." success lines.
func captureErrors(cmd *exec.Cmd) string {
	out, _ := cmd.CombinedOutput()
	s := strings.TrimSpace(string(out))
	if s == "" || strings.HasPrefix(s, "[ok]") {
		return ""
	}
	return s
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
		fmt.Println("Updating knowledge graph...")
		if out := runGraphifyStep(alphaDir, target); out != "" {
			fmt.Println(out)
		} else {
			fmt.Println("  graphify graph updated")
		}
		fmt.Println("Updating understand scan...")
		if out := runUnderstandStep(alphaDir, target, uaBin, "--update"); out != "" {
			fmt.Println(out)
		} else {
			fmt.Println("  understand scan updated")
		}

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
		fmt.Println("Building knowledge graph...")
		if out := runGraphifyStep(alphaDir, target, extraArgs...); out != "" {
			fmt.Println(out)
		} else {
			fmt.Println("  graphify graph built")
		}
		fmt.Println("Running understand scan...")
		if out := runUnderstandStep(alphaDir, target, uaBin, "--start"); out != "" {
			fmt.Println(out)
		} else {
			fmt.Println("  scan complete (run /alpha-understand start for LLM analysis)")
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", subCmd)
		fmt.Fprintln(os.Stderr, "Usage: alpha --knowledge-graph [start|stop|restart|status|logs|update|init]")
		os.Exit(1)
	}
}
