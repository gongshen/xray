import { BarChart, LineChart } from 'echarts/charts'
import {
  DataZoomComponent,
  DataZoomInsideComponent,
  GridComponent,
  LegendComponent,
  TooltipComponent,
} from 'echarts/components'
import { init, use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'

use([
  BarChart,
  LineChart,
  DataZoomComponent,
  DataZoomInsideComponent,
  GridComponent,
  LegendComponent,
  TooltipComponent,
  CanvasRenderer,
])

export { init }
