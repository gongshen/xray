# Responsive Directive Listener Cleanup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix responsive directive resize listener cleanup so unmounted directives remove the exact handler they registered.

**Architecture:** Add a tiny `responsiveDirectiveHandlers.mjs` helper that stores resize handlers per element and directive key. The existing `responsive.js` directives keep their public names and behavior, but use the helper for lifecycle-safe add/remove.

**Tech Stack:** Vue custom directives, browser `window.addEventListener`, Node ESM tests with `node:assert`.

---

### Task 1: Add Listener Helper Tests

**Files:**
- Create: `gin-vue-admin/web/src/directive/responsiveDirectiveHandlers.test.mjs`

- [x] **Step 1: Add failing tests**

Create tests for `bindElementResizeHandler` and `unbindElementResizeHandler` that verify:

- binding registers the provided handler
- unbinding removes the same handler reference
- rebinding the same element/key removes the previous handler before adding the new one
- unbinding an unknown element/key is a no-op

- [x] **Step 2: Run test to verify it fails**

Run: `node src/directive/responsiveDirectiveHandlers.test.mjs`

Expected: FAIL with module-not-found because the helper does not exist yet.

### Task 2: Implement Listener Helper

**Files:**
- Create: `gin-vue-admin/web/src/directive/responsiveDirectiveHandlers.mjs`

- [x] **Step 1: Add helper implementation**

Use a `WeakMap` keyed by element, with an inner `Map` keyed by directive name, to store the handler function.

- [x] **Step 2: Run helper tests**

Run: `node src/directive/responsiveDirectiveHandlers.test.mjs`

Expected: PASS and `responsiveDirectiveHandlers tests passed`.

### Task 3: Refactor Responsive Directives

**Files:**
- Modify: `gin-vue-admin/web/src/directive/responsive.js`

- [x] **Step 1: Import helper**

Import `bindElementResizeHandler` and `unbindElementResizeHandler`.

- [x] **Step 2: Replace anonymous listener registration**

Use named handler closures stored by the helper for both `responsiveTable` and `responsiveForm`.

- [x] **Step 3: Preserve directive behavior**

Keep the existing `makeTableResponsive` and `adaptFormForMobile` logic unchanged.

### Task 4: Verification And Commit

**Files:**
- Test: `gin-vue-admin/web/src/directive/responsiveDirectiveHandlers.test.mjs`
- Review: `gin-vue-admin/web/src/directive/responsive.js`

- [x] **Step 1: Run new helper test**

Run: `node src/directive/responsiveDirectiveHandlers.test.mjs`

Expected: PASS.

- [x] **Step 2: Run existing targeted frontend tests**

Run the existing Node tests for server port reminder, stat traffic, stat chart options, binding share, and system state helpers.

- [x] **Step 3: Run SFC checks**

Run SFC parse/compile checks for pages touched in this optimization branch.

- [x] **Step 4: Skip production build**

Do not run the full production build by default because the user reported it takes about eight minutes.

- [x] **Step 5: Commit**

Commit message: `fix: clean up responsive directive listeners`
