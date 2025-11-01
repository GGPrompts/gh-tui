# gh-tui 🚀

> Interactive Terminal UI for GitHub - Built with [Bubbletea](https://github.com/charmbracelet/bubbletea)

A fast, keyboard-driven terminal interface for GitHub that brings your PRs, issues, repositories, workflow runs, and gists into a beautiful TUI experience.

## ✨ Features

- 📋 **5 Comprehensive Views**
  - Pull Requests - Review and manage PRs
  - Issues - Track and organize issues
  - Repositories - Browse your repos with stats
  - Workflow Runs - Monitor GitHub Actions
  - Gists - Manage your code snippets

- 🎨 **Beautiful UI**
  - GitHub-inspired dark theme
  - Dual-pane layout (list + detail)
  - Smooth keyboard navigation
  - Real-time status indicators with icons

- ⚡ **Fast & Efficient**
  - Async data loading
  - Built on GitHub CLI (`gh`)
  - Lightweight binary (~5MB)
  - Minimal resource usage

- ⌨️ **Keyboard-First**
  - Vim-style navigation (h/j/k/l)
  - Tab switching between views
  - Quick jump shortcuts (1-5)
  - No mouse required

## 📋 Prerequisites

- [GitHub CLI (`gh`)](https://cli.github.com/) - Required for API access
- Go 1.21+ (for building from source)

### Install GitHub CLI

**macOS:**
```bash
brew install gh
```

**Linux:**
```bash
# Debian/Ubuntu
sudo apt install gh

# Fedora/RHEL
sudo dnf install gh

# Arch
sudo pacman -S github-cli
```

**Authenticate:**
```bash
gh auth login
```

## 🚀 Installation

### From Source

```bash
git clone https://github.com/GGPrompts/gh-tui.git
cd gh-tui
go build
sudo mv gh-tui /usr/local/bin/  # Optional: install globally
```

### Using Go Install

```bash
go install github.com/GGPrompts/gh-tui@latest
```

## 💻 Usage

### Basic Usage

Run from any Git repository:
```bash
cd ~/my-github-project
gh-tui
```

Or specify a repository:
```bash
# Note: Repository and Gists views work from anywhere
# PR, Issues, and Actions views require a repo context
cd ~/projects/gh-tui
gh-tui
```

### Views Overview

**1. Pull Requests** (`Tab 1`)
- List all PRs with status indicators
- View PR details (author, branch, reviews, mergeable status)
- See draft/open/merged states
- Quick refresh with `r`

**2. Issues** (`Tab 2`)
- Browse issues with labels and assignees
- View issue details and milestones
- Filter by state (open/closed)
- Track issue activity

**3. Repositories** (`Tab 3`)
- List all your repositories
- View stats (stars ⭐, forks 🍴, open issues 📝)
- See primary language and visibility
- Quick access to repo details

**4. Workflow Runs** (`Tab 4`)
- Monitor GitHub Actions status
- View run conclusions (✓ success, ✗ failure)
- See branch and commit info
- Track workflow run history

**5. Gists** (`Tab 5`)
- Browse your gists
- View public/private status (🌐/🔒)
- See file listings
- Quick access to gist URLs

## ⌨️ Keyboard Shortcuts

### Global Commands
| Key | Action |
|-----|--------|
| `q` | Quit application |
| `?` | Show help |
| `Tab` / `Shift+Tab` | Switch views (next/previous) |
| `1` - `5` | Jump to specific view |
| `r` | Refresh current view |

### Navigation
| Key | Action |
|-----|--------|
| `↑` / `k` | Move up in list |
| `↓` / `j` | Move down in list |
| `g` | Go to top |
| `G` | Go to bottom |

### View-Specific
| Key | Action | Views |
|-----|--------|-------|
| `Enter` | View details | All |
| `o` | Open in browser | All (coming soon) |
| `c` | Create new | PRs, Issues (coming soon) |

## 🏗️ Development

### Project Structure

```
gh-tui/
├── main.go              # Entry point with auth check
├── types.go             # Data types & interfaces
├── model.go             # State management
├── view.go              # Main rendering logic
├── update.go            # Message handlers
├── update_keyboard.go   # Keyboard input
├── update_mouse.go      # Mouse support
├── styles.go            # GitHub theme & styles
├── config.go            # Configuration management
├── github.go            # GitHub CLI integration
├── helpers.go           # Utility functions
├── view_pullrequests.go # PRs view implementation
├── view_issues.go       # Issues view implementation
├── view_repositories.go # Repos view implementation
├── view_actions.go      # Actions view implementation
└── view_gists.go        # Gists view implementation
```

### Building

```bash
go mod tidy
go build
```

### Testing

```bash
# Run the app in development
./gh-tui

# Test with a specific repo context
cd ~/path/to/repo
~/projects/gh-tui/gh-tui
```

### Adding New Features

1. **New View**: Create `view_newfeature.go` implementing the `View` interface
2. **New Data Type**: Add struct to `types.go` and fetch function to `github.go`
3. **New Keyboard Shortcut**: Add to `update_keyboard.go`
4. **New Style**: Define in `styles.go` using Lipgloss

## 🎨 Customization

### Themes

The app uses a GitHub-inspired dark theme by default. Colors are defined in `styles.go`:

- Primary: `#58A6FF` (GitHub blue)
- Success: `#3FB950` (GitHub green)
- Error: `#F85149` (GitHub red)
- Warning: `#D29922` (GitHub yellow)
- Background: `#0D1117` (GitHub dark)

### Configuration

Configuration options are available in `config.go`. Future versions will support:
- Custom themes
- Keybinding customization
- Default view settings
- Refresh intervals

## 🛣️ Roadmap

- [ ] Interactive actions (merge PRs, close issues, etc.)
- [ ] Search and filtering within views
- [ ] Open items in browser (`o` key)
- [ ] Create new PRs/issues from TUI
- [ ] Configuration file support
- [ ] Help screen overlay (`?` key implementation)
- [ ] Mouse support for clicking
- [ ] Pagination for large datasets
- [ ] Custom sorting options
- [ ] Notification support

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details

## 🙏 Acknowledgments

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) - The best TUI framework for Go
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss) - Beautiful terminal styling
- Powered by [GitHub CLI](https://cli.github.com/) - Official GitHub command-line tool
- Inspired by [lazygit](https://github.com/jesseduffield/lazygit) - Simple terminal UI for git

## 📧 Contact

**Author**: GGPrompts
**Repository**: [github.com/GGPrompts/gh-tui](https://github.com/GGPrompts/gh-tui)

---

**Made with ❤️ and Go**
