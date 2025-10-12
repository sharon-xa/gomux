package shared

import (
	"strings"
	"time"
)

// Script represents a script in the system
type Script struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Content     string     `json:"content"`
	Language    string     `json:"language"`
	Tags        []string   `json:"tags"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastRunAt   *time.Time `json:"last_run_at,omitempty"`
	Favorite    bool       `json:"favorite"`
}

// ScriptExecutionResult represents the result of running a script
type ScriptExecutionResult struct {
	ScriptID  string        `json:"script_id"`
	Success   bool          `json:"success"`
	Output    string        `json:"output"`
	Error     string        `json:"error"`
	ExitCode  int           `json:"exit_code"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

// Settings represents application settings
type Settings struct {
	Theme           Theme  `json:"theme"`
	Editor          string `json:"editor"`           // default editor for scripts
	DefaultLanguage string `json:"default_language"` // default script language
	AutoSave        bool   `json:"auto_save"`
	ShowLineNumbers bool   `json:"show_line_numbers"`
	TabSize         int    `json:"tab_size"`
	ScriptDirectory string `json:"script_directory"` // where scripts are stored
}

// Theme settings
type Theme struct {
	Name            string `json:"name"`
	PrimaryColor    string `json:"primary_color"`
	SecondaryColor  string `json:"secondary_color"`
	BackgroundColor string `json:"background_color"`
	TextColor       string `json:"text_color"`
}

// Helper methods for Script
func (s *Script) MarkRun() {
	now := time.Now()
	s.LastRunAt = &now
}

func (s *Script) Update(name, description, content string, tags []string) {
	s.Name = name
	s.Description = description
	s.Content = content
	s.Tags = tags
	s.UpdatedAt = time.Now()
}

func (a *AppState) AddExecutionResult(result ScriptExecutionResult) {
	result.Timestamp = time.Now()
	a.ExecutionHistory = append(a.ExecutionHistory, result)

	// Update script's last run time
	if script := a.GetScriptByID(result.ScriptID); script != nil {
		script.MarkRun()
	}
}

func (a *AppState) GetRecentExecutions(limit int) []ScriptExecutionResult {
	start := max(len(a.ExecutionHistory)-limit, 0)
	return a.ExecutionHistory[start:]
}

func (a *AppState) FilterScripts(query string, tags []string) []Script {
	filtered := []Script{}

	for _, script := range a.Scripts {
		// Check query match
		matchesQuery := query == "" ||
			strings.Contains(script.Name, query) ||
			strings.Contains(script.Description, query)

		// Check tags match
		matchesTags := len(tags) == 0 || hasAllTags(script.Tags, tags)

		if matchesQuery && matchesTags {
			filtered = append(filtered, script)
		}
	}

	return filtered
}

func hasAllTags(scriptTags, filterTags []string) bool {
	tagMap := make(map[string]bool)
	for _, tag := range scriptTags {
		tagMap[tag] = true
	}

	for _, filterTag := range filterTags {
		if !tagMap[filterTag] {
			return false
		}
	}

	return true
}
