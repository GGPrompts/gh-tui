package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// types.go - Type Definitions
// Purpose: All type definitions, structs, enums, and constants
// When to extend: Add new types here when introducing new data structures

// Model represents the application state
type model struct {
	// Configuration
	config Config

	// UI State
	width  int
	height int

	// Focus management
	focusedComponent string

	// Error handling
	err       error
	statusMsg string

	// GitHub data
	pullRequests []PullRequest
	issues       []Issue
	repositories []Repository
	workflowRuns []WorkflowRun
	gists        []Gist

	// View management
	activeView ViewType
	views      map[ViewType]View

	// UI state
	loading  bool
	lastSync time.Time
	showHelp bool
}

// Config holds application configuration
type Config struct {
	// Theme
	Theme       string
	CustomTheme ThemeColors

	// Keybindings
	Keybindings       string
	CustomKeybindings map[string]string

	// Layout
	Layout LayoutConfig

	// UI Elements
	UI UIConfig

	// Performance
	Performance PerformanceConfig

	// Logging
	Logging LogConfig
}

// ThemeColors defines a color theme
type ThemeColors struct {
	Primary    string
	Secondary  string
	Background string
	Foreground string
	Accent     string
	Error      string
}

// LayoutConfig defines layout settings
type LayoutConfig struct {
	Type        string  // single, dual_pane, multi_panel, tabbed
	SplitRatio  float64 // For dual_pane
	ShowDivider bool
}

// UIConfig defines UI element settings
type UIConfig struct {
	ShowTitle       bool
	ShowStatus      bool
	ShowLineNumbers bool
	MouseEnabled    bool
	ShowIcons       bool
	IconSet         string
}

// PerformanceConfig defines performance settings
type PerformanceConfig struct {
	LazyLoading     bool
	CacheSize       int
	AsyncOperations bool
}

// LogConfig defines logging settings
type LogConfig struct {
	Enabled bool
	Level   string
	File    string
}

// Custom message types
// Add your application-specific messages here

type errMsg struct {
	err error
}

type statusMsg struct {
	message string
}

type resizeMsg struct {
	width  int
	height int
}

// ViewType represents different view tabs
type ViewType int

const (
	ViewPullRequests ViewType = iota
	ViewIssues
	ViewRepositories
	ViewActions
	ViewGists
)

// View interface for all view implementations
type View interface {
	Update(tea.Msg) (View, tea.Cmd)
	View(width, height int) string
	Focus()
	Blur()
}

// GitHub data structures
type PullRequest struct {
	Number       int       `json:"number"`
	Title        string    `json:"title"`
	State        string    `json:"state"`
	Author       Author    `json:"author"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	HeadRefName  string    `json:"headRefName"`
	BaseRefName  string    `json:"baseRefName"`
	IsDraft      bool      `json:"isDraft"`
	ReviewDecision string  `json:"reviewDecision"`
	Mergeable    string    `json:"mergeable"`
	URL          string    `json:"url"`
}

type Issue struct {
	Number     int       `json:"number"`
	Title      string    `json:"title"`
	State      string    `json:"state"`
	Author     Author    `json:"author"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Labels     []Label   `json:"labels"`
	Assignees  []Author  `json:"assignees"`
	Milestone  *Milestone `json:"milestone"`
	URL        string    `json:"url"`
}

type Repository struct {
	Name              string   `json:"name"`
	NameWithOwner     string   `json:"nameWithOwner"`
	Description       string   `json:"description"`
	StargazerCount    int      `json:"stargazerCount"`
	ForkCount         int      `json:"forkCount"`
	OpenIssuesCount   int      `json:"openIssuesCount"`
	PrimaryLanguage   *Language `json:"primaryLanguage"`
	Visibility        string   `json:"visibility"`
	URL               string   `json:"url"`
}

type WorkflowRun struct {
	DatabaseId int64     `json:"databaseId"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	Conclusion string    `json:"conclusion"`
	HeadBranch string    `json:"headBranch"`
	HeadSha    string    `json:"headSha"`
	RunNumber  int       `json:"number"`
	CreatedAt  time.Time `json:"createdAt"`
	URL        string    `json:"url"`
}

type Gist struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Public      bool      `json:"public"`
	Files       []GistFile `json:"files"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	URL         string    `json:"url"`
}

// Helper types
type Author struct {
	Login string `json:"login"`
}

type Label struct {
	Name string `json:"name"`
}

type Milestone struct {
	Title string `json:"title"`
}

type Language struct {
	Name string `json:"name"`
}

type GistFile struct {
	Filename string `json:"filename"`
}

// GitHub-specific messages
type prLoadedMsg struct {
	prs []PullRequest
	err error
}

type issuesLoadedMsg struct {
	issues []Issue
	err    error
}

type reposLoadedMsg struct {
	repos []Repository
	err   error
}

type workflowsLoadedMsg struct {
	runs []WorkflowRun
	err  error
}

type gistsLoadedMsg struct {
	gists []Gist
	err   error
}

// Editor-related messages
type editorFinishedMsg struct {
	err error
}

type gistEditorFinishedMsg struct {
	gistID       string
	tempFilePath string
	wasModified  bool
	isNewGist    bool
	err          error
}
