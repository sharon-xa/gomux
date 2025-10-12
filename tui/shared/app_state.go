package shared

import "time"

// AppState holds all shared application state
// This lives in the Root Model
type AppState struct {
	Scripts          []Script                `json:"scripts"`
	ExecutionHistory []ScriptExecutionResult `json:"execution_history"`
	Settings         Settings                `json:"settings"`
	CurrentScriptID  string                  `json:"current_script_id"` // ID of script being edited/run
	SearchQuery      string                  `json:"search_query"`
	FilterTags       []string                `json:"filter_tags"`
	LastError        string                  `json:"last_error"`
	IsLoading        bool                    `json:"is_loading"`
}

// NewDefaultAppState creates a new AppState with sensible defaults
func NewDefaultAppState() AppState {
	return AppState{
		Scripts:          []Script{},
		ExecutionHistory: []ScriptExecutionResult{},
		Settings: Settings{
			Theme: Theme{
				Name:            "default",
				PrimaryColor:    "#00ADD8",
				SecondaryColor:  "#5DC9E2",
				BackgroundColor: "#1a1a1a",
				TextColor:       "#ffffff",
			},
			Editor:          "vim",
			DefaultLanguage: "bash",
			AutoSave:        true,
			ShowLineNumbers: true,
			TabSize:         4,
			ScriptDirectory: "./scripts",
		},
		CurrentScriptID: "",
		SearchQuery:     "",
		FilterTags:      []string{},
		LastError:       "",
		IsLoading:       false,
	}
}

func (a *AppState) GetScriptByID(id string) *Script {
	for i := range a.Scripts {
		if a.Scripts[i].ID == id {
			return &a.Scripts[i]
		}
	}
	return nil
}

func (a *AppState) AddScript(script Script) {
	script.CreatedAt = time.Now()
	script.UpdatedAt = time.Now()
	a.Scripts = append(a.Scripts, script)
}

func (a *AppState) UpdateScript(id string, script Script) bool {
	for i := range a.Scripts {
		if a.Scripts[i].ID == id {
			script.UpdatedAt = time.Now()
			script.CreatedAt = a.Scripts[i].CreatedAt // preserve creation time
			a.Scripts[i] = script
			return true
		}
	}
	return false
}

func (a *AppState) DeleteScript(id string) bool {
	for i := range a.Scripts {
		if a.Scripts[i].ID == id {
			a.Scripts = append(a.Scripts[:i], a.Scripts[i+1:]...)
			return true
		}
	}
	return false
}
