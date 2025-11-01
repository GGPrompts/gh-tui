# gh-tui Development Plan

**Last Updated**: 2025-11-01
**Status**: Core MVP Complete - Enhancement Phase

---

## 🎯 Current Status

### ✅ Completed (Core MVP)

**Foundation (100%)**
- [x] Project scaffolded with Bubbletea
- [x] Tabbed interface with 5 views
- [x] GitHub CLI integration (all fetch functions)
- [x] Configuration system with YAML support
- [x] Keyboard navigation (Tab, arrows, vim keys)
- [x] Mouse support
- [x] GitHub-inspired dark theme styling
- [x] Helper utilities (truncateString, formatTime, etc.)

**UI Components (100%)**
- [x] Animated landing page with GitHub theme
- [x] Table view component (sortable)
- [x] Tree view component (expand/collapse)
- [x] All 5 view implementations:
  - [x] Pull Requests view
  - [x] Issues view
  - [x] Repositories view
  - [x] Actions view
  - [x] Gists view

**Features (100%)**
- [x] Data loading from GitHub CLI
- [x] Dual-pane layout (list + detail)
- [x] Status indicators and icons
- [x] Gist editor integration
- [x] README documentation
- [x] Builds successfully

---

## 📋 Remaining Work

### Phase 1: Interactive Actions (High Priority)

**Pull Requests Tab**
- [x] Open PR in browser (`b` key) ✅ Done
- [ ] View PR diff (`d` key)
- [ ] Merge PR with confirmation (`m` key)
- [ ] Checkout PR branch locally (`c` key)
- [ ] Add comment (`C` key)
- [ ] Approve/Request changes (`a`/`r` keys)
- [ ] Re-run failed checks (`R` key)

**Issues Tab**
- [x] Open issue in browser (`b` key) ✅ Done
- [ ] Create new issue (`n` key)
- [ ] Close/reopen issue (`x`/`r` keys)
- [ ] Add comment (`C` key)
- [ ] Edit issue (`e` key)
- [ ] Add/remove labels (`l` key)
- [ ] Assign/unassign users (`a` key)

**Repositories Tab**
- [x] Open repo in browser (`b` key) ✅ Done
- [ ] Clone repository (`c` key)
- [ ] Star/unstar repo (`s` key)
- [ ] Fork repository (`f` key)

**Actions Tab**
- [x] Open workflow in browser (`b` key) ✅ Done
- [ ] View logs (`l` key)
- [ ] Re-run workflow (`r` key)
- [ ] Cancel running workflow (`c` key)

**Gists Tab**
- [x] View gist in micro (read-only) (`o` key) ✅ Done
- [x] Edit gist with micro editor (`e` key) ✅ Done
- [x] Create new gist (`n` key) ✅ Done
- [x] Open gist in browser (`b` key) ✅ Done
- [ ] Delete gist with confirmation (`d` key)
- [ ] Fork gist (`f` key)

**Progress**: 5/29 actions completed (17%)

**Estimated Time**: 1-2 weeks remaining

---

### Phase 2: Enhanced Gists with Tree View (Medium Priority)

Based on PLUGIN_MANAGER_PLAN.md Phase 1:

- [ ] Add tree view mode to Gists tab (`t` to toggle)
- [ ] Show multi-file gists as expandable trees
- [ ] Navigate files with arrow keys (→ expand, ← collapse)
- [ ] Select individual files for view/edit
- [ ] Visual file hierarchy with tree symbols (├─ └─)
- [ ] "Install as plugin" action for gists with `[Plugin]` prefix

**Estimated Time**: 4-6 hours

---

### Phase 3: Plugin Manager Tab (Medium Priority)

Based on PLUGIN_MANAGER_PLAN.md:

**Core Features**
- [ ] New "Plugins" tab (Tab 6)
- [ ] Display installed plugins from `.claude/` directories
- [ ] List user's gists with `[Plugin]` prefix
- [ ] Tree view with expandable categories:
  - My Gists
  - Installed
  - Marketplaces
  - Local
- [ ] Expand plugins to show files
- [ ] Plugin detail panel

**Management Features**
- [ ] Install plugin from gist to `.claude/` (`i` key)
- [ ] Enable/disable plugins in settings.json (`e`/`x` keys)
- [ ] Uninstall plugins (`d` key)
- [ ] Update check and auto-update (`u` key)
- [ ] Publish local plugin to gist (`p` key)
- [ ] Fork existing plugins (`f` key)

**Estimated Time**: 1-2 weeks

---

### Phase 4: tkan Integration - Kanban Board (Low Priority)

Based on INTEGRATION_PLAN.md:

- [ ] Port GitHub Projects API from tkan (github_projects.go)
- [ ] New "Projects" tab (Tab 7) with Kanban view
- [ ] Columns: BACKLOG, TODO, IN PROGRESS, REVIEW, DONE, ARCHIVE
- [ ] Solitaire-style card stacking
- [ ] Mouse drag & drop between columns
- [ ] Visual drop indicators
- [ ] Ghost cards during drag
- [ ] Detail panel (toggleable with Tab)
- [ ] Table view mode as alternative (`v` to toggle)

**Estimated Time**: 1-2 weeks

---

### Phase 5: TFE Integration - Context Menus & Git Ops (Low Priority)

Based on INTEGRATION_PLAN.md:

**Context Menu System**
- [ ] Create reusable context menu component
- [ ] Right-click support on items
- [ ] Per-view context menu configurations
- [ ] Keyboard navigation in menus

**Git Operations**
- [ ] Pull, push, sync, fetch commands
- [ ] Git status checking
- [ ] Progress indicators
- [ ] Error handling
- [ ] Apply to Repositories tab first

**Table Enhancements**
- [ ] Apply sortable table view to all tabs
- [ ] Click headers to sort
- [ ] Multi-column sorting (Shift+click)
- [ ] View mode toggle (`v` key) for all tabs

**Estimated Time**: 1-2 weeks

---

### Phase 6: UI/UX Polish (Ongoing)

**Core UI**
- [ ] Help screen (`?` key) with keybinding cheat sheet
- [ ] Loading spinners during data fetch
- [ ] Better error messages (user-friendly)
- [ ] Empty state messages ("No PRs found")
- [ ] Confirmation dialogs for destructive actions

**Search & Filter**
- [ ] Global search (`/` key)
- [ ] Per-view filtering
- [ ] Sort options menu (`S` key)
- [ ] Filter by status, author, labels, etc.

**Configuration**
- [ ] Add `--repo` flag for default repository
- [ ] Add `--org` flag for organization
- [ ] Keybinding customization
- [ ] Theme customization (light mode)

**Estimated Time**: Ongoing

---

### Phase 7: Advanced Features (Future)

From ENHANCEMENTS.md (cherry-picked high-value items):

**New Views/Tabs**
- [ ] Notifications tab (Tab 8)
- [ ] Starred repos tab (Tab 9)
- [ ] Releases tab (Tab 10)
- [ ] Discussions tab (Tab 11)
- [ ] Activity feed tab (Tab 12)

**Workflow Features**
- [ ] Auto-refresh with configurable interval
- [ ] Bulk operations (select multiple items)
- [ ] Templates (PR/issue templates)
- [ ] Multi-repo workspace mode

**Integration**
- [ ] Open in VS Code (`code://`)
- [ ] Copy URLs to clipboard
- [ ] Export to CSV/JSON
- [ ] Webhook support for real-time updates

**Performance**
- [ ] Smart caching
- [ ] Lazy loading
- [ ] Virtual scrolling for large lists
- [ ] Background prefetching

**Estimated Time**: 2-3 months (as needed)

---

## 🎯 Recommended Next Steps

### Immediate (This Week)
1. **Interactive Actions** - Start with "open in browser" (`b` key)
   - Easiest to implement (just `gh pr/issue/repo view --web`)
   - Works across all tabs (PR, Issues, Repos, Actions)
   - Note: Gists already uses `o` for view in micro (read-only)
   - High user value
   - **Time**: 1-2 hours

2. **Help Screen** - Add `?` key handler
   - Modal overlay with keybindings
   - Static content, easy to implement
   - Helps discoverability
   - **Time**: 1-2 hours

### Short Term (Next 2 Weeks)
3. **Enhanced Gists with Tree View**
   - Foundation already exists (tree_view.go)
   - High value for multi-file gist management
   - Prepares for plugin system
   - **Time**: 4-6 hours

4. **More Interactive Actions**
   - Merge PRs, create issues, star repos
   - Build on "open in browser" pattern
   - **Time**: 1 week

### Medium Term (Next Month)
5. **Plugin Manager Tab**
   - Leverage enhanced Gists tree view
   - Enables plugin ecosystem
   - **Time**: 1-2 weeks

6. **Context Menus & Git Operations**
   - Right-click context menus
   - Git pull/push/sync
   - **Time**: 1-2 weeks

### Long Term (2-3 Months)
7. **Kanban Board (tkan integration)**
   - Full project management
   - **Time**: 1-2 weeks

8. **Additional Tabs** (Notifications, Stars, etc.)
   - As needed by users
   - **Time**: Varies

---

## 📊 Success Metrics

**MVP Goals (✅ Achieved)**
- [x] App launches without errors
- [x] Shows 5 tabs
- [x] Can switch between tabs
- [x] Loads data from GitHub
- [x] Displays data cleanly
- [x] Keyboard navigation works
- [x] Builds successfully

**Next Level Goals**
- [ ] Users can perform basic GitHub actions without browser
- [ ] Help system guides new users
- [ ] Search/filter makes large datasets manageable
- [ ] Plugin system enables extensibility

**Advanced Goals**
- [ ] Kanban board for project management
- [ ] Context menus for power users
- [ ] Multi-repo workspace support
- [ ] Real-time updates via webhooks

---

## 🗂️ File Structure

```
gh-tui/
├── main.go                      ✅ Complete
├── types.go                     ✅ Complete
├── model.go                     ✅ Complete
├── update.go                    ✅ Complete
├── update_keyboard.go           ✅ Complete
├── update_mouse.go              ✅ Complete
├── view.go                      ✅ Complete
├── config.go                    ✅ Complete
├── github.go                    ✅ Complete
├── styles.go                    ✅ Complete
├── helpers.go                   ✅ Complete
├── landing_page.go              ✅ Complete
├── table_view.go                ✅ Complete
├── tree_view.go                 ✅ Complete
├── gist_editor.go               ✅ Complete
├── view_pullrequests.go         ✅ Complete
├── view_issues.go               ✅ Complete
├── view_repositories.go         ✅ Complete
├── view_actions.go              ✅ Complete
├── view_gists.go                ✅ Complete
│
├── Future Files (Planned)
├── context_menu.go              ⏳ Phase 5
├── git_operations.go            ⏳ Phase 5
├── github_projects.go           ⏳ Phase 4
├── plugin_manager.go            ⏳ Phase 3
├── view_plugins.go              ⏳ Phase 3
├── view_projects.go             ⏳ Phase 4
└── view_notifications.go        ⏳ Phase 7
```

---

## 📚 Documentation Status

- [x] README.md - Installation and usage guide
- [ ] CONTRIBUTING.md - Contribution guidelines
- [ ] ARCHITECTURE.md - Technical architecture doc
- [ ] API.md - GitHub CLI usage patterns
- [ ] PLUGINS.md - Plugin system documentation

---

## 🎓 Learning Resources

**Referenced Projects**
- tkan (~/projects/tkan) - Kanban board, drag & drop patterns
- TFE (~/projects/TFE) - Tree view, context menus, git operations
- TUITemplate (~/projects/TUITemplate) - Base scaffolding

**Dependencies**
- github.com/charmbracelet/bubbletea - TUI framework
- github.com/charmbracelet/lipgloss - Styling
- github.com/charmbracelet/bubbles - Components
- gopkg.in/yaml.v3 - Config parsing

---

## 🤝 Contributing

Priority areas for contribution:
1. Interactive actions (merge, create, comment)
2. Search/filter functionality
3. Help system
4. Plugin system
5. Additional views (notifications, stars)

See individual phase descriptions above for detailed task lists.

---

**Questions or suggestions?** Open an issue or PR!
