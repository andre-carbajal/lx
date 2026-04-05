package executor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/andre-carbajal/lx/internal/detector"
	"github.com/andre-carbajal/lx/internal/parser"
)

type Executor struct{}

func (e *Executor) Execute(
	translated string,
	shell detector.ShellType,
	globalFlags *parser.GlobalFlags,
) int {
	if globalFlags.DryRun {
		fmt.Println(translated)
		return 0
	}

	if globalFlags.Verbose {
		fmt.Printf("Shell: %s\n", shell.String())
		fmt.Printf("Windows: %s\n", translated)
		fmt.Println("─────────────────────────────────────────")
	}

	var cmd *exec.Cmd
	if shell == detector.ShellPowerShell {
		cmd = exec.Command("powershell", "-NoProfile", "-Command", translated)
	} else {
		cmd = exec.Command("cmd", "/c", translated)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return ee.ExitCode()
		}
		// Otros errores retornan 1
		return 1
	}

	return 0
}

func (e *Executor) ExecuteWithResult(
	translated string,
	shell detector.ShellType,
) (string, int) {
	var cmd *exec.Cmd
	if shell == detector.ShellPowerShell {
		cmd = exec.Command("powershell", "-NoProfile", "-Command", translated)
	} else {
		cmd = exec.Command("cmd", "/c", translated)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ee, ok := errors.AsType[*exec.ExitError](err); ok {
			return string(output), ee.ExitCode()
		}
		return string(output), 1
	}

	return string(output), 0
}
