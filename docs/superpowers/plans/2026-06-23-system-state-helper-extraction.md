# System State Helper Extraction Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reduce `system/state.vue` complexity by moving pure formatting, validation, percentage, and traffic target helpers into a tested module.

**Architecture:** Add `src/view/system/stateHelpers.mjs` next to the page. Keep API calls, Element Plus controls, dialog behavior, and styles in `state.vue`; move deterministic data transformation into helper functions covered by Node `assert` tests.

**Tech Stack:** Vue 3 SFC, Element Plus, Node ESM tests with `node:assert`, `@vue/compiler-sfc`.

---

### Task 1: Lock Helper Behavior With Tests

**Files:**
- Create: `gin-vue-admin/web/src/view/system/stateHelpers.test.mjs`

- [x] **Step 1: Add failing helper tests**

Create tests that import these functions from `./stateHelpers.mjs`: `calculatePercent`, `formatBytes`, `formatSize`, `formatUserOption`, `getTrafficAnalysisTotals`, `isServerOnline`, `parseClockMinute`, `rowTrafficAnalysisTargets`, `targetNames`, `todayCompact`, and `validateTrafficAnalysisQuery`.

- [x] **Step 2: Run test to verify it fails**

Run: `node src/view/system/stateHelpers.test.mjs`

Expected: FAIL with module-not-found because `stateHelpers.mjs` does not exist yet.

### Task 2: Implement State Helpers

**Files:**
- Create: `gin-vue-admin/web/src/view/system/stateHelpers.mjs`

- [x] **Step 1: Add minimal helper implementation**

Implement the exported helpers with the same behavior currently embedded in `state.vue`, plus finite-number guards for byte and percentage formatting.

- [x] **Step 2: Run helper tests**

Run: `node src/view/system/stateHelpers.test.mjs`

Expected: PASS and `stateHelpers tests passed`.

### Task 3: Refactor State Page To Use Helpers

**Files:**
- Modify: `gin-vue-admin/web/src/view/system/state.vue`

- [x] **Step 1: Import shared helpers**

Import the helper functions from `./stateHelpers.mjs`.

- [x] **Step 2: Replace inline pure functions**

Remove local definitions for byte/size/time formatting, compact date generation, traffic form validation, totals calculation, user option text, target extraction, online state, and percentage calculation.

- [x] **Step 3: Keep page behavior unchanged**

Leave API calls, loading flags, Element Plus messages, dialog lifecycle, table layout, and styles unchanged.

### Task 4: Verification And Commit

**Files:**
- Review: `gin-vue-admin/web/src/view/system/state.vue`
- Test: `gin-vue-admin/web/src/view/system/stateHelpers.test.mjs`

- [x] **Step 1: Run helper tests**

Run: `node src/view/system/stateHelpers.test.mjs`

Expected: PASS.

- [x] **Step 2: Run existing targeted frontend tests**

Run:

```powershell
node src/view/v2ray_admin/server/serverPortReminder.test.mjs
node src/view/v2ray_admin/stat/statTraffic.test.mjs
node src/view/v2ray_admin/stat/statChartOptions.test.mjs
node src/view/v2ray_admin/binding/bindingShare.test.mjs
```

Expected: all PASS.

- [x] **Step 3: Run SFC checks**

Run SFC parse/compile checks for `src/view/system/state.vue` and the traffic/binding pages touched in previous slices.

Expected: all parse and compile.

- [x] **Step 4: Skip production build**

Do not run the full production build by default because the user reported it takes about eight minutes.

- [x] **Step 5: Commit**

Commit message: `refactor: extract system state helpers`
