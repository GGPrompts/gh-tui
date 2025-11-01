# Next Steps for gh-tui

**Last Updated**: 2025-10-31 (Session 1)
**Status**: Core architecture complete, view integration needed

---

## 🎯 Immediate Next Steps (Session 2)

### 1. Review Subagent Output Files
Check what the subagents created:
```bash
ls -la view_*.go
cat view_pullrequests.go | head -50  # Check structure
cat view_issues.go | head -50
cat view_repositories.go | head -50
cat view_actions.go | head -50
cat view_gists.go | head -50
```

### 2. Ensure View Interface Compliance
Each view file needs these methods:
```go
type View interface {
    Update(tea.Msg) (View, tea.Cmd)
    View(width, height int) string
    Focus()
    Blur()
}
```

**Likely needed fixes**:
- Change `Update(tea.Msg) (tea.Model, tea.Cmd)` to return `(View, tea.Cmd)`
- Ensure `View(width, height int) string` signature (not just `View()`)
- Add `Focus()` and `Blur()` methods if missing

### 3. Add Helper Functions
Create `helpers.go` with utility functions used by views:
```go
// truncateString truncates text to fit width
func truncateString(s string, maxWidth int) string {
    if len(s) <= maxWidth {
        return s
    }
    if maxWidth <= 3 {
        return s[:maxWidth]
    }
    return s[:maxWidth-3] + "..."
}

// formatTime formats timestamps nicely
func formatTime(t time.Time) string {
    return t.Format("Jan 2, 15:04")
}

// formatTimeAgo shows relative time
func formatTimeAgo(t time.Time) string {
    duration := time.Since(t)
    if duration.Hours() < 1 {
        return fmt.Sprintf("%.0fm ago", duration.Minutes())
    }
    if duration.Hours() < 24 {
        return fmt.Sprintf("%.0fh ago", duration.Hours())
    }
    days := int(duration.Hours() / 24)
    if days < 30 {
        return fmt.Sprintf("%dd ago", days)
    }
    return fmt.Sprintf("%dmo ago", days/30)
}
```

### 4. Fix View Constructors in model.go
Update `initialModel()` to match actual constructor signatures:
```go
// Current (might not match):
m.views[ViewPullRequests] = NewPullRequestView()

// May need to be:
m.views[ViewPullRequests] = NewPullRequestView() // Check signature
```

### 5. Try Building
```bash
cd ~/projects/gh-tui
go mod tidy
go build
```

**Common compilation errors to expect**:
- Type mismatches in View interface
- Missing imports (fmt, strings, time, exec)
- Undefined helper functions (truncateString, formatTime, etc.)
- Constructor signature mismatches

### 6. Fix Compilation Errors Systematically
For each error:
1. Read the error message carefully
2. Check which view file has the issue
3. Fix the specific method/import/signature
4. Rebuild and repeat

### 7. Test Run
Once it compiles:
```bash
# Make sure you're authenticated with gh CLI
gh auth status

# Run the app
./gh-tui

# Test keyboard shortcuts:
# - Tab / Shift+Tab: Switch views
# - 1-5: Jump to specific view
# - r: Refresh current view
# - q: Quit
```

---

## 🐛 Expected Issues & Solutions

### Issue 1: View Constructor Signatures Don't Match
**Symptom**: "cannot use NewPullRequestView() as View"
**Solution**: Check what parameters the constructor expects and provide them

### Issue 2: Update() Return Type Wrong
**Symptom**: "method Update has wrong signature"
**Solution**: Change return type from `(tea.Model, tea.Cmd)` to `(View, tea.Cmd)`

### Issue 3: Missing Helper Functions
**Symptom**: "undefined: truncateString"
**Solution**: Create helpers.go with all utility functions

### Issue 4: View Not Rendering
**Symptom**: App runs but shows blank/error
**Solution**: Add debug logging, check if view is nil, verify data is loading

### Issue 5: GitHub Data Not Loading
**Symptom**: "No data found" messages
**Solution**:
- Check `gh` CLI works: `gh pr list --limit 5`
- Check repo context: run from a git repo or specify `--repo owner/repo`
- Add command-line flag support for repo selection

---

## 🎨 Polish Tasks (After Basic Functionality Works)

### UI Improvements
- [ ] Add loading spinners during data fetch
- [ ] Improve error messages (user-friendly)
- [ ] Add empty state messages ("No PRs found")
- [ ] Fine-tune colors and borders
- [ ] Add status icons (✓, ✗, •, etc.)

### Feature Additions
- [ ] Help screen (? key)
- [ ] Search/filter within views
- [ ] Sort options (by date, author, status)
- [ ] Detail pane scrolling
- [ ] Open in browser (o key)
- [ ] Interactive actions (merge PR, close issue)

### Configuration
- [ ] Add `--repo` flag for default repository
- [ ] Add `--org` flag for organization repos
- [ ] Config file (~/.config/gh-tui/config.yaml)
- [ ] Keybinding customization

---

## 📦 File Structure Summary

### ✅ Complete Files
- `main.go` - Entry point with auth check
- `types.go` - All data types and interfaces
- `model.go` - State initialization
- `view.go` - Tabbed rendering system
- `update.go` - Message handlers
- `update_keyboard.go` - Keyboard navigation
- `styles.go` - GitHub theme
- `config.go` - Configuration management
- `github.go` - CLI integration

### 🔄 Need Integration/Fixing
- `view_pullrequests.go`
- `view_issues.go`
- `view_repositories.go`
- `view_actions.go`
- `view_gists.go`

### ❌ Need Creation
- `helpers.go` - Utility functions
- `README.md` - Usage documentation (optional)

---

## 🧪 Testing Checklist

Once app runs:
- [ ] Launches without errors
- [ ] GitHub auth check works
- [ ] Tabs display correctly
- [ ] Tab switching works (Tab/Shift+Tab/1-5)
- [ ] PRs load and display
- [ ] Issues load and display
- [ ] Repositories load and display
- [ ] Actions load and display
- [ ] Gists load and display
- [ ] Refresh works (r key)
- [ ] Navigation works (up/down in lists)
- [ ] Detail panes show correctly
- [ ] Window resize handles gracefully
- [ ] Quit works (q key)

---

## 💡 Quick Debug Commands

```bash
# Check file status
ls -la *.go

# Count lines per file
wc -l *.go

# Check for missing imports
go build 2>&1 | grep "undefined"

# Test GitHub CLI
gh auth status
gh pr list --limit 3
gh issue list --limit 3
gh repo list --limit 3

# Run with verbose errors
go run . 2>&1 | less
```

---

## 🎯 Success Criteria

**Minimum Viable Product (MVP)**:
✅ App launches with auth check
✅ Shows 5 tabs (PRs, Issues, Repos, Actions, Gists)
✅ Can switch between tabs
✅ Loads data from GitHub
✅ Displays data in list format
✅ Can navigate lists with arrow keys
✅ Can quit with q

**Next Level**:
- [ ] Interactive actions (merge, comment, etc.)
- [ ] Search/filter functionality
- [ ] Help screen
- [ ] Mouse support
- [ ] Better error handling

---

Good luck! The hard part (architecture) is done. Now it's just integration and polish! 🚀
