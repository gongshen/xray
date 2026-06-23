const CONFIG_LABELS = {
  config1: 'Shadowrocket',
  config2: 'V2rayN',
}

export function getShareLink(shareInfo, configType) {
  if (!shareInfo) {
    return ''
  }

  if (configType === 'config1') {
    return shareInfo.share1_link || ''
  }

  if (configType === 'config2') {
    return shareInfo.share2_link || ''
  }

  return ''
}

export function getCopySuccessMessage(configType) {
  const label = CONFIG_LABELS[configType]
  return label ? `${label} 配置已复制到剪贴板` : '配置已复制到剪贴板'
}

export function buildQrDownloadName(filename) {
  if (!filename) {
    return 'qrcode.png'
  }

  return filename.endsWith('.png') ? filename : `${filename}.png`
}

export async function createShareDialogInfo(data, toDataUrl) {
  const share1Link = data?.share1 || ''
  const share2Link = data?.share2 || ''
  const [share1, share2] = await Promise.all([
    share1Link ? toDataUrl(share1Link) : '',
    share2Link ? toDataUrl(share2Link) : '',
  ])

  return {
    share1,
    share1_link: share1Link,
    share2,
    share2_link: share2Link,
  }
}
