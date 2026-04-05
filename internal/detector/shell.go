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

// ProcessProvider abstrae la lectura del proceso padre (inyección de dependencias)
type ProcessProvider interface {
	GetParentProcessName(pid int) (string, error)
}

type Detector struct {
	provider ProcessProvider
}

func New(provider ProcessProvider) *Detector {
	return &Detector{provider: provider}
}

// Detect retorna el tipo de shell, con override via LX_SHELL env var
func (d *Detector) Detect() ShellType {
	// Override via env var
	if shell := os.Getenv("LX_SHELL"); shell != "" {
		switch strings.ToLower(shell) {
		case "cmd":
			return ShellCMD
		case "ps", "powershell":
			return ShellPowerShell
		}
	}

	// Lectura via provider (Windows API)
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

// String retorna nombre legible del shell
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
