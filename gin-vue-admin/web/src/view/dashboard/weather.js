
import { ref } from 'vue'

const weatherInfo = ref('本地模式运行中，第三方天气接口已关闭。')

export const useWeatherInfo = () => {
  return weatherInfo
}
