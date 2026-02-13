# Plan: TUI + SQLite

Rough outline of potential tasks. Keep it vague—research and implementation details are up to you.

## SQLite Backend

- [x] Add SQLite dependency
- [x] Define schema / migrations
- [x] Swap out JSON read/write for DB operations
- [x] Make sure existing CRUD still works
- [x] Custom ID (hex string) instead of incremental
- [x] FindShortId for prefix matching (Update, Delete, Toggle)
- [x] Decide what to do with existing JSON data (migrate? ignore?)

## TUI (Bubble Tea)

- [ ] Add Bubble Tea dependency
- [ ] Learn the basic model/update/view pattern
- [ ] Display list of todos
- [ ] Handle keyboard input (navigation, actions)
- [ ] Wire up create / toggle / delete from TUI
- [ ] Figure out how to combine TUI mode with existing CLI flags (separate subcommand? default when no flags?)

## Integration

- [ ] Both CLI and TUI use same SQLite store
- [ ] Test the flow end-to-end

---

## Concurrency & Channels (later, with TUI)

_These fit naturally when you add the TUI—Bubble Tea is already concurrent._

- [ ] **Async commands in Bubble Tea**: When the TUI needs to fetch todos from SQLite, run the DB call in a goroutine so the UI doesn't freeze. Use `tea.Cmd` to spawn it and a channel to send the result back to `Update`. This is the idiomatic way to do I/O in Bubble Tea.
- [ ] **Context for cancellation**: Use `context.Context` when doing longer ops (e.g. future sync/export). Lets you cancel cleanly if the user exits.
- [ ] _(Optional)_ **Single-writer store**: One goroutine owns DB access, others send read/write requests via a channel. Classic Go pattern, but likely overkill for SQLite—only consider if you hit contention or need a clear "all access goes through one place" design.
