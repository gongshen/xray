import { createStatApi } from './statApiFactory'

export const {
  createStat,
  deleteStat,
  deleteStatByIds,
  updateStat,
  findStat,
  getStatList,
  getStatCharts,
} = createStatApi('/v2ray/stat')
