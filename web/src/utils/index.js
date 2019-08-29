/**
 * Created by PanJiaChen on 16/11/18.
 */

/**
 * Parse the time to string
 * @param {(Object|string|number)} time
 * @param {string} cFormat
 * @returns {string}
 */
export function parseTime(time, cFormat) {
  if (arguments.length === 0) {
    return null
  }
  const format = cFormat || '{y}-{m}-{d} {h}:{i}:{s}'
  let date
  if (typeof time === 'object') {
    date = time
  } else {
    if ((typeof time === 'string') && (/^[0-9]+$/.test(time))) {
      time = parseInt(time)
    }
    if ((typeof time === 'number') && (time.toString().length === 10)) {
      time = time * 1000
    }
    date = new Date(time)
  }
  const formatObj = {
    y: date.getFullYear(),
    m: date.getMonth() + 1,
    d: date.getDate(),
    h: date.getHours(),
    i: date.getMinutes(),
    s: date.getSeconds(),
    a: date.getDay()
  }
  const time_str = format.replace(/{(y|m|d|h|i|s|a)+}/g, (result, key) => {
    let value = formatObj[key]
    // Note: getDay() returns 0 on Sunday
    if (key === 'a') { return ['日', '一', '二', '三', '四', '五', '六'][value ] }
    if (result.length > 0 && value < 10) {
      value = '0' + value
    }
    return value || 0
  })
  return time_str
}

/**
 * @param {number} time
 * @param {string} option
 * @returns {string}
 */
export function formatTime(time, option) {
  if (('' + time).length === 10) {
    time = parseInt(time) * 1000
  } else {
    time = +time
  }
  const d = new Date(time)
  const now = Date.now()

  const diff = (now - d) / 1000

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
 * @param {string} url
 * @returns {Object}
 */
export function param2Obj(url) {
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
 * @param {mixed} item
 * @returns {Object}
 */
export function deepClone(item) {
  if (!item) { return item } // null, undefined values check

  var types = [Number, String, Boolean]
  var result

  // normalizing primitives if someone did new String('aaa'), or new Number('444');
  types.forEach(function(type) {
    if (item instanceof type) {
      result = type(item)
    }
  })

  if (typeof result === 'undefined') {
    if (Object.prototype.toString.call(item) === '[object Array]') {
      result = []
      item.forEach(function(child, index, array) {
        result[index] = deepClone(child)
      })
    } else if (typeof item === 'object') {
      // testing that this is DOM
      if (item.nodeType && typeof item.cloneNode === 'function') {
        result = item.cloneNode(true)
      } else if (!item.prototype) { // check that this is a literal
        if (item instanceof Date) {
          result = new Date(item)
        } else {
          // it is an object literal
          result = {}
          for (var i in item) {
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
 * @param {mixed} item
 * @returns {Object}
 */
export function humanSize(bytes) {
  if (bytes === 0) return '0 B'
  var k = 1000 // or 1024
  var sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
  var i = Math.floor(Math.log(bytes) / Math.log(k))

  return (bytes / Math.pow(k, i)).toPrecision(3) + ' ' + sizes[i]
}

