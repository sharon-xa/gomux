package shared

import (
	"time"

	"github.com/sharon-xa/gomux/internal/models"
)

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

func (a *AppState) AddExecutionResult(result ScriptExecutionResult) {
	result.Timestamp = time.Now()
	a.ExecutionHistory = append(a.ExecutionHistory, result)

	// Update script's last run time
	if script := a.GetScriptByID(result.ScriptID); script != nil {
	}
}

func (a *AppState) GetRecentExecutions(limit int) []ScriptExecutionResult {
	start := max(len(a.ExecutionHistory)-limit, 0)
	return a.ExecutionHistory[start:]
}

func (a *AppState) FilterScripts(query string, tags []string) []models.Script {
	return nil
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
