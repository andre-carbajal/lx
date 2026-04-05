package detector

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

type WindowsProcessProvider struct{}

func NewWindowsProcessProvider() *WindowsProcessProvider {
	return &WindowsProcessProvider{}
}

func (w *WindowsProcessProvider) GetParentProcessName(pid int) (string, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return "", fmt.Errorf("failed to create toolhelp32 snapshot: %w", err)
	}
	defer func() {
		_ = windows.CloseHandle(snapshot)
	}()

	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))

	if err := windows.Process32First(snapshot, &pe); err != nil {
		return "", fmt.Errorf("failed to get first process: %w", err)
	}

	for {
		if int(pe.ProcessID) == pid {
			// Encontrado el proceso, retornar su nombre
			return windows.UTF16ToString(pe.ExeFile[:]), nil
		}

		if err := windows.Process32Next(snapshot, &pe); err != nil {
			return "", fmt.Errorf("process %d not found: %w", pid, err)
		}
	}
}

func DefaultProvider() ProcessProvider {
	return NewWindowsProcessProvider()
}
