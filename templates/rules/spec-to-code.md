---
alwaysApply: true
---
## Development Workflow: Spec → Code

THESE INSTRUCTIONS ARE CRITICAL!

They dramatically improve the quality of the work you create.

### Commands Overview

Use these commands for the spec-to-code workflow:

- `/consolidate` - Gather info from codebase, update memory documentation (prefer memory over codebase)
- `/design` - Create feature specification (`design.md`) - replaces asking "Should I create a Spec?"
- `/plan` - Generate tasks (`tasks.md`) from design.md
- `/implement` - Execute tasks following memory documentation patterns
- `/test` - Test with real app, reset DB, update FormX mock scenarios
- `/summarize` - Quick summary of current scope, progress, and next steps (use anytime)

### Handling Misunderstandings

During `/plan` and `/implement`, if user points out misunderstandings:

1. **Acknowledge immediately**: Restate correct understanding
2. **Update design.md**: Add to "Clarifications" section with date, misunderstanding, clarification, impact
3. **Update tasks.md**: Adjust affected tasks, add notes, mark blocked tasks with `[!]`
4. **Update memory**: If pattern-related, update relevant memory file
5. **Confirm**: Restate corrected understanding, ask if correct
6. **Resume**: Continue with corrected understanding

### Phase 1: Requirements First

When asked to implement any feature or make changes:

**Use `/design` command** - This creates the spec automatically. No need to ask "Should I create a Spec?"

The `/design` command will:
- Create `design.md` in `.oursky/scopes/FeatureName/`
- Interview user to clarify: purpose, success criteria, scope, constraints, out of scope
- Present spec for approval

### Phase 2: Review & Refine

After `/design` creates design.md:

**Use `/plan` command** - This generates tasks.md automatically

The `/plan` command will:
- Read design.md
- Generate dependency-ordered tasks in `tasks.md`
- Present tasks for approval
- Wait for "GO!" before implementation

### Phase 3: Implementation

**Use `/implement` command** - Execute tasks following memory-first approach

The `/implement` command will:
- Load tasks.md and design.md
- **Prefer memory documentation over codebase** (memory is source of truth)
- Execute tasks in dependency order
- Update task status as you progress
- Update memory docs when patterns emerge

### Phase 4: Testing & Feedback

**Use `/test` command** - Test with real app, reset DB, update FormX mock

The `/test` command will:
- Reset database for clean testing
- Update FormX mock scenarios
- Test user scenarios from design.md
- Verify success criteria
- Create feedback loop for iterative testing

### Memory Documentation

**CRITICAL**: Always prefer memory documentation over checking codebase.

- Memory location: `.oursky/memory/`
- Constitution: `.oursky/memory/constitution.md` (always read)
- Update memory when patterns change or new patterns emerge
- Memory is source of truth - codebase may contain outdated patterns

### File Organization

\`\`\`

.oursky/
├── scopes/
│   └── FeatureName/
│       ├── design.md      # Feature specification
│       ├── tasks.md       # Dependency-ordered tasks
│       └── test-results.md # Test results (from /test)
└── memory/
    ├── constitution.md    # Core principles
    └── [domain].md        # Domain-specific patterns

\`\`\`

**Remember: Think first, use commands, prefer memory over codebase, _then_ code. The Spec is your north star.**
