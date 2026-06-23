# V2ray Stat Helper Reuse Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor the user-facing `v2ray/stat` page to reuse the tested admin stat helpers while preserving its numeric date-axis formatting.

**Architecture:** Extend the shared `v2ray_admin/stat` chart option helper with tested axis normalization, then wire `v2ray/stat.vue` to `statTraffic.mjs` and `statChartOptions.mjs`. Keep API calls and page layout unchanged.

**Tech Stack:** Vue 3, Element Plus, ECharts, Node `assert`, `@vue/compiler-sfc`.

---

### Task 1: Preserve Numeric Date Axis Behavior

**Files:**
- Modify: `gin-vue-admin/web/src/view/v2ray_admin/stat/statChartOptions.test.mjs`
- Modify: `gin-vue-admin/web/src/view/v2ray_admin/stat/statChartOptions.mjs`

- [x] **Step 1: Add failing test**

Add assertions that `buildTrendChartOptions({ data_axis: [20260621], data: [1024] })` converts the x-axis data to `['2026-06-21']` and that `xAxis.axisLabel.formatter('2026-06-21')` returns `06-21`.

- [x] **Step 2: Run test to verify it fails**

Run: `node src/view/v2ray_admin/stat/statChartOptions.test.mjs`

Expected: FAIL because the current helper does not normalize numeric date axes.

- [x] **Step 3: Implement axis normalization**

Add local helper functions in `statChartOptions.mjs`:

- `normalizeAxisValue(value)`
- `formatAxisLabel(value)`

Use them in `buildTrendChartOptions`.

- [x] **Step 4: Run test to verify it passes**

Run: `node src/view/v2ray_admin/stat/statChartOptions.test.mjs`

Expected: PASS and `statChartOptions tests passed`.

### Task 2: Refactor V2ray Stat Page

**Files:**
- Modify: `gin-vue-admin/web/src/view/v2ray/stat/stat.vue`

- [x] **Step 1: Remove unused imports and refs**

Remove unused `ElMessageBox`, `reactive`, `computed`, `formData`, `elFormRef`, and the inline helper functions replaced by imports.

- [x] **Step 2: Import shared helpers**

Import `formatFlow`, `getDateRangeText`, `getTrafficTagType`, `normalizeDateOnlyToUtcIso`, and `buildTrendChartOptions`.

- [x] **Step 3: Replace inline date and chart code**

Use `dateRangeText` computed from `getDateRangeText(searchInfo.value)`, normalize submit dates via `normalizeDateOnlyToUtcIso`, and replace `setOptions` with `chart.value?.setOption(buildTrendChartOptions(data))`.

- [x] **Step 4: Fix chart lifecycle**

Move resize handler and disposal to top-level hooks so the same handler reference is removed. Remove timer-based refresh paths that are no longer needed because the chart watcher already updates options.

### Task 3: Verification

**Files:**
- Review: `gin-vue-admin/web/src/view/v2ray/stat/stat.vue`
- Review: stat helper files and tests

- [x] **Step 1: Run targeted tests**

Run:

```powershell
node src/view/v2ray_admin/server/serverPortReminder.test.mjs
node src/view/v2ray_admin/stat/statTraffic.test.mjs
node src/view/v2ray_admin/stat/statChartOptions.test.mjs
node src/view/v2ray_admin/binding/bindingShare.test.mjs
```

Expected: all PASS.

- [x] **Step 2: Run SFC checks**

Run SFC parse/compile checks for:

- `src/view/v2ray_admin/stat/stat.vue`
- `src/view/v2ray/stat/stat.vue`
- `src/view/v2ray_admin/binding/binding.vue`
- `src/view/v2ray/binding/binding.vue`

Expected: all parse and compile.

- [x] **Step 3: Inspect git diff**

Run: `git diff --check`

Expected: no whitespace errors.

- [x] **Step 4: Do not run full build by default**

The production build is intentionally skipped unless explicitly requested because the user reported it takes about eight minutes.
