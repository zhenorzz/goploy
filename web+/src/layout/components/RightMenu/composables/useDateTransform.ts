import { reactive, onUnmounted } from 'vue'
import { parseTime } from '@/utils'

export default function useDateTransform() {
  const timeExchange = reactive({
    date: parseTime(new Date().getTime()),
    timestamp: '',
    timer: setInterval(() => {
      timeExchange.placeholder = String(Math.round(Date.now() / 1000))
    }, 1000),
    placeholder: String(Math.round(Date.now() / 1000)),
  })
  const dateExchange = reactive({
    date: '',
    timestamp: Math.round(Date.now() / 1000),
    timer: setInterval(() => {
      dateExchange.placeholder = parseTime(new Date().getTime())
    }, 1000),
    placeholder: parseTime(new Date().getTime()),
  })
  const timestamp = (value: string) => {
    let ts = 0
    switch (value) {
      case 'now':
        ts = Math.round(Date.now() / 1000)
        break
      case 'today':
        ts = Math.round(
          new Date(new Date().setHours(0, 0, 0, 0)).getTime() / 1000
        )
        break
      case 'm1d':
        ts =
          timeExchange.timestamp !== ''
            ? parseInt(timeExchange.timestamp) - 86400
            : Math.round(Date.now() / 1000) - 86400
        break
      case 'p1d':
        ts =
          timeExchange.timestamp !== ''
            ? parseInt(timeExchange.timestamp) + 86400
            : Math.round(Date.now() / 1000) + 86400
        break
      default:
        ts = Math.round(Date.now() / 1000)
    }
    timeExchange.timestamp = String(ts)
    timeExchange.date = parseTime(ts)
  }

  const timestampToDate = () => {
    timeExchange.date = parseTime(Number(timeExchange.timestamp))
  }
  const dateToTimestamp = () => {
    dateExchange.timestamp = new Date(dateExchange.date).getTime() / 1000
  }
  onUnmounted(() => {
    clearTimeout(timeExchange.timer)
    clearTimeout(dateExchange.timer)
  })
  return {
    timeExchange,
    dateExchange,
    timestamp,
    timestampToDate,
    dateToTimestamp,
  }
}
