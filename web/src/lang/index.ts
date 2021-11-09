import { createI18n } from 'vue-i18n'
import Cookies from 'js-cookie'
import elementEnLocale from 'element-plus/lib/locale/lang/en' // element-ui lang
import elementZhLocale from 'element-plus/lib/locale/lang/zh-cn'
import localMessages from '@intlify/vite-plugin-vue-i18n/messages'

export const elementLocale = {
  [elementEnLocale.name]: elementEnLocale,
  [elementZhLocale.name]: elementZhLocale,
}

export function getLanguage(): string {
  const chooseLanguage = Cookies.get('language')
  if (chooseLanguage) return chooseLanguage

  // if has not choose language
  const language = navigator.language.toLowerCase()
  const locales = Object.keys(localMessages)
  for (const locale of locales) {
    if (language.indexOf(locale) > -1) {
      return locale
    }
  }
  return 'en'
}
const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: getLanguage(),
  messages: localMessages,
})
export default i18n
