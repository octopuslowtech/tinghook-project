---
name: let-him-cook
description: Executes beads tasks using parallel sub-agents. Use when asked to execute tasks, run tasks, let him cook, or implement from beads.
---

# Let Him Cook 🍳

You are a task execution orchestrator. Your role is to load ready tasks from beads and execute them **one group at a time** with QA gates between groups.

## Workflow Overview

```
Load Tasks → Group by Dependencies → Execute Group 1 → QA Gate → Execute Group 2 → QA Gate → ...
```

## CRITICAL RULE: Sequential Group Execution

**DO NOT execute the next group until the current group passes QA.**

```
┌─────────────────────────────────────────────────────────────┐
│ ❌ WRONG: Execute all groups, then QA                       │
│    Group 1 → Group 2 → Group 3 → QA                        │
│    Problem: Errors cascade, hard to debug                   │
├─────────────────────────────────────────────────────────────┤
│ ✅ CORRECT: QA after each group                             │
│    Group 1 → QA ✅ → Group 2 → QA ✅ → Group 3 → QA ✅      │
│    Benefit: Catch errors early, easy to fix                 │
└─────────────────────────────────────────────────────────────┘
```

## Phase 1: Load Tasks from Beads

### Get Ready Tasks

Use mcp-beads-village to get tasks ready for execution:

```
ls(status="ready")
```

Or get all open tasks:

```
ls(status="open")
```

### Get Epic Details (if specified)

If user specifies an epic:

```
show(id="bd-[epic-id]")
```

### Output Task List

```markdown
## Ready Tasks

| ID | Title | Role | Priority | Blocked By |
|----|-------|------|----------|------------|
| bd-xxx.1 | BE: Tenant CRUD API | be | 1 | - |
| bd-xxx.2 | BE: Data Isolation | be | 1 | bd-xxx.1 |
| bd-xxx.3 | FE: Tenant Admin UI | fe | 1 | bd-xxx.1 |

**Ready to execute:** bd-xxx.1 (no blockers)
```

## Phase 2: Identify Parallel Groups

Analyze dependencies to find tasks that can run simultaneously:

### Rules

1. **No blockers** = Can run immediately
2. **Same blocker** = Can run in parallel after blocker completes
3. **Different files** = Safe to parallelize
4. **Same file** = Must run sequentially

### Group Example

```
[PARALLEL GROUP 1 - No Dependencies]
- bd-xxx.1: BE: Tenant CRUD API

[PARALLEL GROUP 2 - After Group 1]
- bd-xxx.2: BE: Data Isolation (blocks: bd-xxx.1)
- bd-xxx.3: FE: Tenant Admin UI (blocks: bd-xxx.1)

[PARALLEL GROUP 3 - After Group 2]
- bd-xxx.4: QA: E2E Tests (blocks: bd-xxx.2, bd-xxx.3)
```

## Phase 3: Execute with Sub-Agents

### Before Starting Each Task

1. **Claim the task:**
   ```
   claim()  # Or specify: claim(id="bd-xxx.1")
   ```

2. **Reserve files** (if known):
   ```
   reserve(paths=["src/api/tenants.ts"], reason="bd-xxx.1")
   ```

### Launch Parallel Sub-Agents

For each parallel group, use the **Task tool** to spawn sub-agents:

```markdown
## Executing Parallel Group 1

Launching 3 sub-agents...
```

**Task prompt template (IMPLEMENT ONLY - NO BUILD/TEST):**

```
Implement task bd-xxx.1: [Task Title]

## Context
- Epic: [Epic name]
- Plan: /plans/[feature-name].md

## Requirements
[Copy relevant section from plan]

## Files to Create/Modify
- Create: src/api/tenants.ts
- Modify: src/api/index.ts

## IMPORTANT
- ONLY write code, do NOT run build/test
- Save files and report back
- Build and test will run AFTER all tasks complete

## When Done
Report back with:
1. Files created/modified
2. Summary of changes
3. Any concerns or blockers
```

### Parallel Execution Rules

| Scenario | Action |
|----------|--------|
| 2+ tasks, no shared files | Launch all in parallel |
| Tasks share files | Run sequentially |
| Max parallel agents | 3-5 (avoid overwhelming) |
| Task fails | Stop group, report error |

### Wait for ALL Sub-Agents

```markdown
## Waiting for Group 1 to complete...

- [x] bd-xxx.1: Done - created 3 files
- [x] bd-xxx.2: Done - modified 2 files  
- [x] bd-xxx.3: Done - created 1 file
- [ ] bd-xxx.4: In progress...

[Wait until ALL complete before proceeding to QA]
```

## Phase 4: Mark Tasks Done

After each sub-agent completes successfully:

```
done(id="bd-xxx.1", msg="Implemented Tenant CRUD API with tests")
```

### Release File Locks

```
release(paths=["src/api/tenants.ts"])
```

### Sync Changes

```
sync()
```

## Phase 5: Review & Test (After Group Complete)

**IMPORTANT:** Only run QA after ALL sub-agents in a group complete.

### QA Flow (Sequential)

```
┌─────────────────────────────────────────────────────────────┐
│ Step 1: BUILD                                               │
│   → If FAIL: Stop, fix build errors                         │
├─────────────────────────────────────────────────────────────┤
│ Step 2: TEST (only if build pass)                           │
│   → If FAIL: Stop, fix failing tests                        │
│   → No point reviewing code that doesn't work               │
├─────────────────────────────────────────────────────────────┤
│ Step 3: REVIEW (only if tests pass)                         │
│   → Check quality, security, performance                    │
│   → If critical issues: Stop and fix                        │
├─────────────────────────────────────────────────────────────┤
│ Step 4: Gate Check → Next Group                             │
└─────────────────────────────────────────────────────────────┘
```

**Why sequential?**
- No point reviewing code that doesn't build
- No point reviewing code that doesn't pass tests
- Faster feedback loop (fail early)

### Step 1: Build

Run project build command:

```bash
# Detect and run appropriate build command
npm run build    # JS/TS
dotnet build     # .NET
cargo build      # Rust
go build ./...   # Go
```

**If build fails:**
```markdown
## ❌ Build Failed

```
[Error message]
```

**Affected tasks:** bd-xxx.1, bd-xxx.2

**Options:**
1. Fix build errors now
2. Rollback and retry
3. Stop execution

What would you like to do?
```

### Step 2: Test (Sub-Agent) - Only if Build Pass

Launch tester sub-agent:

```
Run tests for Group 1 changes.

Files changed:
- src/api/tenants.ts
- src/services/tenant.ts
- src/models/tenant.ts

Load skill: shake-it

Report: pass/fail count, coverage, failures
```

**If tests fail:**
```markdown
## ❌ Tests Failed

**Failed tests:**
- `CreateTenant_ShouldReturnId`: Expected 1, got null
- `GetTenant_NotFound`: Timeout

**Coverage:** 72% (below 80% threshold)

**⚠️ Skipping code review - fix tests first**

**Options:**
1. Fix failing tests now
2. Create bug task in beads
3. Stop execution

Choice:
```

### Step 3: Review (Sub-Agent) - Only if Tests Pass

Launch code-reviewer sub-agent:

```
Review Group 1 code changes.

Tasks completed: bd-xxx.1, bd-xxx.2, bd-xxx.3

Files to review:
- src/api/tenants.ts
- src/services/tenant.ts
- src/models/tenant.ts

Load skill: roast-it

Report: critical/high/medium issues
```

### QA Results Summary

```markdown
## ✅ Group 1 QA Complete

### Build
✅ Build successful (12.3s)

### Tests
- Total: 45 | Passed: 45 | Failed: 0
- Coverage: 87%
- Duration: 8.2s

### Code Review
- Critical: 0
- High: 1 (missing null check in tenant.ts:45)
- Medium: 2

### Action Required?
⚠️ 1 High issue found. Fix before continuing? (y/n)
```

### Handle Issues

**If critical/high issues:**
```markdown
## ⚠️ Issues Found in Group 1

**High Priority:**
1. Missing null check in tenant.ts:45

**Options:**
1. Fix now (spawn fix sub-agent)
2. Create follow-up task in beads
3. Continue anyway (not recommended)
4. Stop execution

Choice:
```

**If all good:**
```markdown
## ✅ Group 1 Complete

- Build: ✅
- Tests: ✅ 45/45
- Review: ✅ No critical issues

Proceeding to Group 2...
```

## Phase 6: Gate Check Before Next Group

### MANDATORY: Wait for QA Pass

**Before proceeding to the next group, ALL conditions must be met:**

```markdown
## Gate Check: Group 1 → Group 2

Prerequisites:
- [x] All Group 1 tasks complete
- [x] Build successful
- [x] Tests passing
- [x] No critical/high review issues (or fixed)

✅ Gate PASSED - Proceeding to Group 2
```

### If Gate Fails

```markdown
## ❌ Gate Check Failed

**Blocking issues:**
- Build failed: 2 errors
- Tests: 3 failing

**Cannot proceed to Group 2 until resolved.**

**Options:**
1. Fix issues now
2. Rollback Group 1 changes
3. Stop execution

Choice:
```

### Execution Loop

```
┌─────────────────────────────────────────────────────────────┐
│ FOR EACH GROUP:                                             │
│                                                              │
│   1. Execute all tasks in group (parallel sub-agents)       │
│   2. Wait for ALL sub-agents to complete                    │
│   3. Run BUILD                                              │
│   4. Run TEST + REVIEW (parallel)                           │
│   5. Gate Check:                                            │
│      - If PASS → Continue to next group                     │
│      - If FAIL → Stop and fix                               │
│                                                              │
│ END FOR                                                      │
└─────────────────────────────────────────────────────────────┘
```

## Phase 7: Continue or Complete

### If More Groups Ready

```markdown
## Group 1 Complete ✅

Moving to Parallel Group 2...
- bd-xxx.4: FE: Tenant UI
- bd-xxx.5: FE: Tenant List

Launching 2 sub-agents...
```

### If All Tasks Done

```markdown
## All Tasks Complete 🎉

**Epic:** bd-xxx - Multi-Tenant System
**Tasks Completed:** 5/5
**Time:** ~2 hours

### Summary
| ID | Title | Status |
|----|-------|--------|
| bd-xxx.1 | BE: Tenant CRUD API | ✅ Done |
| bd-xxx.2 | BE: Data Isolation | ✅ Done |
| bd-xxx.3 | FE: Tenant Admin UI | ✅ Done |
| bd-xxx.4 | FE: Tenant List | ✅ Done |
| bd-xxx.5 | QA: E2E Tests | ✅ Done |

### Next Steps
- Review changes: `git diff main`
- Run full test suite: `npm test`
- Create PR: `gh pr create`
```

### If Task Fails

```markdown
## ❌ Task Failed: bd-xxx.2

**Error:** Tests failing - API endpoint not found

**Options:**
1. Fix manually and continue
2. Skip and continue with other tasks
3. Stop execution

What would you like to do?
```

## Execution Best Practices

### Sub-Agent Prompt Guidelines

1. **Include full context** - Plan section, requirements, files
2. **Specify verification** - Test commands, lint, build
3. **Clear acceptance criteria** - What "done" looks like
4. **Error handling** - What to do if something fails

### File Reservation Strategy

```
# Reserve before work
reserve(paths=["src/api/tenants.ts", "src/api/tenants.test.ts"], reason="bd-xxx.1", ttl=600)

# Release after done
release(paths=["src/api/tenants.ts", "src/api/tenants.test.ts"])
```

### Progress Tracking

Update user after each group:

```markdown
## Progress

- [x] Group 1: Backend Core (2/2 tasks)
- [ ] Group 2: Frontend UI (0/2 tasks) ← Current
- [ ] Group 3: Testing (0/1 tasks)
```

## Quick Commands

| User Says | Action |
|-----------|--------|
| "Let him cook" | Execute all ready tasks |
| "Cook bd-xxx" | Execute specific epic |
| "Cook backend only" | Filter by `tags=["be"]` |
| "Cook 1 task" | Execute single next task |
| "Stop cooking" | Pause execution |

## Safety Guidelines

- **Always reserve files** before editing
- **Run tests** after each group (not after each task)
- **Sync frequently** to avoid conflicts
- **Stop on critical errors** - don't continue blindly
- **Max 5 parallel agents** per group
- **Check diagnostics** after each group
- **NEVER skip QA gate** - must pass before next group
