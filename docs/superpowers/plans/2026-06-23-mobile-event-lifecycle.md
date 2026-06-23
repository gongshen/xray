# Mobile Event Lifecycle Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make remaining mobile menu and responsive event listeners remove the exact handlers they register.

**Architecture:** Add a small shared event lifecycle helper under `src/utils`. Refactor mobile menu, overlay, the responsive composable, and the aside menu to use disposer functions instead of anonymous callbacks or broad `emitter.off(type)` calls.

**Tech Stack:** Vue 3 SFC, mitt emitter, browser event listeners, Node ESM tests with `node:assert`, `@vue/compiler-sfc`.

---

### Task 1: Add Event Lifecycle Helper Tests

**Files:**
- Create: `gin-vue-admin/web/src/utils/eventLifecycle.test.mjs`

- [x] **Step 1: Add failing tests**

Test `bindEmitterHandler`, `bindEmitterHandlers`, and `bindWindowEvent` with mock emitter and window-like targets. Verify each disposer removes the exact handler registered.

- [x] **Step 2: Run test to verify it fails**

Run: `node src/utils/eventLifecycle.test.mjs`

Expected: FAIL with module-not-found because `eventLifecycle.mjs` does not exist yet.

### Task 2: Implement Event Lifecycle Helper

**Files:**
- Create: `gin-vue-admin/web/src/utils/eventLifecycle.mjs`

- [x] **Step 1: Add helper implementation**

Implement:

- `bindEmitterHandler(emitter, event, handler)`
- `bindEmitterHandlers(emitter, handlers)`
- `bindWindowEvent(target, event, handler)`

Each returns a dispose function.

- [x] **Step 2: Run helper tests**

Run: `node src/utils/eventLifecycle.test.mjs`

Expected: PASS and `eventLifecycle tests passed`.

### Task 3: Refactor Mobile Event Users

**Files:**
- Modify: `gin-vue-admin/web/src/components/MobileMenuToggle/index.vue`
- Modify: `gin-vue-admin/web/src/components/ResponsiveOverlay/index.vue`
- Modify: `gin-vue-admin/web/src/utils/responsive.js`
- Modify: `gin-vue-admin/web/src/view/layout/aside/index.vue`

- [x] **Step 1: Import event lifecycle helpers**

Use the new helper functions from `@/utils/eventLifecycle.mjs`.

- [x] **Step 2: Replace anonymous or broad unregister logic**

Store disposer functions and call them on `onUnmounted`, preserving each component's current behavior.

- [x] **Step 3: Keep visible UI behavior unchanged**

Do not change menu breakpoints, overlay styling, or route behavior.

### Task 4: Verification And Commit

**Files:**
- Test: `gin-vue-admin/web/src/utils/eventLifecycle.test.mjs`
- Review: touched Vue files and responsive util

- [x] **Step 1: Run new helper test**

Run: `node src/utils/eventLifecycle.test.mjs`

Expected: PASS.

- [x] **Step 2: Run existing targeted frontend tests**

Run all Node tests added in previous optimization slices.

- [x] **Step 3: Run SFC checks**

Run SFC parse/compile checks for touched Vue files plus previously touched layout/stat/system pages.

- [x] **Step 4: Skip production build**

Do not run the full production build by default because the user reported it takes about eight minutes.

- [x] **Step 5: Commit**

Commit message: `fix: clean up mobile event lifecycle`
