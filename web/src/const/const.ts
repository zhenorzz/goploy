const scriptLang = {
  Option: [
    { label: 'sh', value: 'sh', lang: 'sh' },
    { label: 'zsh', value: 'zsh', lang: 'sh' },
    { label: 'bash', value: 'bash', lang: 'sh' },
    { label: 'python', value: 'python', lang: 'python' },
    { label: 'php', value: 'php', lang: 'php' },
    { label: 'bat', value: 'cmd', lang: 'batchfile' },
    { label: 'yaml', value: 'yaml', lang: 'yaml' },
  ],
  getScriptLang: function (mode = '') {
    if (mode !== '') {
      const scriptInfo = scriptLang.Option.find((elem) => elem.value === mode)
      return scriptInfo ? scriptInfo['lang'] : ''
    } else {
      return 'sh'
    }
  },
}
export { scriptLang }
