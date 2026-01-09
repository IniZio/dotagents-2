---
description: Run CI lint and test checks from GitHub Actions workflow and fix issues until all pass. Automatically fixes lint errors, investigates test failures, and ensures codebase passes all CI checks.
---

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Purpose

Run CI lint and test checks from `.github/workflows/job-ci.yaml` and automatically fix issues until all checks pass. This command ensures the codebase is ready for CI/CD pipeline.

## Execution Flow

1. **Load CI workflow configuration**:
   - Read `.github/workflows/job-ci.yaml`
   - Extract lint and test jobs and their commands:
     - API: lint (`make -C api lint`), test (`make -C api test`)
     - Portal: lint (`make -C portal lint`), test (`make -C portal test`)
   - Determine execution order (lint first, then tests)

2. **Initialize CI tracking**:
   - Create `.oursky/ci-results.md` to track progress
   - Set max iterations: `MAX_ITERATIONS = 5` (configurable via arguments)
   - Set current iteration: `iteration = 1`
   - Track status for each check: `{check-name}: {status}` (pending, running, passed, failed, fixed)

3. **Run checks in dependency order**:

   **Phase 1: Lint** (can auto-fix):
   - `api-lint`: `cd api && npm install && make -C api lint`
   - `portal-lint`: `cd portal && npm install && make -C portal lint`

   **Phase 2: Tests** (may need investigation):
   - `api-test`: `cd api && npm install && make -C api test`
   - `portal-test`: `cd portal && npm install && make -C portal test`

4. **For each check**:

   a. **Run the check**:
      - Execute the command
      - Capture stdout, stderr, and exit code
      - Log command and output to `.oursky/ci-results.md`
      - Update status: `{check-name}: running`

   b. **If check passes**:
      - Update status: `{check-name}: passed`
      - Log: "✓ {check-name} passed"
      - Continue to next check

   c. **If check fails**:
      - Update status: `{check-name}: failed`
      - Log: "✗ {check-name} failed"
      - Analyze error output:
        - **Lint errors**: Parse error messages, identify files and issues
        - **Test failures**: Parse test output, identify failing tests and reasons

   d. **Attempt to fix**:
      
      **For lint errors**:
      - Identify fixable issues (formatting, unused imports, etc.)
      - Run auto-fix commands:
        - API: `make -C api format` (if available) or `npm run prettier:fix && npm run eslint:fix` in api/
        - Portal: `make -C portal format` (if available) or `npm run format:fix && npm run lint:fix` in portal/
      - Re-run lint check
      - If still failing, manually fix remaining issues:
        - Read failing files
        - Apply fixes based on error messages
        - Re-run lint check
      - Update status: `{check-name}: fixed` (if successful)

      **For test failures**:
      - Parse test output to identify:
        - Failing test files and test names
        - Error messages and stack traces
        - Assertion failures
      - Read failing test files
      - Analyze test code and implementation
      - Fix issues:
        - Fix implementation bugs
        - Update test expectations if needed
        - Fix test setup/teardown issues
      - Re-run test check
      - Update status: `{check-name}: fixed` (if successful)

   e. **If fix successful**:
      - Update status: `{check-name}: passed`
      - Log: "✓ {check-name} fixed and passed"
      - Continue to next check

   f. **If fix unsuccessful after max attempts**:
      - Update status: `{check-name}: failed (needs manual intervention)`
      - Log detailed error and suggested fixes
      - Continue to next check (don't block other checks)
      - Document in `.oursky/ci-results.md` for user review

5. **Iteration management**:
   - After all checks complete:
     - Count passed vs failed checks
     - If all passed: Exit with success
     - If any failed:
       - If `iteration < MAX_ITERATIONS`:
         - Increment `iteration`
         - Re-run only failed checks
         - Continue fixing
       - If `iteration >= MAX_ITERATIONS`:
         - Report remaining failures
         - Exit with summary

6. **Final summary**:
   - Generate summary report:
     ```markdown
     # CI Results Summary
     
     ## Status: [PASSED/FAILED]
     ## Iterations: [count]
     ## Duration: [time]
     
     ## Checks
     - [✓] api-lint: passed (fixed in iteration 2)
     - [✗] api-test: failed (needs manual intervention)
     - [✓] portal-lint: passed
     - [✓] portal-test: passed
     
     ## Issues Fixed
     - api-lint: Fixed 15 formatting issues, removed unused imports
     - portal-lint: Fixed 8 ESLint warnings
     
     ## Remaining Issues
     - api-test: Test "should handle edge case" failing - needs investigation
       - File: `api/src/modules/invoices/gl-code-recommendation.service.spec.ts`
       - Error: [error message]
       - Suggested fix: [suggestion]
     
     ## Next Steps
     - [Action items for remaining failures]
     ```
   - Save to `.oursky/ci-results.md`
   - Display summary to user

## Error Analysis Patterns

### Lint Error Patterns

**Prettier errors**:
- Pattern: `Replace 'X' with 'Y'` or formatting issues
- Fix: Run `npm run prettier:fix` or `make format`

**ESLint errors**:
- Pattern: `'X' is defined but never used` → Remove unused import/variable
- Pattern: `'X' is not defined` → Add import or fix reference
- Pattern: `Unexpected console.log` → Remove or replace with logger
- Pattern: `Missing return type` → Add return type annotation
- Fix: Run `npm run eslint:fix` or manually fix based on error

**TypeScript errors**:
- Pattern: `Type 'X' is not assignable to type 'Y'` → Fix type mismatch
- Pattern: `Property 'X' does not exist` → Add property or fix reference
- Pattern: `'X' is possibly 'undefined'` → Add null check or default value
- Fix: Fix TypeScript errors in code

### Test Error Patterns

**Test failures**:
- Pattern: `Expected X but got Y` → Fix assertion or implementation
- Pattern: `Cannot read property 'X' of undefined` → Add null checks or mocks
- Pattern: `Timeout` → Fix async handling or increase timeout
- Pattern: `Mock not called` → Fix test setup or mock configuration

**Test setup errors**:
- Pattern: `BeforeAll hook failed` → Fix test setup
- Pattern: `Database connection failed` → Check test database setup

## Auto-Fix Strategies

1. **Formatting fixes** (always try first):
   - Run `make format` or equivalent
   - Re-run lint check

2. **Import fixes**:
   - Remove unused imports
   - Add missing imports
   - Fix import paths

3. **Type fixes**:
   - Add missing type annotations
   - Fix type mismatches
   - Add null checks

4. **Test fixes**:
   - Fix assertions
   - Fix mocks
   - Fix async handling
   - Update test expectations if implementation changed

## Command Arguments

The command accepts arguments in format: `[options]`

Options:
- `--max-iterations <number>` - Maximum fix iterations (default: 5)
- `--skip-lint` - Skip lint checks
- `--skip-test` - Skip test checks
- `--component <name>` - Run checks for specific component only (api, portal)
- `--fix-only` - Only attempt fixes, don't run checks that already pass

Examples:
- `/ci` - Run all checks and fix issues
- `/ci --max-iterations 10` - Allow up to 10 fix iterations
- `/ci --component api` - Only run API checks
- `/ci --skip-lint` - Skip lint checks, run tests only

## CI Results File Format

The `.oursky/ci-results.md` file tracks:

```markdown
# CI Results

## Run Information
- Date: [timestamp]
- Iterations: [count]
- Duration: [time]
- Status: [PASSED/FAILED]

## Checks

### api-lint
- Status: [passed/failed/fixed]
- Iteration: [number]
- Command: `make -C api lint`
- Output: [command output]
- Fixes Applied:
  - Fixed 15 Prettier formatting issues
  - Removed 3 unused imports
  - Fixed 2 TypeScript type errors

[... similar for all checks ...]

## Summary
- Total Checks: [count]
- Passed: [count]
- Failed: [count]
- Fixed: [count]
- Needs Manual Intervention: [count]

## Remaining Issues
[Detailed list of issues that need manual intervention]

## Next Steps
[Action items]
```

## Guidelines

- **Run checks in order**: Lint → Tests
- **Fix automatically when possible**: Lint errors, formatting, simple type errors
- **Investigate test failures**: Read test code, understand failure, fix appropriately
- **Don't block on one check**: Continue with other checks if one needs manual intervention
- **Document everything**: Log all commands, outputs, and fixes
- **Be conservative with fixes**: Don't change logic unless clearly a bug
- **Respect max iterations**: Stop after max iterations to avoid infinite loops
- **Report clearly**: Show what passed, what failed, what was fixed, what needs manual work

## Error Handling

- **Command execution failures**: Log error, mark check as failed, continue
- **Fix application failures**: Log error, mark as needs manual intervention, continue
- **Timeout issues**: Log timeout, suggest increasing timeout or investigating
- **Permission issues**: Log error, suggest checking permissions
- **Missing dependencies**: Attempt to install, log if fails

## Success Criteria

The command succeeds when:
- All checks pass (or all fixable issues are fixed)
- Remaining failures are documented with clear next steps
- CI results file is updated with complete information

The command fails when:
- Critical checks fail and cannot be fixed automatically
- Max iterations reached with remaining failures
- User explicitly stops the process

