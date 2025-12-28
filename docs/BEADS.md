# Beads CLI Guide

> Reference: [beads](https://github.com/steveyegge/beads)

## Quick Start

```bash
# Initialize beads in project
bd init

# Create a task
bd create "Fix login bug" -t task -p 1

# List tasks
bd list
bd ready  # Show ready tasks (no blockers)

# Work on task
bd update bd-xxx --status in_progress

# Complete task
bd close bd-xxx --reason "Fixed"

# Sync with git
bd sync
```

## Common Commands

### Task Creation

```bash
# Create task
bd create "Title" -t task -p 1

# Create with parent (subtask)
bd create "Subtask" -t task -p 1 --parent bd-xxx

# Create with label
bd create "BE: API endpoint" -t task -p 1 --labels be

# Create with description
bd create "Title" -t task -p 1 --desc "Details here"

# Create epic
bd create "Feature Name" -t epic -p 1
```

### Task Types

| Type | Use For |
|------|---------|
| `task` | Regular work item |
| `epic` | Large feature (parent) |
| `bug` | Bug fix |
| `feature` | New feature |
| `chore` | Maintenance |

### Priority Levels

| Priority | Meaning |
|----------|---------|
| 0 | Critical |
| 1 | High |
| 2 | Medium |
| 3 | Low |
| 4 | Backlog |

### Dependencies

```bash
# Add dependency (child blocks on parent)
bd dep add bd-child bd-parent

# Example: Frontend depends on Backend
bd dep add bd-fe-task bd-be-task

# View dependencies
bd show bd-xxx
```

### Status Management

```bash
# Start working
bd update bd-xxx --status in_progress

# Complete
bd close bd-xxx --reason "Done"

# List by status
bd list --status open
bd list --status in_progress
bd ready  # Ready tasks (no blockers)
```

### Listing & Filtering

```bash
# All tasks
bd list

# By parent
bd list --parent bd-epic

# By label
bd list --labels be

# Ready to work
bd ready
```

## Integration with Skills

### Workflow

```
feature-planning → plan-to-beads → let-him-cook
     (PRD)          (bd create)     (execute)
```

### plan-to-beads creates:
```bash
bd create "Feature" -t epic -p 1
bd create "BE: Task 1" --parent bd-xxx --labels be
bd create "FE: Task 2" --parent bd-xxx --labels fe
bd dep add bd-fe bd-be
```

### let-him-cook uses:
```bash
bd ready                           # Find work
bd update bd-xxx --status in_progress  # Claim
# ... implement ...
bd close bd-xxx --reason "Done"    # Complete
bd sync                            # Sync
```

## Role Labels

| Label | Description |
|-------|-------------|
| `be` | Backend |
| `fe` | Frontend |
| `mobile` | Mobile |
| `devops` | DevOps |
| `qa` | QA/Testing |

## Session Checklist

### Start Session
- [ ] `bd ready` - Find available work
- [ ] `bd update bd-xxx --status in_progress` - Claim task

### During Work
- [ ] Implement the task
- [ ] Run tests
- [ ] Commit changes

### End Session
- [ ] `bd close bd-xxx --reason "..."` - Complete task
- [ ] `bd sync` - Sync with git
- [ ] `git push` - Push changes
