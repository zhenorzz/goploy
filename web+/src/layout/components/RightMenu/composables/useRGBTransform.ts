import { reactive } from 'vue'

export default function useRGBTransform() {
  const cHexExchange = reactive({
    hex: '',
    rgb: '',
  })

  const rgbExchange = reactive({
    hex: '',
    rgb: '',
  })

  const hexToRGB = () => {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(
      cHexExchange.hex
    )

    if (result) {
      const r = parseInt(result[1], 16)
      const g = parseInt(result[2], 16)
      const b = parseInt(result[3], 16)
      cHexExchange.rgb = 'rgb(' + r + ', ' + g + ', ' + b + ')'
    } else {
      cHexExchange.rgb = 'rgb(0, 0, 0)'
    }
  }
  const rgbToHex = () => {
    const color = rgbExchange.rgb.replace(/\(|\)|rgb/g, '')
    const rgb = color.split(',')
    const r = parseInt(rgb[0])
    const g = parseInt(rgb[1])
    const b = parseInt(rgb[2])
    const hex =
      '#' + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)
    rgbExchange.hex = hex.toLocaleUpperCase()
  }

  return {
    cHexExchange,
    rgbExchange,
    hexToRGB,
    rgbToHex,
  }
}
