package detector

import (
	"os"
	"strings"
)

type ShellType int

const (
	ShellCMD ShellType = iota
	ShellPowerShell
	ShellUnknown
)

type ProcessProvider interface {
	GetParentProcessName(pid int) (string, error)
}

type Detector struct {
	provider ProcessProvider
}

func New(provider ProcessProvider) *Detector {
	return &Detector{provider: provider}
}

func (d *Detector) Detect() ShellType {
	if shell := os.Getenv("LX_SHELL"); shell != "" {
		switch strings.ToLower(shell) {
		case "cmd":
			return ShellCMD
		case "ps", "powershell":
			return ShellPowerShell
		}
	}

	ppid := os.Getppid()
	parentName, err := d.provider.GetParentProcessName(ppid)
	if err != nil {
		return ShellUnknown
	}

	parentName = strings.ToLower(parentName)
	switch {
	case parentName == "cmd.exe":
		return ShellCMD
	case parentName == "powershell.exe" || parentName == "pwsh.exe":
		return ShellPowerShell
	default:
		return ShellUnknown
	}
}

func (s ShellType) String() string {
	switch s {
	case ShellCMD:
		return "cmd"
	case ShellPowerShell:
		return "powershell"
	default:
		return "unknown"
	}
}
