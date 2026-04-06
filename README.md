# lx — Linux-to-Windows Command Translator

[![Build Status](https://github.com/andre-carbajal/lx/workflows/CI/badge.svg)](https://github.com/andre-carbajal/lx/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Description

**lx** is an open-source CLI tool that translates Linux commands to their Windows equivalents (CMD and PowerShell). It automatically detects the active shell and uses a static dictionary for common commands, with an AI engine as fallback for unknown or complex piped commands.

## Features

- 🔄 Automatic shell detection (CMD vs PowerShell)
- 📚 Static dictionary for 26 core Linux commands
- ⚡ Fast translations (<1ms) for known commands
- 🚩 Support for flags, arguments, and complex commands
- 🖥️ Windows native (CMD, PowerShell, PowerShell Core)
- 🎯 --dry-run and --verbose modes for preview
- 📜 Installation scripts for PowerShell and CMD

## Quick Start

```bash
$ lx ls -la
Get-ChildItem | Format-List -Force

$ lx --dry mkdir test_folder
New-Item -ItemType Directory test_folder

$ lx --verbose pwd
Shell: powershell
Windows: Get-Location
C:\Users\User\Documents

$ lx grep "error" file.txt
Select-String -Pattern error file.txt
```

## Roadmap

- **Phase 0** ✅: Foundation & scaffolding
  - Go module structure
  - Shell detection (CMD vs PowerShell)
  - Command parser (flags, args, pipes)
  - CI/CD automation
  
- **Phase 1** ✅: Dictionary & translation engine
   - Static dictionary (~26 core commands)
   - Command translation engine
   - Executor with --dry-run and --verbose flags
   - Installation scripts for PowerShell and CMD
   
- **Phase 2**: AI integration & expanded support
   - Expand to ~50 commands
   - Anthropic API integration for unknown commands
   - Local JSON cache for AI translations
  
- **Phase 3**: Installers & releases
   - Automated GitHub releases
   - Windows MSI installer
   - Integration with package managers (Scoop, Chocolatey)

## Architecture

```
User Input: "ls -la /home"
        ↓
ShellDetector (CMD vs PowerShell)
        ↓
CommandParser (cmd + flags + args)
        ↓
TranslationRouter
  ├─ [Known] → Static Dictionary → fast path <1ms
  └─ [Unknown] → AI Engine → Cache JSON → slow path ~2-8s
        ↓
Executor (--dry-run | --verbose | exec)
        ↓
Windows Terminal Output
```

## Development

### Prerequisites

- Go 1.26.1 or later
- golangci-lint

### Setup

```bash
git clone https://github.com/andre-carbajal/lx.git
cd lx
make build
```

### Running Tests

```bash
make test        # Run all tests
make lint        # Run linter
make clean       # Clean build artifacts
```

### Building

```bash
make build       # Build binary to ./lx
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License — see [LICENSE](LICENSE)

## Author

Andre Carbajal — [@andre-carbajal](https://github.com/andre-carbajal)
