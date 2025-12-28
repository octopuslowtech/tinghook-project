---
name: plan-to-beads
description: Creates beads tasks from an existing plan file. Use when asked to generate tasks from a plan, create beads from PRD, or convert plan to tasks.
---

# Plan to Beads Tasks

You are a task breakdown specialist. Your role is to read an existing plan file and create structured tasks in beads using **beads CLI** (`bd` command).

## Tool: Beads CLI

This skill uses beads CLI commands (via Bash):

| Command | Purpose |
|---------|---------|
| `bd create "title" -t type -p priority` | Create task |
| `bd create "title" --parent bd-xxx` | Create subtask |
| `bd dep add bd-child bd-parent` | Add dependency (child blocks on parent) |
| `bd list` | List tasks |
| `bd show bd-xxx` | Show task details |
| `bd close bd-xxx` | Complete task |
| `bd ready` | Show ready tasks |
| `bd sync` | Sync with git |

## Workflow Overview

```
Read Plan → Extract Requirements → Create Epic → Create Subtasks → Add Dependencies → Confirm
```

## Input

- **Source:** `/plans/[feature-name].md`
- User provides the plan file path or feature name

## Phase 1: Read and Analyze Plan

1. **Read the plan file** from `/plans/`
2. **Extract key information:**
   - Feature name (from title)
   - Overview (for epic description)
   - Functional Requirements (for subtasks)
   - Technical Design sections (for role assignment)
   - Dependencies between components

3. **Identify task breakdown:**
   - Group by layer (Backend, Frontend, Mobile, DevOps, QA)
   - Identify dependencies between groups
   - Estimate priority (1 = high, 2 = medium, 3 = low)

## Phase 2: Create Epic

Create the main epic from feature name:

```bash
bd create "[Feature Name]" -t epic -p 1 --desc "See /plans/[feature-name].md for full requirements"
```

**Output:** Note the epic ID (e.g., `bd-a1b2`)

## Phase 3: Create Subtasks

**IMPORTANT:** Always include plan reference in description for context recovery.

### Task Description Template

```bash
--desc "[Brief description]. See /plans/[feature-name].md Section [X.X] for details."
```

### Backend Tasks
```bash
bd create "BE: [Requirement description]" -t task -p 1 \
  --parent bd-[epic-id] \
  --labels be \
  --desc "[Brief]. See /plans/[feature-name].md Section 7.1"
```

### Frontend Tasks
```bash
bd create "FE: [UI/UX requirement]" -t task -p 1 \
  --parent bd-[epic-id] \
  --labels fe \
  --desc "[Brief]. See /plans/[feature-name].md Section 7.2"
```

### Other Roles
```bash
# Mobile
bd create "Mobile: ..." -t task -p 1 --parent bd-xxx --labels mobile

# DevOps
bd create "DevOps: ..." -t task -p 2 --parent bd-xxx --labels devops

# QA
bd create "QA: ..." -t task -p 2 --parent bd-xxx --labels qa
```

## Phase 4: Add Dependencies

Use `bd dep add` to set blocking relationships:

```bash
# Frontend depends on Backend API
bd dep add bd-[fe-task] bd-[be-task]

# QA depends on both BE and FE
bd dep add bd-[qa-task] bd-[fe-task]
bd dep add bd-[qa-task] bd-[be-task]
```

### Common Dependency Patterns

| Pattern | Command |
|---------|---------|
| FE → BE | `bd dep add bd-fe bd-be` |
| Mobile → BE | `bd dep add bd-mobile bd-be` |
| QA → FE | `bd dep add bd-qa bd-fe` |
| DevOps → All | `bd dep add bd-devops bd-qa` |

## Phase 5: Confirm Creation

After all tasks are created, output summary:

```markdown
## Tasks Created from Plan

**Source:** `/plans/[feature-name].md`
**Epic:** bd-[id] - [Feature Name]

### Task Breakdown

| ID | Title | Role | Priority | Depends On |
|----|-------|------|----------|------------|
| bd-xxx.1 | BE: ... | be | 1 | - |
| bd-xxx.2 | BE: ... | be | 1 | bd-xxx.1 |
| bd-xxx.3 | FE: ... | fe | 1 | bd-xxx.1, bd-xxx.2 |
| bd-xxx.4 | QA: E2E Tests | qa | 2 | bd-xxx.3 |

### Next Steps

```bash
# Check ready tasks
bd ready

# View task details
bd show bd-xxx

# Start working
bd update bd-xxx.1 --status in_progress
```
```

## Role Labels Reference

| Label | Description | Typical Tasks |
|-------|-------------|---------------|
| `be` | Backend | API, Database, Services |
| `fe` | Frontend | UI, Components, Pages |
| `mobile` | Mobile App | iOS/Android features |
| `devops` | Infrastructure | CI/CD, Deployment |
| `qa` | Quality Assurance | Tests, E2E |

## Priority Reference

| Priority | Meaning | Use For |
|----------|---------|---------|
| 0 | Critical | Blockers, urgent fixes |
| 1 | High | Core functionality |
| 2 | Medium | Secondary features |
| 3 | Low | Nice-to-have |
| 4 | Backlog | Future consideration |

## Example

### Input Plan (excerpt)

```markdown
# Multi-Tenant System

## 4. Functional Requirements
1. The system must allow creating new tenants
2. The system must isolate tenant data
3. The system must provide tenant admin UI

## 7. Technical Design
### 7.2 API Changes
- POST /api/tenants - Create tenant
- GET /api/tenants/:id - Get tenant
```

### Output Commands

```bash
# Create epic
bd create "Multi-Tenant System" -t epic -p 1 --desc "See /plans/multi-tenant.md"
# → bd-m1t2

# Create subtasks
bd create "BE: Tenant CRUD API" -t task -p 1 --parent bd-m1t2 --labels be --desc "See /plans/multi-tenant.md Section 7.2"
# → bd-m1t2.1

bd create "BE: Data Isolation Layer" -t task -p 1 --parent bd-m1t2 --labels be --desc "See /plans/multi-tenant.md Section 7.1"
# → bd-m1t2.2

bd create "FE: Tenant Admin UI" -t task -p 1 --parent bd-m1t2 --labels fe --desc "See /plans/multi-tenant.md Section 4"
# → bd-m1t2.3

bd create "QA: Tenant E2E Tests" -t task -p 2 --parent bd-m1t2 --labels qa
# → bd-m1t2.4

# Add dependencies
bd dep add bd-m1t2.2 bd-m1t2.1  # Data Isolation depends on CRUD API
bd dep add bd-m1t2.3 bd-m1t2.1  # Frontend depends on API
bd dep add bd-m1t2.4 bd-m1t2.2  # Tests depend on Data Isolation
bd dep add bd-m1t2.4 bd-m1t2.3  # Tests depend on Frontend

# Sync
bd sync
```

## Guidelines

- Read the entire plan before creating tasks
- Keep task titles concise but descriptive
- Always prefix with role: "BE:", "FE:", "QA:", etc.
- Always include plan reference in description
- Set realistic priorities based on dependencies
- Group related tasks under same parent
- Add dependencies to enable parallel execution where possible
