package models_test

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite testdata/*.golden.sh from the current output")

// TestStringifyIsValidBash checks every fixture renders to a script that bash
// can parse, with a shebang that points at a real interpreter.
func TestStringifyIsValidBash(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not found in PATH")
	}

	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			script := f.build(goldenDir).Stringify()

			checkShebang(t, script)

			path := filepath.Join(t.TempDir(), "script.sh")
			if err := os.WriteFile(path, []byte(script), 0o644); err != nil {
				t.Fatal(err)
			}

			out, err := exec.Command(bash, "-n", path).CombinedOutput()
			if err != nil {
				t.Errorf("bash -n rejected the script: %v\n%s\n--- script ---\n%s", err, out, numbered(script))
			}
		})
	}
}

func checkShebang(t *testing.T, script string) {
	t.Helper()

	first, _, _ := strings.Cut(script, "\n")
	if !strings.HasPrefix(first, "#!") {
		t.Errorf("first line is not a shebang: %q", first)
		return
	}

	fields := strings.Fields(strings.TrimPrefix(first, "#!"))
	if len(fields) == 0 {
		t.Errorf("empty shebang: %q", first)
		return
	}

	interpreter := fields[0]
	if filepath.Base(interpreter) == "env" && len(fields) > 1 {
		interpreter = fields[1]
		if _, err := exec.LookPath(interpreter); err != nil {
			t.Errorf("shebang %q: env can't find %q", first, interpreter)
		}
		return
	}

	if _, err := os.Stat(interpreter); err != nil {
		t.Errorf("shebang %q: interpreter %q does not exist", first, interpreter)
	}
}

// TestStringifyGolden pins the exact generated output, so later refactors
// (e.g. moving to text/template) can't change scripts by accident.
//
// Golden files are only created on request, since pinning broken output would
// lock the bugs in. Once a fixture's script is correct, run `just update-golden`
// and review the diff before committing.
func TestStringifyGolden(t *testing.T) {
	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			got := f.build(goldenDir).Stringify()
			path := filepath.Join("testdata", f.name+".golden.sh")

			if *update {
				if err := os.MkdirAll("testdata", 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
					t.Fatal(err)
				}
				t.Logf("updated %s", path)
				return
			}

			want, err := os.ReadFile(path)
			if os.IsNotExist(err) {
				t.Skipf("no golden file yet (%s); run `just update-golden` once this output is correct", path)
			}
			if err != nil {
				t.Fatal(err)
			}

			if got != string(want) {
				t.Errorf("output differs from %s (run `just update-golden` if the change is intended)\n--- want ---\n%s\n--- got ---\n%s",
					path, numbered(string(want)), numbered(got))
			}
		})
	}
}

// numbered prefixes each line with its line number, to make failures readable.
func numbered(s string) string {
	var b strings.Builder
	sc := bufio.NewScanner(strings.NewReader(s))
	for i := 1; sc.Scan(); i++ {
		fmt.Fprintf(&b, "%3d  %s\n", i, sc.Text())
	}
	return b.String()
}
