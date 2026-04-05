# lx — Linux-to-Windows Command Translator

[![Build Status](https://github.com/andre-carbajal/lx/workflows/CI/badge.svg)](https://github.com/andre-carbajal/lx/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Description

**lx** is an open-source CLI tool that translates Linux commands to their Windows equivalents (CMD and PowerShell). It automatically detects the active shell and uses a static dictionary for common commands, with an AI engine as fallback for unknown or complex piped commands.

## Features

- 🔄 Automatic shell detection (CMD vs PowerShell)
- 📚 Static dictionary for ~50 common Linux commands
- 🤖 AI-powered fallback for unknown commands (Claude Haiku)
- ⚡ Fast path (<1ms) for known commands
- 💾 Local JSON cache for AI translations
- 🎯 Support for flags, arguments, and piped commands
- 🖥️ Windows native (CMD, PowerShell, PowerShell Core)

## Quick Start

```bash
$ lx ls -la /home
# Output: Get-ChildItem -Path /home -Force | Format-List

$ lx cat file.txt | grep error
# Output: Get-Content file.txt | Select-String error
```

## Roadmap

- **Phase 0** ✅: Foundation & scaffolding
  - Go module structure
  - Shell detection (CMD vs PowerShell)
  - Command parser (flags, args, pipes)
  - CI/CD automation
  
- **Phase 1**: Dictionary & AI integration
  - Static dictionary (~50 commands)
  - Command translation engine
  - Anthropic API integration
  - Local JSON cache
  
- **Phase 2**: Installers & releases
  - PowerShell installer
  - CMD installer
  - Automated GitHub releases
  - Windows MSI installer

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
make test        # Run all tests with coverage
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
