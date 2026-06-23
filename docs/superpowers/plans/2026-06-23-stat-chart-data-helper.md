# Stat Chart Data Helper Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reduce duplicated stat chart response normalization between `v2ray_admin/stat/common.js` and `v2ray/stat/common.js`.

**Architecture:** Add a pure helper module beside the existing admin stat helpers. Keep API calls in each `common.js`, but share trend total calculation, empty-state reset, and rank top-10 slicing through tested functions.

**Tech Stack:** Vue 3 reactive stores, Node ESM tests with `node:assert`.

---

### Task 1: Add Stat Chart Data Tests

**Files:**
- Create: `gin-vue-admin/web/src/view/v2ray_admin/stat/statChartData.test.mjs`

- [x] **Step 1: Add failing tests**

Test `createChartDataState`, `applyTrendChartResponse`, and `applyRankChartResponse` for valid responses, invalid responses, total calculation, and top-10 rank slicing.

- [x] **Step 2: Run test to verify it fails**

Run: `node src/view/v2ray_admin/stat/statChartData.test.mjs`

Expected: FAIL with module-not-found because `statChartData.mjs` does not exist yet.

### Task 2: Implement Stat Chart Data Helper

**Files:**
- Create: `gin-vue-admin/web/src/view/v2ray_admin/stat/statChartData.mjs`

- [x] **Step 1: Add helper implementation**

Implement `createChartDataState({ includeRank = false } = {})`, `applyTrendChartResponse(target, response)`, and `applyRankChartResponse(target, response, { maxRankItems = 10 } = {})`.

- [x] **Step 2: Run helper tests**

Run: `node src/view/v2ray_admin/stat/statChartData.test.mjs`

Expected: PASS and `statChartData tests passed`.

### Task 3: Refactor Stat Common Stores

**Files:**
- Modify: `gin-vue-admin/web/src/view/v2ray_admin/stat/common.js`
- Modify: `gin-vue-admin/web/src/view/v2ray/stat/common.js`

- [x] **Step 1: Import shared helper**

Import `createChartDataState`, `applyTrendChartResponse`, and, for admin only, `applyRankChartResponse`.

- [x] **Step 2: Replace duplicated response mutation**

Use the helper functions to mutate each reactive chart data object after API calls.

- [x] **Step 3: Preserve API behavior**

Keep each file calling its existing API module. Do not change backend request parameters.

### Task 4: Verification And Commit

**Files:**
- Test: `gin-vue-admin/web/src/view/v2ray_admin/stat/statChartData.test.mjs`
- Review: `gin-vue-admin/web/src/view/v2ray_admin/stat/common.js`
- Review: `gin-vue-admin/web/src/view/v2ray/stat/common.js`

- [x] **Step 1: Run new helper test**

Run: `node src/view/v2ray_admin/stat/statChartData.test.mjs`

Expected: PASS.

- [x] **Step 2: Run existing targeted frontend tests**

Run all Node tests added in previous optimization slices.

- [x] **Step 3: Run SFC checks**

Run SFC parse/compile checks for the traffic stat pages and layout pages touched in this branch.

- [x] **Step 4: Skip production build**

Do not run the full production build by default because the user reported it takes about eight minutes.

- [x] **Step 5: Commit**

Commit message: `refactor: share stat chart data helpers`
