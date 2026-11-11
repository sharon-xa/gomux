package store

import (
	"fmt"
	"path/filepath"

	"github.com/sharon-xa/gomux/internal/workspace"
)

// SaveSpec writes the JSON spec to dataState/<type>/<name>.tmux.sh.
//
// TODO(M6): the spec should be <name>.tmux.json (review 4.1), and name/type
// must be sanitized before they become paths (review 4.2).
func (s *Store) SaveSpec(script *workspace.Script) error {
	jsonScript, err := script.Jsonify()
	if err != nil {
		return fmt.Errorf("failed to encode spec: %w", err)
	}

	return s.saveFile(filepath.Join(s.dataStateDir, relPath(script)), jsonScript, false)
}

// SaveScript writes a rendered script, executable, to data/<type>/<name>.tmux.sh.
func (s *Store) SaveScript(script *workspace.Script, content []byte) error {
	return s.saveFile(filepath.Join(s.dataDir, relPath(script)), content, true)
}

func relPath(script *workspace.Script) string {
	filename := script.Session.Name + ".tmux.sh"
	return filepath.Join(script.Session.Type, filename)
}
