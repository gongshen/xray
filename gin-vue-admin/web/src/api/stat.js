import { createStatApi } from './statApiFactory'

export const {
  createStat,
  deleteStat,
  deleteStatByIds,
  updateStat,
  findStat,
  getStatList,
  getStatCharts,
  getStatRank,
} = createStatApi('/v2ray_admin/stat', { includeRank: true })
