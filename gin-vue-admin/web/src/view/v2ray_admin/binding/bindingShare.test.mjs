import assert from 'node:assert/strict'
import {
  buildQrDownloadName,
  createShareDialogInfo,
  getCopySuccessMessage,
  getShareLink,
} from './bindingShare.mjs'

const shareInfo = {
  share1_link: 'vmess://shadowrocket',
  share2_link: 'vmess://v2rayn',
}

assert.equal(getShareLink(shareInfo, 'config1'), 'vmess://shadowrocket')
assert.equal(getShareLink(shareInfo, 'config2'), 'vmess://v2rayn')
assert.equal(getShareLink(shareInfo, 'unknown'), '')
assert.equal(getShareLink(null, 'config1'), '')

assert.equal(getCopySuccessMessage('config1'), 'Shadowrocket 配置已复制到剪贴板')
assert.equal(getCopySuccessMessage('config2'), 'V2rayN 配置已复制到剪贴板')
assert.equal(getCopySuccessMessage('unknown'), '配置已复制到剪贴板')

assert.equal(buildQrDownloadName('shadowrocket-config'), 'shadowrocket-config.png')
assert.equal(buildQrDownloadName('v2ray-config.png'), 'v2ray-config.png')
assert.equal(buildQrDownloadName(''), 'qrcode.png')

const generated = await createShareDialogInfo(
  {
    share1: 'vmess://shadowrocket',
    share2: 'vmess://v2rayn',
  },
  async (value) => `data:image/png;base64,${Buffer.from(value).toString('base64')}`
)

assert.deepEqual(generated, {
  share1: 'data:image/png;base64,dm1lc3M6Ly9zaGFkb3dyb2NrZXQ=',
  share1_link: 'vmess://shadowrocket',
  share2: 'data:image/png;base64,dm1lc3M6Ly92MnJheW4=',
  share2_link: 'vmess://v2rayn',
})

console.log('bindingShare tests passed')
