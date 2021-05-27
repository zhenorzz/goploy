import { createI18n } from 'vue-i18n'
import Cookies from 'js-cookie'
import elementEnLocale from 'element-plus/lib/locale/lang/en' // element-ui lang
import elementZhLocale from 'element-plus/lib/locale/lang/zh-cn'
import localMessages from '@intlify/vite-plugin-vue-i18n/messages'
const messages = {
  [elementEnLocale.name]: {
    // el 这个属性很关键，一定要保证有这个属性，
    el: elementEnLocale.el,
    // 定义您自己的字典，但是请不要和 `el` 重复，这样会导致 ElementPlus 内部组件的翻译失效.
    ...localMessages[elementEnLocale.name],
  },
  [elementZhLocale.name]: {
    el: elementZhLocale.el,
    // 定义您自己的字典，但是请不要和 `el` 重复，这样会导致 ElementPlus 内部组件的翻译失效.
    ...localMessages[elementZhLocale.name],
  },
}

export function getLanguage(): string {
  const chooseLanguage = Cookies.get('language')
  if (chooseLanguage) return chooseLanguage

  // if has not choose language
  const language = navigator.language.toLowerCase()
  const locales = Object.keys(messages)
  for (const locale of locales) {
    if (language.indexOf(locale) > -1) {
      return locale
    }
  }
  return 'en'
}
const i18n = createI18n({
  globalInjection: true,
  locale: getLanguage(),
  messages,
})
export default i18n
