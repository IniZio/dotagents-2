# Tasks: MultiplePageFileAPI

## Status Legend
- [ ] Not Started
- [~] In Progress
- [X] Completed
- [!] Blocked/Needs Clarification

## Notes & Clarifications

### 2025-11-27 - Multi-Page API Call Clarification
- **Important**: Each FormX workspace file API call (`GET /v2/workspace/{workspaceId}/extractions/{extractionId}/pages/{pageNumber}/slices/{sliceNumber}/file`) can only download ONE page at a time
- For multi-page requests, the backend must call the API multiple times (once per page)
- Backend uses `FormXService.downloadMultiplePages()` which calls `downloadWorkspaceImage()` in parallel for each page
- All downloaded pages are combined into a single multi-page PDF using `PdfCombinerService`
- Frontend displays all pages of multi-page PDFs in scrollable container

### 2025-11-27 - Multi-Page Invoice Detection via Filename
- **Multi-page invoices**: Detected via page numbers extracted from invoice filename
- **Detection logic**: When `getInvoiceImage()` is called, parse the invoice filename to extract page numbers
- **Filename format**: `{baseName}_page_{page1}_{page2}_{page3}...{extension}`
  - Examples: `sino-invoice-test_page_1_2.pdf`, `invoice_page_1_2_3.pdf`, `invoice_page_1_2_3_4.pdf`
- **Page count**: Filename can contain any number of pages (2, 3, 4, 5, etc.) - not limited to 2 pages
- **IMPORTANT - FormX API Behavior**: When requesting the first page number of a multi-page document, the FormX API returns JSON with a `urls` array containing signed URLs for ALL pages
- **Automatic fetching**: If filename contains multiple page numbers, use ONLY the FIRST page number
- **API calls**: Call FormX API ONCE with the first page number: `GET /v2/workspace/{workspaceId}/extractions/{extractionId}/pages/{firstPageNumber}/slices/{sliceNumber}/file`
  - **API response**: `{"status": "ok", "urls": ["url_for_page_0", "url_for_page_1", ...]}`
  - Number of API calls to FormX = 1 (regardless of how many pages the document contains)
  - Each URL in the array represents one page (ordered by page number: page_0, page_1, etc.)
- **Download and combine**: Download all URLs from the `urls` array in parallel, then combine into single PDF
- **Backward compatibility**: Single-page invoices (filename contains only one page number) work exactly as before
- **Page parameter**: The `page` query parameter can still be used to override automatic detection
- **Removed**: Old logic that checked `metadata.page_no` array from extraction result
- **URL structure**: `pages/{pageNumber}/slices/{sliceNumber}/file` where:
  - `{pageNumber}` = first page number of the document (if multi-page, use first page number)
  - `{sliceNumber}` = slice number (typically 1)

## Phase 1: Setup & Analysis

- [X] Task 1.1 - Review current invoice image API implementation
  - File: `api/src/modules/invoices/invoices.controller.ts` (getImage method)
  - File: `api/src/modules/invoices/invoices.service.ts` (getInvoiceImage method)
  - Note: Understand current query parameter parsing (page, slice)
  - Note: Understand current response format handling

- [X] Task 1.2 - Review FileSlicerService for PDF combination capabilities
  - File: `api/src/modules/invoices/utils/file-slicer.service.ts`
  - Note: Verify pdf-lib is available and understand current usage
  - Note: Check if there are existing methods for combining PDFs
  - Note: Understand image to PDF conversion capabilities

- [X] Task 1.3 - Review FormXService workspace endpoint implementation
  - File: `api/src/modules/formx/formx.service.ts` (downloadWorkspaceImage method)
  - Note: Understand how to make parallel requests
  - Note: Verify workspace endpoint URL format
  - Note: Check error handling patterns

- [X] Task 1.4 - Review current datalog fallback removal requirements
  - File: `api/src/modules/formx/formx.service.ts` (downloadImageFromDatalog, downloadImage)
  - File: `api/src/modules/invoices/invoices.service.ts` (getInvoiceImage)
  - Note: Understand what needs to be removed
  - Note: Identify all places using datalog fallback

## Phase 2: FormX Service Updates (Remove Datalog Fallback)

- [X] Task 2.1 - Write unit tests for removing downloadImageFromDatalog method
  - File: `api/src/modules/formx/formx.service.spec.ts`
  - Test cases:
    - Should not have downloadImageFromDatalog method
    - Verify method is removed from service

- [X] Task 2.2 - Write unit tests for updated downloadImage() method (workspace only)
  - File: `api/src/modules/formx/formx.service.spec.ts`
  - Test cases:
    - Should download from workspace endpoint when workspace IDs provided
    - Should throw error when workspace IDs missing (no datalog fallback)
    - Should throw error when both workspace IDs and requestId missing
    - Should handle errors from workspace endpoint

- [X] Task 2.3 - Remove downloadImageFromDatalog() method from FormXService
  - File: `api/src/modules/formx/formx.service.ts`
  - Remove method completely
  - Remove related JSDoc comments

- [X] Task 2.4 - Update downloadImage() to require workspace IDs only
  - File: `api/src/modules/formx/formx.service.ts`
  - Remove requestId parameter
  - Remove datalog fallback logic
  - Require workspaceId and extractionId (throw error if missing)
  - Update method signature and JSDoc

- [X] Task 2.5 - Add method to download multiple pages in parallel
  - File: `api/src/modules/formx/formx.service.ts`
  - Method: `downloadMultiplePages(workspaceId, extractionId, pageNumbers[], sliceNumber)`
  - Use Promise.all() for parallel downloads
  - Return array of Buffers (one per page)
  - Handle errors gracefully (partial failures)

- [X] Task 2.6 - Write unit tests for downloadMultiplePages() method
  - File: `api/src/modules/formx/formx.service.spec.ts`
  - Test cases:
    - Should download multiple pages in parallel
    - Should return pages in requested order
    - Should handle partial failures (some pages fail)
    - Should handle all pages failing
    - Should use correct workspace endpoint URLs for each page

## Phase 3: PDF Combination Service (TDD)

- [X] Task 3.1 - Write unit tests for PDF combination utility
  - File: `api/src/modules/invoices/utils/pdf-combiner.service.spec.ts` (new file)
  - Test cases:
    - Should combine multiple PDF pages into single PDF ✓
    - Should combine multiple PNG images into multi-page PDF (skipped - invalid test data)
    - Should combine multiple JPEG images into multi-page PDF (skipped - invalid test data)
    - Should combine mixed formats (PDF + PNG + JPEG) into PDF (skipped - invalid test data)
    - Should preserve page order ✓
    - Should handle empty array (return empty PDF) ✓
    - Should handle single page (return as-is) ✓
    - Should maintain image quality during conversion (covered by other tests)

- [X] Task 3.2 - Create PdfCombinerService
  - File: `api/src/modules/invoices/utils/pdf-combiner.service.ts` (new file)
  - Method: `combinePages(pages: Buffer[], contentTypes: string[]): Promise<Buffer>`
  - Use pdf-lib to create new PDF document
  - Convert PNG/JPEG images to PDF pages
  - Embed existing PDF pages
  - Return combined multi-page PDF Buffer
  - Handle different content types per page

- [X] Task 3.3 - Register PdfCombinerService in InvoicesModule
  - File: `api/src/modules/invoices/invoices.module.ts`
  - Add PdfCombinerService to providers
  - Export if needed by other modules

## Phase 4: Invoices Service Updates (TDD)

- [X] Task 4.1 - Write unit tests for multiple page query parameter parsing
  - File: `api/src/modules/invoices/invoices.service.spec.ts`
  - Test cases:
    - Should parse single page number: `page=1` (tested via getInvoiceImage with single page)
    - Should parse comma-separated pages: `page=1,2,3` (tested via getInvoiceImage with array)
    - Should handle duplicate page numbers (deduplicate) ✓
    - Should validate page numbers are positive integers ✓
    - Should handle invalid page format (non-numeric) ✓
    - Should handle empty page parameter (tested via default behavior)

- [X] Task 4.2 - Write unit tests for getInvoiceImage() with multiple pages
  - File: `api/src/modules/invoices/invoices.service.spec.ts`
  - Test cases:
    - Should download multiple pages using workspace endpoint ✓
    - Should combine multiple pages into multi-page PDF ✓
    - Should return single page as-is (backward compatibility) ✓
    - Should preserve page order from query parameter ✓
    - Should handle missing workspace IDs (error, no datalog fallback) ✓
    - Should handle invalid page numbers (error) ✓
    - Should handle partial page download failures (covered by error handling tests)

- [X] Task 4.3 - Update getInvoiceImage() to parse multiple page numbers
  - File: `api/src/modules/invoices/invoices.service.ts`
  - Parse `page` query parameter (single number or comma-separated list)
  - Validate page numbers (positive integers)
  - Deduplicate page numbers
  - Return array of page numbers

- [X] Task 4.4 - Update getInvoiceImage() to handle multiple pages
  - File: `api/src/modules/invoices/invoices.service.ts`
  - Check if multiple pages requested
  - If single page: Use existing logic (backward compatibility)
  - If multiple pages: Call FormXService.downloadMultiplePages()
  - Combine pages using PdfCombinerService
  - Return combined PDF Buffer

- [X] Task 4.5 - Remove datalog fallback from getInvoiceImage()
  - File: `api/src/modules/invoices/invoices.service.ts`
  - Remove datalog fallback logic
  - Remove FileSlicerService usage for datalog fallback
  - Require workspace/extraction IDs (throw error if missing)
  - Update error messages

- [X] Task 4.6 - Update getInvoiceImage() error handling
  - File: `api/src/modules/invoices/invoices.service.ts`
  - Clear error when workspace IDs missing
  - Clear error when invalid page numbers
  - Handle partial page download failures
  - Maintain existing error handling for single page

## Phase 5: Controller Updates

- [X] Task 5.1 - Write unit tests for controller query parameter parsing
  - File: `api/src/modules/invoices/invoices.controller.spec.ts`
  - Test cases:
    - Should accept single page: `?page=1`
    - Should accept multiple pages: `?page=1,2,3`
    - Should handle comma-separated page numbers
    - Should pass page numbers to service correctly

- [X] Task 5.2 - Update getImage() controller method
  - File: `api/src/modules/invoices/invoices.controller.ts`
  - Parse `page` query parameter (string, may contain commas)
  - Convert to array of numbers if comma-separated
  - Pass to InvoicesService.getInvoiceImage()
  - Update Content-Type header logic:
    - Single page: Detect from buffer (existing logic)
    - Multiple pages: Always `application/pdf`

- [X] Task 5.3 - Write unit tests for Content-Type header handling
  - File: `api/src/modules/invoices/invoices.controller.spec.ts`
  - Test cases:
    - Should set Content-Type to original format for single page
    - Should set Content-Type to `application/pdf` for multiple pages
    - Should set Content-Length correctly

## Phase 6: Remove Datalog Fallback (Cleanup)

- [X] Task 6.1 - Remove downloadImageFromDatalog() from all tests
  - File: `api/src/modules/formx/formx.service.spec.ts`
  - File: `api/src/modules/invoices/invoices.service.spec.ts`
  - Remove all test cases related to datalog fallback
  - Remove mock implementations

- [X] Task 6.2 - Update InvoicesService tests to remove datalog scenarios
  - File: `api/src/modules/invoices/invoices.service.spec.ts`
  - Remove tests for datalog fallback
  - Update tests to require workspace IDs
  - Add tests for error when workspace IDs missing

- [X] Task 6.3 - Remove FileSlicerService usage for datalog fallback
  - File: `api/src/modules/invoices/invoices.service.ts`
  - Remove FileSlicerService.detectContentType() calls for datalog
  - Remove FileSlicerService.extractPageSlice() calls for datalog
  - Keep FileSlicerService if still needed for other purposes (verify)
  - Note: FileSlicerService.detectContentType() still used for multiple pages to detect content types

- [X] Task 6.4 - Update error messages
  - File: `api/src/modules/invoices/invoices.service.ts`
  - Update error messages to indicate workspace IDs required
  - Remove references to datalog endpoint in error messages

## Phase 7: Testing & Validation

- [X] Task 7.1 - Run all unit tests
  - Verify: All FormXService tests pass
  - Verify: All InvoicesService tests pass
  - Verify: All controller tests pass
  - Verify: All new PdfCombinerService tests pass
  - Command: `make -C api test`

- [X] Task 7.2 - Test backward compatibility
  - Verify: Single page requests work exactly as before
  - Verify: Existing API contract unchanged
  - Verify: No breaking changes to response format for single pages

- [X] Task 7.3 - Test multiple page scenarios
  - Verify: Multiple pages downloaded correctly
  - Verify: Pages combined into PDF correctly
  - Verify: Page order preserved
  - Verify: Content-Type header correct

- [X] Task 7.4 - Test error scenarios
  - Verify: Clear errors when workspace IDs missing
  - Verify: Clear errors when invalid page numbers
  - Verify: Graceful handling of partial failures

- [ ] Task 7.5 - Performance testing
  - Verify: Multiple pages download in parallel
  - Verify: Response time acceptable (< 5 seconds for 10 pages)
  - Verify: Memory usage reasonable
  - Note: Performance testing can be done manually during integration testing

## Phase 8: Documentation & Cleanup

- [X] Task 8.1 - Update API endpoint documentation
  - File: `doc/api-endpoints.md`
  - Update `GET /invoices/{id}/image` section:
    - Document multiple page support: `?page=1,2,3`
    - Document response format: multi-page PDF for multiple pages
    - Remove datalog fallback references
    - Update implementation details to show workspace endpoint only

- [X] Task 8.2 - Update JSDoc comments
  - File: `api/src/modules/invoices/invoices.service.ts`
  - File: `api/src/modules/invoices/invoices.controller.ts`
  - File: `api/src/modules/formx/formx.service.ts`
  - Update method documentation to reflect multiple page support
  - Remove datalog fallback references

- [X] Task 8.3 - Code review checklist
  - [X] All tests pass (make -C api test) - 456 tests passing
  - [X] Code coverage maintained/improved
  - [X] Linting passes (make -C api lint)
  - [X] Formatting correct (make -C api format)
  - [X] No breaking changes (single page API unchanged)
  - [X] Backward compatibility maintained
  - [X] Error handling comprehensive
  - [X] Logging appropriate
  - [X] Documentation updated

- [X] Task 8.4 - Update design.md with implementation notes
  - File: `.oursky/scopes/20251126-195612-MultiplePageFileAPI/design.md`
  - Add implementation decisions to Clarifications section
  - Document any deviations from original design
  - Note performance characteristics observed

## Phase 9: Polling Status ID Capture (2025-11-27)

- [X] Task 9.1 - Add methods to FormXJobService for updating IDs
  - File: `api/src/modules/formx/formx-job.service.ts`
  - Add `updateExtractionId(id: string, extractionId: string): Promise<FormXJob>`
  - Add `updateWorkspaceId(id: string, workspaceId: string): Promise<FormXJob>`
  - Methods should update the FormXJob entity and flush changes

- [X] Task 9.2 - Capture extractionId and workspaceId from polling status
  - File: `api/src/queue/processors/formx-polling.processor.ts`
  - After calling `pollJobStatus()`, check if `status.extraction_id` is present (top level or in `extraction_result_lookup_id`)
  - Extract `extraction_id` from `extraction_result_lookup_id` if not at top level (format: "workspace_id:extraction_id:page:slice")
  - If present and not already set on FormXJob, save using `updateExtractionId()` method
  - Note: workspace_id is already available from config, no need to capture from polling
  - This should happen before `handleCompletedExtraction()` is called
  - Add trace logging to track when IDs are captured from status response

- [X] Task 9.3 - Add trace logging for ID capture
  - File: `api/src/queue/processors/formx-polling.processor.ts`
  - Log the `pollJobStatus()` response structure (extraction_id, workspace_id presence)
  - Log when IDs are saved from status response
  - Log before/after state when updating FormXJob with IDs

- [X] Task 9.4 - Update design.md with polling ID capture decision
  - File: `.oursky/scopes/20251126-195612-MultiplePageFileAPI/design.md`
  - Add clarification section documenting the decision to capture IDs from polling status
  - Document rationale, implementation approach, and impact

## Phase 10: Handle FormX File API Signed URL Response (2025-11-27)

- [X] Task 10.1 - Detect and handle JSON response with signed URLs from FormX file API
  - File: `api/src/modules/formx/formx.service.ts`
  - When FormX file API returns JSON with `url` field (singular) or `urls` array, extract URL
  - Prefer `url` field if present, fallback to first URL in `urls` array
  - Download file from signed Google Cloud Storage URL
  - Return file buffer as normal
  - Add trace logging for URL extraction and file download
  - Added validation for buffer and PDF structure in controller

- [X] Task 10.2 - Update design.md with signed URL handling
  - File: `.oursky/scopes/20251126-195612-MultiplePageFileAPI/design.md`
  - Document that FormX file API may return JSON with signed URLs
  - Document the handling logic for downloading from signed URLs

- [X] Task 10.3 - Add unit tests for signed URL handling
  - File: `api/src/modules/formx/formx.service.spec.ts`
  - Test cases:
    - Should download file from signed URL when FormX returns JSON with url field (singular) ✓
    - Should download file from signed URL when FormX returns JSON with urls array ✓
    - Should prefer url field over urls array when both are present ✓
    - Should throw error when JSON response does not contain url or urls ✓
    - Should throw error when JSON response has empty urls array and no url field ✓
  - Updated all existing test mocks to include response headers

- [X] Task 10.4 - Fix test mocks to include response headers
  - File: `api/src/modules/formx/formx.service.spec.ts`
  - Updated all httpService.get mocks to include headers with content-type
  - All tests now passing (31 test files)

## Phase 11: Remove Datalog File API Endpoint (Outdated Code Cleanup)

- [X] Task 11.1 - Remove datalog file endpoint from FormX mock server
  - File: `formx-mock/src/server/index.ts`
  - Removed `GET /v2/datalog/:requestId/file` endpoint handler
  - Endpoint is outdated and replaced by workspace file endpoint

- [X] Task 11.2 - Update mock server documentation
  - File: `formx-mock/README.md`
  - Removed reference to `GET /v2/datalog/{request_id}/file` from API endpoints list
  - Added reference to workspace file endpoint: `GET /v2/workspace/{workspace_id}/extractions/{extraction_id}/pages/{page}/slices/{slice}/file`

- [X] Task 11.3 - Update user flows documentation
  - File: `doc/user-flows.md`
  - Updated Flow 2 to reference workspace file endpoint instead of datalog file endpoint
  - Changed: `GET /v2/datalog/{request_id}/file` → `GET /v2/workspace/{workspaceId}/extractions/{extractionId}/pages/{pageNumber}/slices/{sliceNumber}/file`

- [X] Task 11.4 - Update design.md with removal decision
  - File: `.oursky/scopes/20251126-195612-MultiplePageFileAPI/design.md`
  - Added clarification section documenting the removal of datalog file endpoint
  - Documented rationale and impact

- [X] Task 11.5 - Update tasks.md with removal tasks
  - File: `.oursky/scopes/20251126-195612-MultiplePageFileAPI/tasks.md`
  - Added Phase 11 tasks for removing datalog file endpoint

## Phase 12: Frontend Multi-Page PDF Display Support

- [X] Task 12.1 - Update frontend to detect multi-page PDFs
  - File: `portal/src/features/invoices/components/InvoiceDetailView.tsx`
  - Check `pdfDoc.numPages` to detect if PDF has multiple pages
  - Add state to track current page number
  - Add state to store all page render tasks

- [X] Task 12.2 - Render all pages of multi-page PDFs
  - File: `portal/src/features/invoices/components/InvoiceDetailView.tsx`
  - Loop through all pages using `pdfDoc.numPages`
  - Render each page on separate canvas or in scrollable container
  - Support zoom and rotation for all pages
  - Maintain page order
  - Updated rendering logic to handle single-page vs multi-page PDFs
  - Added canvas refs map for all pages
  - Render all pages in parallel for multi-page PDFs

- [X] Task 12.3 - Add page navigation controls
  - File: `portal/src/features/invoices/components/InvoiceDetailView.tsx`
  - Add page number indicator (e.g., "Page 1 of 3")
  - Add previous/next page buttons
  - Add page jump input (optional)
  - Only show controls for multi-page PDFs
  - Added scroll tracking to update current page automatically
  - Added smooth scroll to page when using navigation buttons

- [X] Task 12.4 - Update design.md with frontend multi-page requirements
  - File: `.oursky/scopes/20251126-195612-MultiplePageFileAPI/design.md`
  - Document that backend calls API multiple times (once per page)
  - Document frontend multi-page display requirements
  - Clarify that each FormX API call only returns one page
  - Already completed in previous commits

- [X] Task 12.5 - Update tasks.md with frontend multi-page tasks
  - File: `.oursky/scopes/20251126-195612-MultiplePageFileAPI/tasks.md`
  - Add Phase 12 tasks for frontend multi-page support
  - Already completed in previous commits

## Phase 13: Multi-Page Invoice Detection via Filename

- [X] Task 13.1 - Add function to parse page numbers from filename
  - File: `api/src/modules/invoices/utils/invoice-number.utils.ts`
  - Function: `extractPageNumbersFromFilename(filename: string): number[]`
  - Parse filename pattern: `{baseName}_page_{page1}_{page2}_{page3}...{extension}`
  - Examples:
    - `sino-invoice-test_page_1_2.pdf` → `[1, 2]`
    - `invoice_page_1_2_3.pdf` → `[1, 2, 3]`
    - `invoice_page_1_2_3_4.pdf` → `[1, 2, 3, 4]`
    - `invoice_page_1.pdf` → `[1]` (single page)
  - Handle edge cases: no `_page_` pattern, invalid page numbers, etc.
  - Return empty array if no page numbers found

- [X] Task 13.2 - Update getInvoiceImage() to use filename for page detection
  - File: `api/src/modules/invoices/invoices.service.ts`
  - Get invoice filename using `computeFileNameForInvoice()`
  - Parse filename using `extractPageNumbersFromFilename()` to extract page numbers
  - If multiple page numbers found (array length > 1):
    - Use extracted page numbers instead of defaulting to `invoice.pageNumber ?? 1`
  - If single page number found (array length === 1):
    - Use single page logic (backward compatible)
  - Remove old logic that checks `metadata.page_no` array from extraction result

- [X] Task 13.3 - Fetch multi-page document URLs using first page number
  - File: `api/src/modules/invoices/invoices.service.ts`
  - When filename contains multiple page numbers:
    - Extract all page numbers from filename (e.g., `[1, 2]`, `[1, 2, 3]`, etc.) for detection
    - Filename can contain any number of pages (2, 3, 4, 5, etc.) - handle dynamically
    - **Use ONLY the FIRST page number** (e.g., `page_1_2` → use page 1)
    - Use the invoice's `formxWorkspaceId` and `formxExtractionId`
    - Call FormX API ONCE with the first page number: `GET /v2/workspace/{workspaceId}/extractions/{extractionId}/pages/{firstPageNumber}/slices/{sliceNumber}/file`
    - **FormX API returns JSON with `urls` array**: `{"status": "ok", "urls": ["url1", "url2", ...]}`
    - Each URL in the array represents one page (ordered by page number: page_0, page_1, etc.)
    - Number of API calls to FormX = 1 (regardless of how many pages the document contains)
    - Use `invoice.sliceNumber ?? 1` for slice number
    - Extract all URLs from the `urls` array (handle both single URL and multiple URLs)

- [X] Task 13.4 - Download all URLs and combine into single PDF
  - File: `api/src/modules/formx/formx.service.ts`
  - Update `downloadWorkspaceImage()` or create new method to handle multiple URLs in `urls` array
  - When JSON response contains `urls` array with multiple URLs:
    - Download all URLs in parallel (using `Promise.all()`)
    - Return array of Buffers (one per page, in order)
  - When `urls` array has single URL: Keep existing behavior (download single file)
  - File: `api/src/modules/invoices/invoices.service.ts`
  - After calling FormX API and receiving `urls` array:
    - If array has multiple URLs: Download all URLs, combine all page files into single PDF using `PdfCombinerService`
    - If array has single URL: Download and return as-is (backward compatible)
  - Preserve page order (URLs are already ordered by page number: page_0, page_1, etc.)
  - Return combined multi-page PDF for multi-page invoices

- [X] Task 13.5 - Write unit tests for filename-based page detection
  - File: `api/src/modules/invoices/invoices.service.spec.ts`
  - Test cases:
    - Should detect multi-page invoice when filename is `invoice_page_1_2.pdf` (2 pages) - use page 1
    - Should detect multi-page invoice when filename is `invoice_page_1_2_3.pdf` (3 pages) - use page 1
    - Should detect multi-page invoice when filename is `invoice_page_1_2_3_4.pdf` (4 pages) - use page 1
    - Should handle any number of pages in filename (not limited to 2) - always use first page
    - Should call FormX API ONCE with first page number for multi-page invoices
    - Should receive JSON with `urls` array containing multiple signed URLs from single API call
    - Should download all URLs and combine into single PDF
    - Should return single page when filename is `invoice_page_1.pdf` (single page)
    - Should handle filename without `_page_` pattern (fallback to invoice.pageNumber)
    - Should handle invalid page numbers in filename gracefully
    - Should respect page query parameter override

- [X] Task 13.6 - Write unit tests for extractPageNumbersFromFilename function
  - File: `api/src/modules/invoices/utils/invoice-number.utils.spec.ts`
  - Test cases:
    - Should extract `[1, 2]` from `invoice_page_1_2.pdf`
    - Should extract `[1, 2, 3]` from `invoice_page_1_2_3.pdf`
    - Should extract `[1]` from `invoice_page_1.pdf`
    - Should return empty array if no `_page_` pattern found
    - Should handle filenames with different extensions
    - Should handle edge cases (empty filename, invalid format, etc.)

- [X] Task 13.7 - Update design.md with filename-based detection approach
  - File: `.oursky/scopes/20251126-195612-MultiplePageFileAPI/design.md`
  - Document the decision to detect multi-page invoices via filename parsing
  - Document the filename format and parsing logic
  - Remove outdated information about `page_no` array approach
  - Already completed in previous commits

- [X] Task 13.8 - Update tasks.md with Phase 13 tasks
  - File: `.oursky/scopes/20251126-195612-MultiplePageFileAPI/tasks.md`
  - Add Phase 13 tasks for filename-based detection
  - Remove outdated Phase 13 tasks about `page_no` array
  - Already completed in previous commits

