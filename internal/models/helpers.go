package models

import (
	"encoding/json"
	"fmt"
)

func LoadScriptFromJsonBytes(jsonFile []byte) (*Script, error) {
	script := &Script{}
	err := json.Unmarshal(jsonFile, script)
	return script, err
}

func shebang() string {
	return `#!/usr/bin/env bash`
}

func checkDirectoryExistence() string {
	return fmt.Sprintf(`if [ ! -d "$%s" ]; then
  echo "Directory $SESSION_DIR does not exist. Exiting..."
  exit 1
fi`, SESSION_DIR_VARIABLE)
}

func checkIfSessionExists() string {
	return fmt.Sprintf(`if tmux has-session -t "$%s" 2>/dev/null; then
  tmux attach -t "$%s"
  exit 0
fi`, SESSION_NAME_VARIABLE, SESSION_NAME_VARIABLE)
}
