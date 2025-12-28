# Skills Catalog

A collection of specialized skills for AI-assisted development workflow.

## Available Skills

| Skill | Purpose | Trigger |
|-------|---------|---------|
| [feature-planning](#feature-planning) | Create PRD + Technical Design | "Plan feature X" |
| [plan-to-beads](#plan-to-beads) | Generate beads tasks from plan | "Create tasks from plan" |
| [let-him-cook](#let-him-cook) | Execute tasks with sub-agents | "Let him cook" |
| [shake-it](#shake-it) | Run tests and coverage | "Shake it" |
| [roast-it](#roast-it) | Code review and quality check | "Roast it" |
| [ship-it](#ship-it) | Git commit and push | "Ship it" |

---

## feature-planning 📋

**Purpose:** Generate Product Requirements Document (PRD) and Technical Design Document (TDD) for new features.

**Triggers:**
```
"Plan a feature for user authentication"
"Create PRD for multi-tenant system"
"Design the payment integration feature"
"I need to architect a new notification system"
```

**Output:** `/plans/[feature-name].md`

**Flow:**
```
Clarifying Questions → PRD → Technical Design (optional) → Save
        ↓ "Go"
   plan-to-beads
        ↓ "Go"
   let-him-cook
```

---

## plan-to-beads 📦

**Purpose:** Convert a plan file into beads tasks with epic, subtasks, and dependencies.

**Triggers:**
```
"Create tasks from /plans/multi-tenant.md"
"Generate beads tasks for this plan"
"Convert plan to tasks"
"Break down this plan into tasks"
```

**Output:** Epic + subtasks in beads

**Example Output:**
```
bd-a1b2: Multi-Tenant System (epic)
├── bd-a1b2.1: BE: Create Tenant entity [be]
├── bd-a1b2.2: BE: Implement TenantContext [be]
├── bd-a1b2.3: FE: Tenant Admin UI [fe]
└── bd-a1b2.4: QA: E2E Tests [qa]
```

---

## let-him-cook 🍳

**Purpose:** Execute beads tasks using parallel sub-agents with QA gates between groups.

**Triggers:**
```
"Let him cook"
"Cook bd-a1b2"
"Execute the tasks"
"Start implementing"
"Cook backend only"
```

**Flow:**
```
Group 1 (parallel) → BUILD → shake-it → roast-it → Gate ✅
                                                      ↓
Group 2 (parallel) → BUILD → shake-it → roast-it → Gate ✅
                                                      ↓
                                                   Done!
```

**Key Features:**
- Parallel task execution within groups
- Sequential QA gates between groups
- Auto-stop on build/test failures

---

## shake-it 🫨

**Purpose:** Run tests, check coverage, validate builds. "Shake the code to see if it breaks."

**Triggers:**
```
"Shake it"
"Run the tests"
"Check if tests pass"
"Validate the build"
"Check coverage"
```

**Supported Languages:**
- JavaScript/TypeScript: `npm test`, `pnpm test`
- Go: `go test ./... -cover`
- .NET: `dotnet test`
- Python: `pytest --cov`
- Rust: `cargo test`

**Output:**
```
## 🫨 Shake Report

| Metric | Value |
|--------|-------|
| Tests  | 45 passed, 0 failed |
| Coverage | 87% |
| Build | ✅ Success |
```

---

## roast-it 🔥

**Purpose:** Code review for quality, security, and performance. "Roast the code like a pro."

**Triggers:**
```
"Roast it"
"Review the code"
"Check code quality"
"Security review"
"What's wrong with this code?"
```

**What It Checks:**
- Code quality and best practices
- Security vulnerabilities (OWASP Top 10)
- Performance issues
- Error handling
- Type safety

**Output:**
```
## 🔥 Roast Report

### Critical Issues
- SQL injection in auth.ts:45

### High Priority
- Missing null check in user.ts:23

### Recommendations
1. Add input validation
2. Use parameterized queries
```

---

## ship-it 🚀

**Purpose:** Git operations - commit, push, and create PRs.

**Triggers:**
```
"Ship it"
"Ship bd-a1b2.1"
"Commit and push"
"Create PR"
"Commit only"
```

**Flow (default):**
```
Stage → Commit → Pull --rebase → Push → Sync beads
```

**Commit Format:**
```
feat(tenant): add multi-tenant support

- Add Tenant entity
- Implement TenantContext
- Add middleware

Closes bd-a1b2.1
```

---

## Complete Workflow Example

```
User: "Plan a feature for user authentication"
      → feature-planning creates /plans/user-auth.md

User: "Go"
      → plan-to-beads creates tasks in beads

User: "Go" 
      → let-him-cook executes tasks:
          Group 1 → BUILD → shake-it → roast-it → ✅
          Group 2 → BUILD → shake-it → roast-it → ✅

User: "Ship it"
      → ship-it commits and pushes all changes
```

---

## Quick Reference

| I want to... | Say... |
|--------------|--------|
| Plan a new feature | "Plan feature X" |
| Create tasks from plan | "Go" or "Create tasks" |
| Execute tasks | "Let him cook" |
| Run tests | "Shake it" |
| Review code | "Roast it" |
| Commit and push | "Ship it" |
| Commit without push | "Commit only" |
| Execute specific epic | "Cook bd-xxx" |
| Execute only backend | "Cook backend only" |
