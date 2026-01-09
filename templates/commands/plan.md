---
description: Create tasks.md from design.md. Break down feature into actionable, dependency-ordered tasks.
handoffs:
  - label: Implement
    agent: implement
    prompt: Start implementation following tasks.md
    send: true
---

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Purpose

Generate `tasks.md` from `design.md` in `.oursky/scopes/FeatureName/`. Create actionable, dependency-ordered tasks organized as a kanban board.

## Execution Flow

1. **Detect active scope and feature directory**:
   - Check if user specified feature name in arguments
   - If not, look for most recent `.oursky/scopes/*/design.md`
   - If multiple exist, ask user which feature
   - Extract feature name from directory:
     - Pattern: `.oursky/scopes/{date-prefix}-{FeatureName}/` (e.g., `20251121-123456-TestCommandFeedbackLoop/`)
     - Extract `FeatureName` by removing date prefix (format: `YYYYMMDD-HHMMSS-`)
     - Example: `20251121-123456-TestCommandFeedbackLoop` → `TestCommandFeedbackLoop`
   - Set `FEATURE_NAME = FeatureName` (PascalCase)
   - Set `FEATURE_DIR = .oursky/scopes/{date-prefix}-{FeatureName}/` (use actual directory name found)
   - **Report active scope**: "Working on feature: [FeatureName]"

2. **Load design.md**:
   - Read `FEATURE_DIR/design.md`
   - Extract: requirements, success criteria, user scenarios, technical considerations
   - If missing: ERROR "design.md not found. Run `/design` first."

3. **Load memory documentation**:
   - Read `.oursky/memory/constitution.md` (mandatory)
   - Check relevant memory files for patterns
   - **Prefer memory over codebase** for patterns

4. **Generate tasks.md**:
   - Location: `FEATURE_DIR/tasks.md`
   - Structure:
     ```markdown
     # Tasks: [FeatureName]
     
     ## Status Legend
     - [ ] Not Started
     - [~] In Progress
     - [X] Completed
     - [!] Blocked/Needs Clarification
     
     ## Notes & Clarifications
     - [Add notes about misunderstandings, clarifications, or decisions here]
     
     ## Phase 1: [Phase Name]
     - [ ] Task 1 - Description
     - [ ] Task 2 - Description [depends on Task 1]
     - [ ] Task 3 - Description [parallel with Task 2]
     
     ## Phase 2: [Phase Name]
     ...
     ```
   
   - Task organization:
     - **By phase**: Setup → Core → Integration → Testing → Polish
     - **Dependencies**: Mark with `[depends on Task X]`
     - **Parallel work**: Mark with `[parallel with Task Y]`
     - **TDD**: Test tasks before implementation tasks
   
   - Task format:
     - Clear, actionable description
     - Include file paths when known
     - Reference design.md requirements
     - Note any memory patterns to follow

5. **Validate completeness**:
   - Each requirement from design.md has corresponding tasks
   - Success criteria can be verified through tasks
   - Dependencies are clear
   - Test tasks included (TDD approach)

6. **Present to user**:
   - Show tasks.md structure
   - Ask: "Does this task breakdown look complete? Any changes needed?"
   - **Handle feedback**: If user points out misunderstandings:
     - Update design.md with clarifications (add to "Clarifications" section)
     - Update tasks.md accordingly
     - Document the misunderstanding and resolution
   - Iterate until approved

7. **After approval**:
   - End with: "Tasks ready! Type 'GO!' or use `/implement` when ready to start implementation"
   - Wait for explicit approval before implementing

## Handling Misunderstandings & Feedback

When user provides feedback during planning:

1. **Acknowledge the misunderstanding**:
   - "I see, let me clarify: [restate understanding]"

2. **Update design.md**:
   - Add a "Clarifications" section if it doesn't exist:
     ```markdown
     ## Clarifications
     
     ### [Date] - [Topic]
     **Misunderstanding**: [What was misunderstood]
     **Clarification**: [Correct understanding]
     **Impact**: [What changed in tasks/requirements]
     ```

3. **Update tasks.md**:
   - Adjust affected tasks
   - Add notes explaining the change
   - Mark any tasks that need re-evaluation

4. **Update memory** (if pattern-related):
   - If misunderstanding reveals a pattern issue, update relevant memory file
   - Document the correct pattern

5. **Confirm understanding**:
   - Restate the corrected understanding
   - Ask: "Does this now match your intent?"

## Guidelines

- **Memory-first**: Reference memory docs for patterns, not codebase
- **Dependency-aware**: Order tasks correctly
- **Test-first**: Include test tasks before implementation
- **Actionable**: Each task should be completable independently
- **Traceable**: Link tasks back to design.md requirements

