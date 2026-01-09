---
description: Execute tasks from tasks.md. Follow memory documentation patterns. Update tasks status as you complete them.
---

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Purpose

Execute implementation following `tasks.md`. Use memory documentation as source of truth for patterns. Update task status as you progress.

## Execution Flow

1. **Detect active scope and feature directory**:
   - Check if user specified feature name in arguments
   - If not, look for most recent `.oursky/scopes/*/tasks.md` or `design.md`
   - Extract feature name from directory:
     - Pattern: `.oursky/scopes/{date-prefix}-{FeatureName}/` (e.g., `20251121-123456-TestCommandFeedbackLoop/`)
     - Extract `FeatureName` by removing date prefix (format: `YYYYMMDD-HHMMSS-`)
     - Example: `20251121-123456-TestCommandFeedbackLoop` → `TestCommandFeedbackLoop`
   - Set `FEATURE_NAME = FeatureName` (PascalCase)
   - Set `FEATURE_DIR = .oursky/scopes/{date-prefix}-{FeatureName}/` (use actual directory name found)
   - **Report active scope**: "Working on feature: [FeatureName]"

2. **Ensure correct feature branch**:
   - **Check current branch**: Run `git branch --show-current` to get current branch name
   - **Convert feature name to branch format**:
     - Convert PascalCase to kebab-case:
       - Example: `TestCommandFeedbackLoop` → `test-command-feedback-loop`
       - Pattern: Insert hyphens before capital letters, convert to lowercase
     - Check for existing branches matching pattern: `{linear-issue}-[kebab-name]`
       - Check remote: `git ls-remote --heads origin | grep -E 'refs/heads/[A-Z]+-[0-9]+-[kebab-name]$'`
       - Check local: `git branch | grep -E '^[* ]*[A-Z]+-[0-9]+-[kebab-name]$'`
     - If branch exists: Use existing branch
     - If no branch exists: **Ask user for Linear issue number**:
       - Prompt: "What is the Linear issue number for this feature? (e.g., ENG-123)"
       - Wait for user input
       - Validate format: Should match pattern `[A-Z]+-[0-9]+` (e.g., `ENG-123`, `FEAT-456`)
   - **Branch naming pattern**: `{linear-issue}-{kebab-case-name}` (e.g., `ENG-123-test-command-feedback-loop`)
   - **Checkout or create branch**:
     - If on correct branch: Continue with implementation
     - If on different branch:
       - Check if feature branch exists
       - If exists: Ask user "Currently on [current-branch]. Should I checkout [feature-branch]?"
       - If doesn't exist: Ask user "Feature branch [feature-branch] doesn't exist. Should I create it?"
     - If branch doesn't exist and user approves: Create new branch from main:
       ```bash
       git switch main
       git pull origin main
       git switch -c {linear-issue}-{kebab-name}
       ```
     - **CRITICAL**: Always ensure you're on the correct feature branch before implementing
   - **Report branch status**: 
     - "Current branch: [branch-name]"
     - "Feature branch: [feature-branch-name]"
     - "Status: [On correct branch / Need to checkout / Need to create]"

3. **Load tasks.md**:
   - Read `FEATURE_DIR/tasks.md`
   - Parse task status and dependencies
   - If missing: ERROR "tasks.md not found. Run `/plan` first."

4. **Load memory documentation** (CRITICAL):
   - Read `.oursky/memory/constitution.md` (mandatory)
   - Check relevant memory files for patterns
   - **ALWAYS prefer memory over codebase** - memory is source of truth
   - If pattern conflicts: Use memory, note discrepancy

5. **Load design.md**:
   - Read `FEATURE_DIR/design.md` for context
   - Reference requirements when making decisions

6. **Detect phases in tasks.md**:
   - Parse `tasks.md` to identify phases (look for headings like `## Phase N: [Name]`)
   - Group tasks by phase
   - Track which phase is currently being worked on
   - Identify when all tasks in a phase are complete

7. **Execute tasks in order**:
   - **Respect dependencies**: Complete prerequisite tasks first
   - **Parallel tasks**: Can work on multiple [parallel] tasks
   - **TDD approach**: Write tests before implementation
   - **Phase-by-phase**: Complete each phase before next
   
8. **For each task**:
   - Check memory docs for patterns/conventions
   - Implement following memory patterns (not codebase patterns)
   - Write tests if test task, implement if implementation task
   - Update task status: `[ ]` → `[~]` → `[X]`
   - **Commit after completing logical units** (atomic commits):
     - Commit when a logical unit of work is complete (e.g., "Add invoice service", "Fix validation bug")
     - Use concise, readable commit messages (no semantic commit format)
     - Commit frequently, not just at phase boundaries

9. **Phase completion checkpoint** (CRITICAL):
   - **After completing all tasks in a phase** (when all tasks in phase marked `[X]`):
     - **Run tests and verify coverage**:
       - Determine which component changed (API or Portal):
         - If tasks involve `api/src/`: API changes
         - If tasks involve `portal/src/`: Portal changes
         - If both: Run tests for both
       - **For API changes**: 
         - Run: `make -C api test`
         - Extract coverage from output (look for "Coverage" section):
           ```
           Coverage summary:
           Lines:    85.23% (1234/1447)
           Functions: 82.45% (567/688)
           Branches:  78.12% (456/584)
           Statements: 85.23% (1234/1447)
           ```
         - Verify thresholds: Lines ≥78%, Functions ≥77%, Branches ≥73%, Statements ≥78%
       - **For Portal changes**:
         - Run: `make -C portal test`
         - Extract coverage similarly
       - **List test files created/modified in this phase**:
         - Use `git diff --name-only` to find new/modified `.spec.ts` files
         - For each test file, extract test case details:
           - Read the test file
           - Parse `it(` or `test(` blocks to get test names
           - Extract setup/arrange code (what's prepared)
           - Extract assertions (what's verified)
       - **Prove test coverage**: Display detailed test case summary:
         ```
         Phase [N] Test Coverage Summary:
         =================================
         
         Test File: api/src/modules/feature/feature.service.spec.ts
         ──────────────────────────────────────────────────────────
         1. "should create feature successfully"
            Arrange: Mock repository with empty find result, valid input DTO
            Assert: Service method called with correct params, returns created entity
         
         2. "should throw error when feature already exists"
            Arrange: Mock repository returning existing feature
            Assert: Throws ConflictException with message "Feature already exists"
         
         3. "should handle invalid input data"
            Arrange: Invalid DTO with missing required fields
            Assert: Throws BadRequestException with validation errors
         
         4. "should return feature by id"
            Arrange: Mock repository returning feature entity
            Assert: Returns correct feature DTO with all fields
         
         5. "should throw NotFoundException when feature not found"
            Arrange: Mock repository returning null
            Assert: Throws NotFoundException with "Feature not found" message
         
         6. "should update feature successfully"
            Arrange: Mock repository with existing feature, update DTO
            Assert: Service update method called, returns updated entity
         
         7. "should handle database errors gracefully"
            Arrange: Mock repository throwing database error
            Assert: Throws InternalServerErrorException with error message
         
         8. "should return empty array when no features exist"
            Arrange: Mock repository returning empty array
            Assert: Returns empty array, no errors thrown
         
         Test File: api/src/modules/feature/feature.controller.spec.ts
         ──────────────────────────────────────────────────────────
         1. "POST /features should create feature"
            Arrange: Valid request body, authenticated user, mock service
            Assert: Returns 201 with created feature, service.create called
         
         2. "POST /features should return 400 for invalid input"
            Arrange: Invalid request body (missing required fields)
            Assert: Returns 400 with validation errors
         
         3. "GET /features/:id should return feature"
            Arrange: Valid feature ID, mock service returning feature
            Assert: Returns 200 with feature data
         
         4. "GET /features/:id should return 404 for non-existent feature"
            Arrange: Invalid feature ID, mock service throwing NotFoundException
            Assert: Returns 404 with error message
         
         5. "PUT /features/:id should update feature"
            Arrange: Valid ID and update body, mock service
            Assert: Returns 200 with updated feature, service.update called
         
         6. "DELETE /features/:id should delete feature"
            Arrange: Valid feature ID, mock service
            Assert: Returns 204, service.delete called
         
         Test File: api/src/queue/processors/feature.processor.spec.ts
         ──────────────────────────────────────────────────────────
         1. "should process feature job successfully"
            Arrange: Valid job data, mock service methods
            Assert: Service methods called in correct order, job completed
         
         2. "should retry on transient errors"
            Arrange: Mock service throwing transient error, retry configured
            Assert: Job retried, eventually succeeds
         
         3. "should fail job on permanent errors"
            Arrange: Mock service throwing permanent error
            Assert: Job fails, error logged
         
         4. "should handle missing job data"
            Arrange: Job with missing required fields
            Assert: Job fails with validation error
         
         Coverage thresholds met: ✓
         All tests passing: ✓
         ```
     - **Atomic git commits**:
       - **Format**: Concise, readable commit messages (no semantic commit format)
       - **Commit frequently**: After completing logical units of work (not per phase)
       - **Message style**: Clear, descriptive, easy to read
       - **Examples**: 
         ```bash
         git commit -m "Add invoice service with CRUD operations"
         git commit -m "Implement invoice validation logic"
         git commit -m "Add tests for invoice service"
         git commit -m "Fix invoice number generation edge case"
         ```
       - **Avoid**: Phase-based commits, semantic prefixes, overly verbose messages
       - **Verify commit**: Show commit hash: `git rev-parse HEAD`
     - **Report phase completion**:
       - Display formatted summary:
         ```
         ✓ Phase [N] completed: [Phase Name]
         ✓ Test cases: [detailed list as shown above]
         ✓ Coverage thresholds met: ✓
         ✓ All tests passing: ✓
         ```
     - **Request user confirmation**:
       - **CRITICAL**: Display clearly:
         ```
         ════════════════════════════════════════════════════════════
         Phase [N] complete. All tests passing and coverage verified.
         Type `/implement` again to proceed to Phase [N+1].
         ════════════════════════════════════════════════════════════
         ```
       - **STOP execution** - Do NOT proceed to next phase
       - **Wait for user to type `/implement`** before continuing
       - If user types `/implement` again, proceed to next phase

9. **Update memory** (when patterns emerge):
   - If you discover new patterns or make decisions:
     - Update relevant memory files immediately
     - Document rationale
   - If codebase differs from memory:
     - Note discrepancy
     - Ask user if codebase should be updated or memory adjusted

10. **Handle feedback and misunderstandings**:
   - If user provides feedback or points out misunderstandings:
     - **Stop current task** if needed
     - **Acknowledge**: "I see, let me clarify: [restate understanding]"
     - **Update design.md**: Add clarification to "Clarifications" section
     - **Update tasks.md**: Adjust affected tasks, mark with `[!]` if blocked
     - **Update memory**: If pattern-related, update relevant memory file
     - **Confirm**: Restate corrected understanding, ask if correct
     - **Resume**: Continue with corrected understanding

11. **Progress reporting**:
   - Report after each completed task
   - Show current phase progress
   - Note any blockers or decisions needed
   - If misunderstandings occur, document in tasks.md notes

12. **Completion**:
   - Verify all tasks marked `[X]`
   - Check that implementation matches design.md
   - Ensure tests pass
   - Update memory if new patterns established
   - Report completion with summary

## Guidelines

- **Memory-first**: Always check memory before codebase
- **Update tasks**: Mark tasks complete as you go
- **Follow constitution**: Adhere to all principles in constitution.md
- **Test-first**: Write tests before implementation
- **Update memory**: Document new patterns immediately
- **Ask before changing scope**: If design needs changes, ask user first
- **Phase checkpoints**: Always commit and verify tests before proceeding to next phase
- **User confirmation**: Wait for user to type `/implement` before starting next phase

## Error Handling

- If task fails: Report error, suggest fix, ask if should continue
- If memory pattern conflicts with codebase: Use memory, note discrepancy
- If design unclear: Reference design.md, ask user for clarification
- **If misunderstanding discovered**: Stop, clarify with user, update design.md and tasks.md, then resume

## Handling Misunderstandings During Implementation

When user provides feedback during implementation:

1. **Acknowledge immediately**:
   - "I understand. Let me clarify: [restate correct understanding]"

2. **Update design.md**:
   - Add to "Clarifications" section:
     ```markdown
     ### [Date] - [Topic]
     **Misunderstanding**: [What was misunderstood]
     **Clarification**: [Correct understanding]
     **Impact**: [What changed in implementation]
     ```

3. **Update tasks.md**:
   - Mark affected tasks with `[!]` if blocked
   - Add notes explaining the change
   - Adjust task descriptions if needed
   - Update status of tasks that need rework

4. **Update code** (if already implemented):
   - Revert or adjust code based on clarification
   - Update tests if needed
   - Document the change

5. **Update memory** (if pattern-related):
   - If misunderstanding reveals pattern issue, update memory
   - Document correct pattern

6. **Confirm and resume**:
   - "Does this now match your intent? Should I continue with [next task]?"

