# gomux roadmap

Goal: describe a tmux workspace once (session → windows → panes → commands),
then open it with one command, from the CLI or the TUI.

**How to use this file**
- Work top to bottom. Each milestone leaves the app *working*, just with less in it.
- A milestone is done when its tests are green. Run them with `just test-fixture <name>`.
- Item numbers like (1.4) point to `project-docs/code-review.md`.
- Run `just check` before every commit.

---

## Working agreements

- **One milestone at a time.** Don't start screens or commands for later milestones while an earlier one is red.
- **Delete instead of commenting out.** Git keeps the history.
- **Stubs fail loudly.** Unfinished code returns `errors.New("not implemented")`, never a silent `return nil`.
- **Small commits, descriptive messages.** One concern per commit; "upload latest code" commits are hard to review or revert.
- **Fix the bug, then add a test.** Every bug you fix gets a test or fixture that would have caught it.

---

## Quick wins (any time, ~minutes each)

- [ ] Handle the `NewScript` error in `saveScript` so `ctrl+s` stops crashing the TUI (3.1)
- [ ] Make the stubbed `gomux create` subcommands return "not implemented", or hide them with `Hidden: true` (2.1)
- [ ] Stop always writing `debug.log` to the working directory; only log when `GOMUX_DEBUG=1` (section 6)

---

## M0 — Guardrails ✅

- [x] `justfile` with `build`, `run`, `test`, `test-unit`, `test-integration`, `test-fixture`, `update-golden`, `check`
- [x] Generator tests: `bash -n` syntax check plus golden files (`internal/render/script_test.go`)
- [x] Real-tmux integration tests with an isolated server and config (`internal/render/tmux_integration_test.go`)
- [x] `actions` tests for validation, defaults and persistence (`pkg/actions/actions_test.go`)
- [x] Restructure step 1: split `internal/models` into `internal/workspace` (spec types) and `internal/render` (bash generation + its tests)

The test cases live in `internal/render/fixtures_test.go`. When the model changes shape, update the `build` funcs there; the expectations describe tmux behaviour and shouldn't need to change.

---

## M1 — One window that actually opens

The smallest useful thing: a valid script that creates a session with one window and runs its command.

- [ ] Fix the section separator and trailing newlines (1.1, 1.2, 1.3)
- [ ] Attach at the end of a fresh session; use `switch-client` when already inside tmux (1.9)
- [ ] Skip attaching when `GOMUX_NO_ATTACH=1` (needed by the tests, which have no terminal)
- [ ] Run `just update-golden` and review the diff; once the output is right, commit `testdata/single-window.golden.sh`

**Done when:** `TestStringifyIsValidBash` and `just test-fixture single-window` are green.

## M2 — Multiple windows

- [ ] Create windows with `tmux new-window`, which isn't called today (1.4)
- [ ] Target windows **by name** (`-t "=$S:=editor"`), not by number. This also makes M5 nearly free.
- [ ] Select the starting window by name

**Done when:** `just test-fixture multi-window` is green.

## M3 — Panes

- [ ] Fix the pane target, which has a stray `$` (1.5)
- [ ] Send each pane's command to the **new** pane. `split-window -P -F '#{pane_id}'` returns its ID.
- [ ] Start panes in the session directory with `-c "$SESSION_DIR"` (1.11)
- [ ] Apply the layout **after** splitting (1.8)
- [ ] Consider nesting panes under windows (`Window.Panes`) instead of `Pane.Location` (5.7)

**Done when:** `just test-fixture panes` is green.

## M4 — Quoting and paths

- [ ] Single-quote every user value in the generated script (1.6) and quote `$SESSION_NAME` everywhere (1.7)
- [ ] Consider `send-keys -l` (literal) followed by a separate `Enter`, so command text is never read as key names
- [ ] Expand a leading `~` in the session directory before saving (1.10)

**Done when:** `just test-fixture quoting` is green.

## M5 — Works with any tmux config

- [ ] Nothing in the script may assume `base-index` or `pane-base-index` (your own `tmux.conf` sets both to 1)
- [ ] Settle the numbering convention in `NewScript`: currently `Location == 0` is rejected and `startingWindow` 0 becomes 1 (5.2)

**Done when:** `just test-fixture user-config-base-index-1` and the whole `just test-integration` are green.
**Then:** run `just update-golden` and commit all golden files.

## M6 — Specs and storage

- [ ] Apply defaults on the real slice elements, not loop copies (5.1) → `TestNewScriptAppliesDefaults`
- [ ] Save the spec as `.tmux.json` (4.1) → `TestNewScriptPersistsSpecAndScript`
- [ ] Sanitize session name and type before using them as paths, and reject `/`, `..` and empty values (4.2)
- [ ] Decide what happens when a script already exists: error or explicit overwrite (4.3)
- [ ] Move validation and defaults onto the spec (`workspace.Script.Validate()`, `ApplyDefaults()`) (5.3)
- [ ] Restructure step 2: rename `internal/database` → `internal/store`, and stop calling `render.Bash` from storage (5.6)
- [ ] Restructure step 2: move errors into the packages that own them (`workspace.ErrNameRequired`, `store.ErrExists`) and delete `gomuxerr`
- [ ] Restructure step 2: move `pkg/actions` → `internal/`, or delete it if `store` + `render` are enough once validation lives on the spec (5.4)
- [ ] Give scripts an ID; `<type>/<name>` matches the on-disk layout (3.4)
- [ ] Add `List()`, `Get(id)` and `Delete(id)` to storage (4.9)
- [ ] Decide whether to store the `.sh` at all, or generate it from the spec when running (4.4)

## M7 — TUI end to end

- [ ] Restructure step 3: move `tui/` → `internal/tui/`; rename `commenui` → `style` and `shared` → `msg` (messages only)
- [ ] Rebuild the create form around session (name, dir, type), windows and panes; drop language, tags and content (3.2)
- [ ] Show a live preview of the generated script in the form
- [ ] Load scripts from disk on startup (3.3) and render them in the list (3.5)
- [ ] Run the selected script from the list
- [ ] Show validation errors in the form; fix the value-receiver bug (3.7, 3.11)
- [ ] Reset the form after saving (3.8); fix menu filtering (3.9)
- [ ] Hide the Edit and Settings menu items until those screens exist (3.6)
- [ ] Create one `actions` instance in the root model and pass it to the screens (3.12)

## M8 — CLI

- [ ] `gomux create -f spec.json` (or `gomux new <name>`, which opens `$EDITOR` on a template)
- [ ] `gomux list`
- [ ] `gomux run <id>`
- [ ] `gomux delete <id>`
- [ ] Remove the old step-by-step `create session|window|pane` subcommands (2.1–2.3)

## M9 — Cleanup and polish

- [ ] Read `$XDG_DATA_HOME`; merge `InitDB`/`NewDB`; use `path/filepath` (4.5–4.8)
- [ ] Delete unused template state and messages (section 6)
- [ ] Lowercase error strings in `gomuxerr` (5.5)
- [ ] README with install instructions and an example spec
- [ ] CI: run `just check` on every push (GitHub Actions)

---

## Later / ideas

- Edit screen and `gomux edit <id>`
- Per-window environment variables and hooks (see `project-docs/example-script.sh`)
- Apply layouts directly with `exec.Command("tmux", ...)` instead of generating bash
- Import an existing running session into a spec
