---
description: Consolidate information from codebase and update memory documentation. Prefer memory docs over checking codebase for patterns.
---

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Purpose

Gather information from the codebase, consolidate findings, and update memory documentation. Memory documentation is the source of truth - prefer it over checking codebase which may contain outdated patterns.

## Execution Flow

1. **Detect active scope** (if applicable):
   - Check for active feature scope: Look for most recent `.oursky/scopes/*/design.md` or `tasks.md`
   - If feature scope found: Extract feature name and note it
   - **Report active scope**: "Consolidating for feature: [FeatureName]" (if applicable)

2. **Identify what to consolidate**:
   - If user provided context: Focus on that area (e.g., "authentication patterns", "database schema", "API endpoints")
   - If no context: Scan for recent changes or ask user what area to consolidate

3. **Check memory documentation first**:
   - Read `.oursky/memory/constitution.md` (always)
   - Check for existing memory files in `.oursky/memory/` related to the topic
   - Note what's already documented

4. **Gather from codebase** (only if memory is incomplete):
   - Search for relevant patterns, implementations, or examples
   - Extract key decisions, patterns, and conventions
   - Note any inconsistencies or outdated patterns

5. **Update memory documentation**:
   - **CRITICAL**: Always use `date +%Y-%m-%d` command to get the real current date for "Last Updated" fields
   - Create or update memory files in `.oursky/memory/`
   - Document:
     - Patterns and conventions
     - Architectural decisions
     - Common pitfalls and solutions
     - Testing approaches
     - Integration patterns
   - Format: Use clear sections, examples, and rationale
   - Always update "Last Updated" field with current date from `date` command

6. **Report consolidation**:
   - List what was consolidated
   - Highlight any patterns that differ from memory (may indicate need for code updates)
   - Note any gaps that need clarification

## Memory Documentation Structure

Memory files should be organized by domain/concern:
- `.oursky/memory/constitution.md` - Core principles (already exists)
- `.oursky/memory/architecture.md` - System architecture patterns
- `.oursky/memory/api-patterns.md` - API design patterns
- `.oursky/memory/database-patterns.md` - Database and migration patterns
- `.oursky/memory/testing-patterns.md` - Testing approaches
- `.oursky/memory/integration-patterns.md` - External integrations (FormX, etc.)

## Guidelines

- **Memory-first**: Always check memory before codebase
- **Update memory**: When patterns change, update memory docs immediately
- **Be concise**: Focus on patterns and decisions, not implementation details
- **Include rationale**: Explain why patterns exist
- **Version awareness**: Note when patterns were established or changed
- **Automated testing**: When testing patterns or workflows need automation, **ask the user for help** to create automated test flows instead of falling back to manual testing guidance. The user can provide assistance with browser automation, test data setup, or other automation needs.

