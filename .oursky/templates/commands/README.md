# Spec-to-Code Commands

Concise workflow commands aligned with the spec-to-code development process. These commands replace the need to manually ask "Should I create a Spec?" and streamline the entire development cycle.

## Command Flow

```
/consolidate → /design → /plan → /implement → /test → /ci
/summarize (anytime - get quick context)
/ci (anytime - ensure CI passes)
```

### `/consolidate`

Gather information from codebase and update memory documentation.

**When to use:**
- Starting a new feature area
- Patterns have changed
- Need to document decisions

**What it does:**
- Checks memory documentation first (preferred over codebase)
- Gathers patterns from codebase if memory incomplete
- Updates `.oursky/memory/` files
- Documents patterns, decisions, and conventions

**Memory files:**
- `.oursky/memory/constitution.md` - Core principles
- `.oursky/memory/architecture.md` - System architecture
- `.oursky/memory/api-patterns.md` - API design patterns
- `.oursky/memory/database-patterns.md` - DB patterns
- `.oursky/memory/testing-patterns.md` - Testing approaches
- `.oursky/memory/integration-patterns.md` - External integrations

---

### `/design`

Create feature specification (`design.md`). **Replaces asking "Should I create a Spec?"**

**When to use:**
- Starting any new feature
- Need to clarify requirements before coding

**What it does:**
- Creates `.oursky/scopes/{date-prefix}-FeatureName/design.md` (date format: `YYYYMMDD-HHMMSS`)
- Interviews user for: purpose, success criteria, scope, constraints
- Presents spec for approval
- Keeps spec concise and technology-agnostic

**Output:**
- `design.md` with: purpose, success criteria, scope, requirements, user scenarios

---

### `/plan`

Generate dependency-ordered tasks from design.md.

**When to use:**
- After `/design` creates design.md
- Need to break down feature into actionable tasks

**What it does:**
- Reads `design.md`
- Generates `tasks.md` with dependency-ordered tasks
- Organizes by phases: Setup → Core → Integration → Testing → Polish
- Marks dependencies and parallel work
- Follows TDD (test tasks before implementation)

**Output:**
- `tasks.md` with kanban-style task list

---

### `/implement`

Execute tasks following memory-first approach.

**When to use:**
- After `/plan` creates tasks.md
- Ready to start coding

**What it does:**
- **Detects active scope**: Identifies feature from `.oursky/scopes/{date-prefix}-FeatureName/`
- **Ensures correct branch**: Asks for Linear issue number, checks/creates feature branch (`{linear-issue}-{kebab-name}`)
- Loads `tasks.md` and `design.md`
- **Prefers memory documentation over codebase** (memory is source of truth)
- Executes tasks in dependency order
- Updates task status: `[ ]` → `[~]` → `[X]`
- Updates memory when new patterns emerge
- Follows constitution principles

**Branch management:**
- Feature name (PascalCase) → Branch name (kebab-case): `TestCommandFeedbackLoop` → `test-command-feedback-loop`
- Branch format: `{linear-issue}-{kebab-name}` (e.g., `ENG-123-test-command-feedback-loop`)
- Commits: Atomic, concise messages (no semantic commit format)
- Automatically checks if correct branch exists
- Creates branch if needed (with user approval)
- Ensures you're on correct branch before implementing

**Key principle:**
- Always check memory docs before codebase
- Memory is authoritative - codebase may have outdated patterns
- Always work on correct feature branch

---

### `/test`

Test feature with real app, reset database, update FormX mock scenarios.

**When to use:**
- After implementation or during development
- Need to verify functionality end-to-end

**What it does:**
- Resets database (full reset or clear data)
- Updates FormX mock scenarios for testing
- Tests user scenarios from design.md
- Verifies success criteria
- Creates `test-results.md` with findings
- Creates feedback loop for iterative testing

**Test flow:**
1. Reset DB → Start services → Update FormX mock
2. Test scenarios → Verify success criteria
3. Document results → Fix issues → Re-test

---

## Workflow Example

```bash
# 1. Consolidate patterns (optional, but recommended)
/consolidate authentication patterns

# 2. Create specification
/design Add user authentication with OAuth2

# 3. Generate tasks
/plan

# 4. Implement (after approval)
/implement

# 5. Test with real app
/test

# Get quick context anytime
/summarize
```

---

### `/ci`

Run CI lint and test checks from GitHub Actions workflow and fix issues until all pass.

**When to use:**
- Before pushing code to ensure CI will pass
- After making changes that might break CI
- Need to fix lint or test errors automatically

**What it does:**
- Runs lint and test checks from `.github/workflows/job-ci.yaml`:
  - API: lint, test
  - Portal: lint, test
- Automatically fixes issues:
  - Lint errors: Auto-fixes formatting, removes unused imports, fixes simple issues
  - Test failures: Investigates and fixes failing tests
- Iterates until all checks pass or max iterations reached
- Creates `.oursky/ci-results.md` with detailed results

**Options:**
- `--max-iterations <number>` - Maximum fix iterations (default: 5)
- `--component <name>` - Run checks for specific component (api, portal)
- `--skip-lint` - Skip lint checks
- `--skip-test` - Skip test checks

**Output:**
- `.oursky/ci-results.md` with check status, fixes applied, remaining issues
- Summary of what passed, what failed, what was fixed

---

### `/summarize`

Generate concise summary of current feature scope, progress, and next steps.

**When to use:**
- Handing off to a colleague
- Returning to a feature after time away
- Need quick context without reading all files
- Before starting work session

**What it does:**
- Reads design.md, tasks.md, test-results.md
- Summarizes: what's built, progress, decisions, next steps
- Highlights blockers and clarifications
- Provides actionable next steps

**Output:**
- Concise summary (< 200 lines)
- Self-contained (no questions needed)
- Actionable next steps
- Optionally saves to `summary.md`

---

## Active Scope Detection

All commands detect the active feature scope:

- **Scope location**: `.oursky/scopes/{date-prefix}-FeatureName/` (e.g., `.oursky/scopes/20251121-123456-FeatureName/`)
- **Date prefix format**: `YYYYMMDD-HHMMSS` (generated when scope is created)
- **Detection**: Commands look for most recent `design.md` or `tasks.md`
- **Feature name**: Extracted from directory name (e.g., `TestCommandFeedbackLoop` from `20251121-123456-TestCommandFeedbackLoop`)
- **Reporting**: Commands report "Working on feature: [FeatureName]"

**Branch naming**:
- Feature name: PascalCase (e.g., `TestCommandFeedbackLoop`)
- Branch name: `{linear-issue}-{kebab-case}` (e.g., `ENG-123-test-command-feedback-loop`)
- Linear issue format: `[A-Z]+-[0-9]+` (e.g., `ENG-123`, `FEAT-456`)
- Conversion: Insert hyphens before capitals, convert to lowercase
- **User prompt**: `/implement` asks for Linear issue number before creating branch

**Branch management** (`/implement` only):
- Checks current branch vs feature branch
- Asks user for Linear issue number if branch doesn't exist
- Creates branch if doesn't exist (with user approval)
- Ensures correct branch before implementing

---

## Handling Misunderstandings & Feedback

During `/plan` and `/implement`, if you discover misunderstandings:

1. **Point out the issue**: "Actually, I meant X, not Y"
2. **AI acknowledges**: Clarifies understanding
3. **Updates artifacts**:
   - Adds clarification to `design.md` → "Clarifications" section
   - Updates affected tasks in `tasks.md` → "Notes & Clarifications" section
   - Updates memory if pattern-related
4. **Confirms**: "Does this now match your intent?"
5. **Resumes**: Continues with corrected understanding

**Documentation**:
- Misunderstandings tracked in `design.md` → "Clarifications" section
- Task notes in `tasks.md` → "Notes & Clarifications" section
- Pattern updates in `.oursky/memory/` files

---

## Key Principles

1. **Memory-first**: Always prefer memory documentation over codebase
2. **Spec-first**: Create design.md before coding
3. **Test-first**: Write tests before implementation
4. **Iterative**: Test → Fix → Re-test until passing
5. **Document**: Update memory when patterns change

## File Structure

```
.oursky/
├── scopes/
│   └── FeatureName/
│       ├── design.md       # Feature specification
│       ├── tasks.md        # Dependency-ordered tasks
│       └── test-results.md # Test results
└── memory/
    ├── constitution.md      # Core principles (always read)
    └── [domain].md         # Domain-specific patterns
```

## Memory Documentation

Memory documentation is the **source of truth**. Always check memory before codebase because:

- Codebase may contain outdated patterns
- Memory documents current best practices
- Memory includes rationale for decisions
- Memory is easier to maintain than scattered code

When patterns change:
1. Update memory immediately
2. Note if codebase needs updating
3. Document rationale for changes

