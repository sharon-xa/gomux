package messages

import (
	"time"
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
