---
description: Create or update feature specification (design.md) following spec-to-code workflow. Interview user to clarify requirements before coding.
handoffs:
  - label: Create Tasks
    agent: plan
    prompt: Create tasks.md from this design
    send: true
---

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Purpose

Create a feature specification (`design.md`) in `.oursky/scopes/FeatureName/` following the spec-to-code workflow. This replaces asking "Should I create a Spec?" - this command IS the spec creation.

## Execution Flow

1. **Determine feature name and scope**:
   - Extract from user input or ask if unclear
   - Use PascalCase (e.g., "UserAuthentication", "InvoiceProcessing")
   - **Generate date prefix**: Use format `YYYYMMDD-HHMMSS` (e.g., `20251121-123456`)
     - Get current date/time: `date +%Y%m%d-%H%M%S`
   - Create directory: `.oursky/scopes/{date-prefix}-{FeatureName}/`
     - Example: `.oursky/scopes/20251121-123456-UserAuthentication/`
   - Set `FEATURE_NAME = FeatureName` (PascalCase)
   - Set `FEATURE_DIR = .oursky/scopes/{date-prefix}-{FeatureName}/`
   - **Report active scope**: "Creating feature: [FeatureName] in [FEATURE_DIR]"

2. **Load memory documentation**:
   - Read `.oursky/memory/constitution.md` (mandatory)
   - Check relevant memory files for patterns/constraints
   - **Prefer memory over codebase** - memory is source of truth

3. **Interview user** (if needed):
   - Purpose & user problem
   - Success criteria (measurable, technology-agnostic)
   - Scope & constraints
   - Technical considerations
   - Out of scope items
   - Only ask critical questions (max 3-5)

4. **Create design.md**:
   - Location: `FEATURE_DIR/design.md` (e.g., `.oursky/scopes/20251121-123456-FeatureName/design.md`)
   - Structure:
     ```markdown
     # Feature: [FeatureName]
     
     ## Purpose & User Problem
     [What problem does this solve?]
     
     ## Success Criteria
     [Measurable, technology-agnostic outcomes]
     
     ## Scope
     ### In Scope
     - [What's included]
     
     ### Out of Scope
     - [What's explicitly excluded]
     
     ## Requirements
     ### Functional
     - [Testable requirements]
     
     ### Non-Functional
     - [Performance, security, etc.]
     
     ## Technical Considerations
     - [Constraints, dependencies, integrations]
     
     ## User Scenarios
     - [Primary user flows]
     
     ## Assumptions
     - [Documented assumptions]
     
     ## Clarifications
     [Added during planning/implementation when misunderstandings are resolved]
     ### [Date] - [Topic]
     **Misunderstanding**: [What was misunderstood]
     **Clarification**: [Correct understanding]
     **Impact**: [What changed in requirements/tasks]
     ```
   - Keep it concise and focused on WHAT/WHY, not HOW

5. **Present to user**:
   - Show the design.md
   - Ask: "Does this capture your intent? Any changes needed?"
   - Iterate until approved

6. **After approval**:
   - Suggest: "Ready to create tasks? Use `/plan` or I can create tasks.md now"
   - Wait for explicit approval before proceeding to tasks

## Guidelines

- **Requirements first**: Focus on user needs, not implementation
- **Memory-aware**: Reference memory docs for patterns/constraints
- **Concise**: Avoid over-specification
- **Testable**: All requirements must be verifiable
- **Technology-agnostic**: Success criteria shouldn't mention frameworks/tools

