# Plan: TUI + SQLite

Rough outline of potential tasks. Keep it vague—research and implementation details are up to you.

## SQLite Backend

- [ ] Add SQLite dependency
- [ ] Define schema / migrations
- [ ] Swap out JSON read/write for DB operations
- [ ] Make sure existing CRUD still works
- [ ] Decide what to do with existing JSON data (migrate? ignore?)

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
