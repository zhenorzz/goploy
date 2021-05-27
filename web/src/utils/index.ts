/**
 * Created by PanJiaChen on 16/11/18.
 */

/**
 * 判断是否为空
 * @param {any} mixedVar
 * @return {boolean}
 */
export function empty(mixedVar: any): boolean {
  //  discuss at: https://locutus.io/php/empty/
  // original by: Philippe Baumann
  //    input by: Onno Marsman (https://twitter.com/onnomarsman)
  //    input by: LH
  //    input by: Stoyan Kyosev (https://www.svest.org/)
  // bugfixed by: Kevin van Zonneveld (https://kvz.io)
  // improved by: Onno Marsman (https://twitter.com/onnomarsman)
  // improved by: Francesco
  // improved by: Marc Jansen
  // improved by: Rafał Kukawski (https://blog.kukawski.pl)
  //   example 1: empty(null)
  //   returns 1: true
  //   example 2: empty(undefined)
  //   returns 2: true
  //   example 3: empty([])
  //   returns 3: true
  //   example 4: empty({})
  //   returns 4: true
  //   example 5: empty({'aFunc' : function () { alert('humpty'); } })
  //   returns 5: false

  let undef
  const emptyValues = [undef, null, false, 0, '', '0']

  for (let i = 0, len = emptyValues.length; i < len; i++) {
    if (mixedVar === emptyValues[i]) {
      return true
    }
  }

  if (typeof mixedVar === 'object') {
    for (const key in mixedVar) {
      if (Object.prototype.hasOwnProperty.call(mixedVar, key)) {
        return false
      }
    }
    return true
  }

  return false
}

/**
 * Parse the time to string
 */
export function parseTime(time: number, cFormat?: string): string {
  const format = cFormat || '{y}-{m}-{d} {h}:{i}:{s}'
  let date
  if (typeof time === 'object') {
    date = time
  } else {
    if (typeof time === 'string' && /^[0-9]+$/.test(time)) {
      time = parseInt(time)
    }
    if (typeof time === 'number' && time.toString().length === 10) {
      time = time * 1000
    }
    date = new Date(time)
  }
  const formatObj: { [key: string]: number } = {
    y: date.getFullYear(),
    m: date.getMonth() + 1,
    d: date.getDate(),
    h: date.getHours(),
    i: date.getMinutes(),
    s: date.getSeconds(),
    a: date.getDay(),
  }
  const time_str = format.replace(
    /{(y|m|d|h|i|s|a)+}/g,
    (result, key): string => {
      const value = formatObj[key]
      // Note: getDay() returns 0 on Sunday
      if (key === 'a') {
        return ['日', '一', '二', '三', '四', '五', '六'][value]
      }
      if (result.length > 0 && value < 10) {
        return '0' + value
      }
      return value ? '' + value : '0'
    }
  )
  return time_str
}

/**
 * Parse the string to time
 */
export function parseDate(date: string): number {
  date = date.substring(0, 19)
  date = date.replace(/-/g, '/')
  const timestamp = new Date(date).getTime()
  return timestamp / 1000
}

/**
 * @param {number} time
 * @param {string} option
 * @returns {string}
 */
export function formatTime(time: number, option: string): string {
  if (('' + time).length === 10) {
    time = parseInt('' + time) * 1000
  } else {
    time = +time
  }
  const d = new Date(time)
  const now = Date.now()

  const diff = (now - d.getTime()) / 1000

  if (diff < 30) {
    return '刚刚'
  } else if (diff < 3600) {
    // less 1 hour
    return Math.ceil(diff / 60) + '分钟前'
  } else if (diff < 3600 * 24) {
    return Math.ceil(diff / 3600) + '小时前'
  } else if (diff < 3600 * 24 * 2) {
    return '1天前'
  }
  if (option) {
    return parseTime(time, option)
  } else {
    return (
      d.getMonth() +
      1 +
      '月' +
      d.getDate() +
      '日' +
      d.getHours() +
      '时' +
      d.getMinutes() +
      '分'
    )
  }
}

/**
 */
export function param2Obj(url: string): any {
  const search = url.split('?')[1]
  if (!search) {
    return {}
  }
  return JSON.parse(
    '{"' +
      decodeURIComponent(search)
        .replace(/"/g, '\\"')
        .replace(/&/g, '","')
        .replace(/=/g, '":"')
        .replace(/\+/g, ' ') +
      '"}'
  )
}

/**
 * 深拷贝
 * 这个函数有个小问题，一般情况不会遇到，如果遇到请查看：
 * https://stackoverflow.com/a/4460624
 */
export function deepClone(item: any): any {
  if (!item) {
    return item
  } // null, undefined values check

  const types = [Number, String, Boolean]
  let result

  // normalizing primitives if someone did new String('aaa'), or new Number('444');
  types.forEach(function (type) {
    if (item instanceof type) {
      result = type(item)
    }
  })

  if (typeof result === 'undefined') {
    if (Object.prototype.toString.call(item) === '[object Array]') {
      result = []
      item.forEach(function (child, index) {
        result[index] = deepClone(child)
      })
    } else if (typeof item === 'object') {
      // testing that this is DOM
      if (item.nodeType && typeof item.cloneNode === 'function') {
        result = item.cloneNode(true)
      } else if (!item.prototype) {
        // check that this is a literal
        if (item instanceof Date) {
          result = new Date(item)
        } else {
          // it is an object literal
          result = {}
          for (const i in item) {
            result[i] = deepClone(item[i])
          }
        }
      } else {
        result = item
      }
    } else {
      result = item
    }
  }

  return result
}

/**
 * 人类可读的字节大小
 */
export function humanSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024 // or 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return (bytes / Math.pow(k, i)).toPrecision(3) + ' ' + sizes[i]
}

export function parseGitURL(url: string): string {
  const lastDotGitIndex = url.lastIndexOf('.git')
  if (lastDotGitIndex !== -1) {
    url = url.substring(0, lastDotGitIndex)
  }
  const lastAtIndex = url.lastIndexOf('@')
  if (lastAtIndex === -1) {
    return url
  } else {
    return 'http://' + url.substring(lastAtIndex + 1).replace(':', '/')
  }
}

export function copy(content: string): void {
  const el = document.createElement('textarea')
  el.value = content
  el.setAttribute('readonly', '')
  el.style.position = 'absolute'
  el.style.left = '-9999px'
  document.body.appendChild(el)
  const selected =
    document.getSelection().rangeCount > 0
      ? document.getSelection().getRangeAt(0)
      : false
  el.select()
  document.execCommand('copy')
  document.body.removeChild(el)
  if (selected) {
    document.getSelection().removeAllRanges()
    document.getSelection().addRange(selected)
  }
}
