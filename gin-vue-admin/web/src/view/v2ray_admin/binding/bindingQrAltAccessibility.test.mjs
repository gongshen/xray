import assert from 'node:assert/strict'
import fs from 'node:fs'

const files = [
  'src/view/v2ray_admin/binding/binding.vue',
  'src/view/v2ray/binding/binding.vue',
]

for (const file of files) {
  const source = fs.readFileSync(file, 'utf8')
  const qrImages = [...source.matchAll(/<img[^>]*class="qr-image"[^>]*>/g)].map((match) => match[0])

  assert.equal(qrImages.length, 2, `${file} should render two QR code images`)
  assert.match(
    qrImages[0],
    /alt="Shadowrocket \/ Qv2ray \/ V2rayXS 配置二维码"/,
    `${file} first QR code should identify the compatible clients`
  )
  assert.match(
    qrImages[1],
    /alt="V2rayN \/ V2rayNG \/ V2rayXS 配置二维码"/,
    `${file} second QR code should identify the compatible clients`
  )
}

console.log('bindingQrAltAccessibility tests passed')
