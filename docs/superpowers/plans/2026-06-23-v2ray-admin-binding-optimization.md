# V2ray Admin Binding Optimization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor the `v2ray_admin/binding` page by extracting testable share/copy/download logic and adding clearer loading behavior without changing backend API contracts.

**Architecture:** Keep the existing Vue 3 and Element Plus page. Move pure share-state behavior into a local `.mjs` helper with Node tests, while leaving browser APIs such as clipboard and anchor download in the Vue component.

**Tech Stack:** Vue 3, Element Plus, qrcode, Node `assert`, existing Vite aliases.

---

### Task 1: Binding Share Helpers

**Files:**
- Create: `gin-vue-admin/web/src/view/v2ray_admin/binding/bindingShare.mjs`
- Create: `gin-vue-admin/web/src/view/v2ray_admin/binding/bindingShare.test.mjs`

- [ ] **Step 1: Write the failing test**

Create `bindingShare.test.mjs` with assertions for:

- resolving `config1` and `config2` share links
- building copy success messages
- building `.png` download names
- converting backend share data into dialog state with injected QR-code generator

- [ ] **Step 2: Run test to verify it fails**

Run: `node src/view/v2ray_admin/binding/bindingShare.test.mjs`

Expected: FAIL with module not found for `bindingShare.mjs`.

- [ ] **Step 3: Write minimal implementation**

Create `bindingShare.mjs` with:

- `getShareLink(shareInfo, configType)`
- `getCopySuccessMessage(configType)`
- `buildQrDownloadName(filename)`
- `createShareDialogInfo(data, toDataUrl)`

- [ ] **Step 4: Run test to verify it passes**

Run: `node src/view/v2ray_admin/binding/bindingShare.test.mjs`

Expected: PASS and `bindingShare tests passed`.

### Task 2: Wire Admin Binding Page To Helpers

**Files:**
- Modify: `gin-vue-admin/web/src/view/v2ray_admin/binding/binding.vue`

- [ ] **Step 1: Remove unused imports and unused ClipboardJS instance**

Remove `ClipboardJS`, unused dictionary imports, and unused `onMounted`.

- [ ] **Step 2: Import binding share helpers**

Import helper functions from `./bindingShare.mjs`.

- [ ] **Step 3: Add loading refs**

Add `tableLoading` for list loading and `shareLoading` for share QR generation.

- [ ] **Step 4: Refactor `getTableData`**

Wrap the list API call in `try/finally`, bind `v-loading="tableLoading"` to the table, and show Element Plus error feedback on request failure or non-zero API code.

- [ ] **Step 5: Refactor `shareBindingFunc`**

Use `createShareDialogInfo(res.data, QRCode.toDataURL)` and set `shareLoading` around the async work.

- [ ] **Step 6: Refactor `handleCopy` and `downloadQR`**

Use `getShareLink`, `getCopySuccessMessage`, and `buildQrDownloadName` while preserving existing browser clipboard fallback and anchor download behavior.

### Task 3: Verification

**Files:**
- Review: `gin-vue-admin/web/src/view/v2ray_admin/binding/binding.vue`
- Review: `gin-vue-admin/web/src/view/v2ray_admin/binding/bindingShare.mjs`
- Review: `gin-vue-admin/web/src/view/v2ray_admin/binding/bindingShare.test.mjs`

- [ ] **Step 1: Run targeted tests**

Run:

```powershell
node src/view/v2ray_admin/server/serverPortReminder.test.mjs
node src/view/v2ray_admin/stat/statTraffic.test.mjs
node src/view/v2ray_admin/stat/statChartOptions.test.mjs
node src/view/v2ray_admin/binding/bindingShare.test.mjs
```

Expected: all PASS.

- [ ] **Step 2: Run SFC parse checks**

Run Node with `@vue/compiler-sfc` against:

- `src/view/v2ray_admin/stat/stat.vue`
- `src/view/v2ray_admin/binding/binding.vue`

Expected: both parse and compile without SFC errors.

- [ ] **Step 3: Inspect git diff**

Run: `git diff --check`

Expected: no whitespace errors.

- [ ] **Step 4: Do not run full build by default**

The user reported production build takes about eight minutes. Record build verification as not run unless the user provides output or explicitly asks Codex to run it.
