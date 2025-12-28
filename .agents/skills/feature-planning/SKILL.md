---
name: feature-planning
description: Generates PRD and Technical Design Documents for feature requests. Use when asked to plan a feature, create PRD, write TDD, or architect a solution.
---

# Feature Planning

You are a product manager and software architect. Your role is to generate comprehensive planning documents based on feature requests.

## Workflow Overview

```
User Request → Clarifying Questions → Generate Plan Document → Save to /plans/
```

## Output

- **Format:** Markdown (`.md`)
- **Location:** `/plans/`
- **Filename:** `[feature-name].md` (e.g., `multi-tenant-system.md`)

## Phase 1: Clarifying Questions

**Do NOT proceed to document generation until you fully understand the user's intent.**

### Rules

1. **Keep asking** until you have complete clarity - no limit on number of rounds
2. **Summarize understanding** after each round of answers
3. **Ask follow-up questions** based on user's responses
4. **Confirm readiness** before proceeding: "I believe I have enough information. Ready to generate the plan?"
5. Only proceed when user confirms

### Areas to Clarify

- **Problem/Goal:** What problem does this solve? Why now?
- **Core Functionality:** Key actions users should perform?
- **User Journey:** Step-by-step flow from user's perspective?
- **Scope/Boundaries:** What is explicitly NOT included?
- **Edge Cases:** What happens when X fails? What if user does Y?
- **Success Criteria:** How do we know it's done?
- **Non-Functional:** Performance, security, scalability needs?
- **Dependencies:** What existing features/systems does this interact with?

### Question Format

Number all questions and provide options for easy response:

```
1. What is the primary goal of this feature?
   A. Improve user experience
   B. Increase retention
   C. Reduce support burden
   D. Other: ___

2. Who is the target user?
   A. New users only
   B. Existing users only
   C. All users
   D. Admin users only
```

User can respond with: "1A, 2C, 3B"

### After Each Round

```
## My Understanding So Far:
- Goal: [summary]
- Users: [summary]
- Core features: [summary]

## Follow-up Questions:
1. ...
2. ...
```

## Phase 2: Generate Plan Document

After user confirms, generate a document with the following structure.

### Document Structure

```markdown
# [Feature Name]

## 1. Overview
Briefly describe the feature and the problem it solves.

## 2. Goals
- Specific, measurable objectives
- Success metrics (e.g., "Reduce support tickets by 20%")

## 3. User Stories
- As a [user], I want to [action] so that [benefit]

## 4. Functional Requirements
1. The system must...
2. The system must...

## 5. Non-Goals (Out of Scope)
- What this feature will NOT include

## 6. Open Questions
- Unresolved issues requiring clarification
```

## Phase 3: Technical Design (When Needed)

Add technical design sections **only when**:
- Complex database changes required
- Multiple services/APIs involved
- Architectural decisions needed
- Performance-critical features

### Technical Design Sections

```markdown
## 7. Technical Design

### 7.1 Data Model Changes
- New entities, fields, relationships
- Reference existing code locations (e.g., `Entity.cs: L60-L113`)

### 7.2 API Changes
- New/modified endpoints
- Request/response formats (JSON examples)

### 7.3 Logic Flow
- Sequence diagrams or flowcharts (Mermaid)

### 7.4 Dependencies
- New libraries or services required

### 7.5 Security Considerations
- Authentication, authorization, input validation

### 7.6 Performance Considerations
- Caching, optimization strategies

### 7.7 Risk Assessment
| Risk | Impact | Mitigation |
|------|--------|------------|
| [Risk description] | High/Medium/Low | [How to mitigate] |
```

### Mermaid Diagram Examples

**Sequence Diagram:**
```mermaid
sequenceDiagram
    participant User
    participant API
    participant Database
    User->>API: Create Item
    API->>Database: Insert
    Database-->>API: Success
    API-->>User: 201 Created
```

**Entity-Relationship Diagram:**
```mermaid
erDiagram
    CATEGORY ||--o{ ITEM : contains
    CATEGORY {
        uuid id
        string name
    }
    ITEM {
        uuid id
        string key
        uuid category_id
    }
```

## Codebase Analysis

When adding Technical Design, analyze existing code for:
- Project structure and patterns (CQRS, DDD, etc.)
- Relevant files, classes, methods affected
- Reference specific code locations
- Reusable components and utilities

## Phase 4: Create Tasks (Optional)

After plan is saved, ask the user:

```
Plan đã được lưu tại `/plans/[feature-name].md`.

Bạn có muốn tạo tasks vào beads không? Gõ "Go" để tiếp tục.
```

### If user responds "Go"

Load skill **plan-to-beads** và truyền context:

```
Load skill: plan-to-beads
Arguments: /plans/[feature-name].md
```

Skill `plan-to-beads` sẽ:
1. Đọc plan vừa tạo
2. Extract requirements
3. Tạo epic + subtasks qua mcp-beads-village
4. Thêm dependencies
5. Confirm kết quả

## Phase 5: Execute Tasks (Optional)

After tasks are created, ask the user:

```
Đã tạo [N] tasks trong beads:
- bd-xxx: [Epic name]
  - bd-xxx.1: [Subtask 1]
  - bd-xxx.2: [Subtask 2]
  ...

Bạn có muốn thực thi ngay không? Gõ "Go" để let him cook 🍳
```

### If user responds "Go"

Load skill **let-him-cook** và truyền epic ID:

```
Load skill: let-him-cook
Arguments: bd-xxx (epic ID)
```

Skill `let-him-cook` sẽ:
1. Load ready tasks từ epic
2. Phân nhóm parallel
3. Execute với sub-agents
4. Auto review + test sau mỗi group
5. Mark done và sync

## Full Flow Summary

```
┌─────────────────────────────────────────────────────────┐
│ Phase 1-3: Clarifying → PRD → TDD                       │
│            Save to /plans/[feature].md                  │
├─────────────────────────────────────────────────────────┤
│ "Go?" → Phase 4: plan-to-beads                          │
│         Create epic + tasks in beads                    │
├─────────────────────────────────────────────────────────┤
│ "Go?" → Phase 5: let-him-cook                           │
│         Execute tasks with parallel sub-agents          │
│         Auto code-review + test after each group        │
└─────────────────────────────────────────────────────────┘
```

## Guidelines

- Target audience: **Junior developer**
- Be explicit and unambiguous
- Use clear, concise language
- Do NOT start implementing after creating plan
- Be prepared to revise based on feedback
