# History Listener Cleanup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the layout history tab component clean up body click and event bus listeners with the same handler references used during registration.

**Architecture:** Add a tiny `historyEvents.mjs` helper next to `history.vue` for DOM and mitt-style event binding disposal. Refactor `history.vue` to keep named callbacks and dispose functions while preserving existing tab behavior.

**Tech Stack:** Vue 3 SFC, mitt emitter, browser DOM events, Node ESM tests with `node:assert`.

---

### Task 1: Add Event Helper Tests

**Files:**
- Create: `gin-vue-admin/web/src/view/layout/aside/historyComponent/historyEvents.test.mjs`

- [x] **Step 1: Add failing tests**

Create tests for `bindBodyClickHandler` and `bindEmitterHandlers` that verify returned dispose functions remove the exact handlers registered.

- [x] **Step 2: Run test to verify it fails**

Run: `node src/view/layout/aside/historyComponent/historyEvents.test.mjs`

Expected: FAIL with module-not-found because `historyEvents.mjs` does not exist yet.

### Task 2: Implement Event Helper

**Files:**
- Create: `gin-vue-admin/web/src/view/layout/aside/historyComponent/historyEvents.mjs`

- [x] **Step 1: Add helper implementation**

Implement `bindBodyClickHandler(handler, target = document.body)` and `bindEmitterHandlers(emitter, handlers)`; each should return a dispose function.

- [x] **Step 2: Run helper tests**

Run: `node src/view/layout/aside/historyComponent/historyEvents.test.mjs`

Expected: PASS and `historyEvents tests passed`.

### Task 3: Refactor History Component

**Files:**
- Modify: `gin-vue-admin/web/src/view/layout/aside/historyComponent/history.vue`

- [x] **Step 1: Import helper**

Import `bindBodyClickHandler` and `bindEmitterHandlers`.

- [x] **Step 2: Replace anonymous body click handlers**

Use a named `closeContextMenu` function and a stored dispose function so the click listener is removed correctly.

- [x] **Step 3: Replace anonymous emitter handlers**

Use named handlers for `closeThisPage`, `closeAllPage`, `mobile`, and `collapse`, then dispose them on component unmount.

### Task 4: Verification And Commit

**Files:**
- Test: `gin-vue-admin/web/src/view/layout/aside/historyComponent/historyEvents.test.mjs`
- Review: `gin-vue-admin/web/src/view/layout/aside/historyComponent/history.vue`

- [x] **Step 1: Run new helper test**

Run: `node src/view/layout/aside/historyComponent/historyEvents.test.mjs`

Expected: PASS.

- [x] **Step 2: Run existing targeted frontend tests**

Run the existing Node tests added in previous optimization slices.

- [x] **Step 3: Run SFC checks**

Run SFC parse/compile checks for `history.vue` and the pages touched in this optimization branch.

- [x] **Step 4: Skip production build**

Do not run the full production build by default because the user reported it takes about eight minutes.

- [x] **Step 5: Commit**

Commit message: `fix: clean up history component listeners`
