<template>
  <div class="fab">
    <el-popover
      placement="right-end"
      trigger="hover"
      popper-class="fab-popper"
    >
      <div style="margin: 0 30px">
        <el-row class="fab-item"><el-link :underline="false" @click="showTransformDialog('time')">时间转换</el-link></el-row>
        <el-row class="fab-item"><el-link :underline="false" @click="showTransformDialog('password')">随机密码</el-link></el-row>
        <el-row class="fab-item"><el-link :underline="false" @click="showTransformDialog('unicode')">unicode</el-link></el-row>
        <el-row class="fab-item"><el-link :underline="false" @click="showTransformDialog('decodeURI')">decodeURI</el-link></el-row>
        <el-row class="fab-item"><el-link :underline="false" @click="showTransformDialog('md5')">md5</el-link></el-row>
        <el-row class="fab-item"><el-link :underline="false" @click="showTransformDialog('cron')">cron生成</el-link></el-row>
        <el-row class="fab-item"><el-link :underline="false" @click="showTransformDialog('qrcode')">二维码</el-link></el-row>
        <el-row class="fab-item"><el-link :underline="false" @click="showTransformDialog('byte')">字节转换</el-link></el-row>
        <el-row class="fab-item"><el-link :underline="false" @click="showTransformDialog('color')">颜色转换</el-link></el-row>
      </div>
      <div slot="reference" class="fab-cell">
        <i class="el-icon-s-cooperation fab-icon" />
      </div>
    </el-popover>
    <el-dialog
      :visible.sync="transformVisible"
      width="600px"
      append-to-body
      :close-on-click-modal="false"
    >
      <el-row class="transform-content">
        <el-row v-show="transformType === 'time'">
          <el-row>
            <span style="display:inline-block;width:60px;font-size:14px;margin-right:10px">时间戳</span>
            <el-input v-model="timeExchange.timestamp" style="width: 200px" :placeholder="timeExchange.placeholder" clearable />
            <el-button type="primary" @click="timestampToDate">转换>></el-button>
            <el-input v-model="timeExchange.date" style="width: 200px" />
          </el-row>
          <el-row style="margin-top: 10px">
            <span style="display:inline-block;width:60px;font-size:14px;margin-right:10px">时间</span>
            <el-input v-model="dateExchange.date" style="width: 200px" :placeholder="dateExchange.placeholder" clearable />
            <el-button type="primary" @click="dateToTimestamp">转换>></el-button>
            <el-input v-model="dateExchange.timestamp" style="width: 200px" />
          </el-row>
        </el-row>
        <el-row v-show="transformType === 'cron'">
          <el-input v-model="cron.expression" placeholder="* * * * ?" style="width: 450px; margin-bottom:10px" />
          <el-button type="primary" @click="crontabTranslate">转换>></el-button>
          <el-row style="margin-left: 5px;">{{ cron.chinese }}</el-row>
          <pre>
    *    *    *    *    *
    -    -    -    -    -
    |    |    |    |    |
    |    |    |    |    +----- 星期中星期几 (0 - 7) (星期天为0)
    |    |    |    +---------- 月份 (1 - 12)
    |    |    +--------------- 一个月中的第几天 (1 - 31)
    |    +-------------------- 小时 (0 - 23)
    +------------------------- 分钟 (0 - 59)
          </pre>
          <el-row style="padding: 0 5px;">
            在以上各个字段中，还可以使用以下特殊字符：
            <p>星号( * )：代表所有可能的值，例如month字段如果是星号，则表示在满足其它字段的制约条件后每月都执行该命令操作。</p>
            <p>逗号( , )：可以用逗号隔开的值指定一个列表范围，例如，"1,2,5,7,8,9"</p>
            <p>中杠( - )：可以用整数之间的中杠表示一个整数范围，例如"2-6"表示"2,3,4,5,6"</p>
            <p>正斜线( / )：可以用正斜线指定时间的间隔频率，例如"0-23/2"表示每两小时执行一次。同时正斜线可以和星号一起使用，例如*/10，如果用在minute字段，表示每十分钟执行一次。</p>
          </el-row>
        </el-row>
        <el-row v-show="transformType === 'unicode'">
          <el-input
            v-model="unicode.escape"
            type="textarea"
            :autosize="{ minRows: 2}"
            placeholder="请输入未转义的unicode编码"
          />
          <el-input
            :value="unicodeUnescape"
            style="margin-top:10px"
            type="textarea"
            :autosize="{ minRows: 2}"
            readonly
          />
        </el-row>
        <el-row v-show="transformType === 'decodeURI'">
          <el-input
            v-model="decodeURI.escape"
            type="textarea"
            :autosize="{ minRows: 2}"
            placeholder="请输入未转义的URI"
          />
          <el-input
            :value="decodeURI.escape ? decodeURI(decodeURI.escape) : ''"
            style="margin-top:10px"
            type="textarea"
            :autosize="{ minRows: 2}"
            readonly
          />
        </el-row>
        <el-row v-show="transformType === 'password'">
          <el-checkbox-group v-model="password.checkbox">
            <el-checkbox :label="1">A-Z</el-checkbox>
            <el-checkbox :label="2">a-z</el-checkbox>
            <el-checkbox :label="3">0-9</el-checkbox>
            <el-checkbox :label="4">!@#$%^&*</el-checkbox>
          </el-checkbox-group>
          <el-row style="margin-top: 10px">
            <span style="display:inline-block;width:60px;font-size:14px;margin-right:5px">密码长度</span>
            <el-input-number
              v-model="password.length"
              :min="1"
              placeholder="请输入密码长度"
            />
            <el-button type="primary" @click="createPassword">生成>></el-button>
          </el-row>

          <el-input
            :value="password.text"
            style="margin-top:10px"
            readonly
          />
        </el-row>
        <el-row v-show="transformType === 'md5'">
          <el-input
            v-model="md5.text"
            type="textarea"
            :autosize="{ minRows: 3}"
            placeholder="请输入....."
          />
          <el-input
            :value="md5.text|md5"
            style="margin-top:10px"
            readonly
          />
        </el-row>
        <el-row v-show="transformType === 'color'">
          <el-row>
            <span style="display:inline-block;width:60px;font-size:14px;margin-right:10px">十六进制</span>
            <el-input v-model="cHexExchange.hex" style="width: 200px" placeholder="#FFFFFF" clearable />
            <el-button type="primary" @click="hexToRGB">转换>></el-button>
            <el-input v-model="cHexExchange.rgb" style="width: 200px" />
          </el-row>
          <el-row style="margin-top: 10px">
            <span style="display:inline-block;width:60px;font-size:14px;margin-right:10px">RGB</span>
            <el-input v-model="RGBExchange.rgb" style="width: 200px" placeholder="(255,255,255)" clearable />
            <el-button type="primary" @click="rgbToHex">转换>></el-button>
            <el-input v-model="RGBExchange.hex" style="width: 200px" />
          </el-row>
        </el-row>
        <el-row v-show="transformType === 'byte'">
          <span style="display:inline-block;width:60px;font-size:14px;margin-right:10px">字节</span>
          <el-input v-model="bytes" style="width: 130px" />
          <el-select v-model="bytesUnit" style="width: 70px">
            <el-option :value="1" label="B" />
            <el-option :value="1*1024" label="KB" />
            <el-option :value="1024*1024" label="MB" />
          </el-select>
          <el-button type="primary" @click="bytesToHumanSize">转换>></el-button>
          <el-input v-model="humanSize" style="width: 200px" />
        </el-row>
        <el-row v-show="transformType === 'qrcode'">
          <el-input
            v-model="qrcode.text"
            type="textarea"
            :autosize="{ minRows: 2}"
            placeholder="请输入内容"
          />
          <el-row style="margin-top: 10px">
            <span style="display:inline-block;width:30px;font-size:14px;margin-right:10px">大小</span>
            <el-input-number v-model="qrcode.width" />
          </el-row>
          <vue-qrcode class="text-align:center" :value="qrcode.text" :width="qrcode.width" />
        </el-row>
      </el-row>
    </el-dialog>
  </div>
</template>

<script>
import VueQrcode from 'vue-qrcode'
import cronstrue from 'cronstrue/i18n'
import { parseTime, humanSize } from '@/utils'
import { md5 } from '@/utils/md5'
export default {
  components: {
    VueQrcode
  },
  filters: { md5 },
  data() {
    return {
      transformVisible: false,
      transformType: '',
      timeExchange: {
        date: parseTime(new Date()),
        timestamp: '',
        timer: null,
        placeholder: Date.parse(new Date()) / 1000
      },
      dateExchange: {
        date: '',
        timestamp: Date.parse(new Date()) / 1000,
        timer: null,
        placeholder: parseTime(new Date())
      },
      cHexExchange: {
        hex: '',
        rgb: ''
      },
      RGBExchange: {
        hex: '',
        rgb: ''
      },
      qrcode: {
        text: 'https://github.com/zhenorzz/goploy',
        width: 200
      },
      unicode: {
        escape: ''
      },
      decodeURI: {
        escape: ''
      },
      password: {
        checkbox: [1, 2, 3],
        length: 8,
        text: ''
      },
      md5: {
        text: ''
      },
      cron: {
        expression: '',
        chinese: ''
      },
      bytes: '',
      bytesUnit: 1,
      humanSize: ''
    }
  },
  computed: {
    // 计算属性的 getter
    unicodeUnescape: function() {
      // `this` 指向 vm 实例
      return unescape(this.unicode.escape.replace(/\\u/g, '%u'))
    }
  },
  created() {
    this.timeExchange.timer = setInterval(() => {
      this.timeExchange.placeholder = Date.parse(new Date()) / 1000
    }, 1000)
    this.dateExchange.timer = setInterval(() => {
      this.dateExchange.placeholder = parseTime(new Date())
    }, 1000)
  },
  beforeDestroy() {
    clearTimeout(this.timeExchange.timer)
    clearTimeout(this.dateExchange.timer)
  },
  methods: {
    showTransformDialog(type) {
      this.transformVisible = true
      this.transformType = type
    },
    timestampToDate() {
      this.timeExchange.date = parseTime(this.timeExchange.timestamp)
    },
    dateToTimestamp() {
      this.dateExchange.timestamp = Date.parse(new Date(this.dateExchange.date)) / 1000
    },
    hexToRGB() {
      const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(this.cHexExchange.hex)

      if (result) {
        const r = parseInt(result[1], 16)
        const g = parseInt(result[2], 16)
        const b = parseInt(result[3], 16)
        this.cHexExchange.rgb = 'rgb(' + r + ', ' + g + ', ' + b + ')'
      } else {
        this.cHexExchange.rgb = 'rgb(0, 0, 0)'
      }
    },
    rgbToHex() {
      const color = this.RGBExchange.rgb.replace(/\(|\)|rgb/g, '')
      const rgb = color.split(',')
      const r = parseInt(rgb[0])
      const g = parseInt(rgb[1])
      const b = parseInt(rgb[2])
      const hex = '#' + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)
      this.RGBExchange.hex = hex.toLocaleUpperCase()
    },
    bytesToHumanSize() {
      this.humanSize = humanSize(this.bytes * this.bytesUnit)
    },
    crontabTranslate() {
      try {
        this.cron.chinese = cronstrue.toString(this.cron.expression, { locale: 'zh_CN' })
      } catch (error) {
        this.$message.error(error)
      }
    },
    createPassword() {
      let randArr = []
      for (const num of this.password.checkbox) {
        if (num === 1) {
          for (let i = 0; i < 26; i++) {
            randArr.push(String.fromCharCode(65 + i))
          }
        } else if (num === 2) {
          for (let i = 0; i < 26; i++) {
            randArr.push(String.fromCharCode(97 + i))
          }
        } else if (num === 3) {
          for (let i = 0; i < 10; i++) {
            randArr.push(i)
          }
        } else {
          randArr = randArr.concat(['!', '@', '#', '$', '%', '^', '&', '*'])
        }
      }
      let password = ''
      for (let i = 0; i < this.password.length; i++) {
        password += randArr[Math.floor(Math.random() * randArr.length)]
      }
      this.password.text = password
    }
  }
}
</script>

<style lang="scss" scoped>
@import "@/styles/mixin.scss";
.fab {
  z-index: 20;
  position: fixed;
  right: 0;
  bottom: 40px;
  width: 36px;
  border-radius: 4px 0 0 4px;
  background-color: #fff;
  box-shadow: 0 2px 8px 0 rgba(0,0,0,.18);
  text-align: center;
  &-cell {
    display: block;
    position: relative;
    width: 36px;
    height: 36px;
    transform: translateZ(0);
    color: #999;
    cursor: pointer;
  }
  &-icon {
    line-height: 36px;
    color: #999;
    cursor: pointer;
    width: 20px;
    height: 20px;
  }
}

.transform {
  &-content {
    max-height: 500px;
    overflow-y: auto;
    @include scrollBar();
  }
}

</style>
<style lang="scss">
.fab-popper {
  padding: 0;
  min-width: 120px !important;
}
.fab-item {
  margin: 10px 0;
}
</style>
