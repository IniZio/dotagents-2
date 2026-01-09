---
description: Summarize current feature scope, progress, and next steps. Help colleagues quickly pick up context without questions.
---

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Purpose

Generate a concise summary of the current feature scope, progress, key decisions, and next steps. Enables colleagues to quickly understand context and continue work without asking questions.

## Execution Flow

1. **Detect active scope and feature directory**:
   - Check if user specified feature name in arguments
   - If not, look for most recent `.oursky/scopes/*/design.md` or `tasks.md`
   - If multiple exist, use most recently modified
   - Extract feature name from directory:
     - Pattern: `.oursky/scopes/{date-prefix}-{FeatureName}/` (e.g., `20251121-123456-TestCommandFeedbackLoop/`)
     - Extract `FeatureName` by removing date prefix (format: `YYYYMMDD-HHMMSS-`)
     - Example: `20251121-123456-TestCommandFeedbackLoop` → `TestCommandFeedbackLoop`
   - Set `FEATURE_NAME = FeatureName` (PascalCase)
   - Set `FEATURE_DIR = .oursky/scopes/{date-prefix}-{FeatureName}/` (use actual directory name found)
   - **Report active scope**: "Summarizing feature: [FeatureName]"

2. **Load all artifacts**:
   - Read `FEATURE_DIR/design.md` (if exists)
   - Read `FEATURE_DIR/tasks.md` (if exists)
   - Read `FEATURE_DIR/test-results.md` (if exists)
   - Check for any notes or clarifications files

3. **Load memory documentation**:
   - Read `.oursky/memory/constitution.md` (for context)
   - Check relevant memory files if needed

4. **Generate summary**:
   - Create concise, actionable summary
   - Structure:
     ```markdown
     # Summary: [FeatureName]
     
     **Last Updated**: [Date from most recent file]
     **Status**: [In Progress / Planning / Testing / Complete]
     
     ## What We're Building
     [1-2 sentence description from design.md]
     
     ## Current Progress
     - **Phase**: [Current phase]
     - **Tasks Completed**: X / Y
     - **Tasks In Progress**: [List]
     - **Tasks Blocked**: [List with reasons]
     
     ## Key Decisions Made
     - [Decision 1]: [Rationale]
     - [Decision 2]: [Rationale]
     
     ## Clarifications & Notes
     [From design.md clarifications and tasks.md notes]
     
     ## Next Steps
     1. [Next immediate task]
     2. [Following task]
     3. [Upcoming milestone]
     
     ## Success Criteria Status
     - [ ] Criterion 1: [Status]
     - [ ] Criterion 2: [Status]
     
     ## Blockers & Questions
     - [Any blockers or open questions]
     
     ## Quick Context
     - **Files Changed**: [Key files modified]
     - **Patterns Used**: [Memory patterns referenced]
     - **Dependencies**: [External dependencies or integrations]
     ```
   
   - Keep summary concise (aim for < 200 lines)
   - Focus on actionable information
   - Highlight what's done vs what's next

5. **Output summary**:
   - Display summary in chat
   - Optionally save to `FEATURE_DIR/summary.md` if user requests

## Guidelines

- **Be concise**: Focus on essential information
- **Be actionable**: Clear next steps
- **Be current**: Reflect latest state
- **Be complete**: Include all critical context
- **No questions**: Summary should be self-contained

## What to Include

**Must include**:
- Feature purpose (1-2 sentences)
- Current phase and progress
- Next 2-3 immediate tasks
- Any blockers or clarifications
- Key decisions and rationale

**Should include**:
- Success criteria status
- Files/modules affected
- Patterns or conventions used
- Dependencies

**Optional**:
- Historical context (if relevant)
- Related features or connections
- Testing status

## Example Output

```markdown
# Summary: UserAuthentication

**Last Updated**: 2025-01-27
**Status**: In Progress (Phase 2: Core Implementation)

## What We're Building
Add OAuth2 authentication to allow users to sign in with Google and GitHub accounts. Users can link multiple providers to one account.

## Current Progress
- **Phase**: Core Implementation
- **Tasks Completed**: 8 / 15
- **Tasks In Progress**: 
  - [~] Implement OAuth2 callback handler
  - [~] Create user-provider linking logic
- **Tasks Blocked**: None

## Key Decisions Made
- Use Passport.js for OAuth2: Standard library, well-maintained
- Store provider tokens encrypted: Security requirement from constitution
- Support multiple providers per user: Requirement from design.md

## Clarifications & Notes
- 2025-01-27: Clarified that users can link multiple providers (not just one)
- Email verification optional for OAuth providers (they verify email)

## Next Steps
1. Complete OAuth2 callback handler
2. Implement provider linking UI
3. Add tests for multi-provider scenarios

## Success Criteria Status
- [✓] Users can sign in with Google (completed)
- [~] Users can sign in with GitHub (in progress)
- [ ] Users can link multiple providers (pending)
- [ ] 95% test coverage (pending)

## Blockers & Questions
- None currently

## Quick Context
- **Files Changed**: `api/src/auth/`, `api/src/users/`, `portal/src/auth/`
- **Patterns Used**: NestJS modules, DTOs for validation, encrypted storage
- **Dependencies**: Passport.js, @nestjs/passport, crypto for encryption
```

