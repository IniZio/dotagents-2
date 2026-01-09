---
description: Test feature with real app. Reset database, update FormX mock scenarios, verify functionality. Create feedback loop for iterative testing with browser automation.
---

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Purpose

Test the implemented feature with the real application using browser automation. Reset database, configure FormX mock scenarios programmatically, verify UI layouts with screenshots and snapshots, and create an automated feedback loop for iterative testing.

## Execution Flow

1. **Detect active scope and feature directory**:
   - Check if user specified feature name in arguments
   - If not, look for most recent `.oursky/scopes/*/design.md`
   - Extract feature name from directory:
     - Pattern: `.oursky/scopes/{date-prefix}-{FeatureName}/` (e.g., `20251121-123456-TestCommandFeedbackLoop/`)
     - Extract `FeatureName` by removing date prefix (format: `YYYYMMDD-HHMMSS-`)
     - Example: `20251121-123456-TestCommandFeedbackLoop` → `TestCommandFeedbackLoop`
   - Set `FEATURE_NAME = FeatureName` (PascalCase)
   - Set `FEATURE_DIR = .oursky/scopes/{date-prefix}-{FeatureName}/` (use actual directory name found)
   - **Report active scope**: "Testing feature: [FeatureName]"
   - Create `FEATURE_DIR/screenshots/` directory if it doesn't exist

2. **Load design.md**:
   - Read `FEATURE_DIR/design.md` for test scenarios
   - Extract user scenarios and success criteria
   - Understand what to test and expected UI layouts

3. **Check browser tools availability**:
   - Verify MCP browser tools are available (browser_navigate, browser_snapshot, browser_take_screenshot, etc.)
   - If unavailable, log warning but continue with manual testing guidance
   - Browser tools enable automated visual verification

4. **Prepare test environment**:
   - **Reset database** (automated):
     - Determine reset type from arguments or default to "full" (option 1)
     - Execute `./scripts/reset-db.sh` non-interactively:
       - For "full": echo "1" and "yes" to script
       - For "clear-data": echo "2" and "yes" to script
       - For "migrations": echo "3" and "yes" to script
     - Verify reset completion by checking database connection
     - Log reset type and timestamp in test results
   
   - **Start services**:
     - Ensure docker-compose services are running (`docker compose ps`)
     - Check API is accessible (default: `http://localhost:8080`)
     - Check Portal is accessible (default: `http://localhost:3000`)
     - Check FormX mock is accessible (default: `http://localhost:3001`)
     - If services not running, start them with `docker compose up -d`
   
   - **Configure FormX mock scenarios** (programmatic):
     - Use FormX mock API to configure scenarios:
       - Base URL: `http://localhost:3001` (or from config)
       - Clear existing matchers: `GET /api/matchers`, then `DELETE /api/matchers/:id` for each
       - Create new matchers: `POST /api/matchers` with scenario configuration
     - Common scenarios:
       - **Happy Path**: Successful extraction (200, immediate data)
       - **Error**: Rate limit (429), server error (500), etc.
       - **Async Job**: Pending → Processing → Completed progression
       - **Timeout**: Delayed response
     - Document scenario configuration in test results

5. **Test scenarios with browser automation**:
   - For each user scenario from design.md:
     - **Navigate to Portal**:
       - Detect Portal URL (default: `http://localhost:3000` or from config)
       - Use `browser_navigate` to go to Portal URL
       - Wait for page load using `browser_wait_for` (wait for key text like "Dashboard" or "Invoices")
     
     - **Capture visual evidence**:
       - Take screenshot using `browser_take_screenshot`:
         - **Use relative path** (not absolute): `.oursky/scopes/{FeatureName}/screenshots/{scenario-name}-{timestamp}.png`
         - Example call: `browser_take_screenshot({ fullPage: true, filename: ".oursky/scopes/TestCommandFeedbackLoop/screenshots/scenario-1-happy-path-initial-20251121.png" })`
         - **Do NOT use**: Absolute paths like `/Users/.../formx-sino/.oursky/scopes/...` or variables like `FEATURE_DIR/screenshots/...`
         - Use fullPage: true for complete page capture
         - Generate unique filename with scenario name and timestamp
       
       - **Retrieve screenshot from temp directory** (CRITICAL - do immediately):
         - **Problem**: Browser extension saves to temp directory that may be cleaned up quickly
         - **Solution**: Try multiple retrieval methods, document if fails:
           ```bash
           SCREENSHOT_FILENAME="{scenario-name}-{timestamp}.png"
           SCREENSHOT_PATH=".oursky/scopes/{FeatureName}/screenshots/$SCREENSHOT_FILENAME"
           TEMP_PATH="/var/folders/.../cursor-browser-extension/.../.oursky/scopes/{FeatureName}/screenshots/$SCREENSHOT_FILENAME"
           
           # Method 1: Try exact path from tool result
           if [ -f "$TEMP_PATH" ]; then
             cp "$TEMP_PATH" "$SCREENSHOT_PATH" && echo "✓ Screenshot copied: $SCREENSHOT_FILENAME"
           else
             # Method 2: Search for screenshot by filename
             FOUND=$(find /var/folders -path "*cursor-browser-extension*" -name "$SCREENSHOT_FILENAME" -type f 2>/dev/null | head -1)
             if [ -n "$FOUND" ] && [ -f "$FOUND" ]; then
               cp "$FOUND" "$SCREENSHOT_PATH" && echo "✓ Screenshot found and copied: $SCREENSHOT_FILENAME"
             else
               # Method 3: Search by pattern
               FOUND=$(find /var/folders -path "*cursor-browser-extension*" -path "*scopes/{FeatureName}/screenshots/*.png" -type f 2>/dev/null | grep -E "$SCREENSHOT_FILENAME" | head -1)
               if [ -n "$FOUND" ] && [ -f "$FOUND" ]; then
                 cp "$FOUND" "$SCREENSHOT_PATH" && echo "✓ Screenshot found by pattern: $SCREENSHOT_FILENAME"
               else
                 echo "⚠ Screenshot could not be retrieved: $SCREENSHOT_FILENAME"
                 echo "⚠ Test will continue with snapshots only - screenshots are optional evidence"
                 # Document in test results that screenshot retrieval failed
               fi
             fi
           fi
           
           # Always verify final screenshot exists
           if [ -f "$SCREENSHOT_PATH" ] && [ -s "$SCREENSHOT_PATH" ]; then
             echo "✓ Screenshot verified: $SCREENSHOT_PATH ($(stat -f%z "$SCREENSHOT_PATH" 2>/dev/null || stat -c%s "$SCREENSHOT_PATH" 2>/dev/null) bytes)"
           else
             echo "✗ Screenshot not available: $SCREENSHOT_PATH"
             echo "⚠ Continuing test with snapshots only"
           fi
           ```
           **Note**: If screenshot retrieval fails, test can continue - snapshots provide layout verification evidence.
       - Get accessibility snapshot using `browser_snapshot`:
         - Parse snapshot JSON to extract layout structure
         - Save snapshot to `.oursky/scopes/{FeatureName}/screenshots/{scenario-name}-{timestamp}-snapshot.json`
         - Extract key UI elements (buttons, inputs, lists, headings)
         - **Note**: Snapshots (JSON) are saved correctly, only PNG screenshots need copying
     
     - **Verify layout structure**:
       - Compare snapshot with expected layout from design.md
       - Verify key elements are present (e.g., "Upload" button, invoice list, navigation)
       - Verify element hierarchy is correct
       - Document any missing or unexpected elements
       - Mark layout verification as pass/fail
     
     - **Execute scenario actions** (if applicable):
       - Use `browser_click` to interact with UI elements
       - Use `browser_type` to fill form fields
       - Use `browser_wait_for` to wait for state changes
       - Take additional screenshots at key interaction points
     
     - **Verify success criteria**:
       - Check each success criterion from design.md
       - Use browser snapshot to verify UI state
       - Take screenshot of final state
       - Document results (pass/fail, issues found)
   
   - **Edge cases**:
     - Test error conditions (configure FormX mock for errors)
     - Test boundary conditions (empty states, large datasets)
     - Test integration points (FormX mock responses)
     - Capture screenshots of error states

6. **Verify success criteria**:
   - Check each success criterion from design.md
   - Measure quantitative metrics if applicable
   - Document qualitative observations
   - Include screenshots as evidence for each criterion

7. **Create feedback loop**:
   - **Document findings**:
     - Create or update `FEATURE_DIR/test-results.md`:
       ```markdown
       # Test Results: [FeatureName]
       
       ## Test Date: [DATE]
       ## Test Duration: [DURATION]
       ## Database Reset: [TYPE] - [STATUS]
       ## FormX Mock Configuration: [SCENARIOS]
       
       ## Scenarios Tested
       - [X] Scenario 1: Happy Path - [PASS/FAIL]
         - Screenshot: `screenshots/scenario-1-{timestamp}.png`
         - Snapshot: `screenshots/scenario-1-{timestamp}-snapshot.json`
         - Layout Verification: [PASS/FAIL]
         - Issues: [None or description]
       
       - [ ] Scenario 2: Error Handling - [PASS/FAIL]
         - [Similar structure]
       
       ## Success Criteria Verification
       - [X] Criterion 1: [Status] - [PASS/FAIL]
         - Evidence: [Screenshot/snapshot reference]
       - [ ] Criterion 2: [Status] - [PASS/FAIL]
       
       ## Layout Verification Results
       - Home Page: [PASS/FAIL] - [Details]
       - Invoice List: [PASS/FAIL] - [Details]
       - [Other pages tested]
       
       ## Issues Found
       - [Issue description with screenshot reference]
       
       ## Next Steps
       - [What needs fixing]
       - [Suggested fixes]
       ```
   
   - **Failure detection**:
     - If any scenario fails, document failure details
     - Include screenshots showing failure state
     - Include snapshot showing layout issues
     - Suggest fixes based on failure type:
       - Layout issues → UI component fixes
       - API errors → Backend fixes
       - Integration issues → FormX mock configuration fixes
   
   - **If issues found**:
     - Document in test-results.md with evidence
     - Update tasks.md with new tasks if needed
     - Suggest specific fixes with file paths and code changes
     - Ask user how to proceed (fix now or later)
   
   - **If all tests pass**:
     - Mark feature complete in test-results.md
     - Update memory if patterns validated
     - Suggest next steps (deploy, document, etc.)

8. **Iterative testing**:
   - After fixes, re-run failed scenarios:
     - Re-configure FormX mock if needed
     - Re-run database reset if needed
     - Re-execute browser automation for failed scenarios
     - Update test-results.md with new results
   - Continue until all scenarios pass or max iterations reached
   - Track iteration count in test-results.md

9. **Generate summary report**:
   - Summary of test results:
     - Total scenarios: X
     - Passed: Y
     - Failed: Z
     - Success criteria met: A/B
   - List all issues found with priorities
   - Provide next steps for fixing failures
   - Include links to screenshots and snapshots

## FormX Mock API Configuration

### API Endpoints

Base URL: `http://localhost:3001` (or from config)

- `GET /api/matchers` - List all response matchers
- `POST /api/matchers` - Create new response matcher
- `PUT /api/matchers/:id` - Update response matcher
- `DELETE /api/matchers/:id` - Delete response matcher
- `GET /api/jobs` - List active jobs
- `POST /api/jobs/:jobId/status` - Update job status

### Scenario Configuration Helpers

**Happy Path Scenario**:
```json
{
  "name": "Happy Path - Successful Extraction",
  "enabled": true,
  "pattern": {
    "endpoint": "/v2/workspace",
    "method": "POST"
  },
  "response": {
    "statusCode": 200,
    "body": {
      "status": "ok",
      "extraction_id": "ext_123",
      "workspace_id": "ws_abc",
      "data": { /* extraction data */ }
    }
  }
}
```

**Error Scenario (Rate Limit)**:
```json
{
  "name": "Rate Limit Error",
  "enabled": true,
  "pattern": {
    "endpoint": "/v2/workspace",
    "method": "POST"
  },
  "response": {
    "statusCode": 429,
    "body": {
      "status": "error",
      "error": {
        "code": "RATE_LIMIT",
        "message": "Rate limit exceeded"
      }
    }
  }
}
```

**Async Job Scenario**:
- Create matcher for `/v2/workspace` returning job_id
- Create matcher for `/v2/extract/jobs/:job_id` returning status progression
- Use `POST /api/jobs/:jobId/status` to update job status during test

### Programmatic Configuration

1. **Clear existing matchers**:
   - GET `/api/matchers`
   - DELETE each matcher by ID

2. **Create scenario matchers**:
   - POST `/api/matchers` with scenario configuration

3. **Update job states** (for async scenarios):
   - POST `/api/jobs/:jobId/status` with new status

## Database Reset Options

### Reset Types

- **"full"** (option 1): Drop and recreate database (use for clean slate)
  - Stops services, removes volume, recreates, runs migrations
- **"clear-data"** (option 2): Truncate all tables (use for quick reset)
  - Keeps schema, clears all data
- **"migrations"** (option 3): Rollback and reapply migrations
  - Use for schema changes testing
- **"volumes"** (option 4): Clear all Docker volumes
  - Removes DB, Redis, MinIO data

### Non-Interactive Execution

To execute reset non-interactively:
```bash
echo -e "1\nyes" | ./scripts/reset-db.sh  # Full reset
echo -e "2\nyes" | ./scripts/reset-db.sh  # Clear data
```

## Scenario Implementations

**Important**: After each `browser_take_screenshot` call, screenshots are saved to a temp directory. **You MUST copy them to the feature directory**.

**Use relative paths** (not absolute) when calling `browser_take_screenshot`:
- Format: `.oursky/scopes/{FeatureName}/screenshots/{scenario-name}-{timestamp}.png`
- Example call: `browser_take_screenshot({ fullPage: true, filename: ".oursky/scopes/TestCommandFeedbackLoop/screenshots/scenario-1-happy-path-initial-20251121.png" })`
- **Do NOT use absolute paths** like `/Users/.../formx-sino/.oursky/scopes/...` or variables like `FEATURE_DIR/screenshots/...`

**Copy screenshots from temp directory** (verify file exists and copy immediately after taking screenshot):
```bash
# Method 1: Use exact path from browser_take_screenshot tool result (RECOMMENDED)
# Extract temp path immediately after calling browser_take_screenshot
# The tool result shows: "/var/folders/.../cursor-browser-extension/.../.oursky/scopes/{FeatureName}/screenshots/{filename}.png"
TEMP_PATH="/var/folders/.../cursor-browser-extension/.../.oursky/scopes/{FeatureName}/screenshots/{filename}.png"

# Verify file exists before copying
if [ -f "$TEMP_PATH" ]; then
  cp "$TEMP_PATH" ".oursky/scopes/{FeatureName}/screenshots/" && echo "✓ Copied: $(basename "$TEMP_PATH")"
  
  # Verify screenshot exists in feature directory after copying
  FINAL_PATH=".oursky/scopes/{FeatureName}/screenshots/$(basename "$TEMP_PATH")"
  if [ -f "$FINAL_PATH" ]; then
    echo "✓ Screenshot verified: $FINAL_PATH"
    ls -lh "$FINAL_PATH"
  else
    echo "✗ Screenshot not found in feature directory: $FINAL_PATH"
  fi
else
  echo "⚠ Screenshot not found at temp path: $TEMP_PATH"
  echo "Attempting to find screenshot in temp directory..."
  # Fall back to Method 2
fi

# Method 2: Find and copy all screenshots (if temp directory still exists)
find /var/folders -path "*cursor-browser-extension*" -path "*scopes/{FeatureName}/screenshots/*.png" -type f 2>/dev/null | \
  while read file; do 
    if [ -f "$file" ]; then
      cp "$file" ".oursky/scopes/{FeatureName}/screenshots/$(basename "$file")"
      FINAL_PATH=".oursky/scopes/{FeatureName}/screenshots/$(basename "$file")"
      if [ -f "$FINAL_PATH" ]; then
        echo "✓ Copied and verified: $(basename "$file")"
      else
        echo "✗ Failed to copy: $(basename "$file")"
      fi
    fi
  done
```

**Important**: 
- Temp directory paths are session-specific and may be cleaned up quickly
- Always verify file exists before copying
- Always verify screenshot exists in feature directory after copying
- Use `[ -f "$PATH" ]` to check file existence

### Scenario 1: Happy Path Testing

**Purpose**: Verify successful document upload and extraction flow

**Steps**:
1. **Reset database**:
   - Execute: `echo -e "1\nyes" | ./scripts/reset-db.sh` (full reset)
   - Verify: Check database connection and empty state
   - Log: "Database reset: Full reset completed"

2. **Configure FormX mock for successful extraction**:
   - Clear all existing matchers: `GET /api/matchers`, then `DELETE /api/matchers/:id` for each
   - Create happy path matcher:
     ```json
     POST /api/matchers
     {
       "name": "Happy Path - Successful Extraction",
       "enabled": true,
       "pattern": { "endpoint": "/v2/workspace", "method": "POST" },
       "response": {
         "statusCode": 200,
         "body": {
           "status": "ok",
           "extraction_id": "ext_happy_123",
           "workspace_id": "ws_happy_abc",
           "job_id": "job_happy_456",
           "data": { /* sample extraction data */ }
         }
       }
     }
     ```
   - Log: "FormX mock configured: Happy Path scenario"

3. **Navigate to Portal home page**:
   - Use `browser_navigate` to go to `http://localhost:3000` (or configured Portal URL)
   - Wait for page load: `browser_wait_for({ text: "Jobs" })` or similar key text
   - Log: "Navigated to Portal home page"

4. **Take screenshot and snapshot**:
   - Screenshot: `browser_take_screenshot({ fullPage: true, filename: ".oursky/scopes/{FeatureName}/screenshots/scenario-1-happy-path-initial-{timestamp}.png" })`
   - **Verify and copy screenshot immediately**:
     - Extract temp path from tool result
     - Verify file exists: `[ -f "$TEMP_PATH" ]`
     - Copy to feature directory if exists
     - Verify screenshot exists in feature directory after copying
     - Log success or failure
   - Snapshot: `browser_snapshot()` → Save to `.oursky/scopes/{FeatureName}/screenshots/scenario-1-happy-path-initial-{timestamp}-snapshot.json`
   - Log: "Captured initial state screenshot and snapshot"

5. **Verify layout structure**:
   - Parse snapshot JSON to extract UI elements
   - Verify key elements present:
     - Navigation menu/header
     - Main content area
     - Upload button or upload area (if applicable)
     - Invoice/collection list (if applicable)
   - Compare with expected layout from design.md
   - Document: "Layout verification: [PASS/FAIL] - [Details]"

6. **Upload document** (if applicable):
   - Locate upload button/area in snapshot
   - Use `browser_click` to click upload button
   - Wait for upload dialog: `browser_wait_for({ text: "Upload" })`
   - Use `browser_type` or file upload to select document
   - Click submit/upload button
   - Wait for upload completion: `browser_wait_for({ text: "Success" })` or similar
   - Take screenshot: `scenario-1-happy-path-upload-{timestamp}.png`
   - Log: "Document uploaded successfully"

7. **Verify extraction completes**:
   - Wait for extraction status: `browser_wait_for({ text: "Completed" })` or check job status
   - Take snapshot: `scenario-1-happy-path-extraction-{timestamp}-snapshot.json`
   - Verify extraction data appears in UI
   - Log: "Extraction completed successfully"

8. **Navigate to invoice list**:
   - Use `browser_navigate` to go to `/invoices` or click invoice list link
   - Wait for page load: `browser_wait_for({ text: "Invoices" })`
   - Take screenshot: `scenario-1-happy-path-invoice-list-{timestamp}.png`
   - Take snapshot: `scenario-1-happy-path-invoice-list-{timestamp}-snapshot.json`

9. **Verify invoice appears in list**:
   - Parse snapshot to find invoice list items
   - Verify uploaded invoice appears in list
   - Verify invoice data is correct (if visible)
   - Document: "Invoice list verification: [PASS/FAIL]"

10. **Document results**:
    - Update `test-results.md`:
      ```markdown
      - [X] Scenario 1: Happy Path - [PASS/FAIL]
        - Screenshots: 
          - Initial: `screenshots/scenario-1-happy-path-initial-{timestamp}.png`
          - Upload: `screenshots/scenario-1-happy-path-upload-{timestamp}.png`
          - Invoice List: `screenshots/scenario-1-happy-path-invoice-list-{timestamp}.png`
        - Snapshots:
          - Initial: `screenshots/scenario-1-happy-path-initial-{timestamp}-snapshot.json`
          - Invoice List: `screenshots/scenario-1-happy-path-invoice-list-{timestamp}-snapshot.json`
        - Layout Verification: [PASS/FAIL]
        - Issues: [None or description]
      ```

### Scenario 2: Error Handling Testing

**Purpose**: Verify error messages display correctly and layout handles errors gracefully

**Steps**:
1. **Reset database**:
   - Execute: `echo -e "2\nyes" | ./scripts/reset-db.sh` (clear data)
   - Log: "Database reset: Data cleared"

2. **Configure FormX mock for error response**:
   - Clear all existing matchers
   - Create rate limit error matcher:
     ```json
     POST /api/matchers
     {
       "name": "Rate Limit Error",
       "enabled": true,
       "pattern": { "endpoint": "/v2/workspace", "method": "POST" },
       "response": {
         "statusCode": 429,
         "body": {
           "status": "error",
           "error": {
             "code": "RATE_LIMIT",
             "message": "Rate limit exceeded. Please try again later."
           }
         }
       }
     }
     ```
   - Log: "FormX mock configured: Rate Limit Error scenario"

3. **Navigate to Portal upload page**:
   - Use `browser_navigate` to go to `/uploads` or upload page
   - Wait for page load: `browser_wait_for({ text: "Upload" })`
   - Take screenshot: `scenario-2-error-handling-initial-{timestamp}.png`

4. **Attempt to upload document**:
   - Click upload button
   - Select document and submit
   - Wait for error response: `browser_wait_for({ text: "Rate limit" })` or error indicator
   - Log: "Upload attempted, error response received"

5. **Verify error message displays correctly**:
   - Take screenshot: `scenario-2-error-handling-error-state-{timestamp}.png`
   - Take snapshot: `scenario-2-error-handling-error-state-{timestamp}-snapshot.json`
   - Parse snapshot to find error message element
   - Verify error message text contains "Rate limit" or similar
   - Verify error message is visible and readable
   - Document: "Error message verification: [PASS/FAIL]"

6. **Verify layout handles error gracefully**:
   - Check snapshot for layout structure
   - Verify error doesn't break page layout
   - Verify error message is positioned appropriately
   - Verify other UI elements remain functional
   - Document: "Layout error handling: [PASS/FAIL]"

7. **Document results**:
   - Update `test-results.md`:
     ```markdown
     - [X] Scenario 2: Error Handling - [PASS/FAIL]
       - Screenshots:
         - Initial: `screenshots/scenario-2-error-handling-initial-{timestamp}.png`
         - Error State: `screenshots/scenario-2-error-handling-error-state-{timestamp}.png`
       - Snapshot:
         - Error State: `screenshots/scenario-2-error-handling-error-state-{timestamp}-snapshot.json`
       - Error Message Verification: [PASS/FAIL]
       - Layout Error Handling: [PASS/FAIL]
       - Issues: [None or description]
     ```

### Scenario 3: Async Job Testing

**Purpose**: Verify async job status progression and UI updates

**Steps**:
1. **Reset database**:
   - Execute: `echo -e "1\nyes" | ./scripts/reset-db.sh` (full reset)
   - Log: "Database reset: Full reset completed"

2. **Configure FormX mock for async job progression**:
   - Clear all existing matchers
   - Create async job matchers:
     ```json
     # Workspace upload matcher (returns job_id)
     POST /api/matchers
     {
       "name": "Async Job - Workspace Upload",
       "enabled": true,
       "pattern": { "endpoint": "/v2/workspace", "method": "POST" },
       "response": {
         "statusCode": 202,
         "body": {
           "status": "ok",
           "job_id": "job_async_123",
           "request_id": "req_async_456"
         }
       }
     }
     
     # Job status matcher (initially pending)
     POST /api/matchers
     {
       "name": "Async Job - Job Status",
       "enabled": true,
       "pattern": { "endpoint": "/v2/extract/jobs/*", "method": "GET" },
       "response": {
         "statusCode": 200,
         "body": {
           "status": "pending",
           "job_id": "job_async_123"
         }
       }
     }
     ```
   - Log: "FormX mock configured: Async Job scenario"

3. **Navigate to Portal**:
   - Use `browser_navigate` to go to Portal home
   - Wait for page load
   - Log: "Navigated to Portal"

4. **Upload document**:
   - Click upload button
   - Select document and submit
   - Wait for job creation: `browser_wait_for({ text: "Processing" })` or job ID
   - Take screenshot: `scenario-3-async-job-upload-{timestamp}.png`
   - Take snapshot: `scenario-3-async-job-pending-{timestamp}-snapshot.json`
   - Log: "Document uploaded, job created"

5. **Verify job status updates**:
   - **Pending state**: Verify UI shows "Pending" or "Processing" status
     - Take screenshot: `scenario-3-async-job-pending-{timestamp}.png`
     - Document: "Pending state verified: [PASS/FAIL]"
   
   - **Update to Processing**:
     - Update job status: `POST /api/jobs/job_async_123/status` with `{ "status": "processing" }`
     - Wait for UI update: `browser_wait_for({ text: "Processing" })`
     - Take screenshot: `scenario-3-async-job-processing-{timestamp}.png`
     - Take snapshot: `scenario-3-async-job-processing-{timestamp}-snapshot.json`
     - Document: "Processing state verified: [PASS/FAIL]"
   
   - **Update to Completed**:
     - Update job status: `POST /api/jobs/job_async_123/status` with `{ "status": "completed", "data": { /* extraction data */ } }`
     - Wait for completion: `browser_wait_for({ text: "Completed" })`
     - Take screenshot: `scenario-3-async-job-completed-{timestamp}.png`
     - Take snapshot: `scenario-3-async-job-completed-{timestamp}-snapshot.json`
     - Document: "Completed state verified: [PASS/FAIL]"

6. **Verify final state shows completed extraction**:
   - Parse final snapshot to verify extraction data appears
   - Verify UI shows completed status
   - Verify extraction results are displayed
   - Document: "Final state verification: [PASS/FAIL]"

7. **Document results**:
   - Update `test-results.md`:
     ```markdown
     - [X] Scenario 3: Async Job Testing - [PASS/FAIL]
       - Screenshots:
         - Upload: `screenshots/scenario-3-async-job-upload-{timestamp}.png`
         - Pending: `screenshots/scenario-3-async-job-pending-{timestamp}.png`
         - Processing: `screenshots/scenario-3-async-job-processing-{timestamp}.png`
         - Completed: `screenshots/scenario-3-async-job-completed-{timestamp}.png`
       - Snapshots:
         - Pending: `screenshots/scenario-3-async-job-pending-{timestamp}-snapshot.json`
         - Processing: `screenshots/scenario-3-async-job-processing-{timestamp}-snapshot.json`
         - Completed: `screenshots/scenario-3-async-job-completed-{timestamp}-snapshot.json`
       - Status Progression: Pending → Processing → Completed
       - Issues: [None or description]
     ```

### Scenario 4: Multiple Scenarios Iteration

**Purpose**: Run all scenarios in sequence and support iterative testing

**Steps**:
1. **Initialize test results**:
   - Create or reset `test-results.md`
   - Set iteration count: `iteration = 1`
   - Set max iterations: `maxIterations = 3` (configurable)

2. **Run Scenario 1 (Happy Path)**:
   - Execute Scenario 1 steps
   - Record result: `scenario1Result = [PASS/FAIL]`
   - If FAIL: Document issues in test-results.md

3. **Run Scenario 2 (Error Handling)**:
   - Execute Scenario 2 steps
   - Record result: `scenario2Result = [PASS/FAIL]`
   - If FAIL: Document issues in test-results.md

4. **Run Scenario 3 (Async Job)**:
   - Execute Scenario 3 steps
   - Record result: `scenario3Result = [PASS/FAIL]`
   - If FAIL: Document issues in test-results.md

5. **Check if any scenarios failed**:
   - If all pass: Proceed to final summary
   - If any fail:
     - Document all failures in test-results.md
     - Suggest fixes based on failure types
     - Ask user: "Some scenarios failed. Should I re-run failed scenarios after fixes?"
     - If user approves and `iteration < maxIterations`:
       - Increment iteration count
       - Re-run only failed scenarios
       - Update test-results.md with new results
       - Repeat until all pass or max iterations reached

6. **Generate final summary report**:
   - Calculate totals:
     - Total scenarios: 3
     - Passed: [count]
     - Failed: [count]
     - Success criteria met: [count]/[total]
   - List all issues found with priorities
   - Provide next steps for fixing failures
   - Include links to all screenshots and snapshots
   - Update test-results.md with summary

7. **Document results**:
   - Update `test-results.md`:
     ```markdown
     ## Iteration Summary
     - Iteration 1: Scenario 1: [PASS/FAIL], Scenario 2: [PASS/FAIL], Scenario 3: [PASS/FAIL]
     - Iteration 2: [If re-run] Scenario X: [PASS/FAIL]
     - Final Status: All scenarios [PASS/FAIL]
     
     ## Final Summary
     - Total Scenarios: 3
     - Passed: [count]
     - Failed: [count]
     - Success Criteria Met: [count]/[total]
     - Issues Found: [list]
     - Next Steps: [recommendations]
     ```

## Browser Tool Usage

### Available Tools

- `browser_navigate(url)` - Navigate to URL
- `browser_snapshot()` - Get accessibility snapshot (JSON structure)
- `browser_take_screenshot(options)` - Capture screenshot
  - Options: `{ fullPage: true, filename: "path.png" }`
- `browser_wait_for({ text, time })` - Wait for text or time
- `browser_click({ element, ref })` - Click element
- `browser_type({ element, ref, text })` - Type text
- `browser_fill_form({ fields })` - Fill multiple form fields

### Layout Verification

1. **Get snapshot**: Use `browser_snapshot()` to get page structure
2. **Parse snapshot**: Extract key elements (buttons, inputs, headings)
3. **Compare with expected**: Check against expected layout from design.md
4. **Document differences**: Log missing or unexpected elements

### Screenshot Naming

Format: `{scenario-name}-{timestamp}-{step}.png`

Examples:
- `happy-path-20250127-143022-initial.png`
- `error-handling-20250127-143045-error-state.png`
- `async-job-20250127-143100-pending.png`

### Screenshot Path Handling

**CRITICAL**: Browser extension saves screenshots to temp directories that may be cleaned up quickly, making retrieval unreliable. **Use Method 1 (browser_evaluate) for reliable screenshot capture**.

**Method 1: Capture screenshot as base64 (RECOMMENDED - Most Reliable)**:
- Use `browser_evaluate` to capture screenshot as base64 data directly from page
- Save base64 data directly to feature directory
- Avoids temp directory issues entirely
- Example:
  ```javascript
  // Capture screenshot as base64
  const result = await browser_evaluate({
    function: "() => {
      return new Promise((resolve) => {
        // Use html2canvas if available, or canvas API
        const canvas = document.createElement('canvas');
        canvas.width = window.innerWidth;
        canvas.height = document.body.scrollHeight;
        // ... capture logic ...
        resolve(canvas.toDataURL('image/png'));
      });
    }"
  });
  
  // Save base64 to file
  const base64Data = result.replace(/^data:image\/png;base64,/, '');
  const screenshotPath = '.oursky/scopes/{FeatureName}/screenshots/{filename}.png';
  // Write base64 data to file using base64 -d
  ```
  
  ```bash
  # Save base64 screenshot data
  BASE64_DATA="[from browser_evaluate result]"
  SCREENSHOT_PATH=".oursky/scopes/{FeatureName}/screenshots/{filename}.png"
  echo "$BASE64_DATA" | sed 's/^data:image\/png;base64,//' | base64 -d > "$SCREENSHOT_PATH"
  
  # Verify
  if [ -f "$SCREENSHOT_PATH" ]; then
    echo "✓ Screenshot saved: $SCREENSHOT_PATH"
    ls -lh "$SCREENSHOT_PATH"
  fi
  ```

**Method 2: Use browser_take_screenshot with immediate copy (Fallback)**:
- **Use relative paths from workspace root** when calling `browser_take_screenshot`:
  - Format: `.oursky/scopes/{FeatureName}/screenshots/{filename}.png`
  - Example: `.oursky/scopes/TestCommandFeedbackLoop/screenshots/scenario-1-happy-path-initial-20251121.png`

- **Copy immediately after taking screenshot**:
  ```bash
  # Extract temp path from browser_take_screenshot tool result
  TEMP_PATH="/var/folders/.../cursor-browser-extension/.../.oursky/scopes/{FeatureName}/screenshots/{filename}.png"
  SCREENSHOT_PATH=".oursky/scopes/{FeatureName}/screenshots/{filename}.png"
  
  # Try to copy immediately (temp directory may be cleaned up quickly)
  if [ -f "$TEMP_PATH" ]; then
    cp "$TEMP_PATH" "$SCREENSHOT_PATH" && echo "✓ Copied: $(basename "$TEMP_PATH")"
  else
    echo "⚠ Screenshot not found at temp path - temp directory may have been cleaned up"
    echo "⚠ Consider using Method 1 (browser_evaluate) for more reliable capture"
  fi
  
  # Always verify screenshot exists
  if [ -f "$SCREENSHOT_PATH" ]; then
    echo "✓ Screenshot verified: $SCREENSHOT_PATH"
    ls -lh "$SCREENSHOT_PATH"
  else
    echo "✗ Screenshot not found: $SCREENSHOT_PATH"
  fi
  ```

**Note**: 
- Snapshots (JSON files) are saved correctly to the feature directory - no copying needed
- Method 1 (browser_evaluate) is more reliable as it avoids temp directory issues
- Method 2 may fail if temp directory is cleaned up before copying
- Always verify screenshot exists in feature directory after saving

## Configuration Options

### Environment Variables

- `PORTAL_URL` - Portal URL (default: `http://localhost:3000`)
- `FORMX_MOCK_URL` - FormX mock URL (default: `http://localhost:3001`)
- `API_URL` - API URL (default: `http://localhost:8080`)
- `DB_RESET_TYPE` - Database reset type (default: `full`)
- `SCREENSHOT_DIR` - Screenshot directory (default: `FEATURE_DIR/screenshots/`)

### Command Arguments

The command accepts arguments in format: `[feature-name] [options]`

Options:
- `--reset-type <type>` - Database reset type (full, clear-data, migrations, volumes)
- `--portal-url <url>` - Override Portal URL
- `--formx-mock-url <url>` - Override FormX mock URL
- `--scenarios <list>` - Comma-separated list of scenarios to run
- `--no-browser` - Skip browser automation (manual testing only)

## Guidelines

- **Test real app**: Don't just run unit tests, test the full stack
- **Reset cleanly**: Start with clean database for consistent results
- **Use browser automation**: Automate visual verification when possible
- **Capture evidence**: Always take screenshots and snapshots
- **Document scenarios**: Keep FormX mock scenarios documented in test results
- **Iterate**: Fix issues and re-test until passing
- **Update memory**: Document testing patterns that work

## Common Test Scenarios

- **Happy path**: Primary user flow works end-to-end
  - Configure FormX mock for success
  - Navigate to Portal, verify layout
  - Execute primary action
  - Verify success state with screenshot

- **Error handling**: Appropriate errors for invalid inputs
  - Configure FormX mock for error (429, 500, etc.)
  - Navigate to Portal, attempt action
  - Verify error message displays correctly
  - Capture screenshot of error state

- **Async jobs**: Job status progression
  - Configure FormX mock for async job
  - Upload document, verify pending state
  - Update job status, verify processing state
  - Complete job, verify final state
  - Capture screenshots at each transition

- **Edge cases**: Boundary conditions, empty states
  - Test with empty database
  - Test with large datasets
  - Test with missing data
  - Verify UI handles edge cases gracefully

- **Integration**: External services (FormX) work correctly
  - Test different FormX mock scenarios
  - Verify API integration points
  - Verify error propagation

## Error Handling

- **Browser tool failures**: Log warning, continue with manual testing guidance
- **Database reset failures**: Log error, suggest manual reset, continue if possible
- **FormX mock API failures**: Log error, suggest manual configuration, continue if possible
- **Service unavailability**: Log error, suggest starting services, abort test
- **Layout verification failures**: Document in test results with evidence, continue testing

## Test Results Format

The `test-results.md` file should include:

1. **Header**: Feature name, test date, duration
2. **Environment**: Database reset type, FormX mock scenarios, service URLs
3. **Scenarios**: Each scenario with:
   - Status (pass/fail)
   - Screenshot references
   - Snapshot references
   - Layout verification results
   - Issues found
4. **Success Criteria**: Each criterion with status and evidence
5. **Layout Verification**: Per-page results
6. **Issues**: Detailed issue descriptions with evidence
7. **Next Steps**: Suggested fixes and actions

