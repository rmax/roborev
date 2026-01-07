package agent

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// OpenCodeAgent runs code reviews using the OpenCode CLI
type OpenCodeAgent struct {
	Command string // The opencode command to run (default: "opencode")
}

// NewOpenCodeAgent creates a new OpenCode agent
func NewOpenCodeAgent(command string) *OpenCodeAgent {
	if command == "" {
		command = "opencode"
	}
	return &OpenCodeAgent{Command: command}
}

func (a *OpenCodeAgent) Name() string {
	return "opencode"
}

func (a *OpenCodeAgent) CommandName() string {
	return a.Command
}

func (a *OpenCodeAgent) Review(ctx context.Context, repoPath, commitSHA, prompt string) (string, error) {
	// Use opencode run for non-interactive mode
	args := []string{
		"run",
		prompt,
	}

	cmd := exec.CommandContext(ctx, a.Command, args...)
	cmd.Dir = repoPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("opencode failed: %w\nstderr: %s", err, stderr.String())
	}

	output := strings.TrimSpace(stdout.String())
	if len(output) == 0 {
		return "No review output generated", nil
	}

	return output, nil
}

func init() {
	Register(NewOpenCodeAgent(""))
}
