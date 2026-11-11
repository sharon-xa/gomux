# gomux

A Go CLI/TUI for tmux workspaces: describe a session once (session → windows → panes → commands), store it as a JSON spec, generate a bash script from it, and open it with one command.

**Early development.** Much of the app is stubbed or unfinished. Before working on anything, read:
- `project-docs/todo.md`: the roadmap. Work happens **one milestone at a time, top to bottom**.
- `project-docs/code-review.md`: the review the roadmap refers to by item number, e.g. (1.4).

## Commands

Use `just` (no Makefile):

- `just check`: gofmt, vet, staticcheck (if installed), all tests. Run before every commit.
- `just test-unit`: fast tests, no tmux.
- `just test-fixture <name>`: every test for one milestone fixture (e.g. `single-window`).
- `just test-integration`: real tmux tests, verbose and uncached.
- `just update-golden`: rewrite golden scripts. Only run when the output is known correct, and review the diff.
- `just build` / `just run <args>`

Many tests are **expected to fail** until their milestone is done. Don't "fix" a test to match broken output; fix the code. Don't disable or skip failing tests to get `just check` green.

## Layout

- `main.go`, `cmd/`: Cobra commands. No args launches the TUI. `create` subcommands are stubs.
- `internal/workspace/`: the spec types (`Script`, `Session`, `Window`, `Pane`), JSON (de)serialization, and the rules: `ApplyDefaults()`, `Validate()` and the `Err*` validation errors (compare with `errors.Is`). No bash, no storage.
- `internal/render/`: `render.Bash(*workspace.Script)` turns a spec into a bash script. The most broken part (review section 1), and where milestones M1–M5 happen. It holds the generator tests.
- `internal/store/`: file storage under `~/.local/share/gomux/` (`dataState/` = JSON specs, `data/` = generated scripts). Saves bytes it's given; never renders.
- `internal/actions/`: operations shared by the CLI and the TUI, combining workspace, render and store. `NewScript(...)` applies defaults, validates, saves the spec, renders, then saves the script.
- `internal/tui/`: Bubble Tea. The root model (`model.go`) owns `AppState` (`state.go`) and routes messages to the screens in `screens/`. Screens talk to the root only through the message types in `messages/` (plus `ScreenState` and message payloads). Styles and the help bar are in `style/`.
  - The package is `messages`, not `msg`, so it isn't hidden by the usual `Update(msg tea.Msg)` parameter.

## Testing

- The test cases are in `internal/render/fixtures_test.go`. Expectations describe tmux behaviour, not model fields. When the model changes shape, update only the `build` funcs.
- The integration tests start a private tmux server per case: own `TMUX_TMPDIR`, `HOME`, `XDG_CONFIG_HOME`, `SHELL=/bin/sh`, with `TMUX` unset. Keep it that way. Tests must never touch the user's running tmux server, `~/.config/tmux/tmux.conf`, or `~/.local/share/gomux`.
- **Generator contract:** when `GOMUX_NO_ATTACH=1` is set, generated scripts must not `attach` or `switch-client`. Tests have no terminal.
- The user's own tmux config sets `base-index 1` and `pane-base-index 1`. Generated scripts must not assume any index base; target windows by name.
- Tests in `internal/actions` set `HOME` to a temp dir before calling `actions.New()`. Do the same in any test that touches storage.

## Conventions

- Generated bash: single-quote every user-supplied value. Never interpolate raw strings into `"..."`.
- Delete dead code instead of commenting it out.
- Unfinished code returns an error (`errors.New("not implemented")`), never a silent `return nil`.
- Errors live in the package that owns them (`workspace.ErrSessionNameRequired`), not in a shared errors package. Messages are lowercase.
- Every bug fix comes with a test or fixture that would have caught it.
- Commits: small, one concern each, Conventional Commit prefixes (`feat:`, `fix:`, `test:`, `docs:`, `build:`, `chore:`).
- When a roadmap item is done, tick it in `project-docs/todo.md` in the same commit, and update its **Status** line when a milestone finishes.
- Git workflow:
  - Only commit when the user asks.
  - Commit on a short-lived branch, then `git merge --ff-only` into `main` and delete the branch.
  - Never push; the user pushes.
- Before committing a behaviour-preserving refactor, compare before and after: test outcomes, generated script output, and (for TUI changes) screen captures from an isolated tmux server.
