# Layout Event Lifecycle Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `layout/index.vue` responsive behavior and event bus listeners lifecycle-safe without changing the visible layout breakpoints.

**Architecture:** Add `src/view/layout/layoutEvents.mjs` with pure breakpoint state calculation and a small event binding helper. Refactor `layout/index.vue` to use `window.addEventListener` plus disposer cleanup instead of assigning `window.onresize`, and to remove the exact mitt handlers registered on mount.

**Tech Stack:** Vue 3 SFC, mitt emitter, browser resize events, Node ESM tests with `node:assert`, `@vue/compiler-sfc`.

---

### Task 1: Add Layout Helper Tests

**Files:**
- Create: `gin-vue-admin/web/src/view/layout/layoutEvents.test.mjs`

- [x] **Step 1: Add failing tests**

Test `getLayoutState` for widths `999`, `1000`, `1199`, and `1200`. Test `bindLayoutEventHandlers` registers `reload`, `showLoading`, `closeLoading`, and `resize`, then removes those same handler references from the emitter and window-like target.

- [x] **Step 2: Run test to verify it fails**

Run: `node src/view/layout/layoutEvents.test.mjs`

Expected: FAIL with module-not-found because `layoutEvents.mjs` does not exist yet.

### Task 2: Implement Layout Events Helper

**Files:**
- Create: `gin-vue-admin/web/src/view/layout/layoutEvents.mjs`

- [x] **Step 1: Add helper implementation**

Implement `getLayoutState(screenWidth)` and `bindLayoutEventHandlers({ emitter, target, onReload, onShowLoading, onCloseLoading, onResize })`.

- [x] **Step 2: Run helper tests**

Run: `node src/view/layout/layoutEvents.test.mjs`

Expected: PASS and `layoutEvents tests passed`.

### Task 3: Refactor Layout Component

**Files:**
- Modify: `gin-vue-admin/web/src/view/layout/index.vue`

- [x] **Step 1: Import helper and `onUnmounted`**

Import `getLayoutState` and `bindLayoutEventHandlers`; add `onUnmounted` to Vue imports.

- [x] **Step 2: Replace breakpoint logic**

Make `initPage` derive `isMobile`, `isSider`, and `isCollapse` from `getLayoutState(document.body.clientWidth)`.

- [x] **Step 3: Replace event registration**

Use named handlers for reload, loading, close loading, and resize. Store the disposer returned by `bindLayoutEventHandlers`.

- [x] **Step 4: Add cleanup**

On unmount, dispose layout handlers and clear any pending reload timer.

### Task 4: Verification And Commit

**Files:**
- Test: `gin-vue-admin/web/src/view/layout/layoutEvents.test.mjs`
- Review: `gin-vue-admin/web/src/view/layout/index.vue`

- [x] **Step 1: Run new helper test**

Run: `node src/view/layout/layoutEvents.test.mjs`

Expected: PASS.

- [x] **Step 2: Run existing targeted frontend tests**

Run the Node tests added in previous optimization slices.

- [x] **Step 3: Run SFC checks**

Run SFC parse/compile checks for `layout/index.vue`, `history.vue`, `system/state.vue`, and the traffic/binding pages touched in this branch.

- [x] **Step 4: Skip production build**

Do not run the full production build by default because the user reported it takes about eight minutes.

- [x] **Step 5: Commit**

Commit message: `fix: clean up layout event lifecycle`
