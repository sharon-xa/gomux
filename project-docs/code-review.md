# gomux — Codebase Review

_Reviewed: 2026-11-11 · commit `854153b` ("upload latest code") · ~2.1k lines of Go_
_First pass was on `035f681`; see [What changed since the first pass](#what-changed-since-the-first-pass)._
_Since then, `internal/models` has been split. Rendering (`Stringify` → `render.Bash`, the constants and helpers) is now in `internal/render/bash.go`, and the types are in `internal/workspace/`. Restructure step 2 then renamed `internal/database` → `internal/store`, moved `pkg/actions` → `internal/actions`, moved validation and its errors onto `workspace` (deleting `gomuxerr`), and fixed 5.1 and 5.5. File and line references below are for `854153b`._

**Method:** I read every Go file. `go build` and `go vet` both pass. I also:
- rendered a sample script from `internal/models` and ran it through `bash -n`
- ran the CLI with `HOME` pointed at a sandbox directory
- ran throwaway Go tests against a scratch copy of the repo (nothing was added to your tree)

Every finding marked **[verified]** was reproduced, not just inferred from the code.

> **Context:** the project is in early development, so stubs, unfinished screens and leftover scaffolding are expected, and sections 3.6–3.12 and 6 are listed for tracking, not as criticism. What matters at this stage is sections 1, 3.1–3.5, 4.2, 5.1–5.3 and 7. They are foundational: the code built on top of them gets harder to change the longer they stay as they are.

---

## TL;DR

1. **The generated bash script is still broken.** The generator didn't change in this commit, and its output fails `bash -n`. No generated script can run yet.
2. **The TUI crashes on save** **[verified]**. `ctrl+s` in the create form calls `NewScript` with empty structs, ignores the error, and dereferences a nil `*Script`.
3. **The architecture moved in the right direction.** The TUI now uses `models.Script` and goes through `actions`, and the stateful step-by-step builder is gone. But the create form still collects the old generic fields (language, tags, content), which map to nothing, and scripts are never loaded back from disk.
4. **`NewScript` silently drops its defaults** **[verified]**. It sets `Layout` and `Direction` on loop copies, and its window-numbering rules contradict the generator's.
5. **There are still no tests.** A golden-file test for `Stringify()` and a table test for `NewScript` would cover the two most broken areas.

---

## What changed since the first pass

| Area | Change in `854153b` | Effect on the review |
|------|---------------------|----------------------|
| `pkg/actions` | Stateful builder (`AddSession`/`AddWindow`/`AddPane` + `CurrentScript`) replaced by one all-in-one `NewScript(...)` that validates and saves | Good move. Old 2.1–2.3, 2.5 and 5.3 resolved. New bugs in 5.1–5.2. |
| `internal/gomuxerr` | New sentinel errors package | Good idea; small notes in 5.5. |
| CLI `create` | All subcommands stubbed (`return nil`) | Old "fails / writes nameless file" bugs gone, replaced by "silently does nothing" (2.1). |
| TUI types | `shared.Script` deleted; `AppState`, messages and screens use `models.Script` | **Addresses the core "two products" problem** at the type level. New issues 3.1–3.5. |
| TUI styling | Styles and help moved into `commenui` (`PageStyle`, `LabelStyle`, `GetHelpMessage`, …) and used by main menu and create | Old 6 (duplicated inline styles) mostly resolved; list screen still inline. |
| Main menu | Per-render `log.Printf` removed | Old 3.10 (`debug.log` growth) resolved. |
| `Constants.go` | Renamed to `constants.go` | Filename nit resolved. |
| Generator, storage, `root.go`, `list.go`, `main.go` | **Unchanged** | Sections 1 and 4 still apply in full. |

---

## 1. Critical: generated script output _(unchanged)_

Here is the output for a session with 2 windows and 1 pane **[verified]**:

```bash
#!/usr/bin/env bash#            <- shebang now asks env for "bash#"
--------                        <- executed as a command: "command not found"
SESSION_NAME="dev"
SESSION_DIR="/tmp"#
--------
if [ ! -d "$SESSION_DIR" ]; then
  ...
fi#                             <- "fi#" is not "fi": the if never closes -> SYNTAX ERROR
...
tmux new-session -d -s $SESSION_NAME -c "$SESSION_DIR"#    <- dir becomes "/tmp#"
--------
# Window Number 0tmux rename-window -t $SESSION_NAME:0 "editor"   <- command commented out
tmux send-keys -t $SESSION_NAME:1 "echo "hi"" C-m                  <- quoting broken
## Pane Number 0tmux split-window -h -t $SESSION_NAME:$1.0         <- commented out + "$1"
```

`bash -n` reports: `syntax error: unexpected end of file from 'if' command on line 11`.

| # | Issue | Location |
|---|-------|----------|
| 1.1 | `SECTION_SEPARATOR = "#\n--------\n"`. The `#` sticks to the end of the previous line and corrupts it (`bash#`, `fi#`, `"$SESSION_DIR"#`). `--------` sits on its own uncommented line, so bash tries to run it. Use something like `"\n\n# --------\n"`. | `internal/models/constants.go:4` |
| 1.2 | `SCRIPT_TEMPLATE_SECTIONS = 10`, but only 8 slots are filled, so two empty separator sections get appended. Build with `append` or `strings.Builder` instead of a fixed-size slice. | `constants.go:7`, `script.go:32-42` |
| 1.3 | The window and pane header comments have no trailing `\n`, so the first tmux command of every window and pane is **commented out**. | `script.go:58`, `script.go:74` |
| 1.4 | Windows are never created. `createWindow` only runs `rename-window` and `send-keys` against `$SESSION:N`. Only the first window exists after `new-session`, so every other window fails with "can't find window". | `tmux.go:38-76` |
| 1.5 | The pane target format `$%s:$%d.%d` has a stray `$`, producing `$SESSION_NAME:$1.0`. `$1` is bash's first positional argument (empty), so the target becomes `dev:.0`. | `tmux.go:86` |
| 1.6 | No shell quoting or escaping. `Name`, `Directory` and `Command` go into `"..."` as-is. Any `"`, `$`, `` ` `` or `\` breaks the script or runs at generation time (for example `$PORT` expands in the script, not in the pane). Use single-quote escaping (`'` → `'\''`) or `%q`-style quoting. | `script.go:46`, `tmux.go:46-53`, `tmux.go:92-97` |
| 1.7 | `$SESSION_NAME` is unquoted in `-s` and `-t`, so names with spaces break. tmux also rewrites `.` and `:` in session names, so later targets won't match. | `tmux.go:24`, `tmux.go:39`, `tmux.go:86`, `script.go:67` |
| 1.8 | `select-layout` runs while creating the window, **before** any panes exist, so it has no effect. Apply the layout after the panes are split. | `tmux.go:56-64` vs `script.go:39` |
| 1.9 | On a fresh session the script never attaches (the session is `-d` and nothing calls `attach` at the end). It also uses `tmux attach` even when already inside tmux, which fails. Use `switch-client` when `$TMUX` is set. | `helpers.go:25-30`, `script.go:66-68` |
| 1.10 | `~` in `Directory` isn't expanded inside double quotes, so `SESSION_DIR="~/x"` fails the existence check. Expand `~` in Go before saving. This matters more now that the TUI is the main entry point, because no shell expands `~` there. | `script.go:46`, `helpers.go:18` |
| 1.11 | Panes split without `-c "$SESSION_DIR"`, so they may open in the wrong directory depending on tmux config. | `tmux.go:92` |

**Recommendation:** switch to `text/template` with a quoting helper (`{{ q .Command }}`), and add golden-file tests that run `bash -n` on the output. Alternatively, skip script generation entirely and run `tmux` directly with `exec.Command("tmux", args...)`. That removes the whole quoting problem, because arguments never pass through a shell. Also consider **targeting windows by name** (`-t "$S:editor"`) instead of by number, which sidesteps the `base-index` problem in 5.2.

---

## 2. CLI

| # | Issue | Location |
|---|-------|----------|
| 2.1 | Every `create` subcommand is now a stub returning `nil`, but the flags are still marked required and the help still shows working examples. `gomux create session --name dev --dir /tmp` **exits 0 and does nothing**. Until the redesign, set `Hidden: true` on `createCmd` or return `errors.New("not implemented yet")`. | `cmd/create.go:43-131` |
| 2.2 | `PersistentPreRun` still builds an `actionsInstance` that nothing uses. | `cmd/create.go:33-39` |
| 2.3 | `gomux create all ...` appears in the help text but doesn't exist. | `cmd/create.go:31-32` |
| 2.4 | `gomux list` prints `list called`. It's still the Cobra scaffold. | `cmd/list.go:17` |

**Direction:** the new `NewScript(session, windows, start, panes)` shape fits a one-shot CLI much better than the old builder did. The natural next step is `gomux create -f spec.yaml` (or `gomux new <name>`, which opens `$EDITOR` on a template spec), feeding straight into `NewScript`.

---

## 3. TUI

| # | Issue | Location |
|---|-------|----------|
| 3.1 | **Crash on save** **[verified]**. `saveScript` passes empty `models.Session{}`, `models.Window{}` and `models.Pane{}` to `NewScript`. That returns `(nil, ErrNoSessionNameFound)`, the `if err != nil { // handle error }` block is empty, and the returned command does `*script`, which is a nil dereference. Fill in name and content, press `ctrl+s`, and the program crashes. At minimum: `if err != nil { return func() tea.Msg { return shared.NewError(err) } }`. | `tui/screens/create_script.go:210-228` |
| 3.2 | **The form doesn't match the model.** It still collects Name, Description, Language, Tags and Content. None of those exist on `models.Script`, and the values are thrown away (`tags` and `language` are computed and never used). The form needs rebuilding around session (name, dir, type), a window list, and per-window panes. A live preview of `Stringify()` would make it far more useful. | `create_script.go:17-31`, `create_script.go:192-213` |
| 3.3 | **Scripts are never loaded from disk.** `AppState.Scripts` starts empty on every launch and only grows through `ScriptCreatedMsg`, so the list shows nothing from previous runs even though `NewScript` now saves to disk. This needs `database.List()`/`Load()` (using `LoadScriptFromJsonBytes`, which already exists) and a load command in the root `Init()`. | `tui/shared/app_state.go:21-46`, `tui/model.go:51-53` |
| 3.4 | **`models.Script` has no identity.** `GetScriptByID` always returns `nil`. The list sends the whole `models.Script` as `ScreenChangeMsg.Data`, but the root does `msg.Data.(string)`, which never matches. Delete sends `ScriptID: ""`. Pick a key, such as `type/name` (already the on-disk path) or a UUID field, and use it everywhere. | `app_state.go:48-62`, `list_scripts.go:67,75,84,110,119,135`, `model.go:83` |
| 3.5 | **The list and detail views render nothing.** The `else` branch for a non-empty list is empty, and the detail view is fully commented out. You see "Scripts (N)" with no rows. | `list_scripts.go:189-193`, `list_scripts.go:206-306` |
| 3.6 | **Run Script**, **Settings** and **edit** still lead to a blank or "Unknown screen" view where only `ctrl+c` works. **Edit Script** in the main menu does nothing, and **delete** is a no-op. | `model.go:137-187`, `run_script.go:31-33`, `main_menu.go:108-133` |
| 3.7 | Validation errors still never show. `saveScript` has a value receiver, so `m.err = ...` sets a copy. | `create_script.go:180-190` |
| 3.8 | The create form is still built once and never reset, so the second "Create" shows the previous values. | `model.go:43`, `model.go:99-103` |
| 3.9 | Main menu filtering is still broken. `enter`/space/`q` are intercepted before the list sees them, so typing `q` in the filter box quits the app. New: `View()` now sets `Title`, `Styles`, `SetShowStatusBar` and `SetFilteringEnabled` **on a copy, every frame**. Move that into `NewMainMenuModel`. | `main_menu.go:91-106`, `main_menu.go:135-141` |
| 3.10 | `WindowSizeMsg` still only reaches the main menu (it now accounts for `PageStyle` margins, which is good). | `model.go:59-64` |
| 3.11 | `ErrorMsg` is stored in `state.LastError` but never rendered. Child `Init()` is never called, so the cursor blink never starts. | `model.go:132-134`, `model.go:51-53` |
| 3.12 | `CreateScriptModel` builds its own `actions.NewActions()` (and so its own `DB`). Create one in the root and pass it in, so all screens share it and tests can substitute a fake. | `create_script.go:76` |

---

## 4. Medium: storage (`internal/database`) _(unchanged)_

| # | Issue | Location |
|---|-------|----------|
| 4.1 | The JSON state file is saved with a **`.tmux.sh`** extension (the same filename as the script). The doc comment says `.tmux.json`. | `data.go:17-20` |
| 4.2 | **Path traversal and unsanitized filenames.** `Session.Name` and `Session.Type` go straight into `filepath.Join`, and the result is written with mode `0755`. A name like `../../../bin/foo` writes an executable outside the data dir. This matters more now that the TUI takes free-text names. | `data.go:17-29`, `tmux.go:9,11` |
| 4.3 | Saving silently **overwrites** an existing script with the same name. | `data.go:21`, `data.go:29` |
| 4.4 | The JSON and `.sh` are written non-atomically. If the second write fails, they drift apart. | `data.go:21-32` |
| 4.5 | `InitDB` and `NewDB` duplicate the path logic in different ways (string concatenation vs `path.Join`). `InitDB` panics on error; `NewDB` ignores errors. | `db.go:40-75` |
| 4.6 | It uses `path` (for URL-style slash paths) instead of `path/filepath` for OS paths. | `db.go:7,66,78` |
| 4.7 | `~/.local/share` is hard-coded, so `$XDG_DATA_HOME` is ignored. `GOMUX_CONFIG_PATH` and `GOMUX_CACHE_PATH` are unused. | `db.go:10-14` |
| 4.8 | `InitDB()` runs before **every** command, even `--help`, creating directories as a side effect. | `main.go:9` |
| 4.9 | No read side: there is no `List`, `Get` or `Delete` (which 3.3 and 3.4 need). | `internal/database/` |

---

## 5. `pkg/actions` and `gomuxerr`

| # | Issue | Location |
|---|-------|----------|
| 5.1 | **Defaults are silently lost** **[verified]**. `for _, window := range windows { window.Layout = ... }` sets the default on a **copy**; same for `pane.Direction`. A test with no layout or direction set produced `layout=""` and `direction=""`. Use `for i := range windows { if windows[i].Layout == "" { windows[i].Layout = ... } }`. | `pkg/actions/actions.go:50-58`, `actions.go:64-71` |
| 5.2 | **Inconsistent window numbering.** `startingWindow <= 0` becomes `1`, and `pane.Location == 0` is rejected as "required" **[verified]**. That assumes `base-index 1`, but tmux defaults to `0`, and the generator targets `$S:0` for the first window. As written, you can't put a pane on the first window of a default tmux setup. Choose one convention: positional indices in the spec with the renderer mapping them, or window names as targets. Use a pointer or `-1` if you need "unset". | `actions.go:60-62`, `actions.go:68-70` |
| 5.3 | **Validation lives in the wrong place.** `NewScript` both validates and saves, so loading a spec from disk or a future `edit` path can't reuse the rules. Move them to `func (s *Script) Validate() error` in `models` and have `actions` call it. Also, "NewScript" reads like a constructor; `CreateScript` or `SaveScript` would say what it does. | `actions.go:29-80` |
| 5.4 | Still under `pkg/` while exposing `internal/models` types, so outside importers can't use it. Move to `internal/`. `IActions`: the `I` prefix isn't idiomatic, and the interface is defined next to its only implementation. | `pkg/actions/actions.go` |
| 5.5 | `gomuxerr`: sentinel errors are a good choice and work with `errors.Is`. Small notes: Go error strings should be lowercase without punctuation (`"session name is required"`); `ErrNo…Found` reads like a lookup failure (`ErrSessionNameRequired` is clearer); and the custom `GomuxError` type adds nothing over `errors.New` until it carries extra fields like `Field`. | `internal/gomuxerr/errors.go` |
| 5.6 | `models` still mixes data with bash rendering, and storage still calls `Stringify()`. | `internal/models/tmux.go`, `internal/database/data.go:26` |
| 5.7 | `Panes` is still a flat list on `Script` pointing at windows by `Location int`. Nesting as `Window{ Panes []Pane }` would make 5.2 disappear and invalid states impossible. | `internal/models/script.go:13` |

---

## 6. Low: dead code, scaffolding, polish

- **Empty files:** `pkg/actions/{create,delete,get,list,update}.go` contain only `package actions`.
- **Unused or stub code:**
  - `AppState`: `GetScriptByID`, `UpdateScript`, `DeleteScript` and `FilterScripts` are stubs returning `nil`/`false`.
  - `AddExecutionResult` has an empty `if` body.
  - `hasAllTags`, `GetRecentExecutions`, `LoadScriptFromJsonBytes` and `NewSession` are never called.
  - `formatTime` is only called from commented-out code.
  - `commenui`: `HeaderStyle`, `SuccessStyle`, `QuitInstruction`, `RefreshInstruction`, `SelectInstuction`, and the up/down/left/right instructions are unused.
  - `shared`: `ConfirmActionMsg`, `StatusMsg`, `LoadingMsg`, `RefreshMsg`, `SearchQueryMsg` and `FilterTagsMsg` are unused.
- **Template-shaped state:** `AppState.Settings` (`Editor`, `TabSize`, `ShowLineNumbers`, `AutoSave`, `DefaultLanguage`, `ScriptDirectory: "./scripts"`), `ExecutionHistory`, `FilterTags`, `SearchQuery` and `IsLoading` still don't apply to a tmux session manager. Now that `shared.Script` is gone, these are the last leftovers from the generic template.
- **Large blocks of commented-out code** in `tui/model.go`, `cmd/create.go`, and most of `list_scripts.go`. Git already keeps that history.
- **Typos:** package `commenui` (→ `commonui`), `SelectInstuction`/`UpInstuction`/`DownInstuction`, "direcotories", "sanitaize".
- **`cmd/root.go`:**
  - `len(os.Args) == 1` bypasses Cobra.
  - `log.Fatalln` followed by `os.Exit(1)`: the exit is unreachable.
  - `debug.log` is always written to the current working directory. Gate it behind something like `GOMUX_DEBUG=1`.
- `RunScriptModel` sets `cursor: 1` for a 1-element slice.
- The list screen still declares its styles inline; move them to `commenui` like the other screens.
- No README, no tests, no CI.

---

## 7. Architecture assessment (updated)

**What improved in `854153b`**
- **One domain type.** The TUI now speaks `models.Script`, and the generic `shared.Script` is gone. This was the biggest structural problem in the first pass, and it's now resolved at the type level.
- **One entry point.** The TUI goes through `actions` → `database` instead of keeping its own in-memory world.
- **Stateless service.** Replacing the per-process builder with an all-in-one `NewScript` fits both a one-shot CLI and the TUI.
- **Shared UI kit.** The `commenui` styles and help helpers are taking shape.

**What still needs work**
1. **The wiring is only half done.** The types line up, but data still doesn't flow through. The create form doesn't produce a `Script`, the list doesn't read from disk, and nothing has an ID. Next steps in order: an ID scheme, then `database.List`, then a form that builds a real `Script`.
2. **Rules are spread across layers.** Validation sits in `actions`, defaults are (meant to be) in `actions`, numbering assumptions are in `actions` and the generator (and they disagree), and rendering is in `models`. Put validation and defaults on the model, and give rendering its own package.
3. **Bash is assembled by string concatenation.** That is still the most broken part of the code and still has no tests.

**Suggested shape** (unchanged, but closer now):

```
main.go
cmd/                  thin cobra commands: new, list, run, edit, delete, (root -> TUI)
internal/
  config/             XDG paths, settings
  layout/             domain: Session{Name, Dir, Windows[]{Name, Cmd, Layout, Panes[]{Split, Cmd}}}
                      + Validate() + ApplyDefaults()   (single source of truth for rules)
  store/              Load/Save/List/Delete specs; ID = "<type>/<name>"; sanitize names; no rendering
  render/             layout -> bash via text/template + shell quoting   (golden tests + `bash -n`)
  tmux/               optional: apply layout directly via exec.Command("tmux", ...)
  service/            (today's actions) orchestration only: validate -> store -> render/run
  tui/                screens receive the service from the root model
```

---

## Suggested order of work

1. **Stop the crash (3.1)** by handling the `NewScript` error. This takes about two minutes.
2. **Fix the generator output (1.1–1.3)** and add a golden test with `bash -n`.
3. **Fix the loop-copy defaults (5.1)** and settle the window-numbering convention (5.2, 1.4); targeting windows by name is the simplest route.
4. **Add identity and a read path:** an ID scheme (3.4), `database.List`/`Load` (4.9), and loading on TUI start (3.3). Then render the list (3.5).
5. **Rebuild the create form** around session, windows and panes (3.2), ideally with a live script preview.
6. **Sanitize names (4.2)**, fix the JSON extension (4.1), and move validation onto the model (5.3).
7. **Clean up:** hide or error the stubbed CLI (2.1), delete the leftover template state and commented code, move `pkg/actions` → `internal/`, read XDG paths, and gate the debug log.
