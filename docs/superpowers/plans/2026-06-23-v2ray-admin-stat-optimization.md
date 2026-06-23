# V2ray Admin Stat Optimization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor the `v2ray_admin/stat` page by extracting testable traffic formatting and ECharts option logic, then wire the page to those helpers with clearer loading and lifecycle handling.

**Architecture:** Keep the existing Vue 3, Element Plus, and ECharts page. Move pure behavior into local `.mjs` modules with Node-based tests, and keep DOM/ECharts lifecycle inside the Vue component.

**Tech Stack:** Vue 3, Element Plus, ECharts, Node `assert`, existing Vite aliases.

---

### Task 1: Traffic Formatting Helpers

**Files:**
- Create: `gin-vue-admin/web/src/view/v2ray_admin/stat/statTraffic.mjs`
- Create: `gin-vue-admin/web/src/view/v2ray_admin/stat/statTraffic.test.mjs`

- [ ] **Step 1: Write the failing test**

Create `statTraffic.test.mjs` with assertions for byte formatting, text parsing, tag thresholds, UTC date normalization, and date range labels.

- [ ] **Step 2: Run test to verify it fails**

Run: `node src/view/v2ray_admin/stat/statTraffic.test.mjs`

Expected: FAIL with module not found for `statTraffic.mjs`.

- [ ] **Step 3: Write minimal implementation**

Create `statTraffic.mjs` with:

- `formatFlow(value)`
- `parseTrafficBytes(value)`
- `getTrafficTagType(value)`
- `normalizeDateOnlyToUtcIso(value)`
- `getDateRangeText(searchInfo)`

- [ ] **Step 4: Run test to verify it passes**

Run: `node src/view/v2ray_admin/stat/statTraffic.test.mjs`

Expected: PASS and `statTraffic tests passed`.

### Task 2: Chart Option Builders

**Files:**
- Create: `gin-vue-admin/web/src/view/v2ray_admin/stat/statChartOptions.mjs`
- Create: `gin-vue-admin/web/src/view/v2ray_admin/stat/statChartOptions.test.mjs`

- [ ] **Step 1: Write the failing test**

Create `statChartOptions.test.mjs` with assertions for trend line options, rank bar options, empty data fallback, and tooltip formatter output.

- [ ] **Step 2: Run test to verify it fails**

Run: `node src/view/v2ray_admin/stat/statChartOptions.test.mjs`

Expected: FAIL with module not found for `statChartOptions.mjs`.

- [ ] **Step 3: Write minimal implementation**

Create `statChartOptions.mjs` with:

- `buildTrendChartOptions(data)`
- `buildRankChartOptions(data)`

Both functions should accept the existing `chartData` shape: `{ data, data_axis, rank, rank_axis }`.

- [ ] **Step 4: Run test to verify it passes**

Run: `node src/view/v2ray_admin/stat/statChartOptions.test.mjs`

Expected: PASS and `statChartOptions tests passed`.

### Task 3: Wire Stat Page To Helpers

**Files:**
- Modify: `gin-vue-admin/web/src/view/v2ray_admin/stat/stat.vue`

- [ ] **Step 1: Replace inline pure helpers**

Import helper functions from `statTraffic.mjs` and `statChartOptions.mjs`.

- [ ] **Step 2: Simplify chart option application**

Replace inline `setOptions` contents with:

```js
const setOptions = (data) => {
  chart.value?.setOption(buildTrendChartOptions(data))
  rankChart.value?.setOption(buildRankChartOptions(data))
}
```

- [ ] **Step 3: Fix chart lifecycle**

Move resize handler registration and cleanup to top-level `onMounted` and `onUnmounted` hooks so the same function reference is removed.

- [ ] **Step 4: Add table loading state**

Wrap `getTableData` in a `try/finally`, set `tableLoading`, and bind it to the table with `v-loading`.

- [ ] **Step 5: Run helper tests**

Run:

```powershell
node src/view/v2ray_admin/stat/statTraffic.test.mjs
node src/view/v2ray_admin/stat/statChartOptions.test.mjs
```

Expected: both PASS.

### Task 4: Verification And Review

**Files:**
- Review: `gin-vue-admin/web/src/view/v2ray_admin/stat/stat.vue`
- Review: helper modules and tests from Tasks 1-2

- [ ] **Step 1: Run available targeted tests**

Run:

```powershell
node src/view/v2ray_admin/server/serverPortReminder.test.mjs
node src/view/v2ray_admin/stat/statTraffic.test.mjs
node src/view/v2ray_admin/stat/statChartOptions.test.mjs
```

Expected: all PASS.

- [ ] **Step 2: Inspect git diff**

Run: `git diff --check`

Expected: no whitespace errors.

- [ ] **Step 3: Do not run full build by default**

The user reported production build takes about eight minutes. Record build verification as not run unless the user provides output or explicitly asks Codex to run it.
