<template>
  <div class="fab">
    <el-popover placement="right-end" trigger="hover" popper-class="fab-popper">
      <div style="margin: 0 10px">
        <el-button
          type="text"
          size="medium"
          @click="showTransformDialog('time')"
        >
          Date Transform
        </el-button>
        <el-button
          type="text"
          size="medium"
          @click="showTransformDialog('password')"
        >
          Random PWD
        </el-button>
        <el-button
          type="text"
          size="medium"
          @click="showTransformDialog('unicode')"
        >
          Unicode
        </el-button>
        <el-button
          type="text"
          size="medium"
          @click="showTransformDialog('decodeURI')"
        >
          DecodeURI
        </el-button>
        <div>
          <el-button
            type="text"
            size="medium"
            @click="showTransformDialog('md5')"
          >
            MD5
          </el-button>
        </div>
        <el-button
          type="text"
          size="medium"
          @click="showTransformDialog('cron')"
        >
          Crontab
        </el-button>
        <el-button
          type="text"
          size="medium"
          @click="showTransformDialog('qrcode')"
        >
          QRcode
        </el-button>
        <el-button
          type="text"
          size="medium"
          @click="showTransformDialog('byte')"
        >
          Byte Transform
        </el-button>
        <el-button
          type="text"
          size="medium"
          @click="showTransformDialog('color')"
        >
          Color Transform
        </el-button>
      </div>
      <template #reference>
        <div class="fab-cell">
          <i class="el-icon-s-cooperation fab-icon" />
        </div>
      </template>
    </el-popover>
    <el-dialog
      v-model="transformVisible"
      width="600px"
      append-to-body
      :close-on-click-modal="false"
    >
      <el-row class="transform-content">
        <el-row v-show="transformType === 'time'">
          <el-button
            style="margin-left: 80px"
            type="primary"
            @click="timestamp('now')"
          >
            {{ $t('now') }}
          </el-button>
          <el-button type="primary" @click="timestamp('today')">
            {{ $t('today') }}
          </el-button>
          <el-button type="primary" @click="timestamp('m1d')">
            {{ $t('m1d') }}
          </el-button>
          <el-button type="primary" @click="timestamp('p1d')">
            {{ $t('p1d') }}
          </el-button>
          <el-row style="margin-top: 10px" type="flex" align="middle">
            <span style="width: 70px; font-size: 14px; margin-right: 10px">
              Timestamp
            </span>
            <el-input
              v-model="timeExchange.timestamp"
              style="width: 200px"
              :placeholder="timeExchange.placeholder"
              clearable
              @keyup.enter="timestampToDate"
            />
            <el-button type="primary" @click="timestampToDate">>></el-button>
            <el-input v-model="timeExchange.date" style="width: 200px" />
          </el-row>
          <el-row style="margin-top: 10px" type="flex" align="middle">
            <span style="width: 70px; font-size: 14px; margin-right: 10px">
              Date
            </span>
            <el-input
              v-model="dateExchange.date"
              style="width: 200px"
              :placeholder="dateExchange.placeholder"
              clearable
              @keyup.enter="dateToTimestamp"
            />
            <el-button type="primary" @click="dateToTimestamp">>></el-button>
            <el-input v-model="dateExchange.timestamp" style="width: 200px" />
          </el-row>
        </el-row>
        <el-row v-show="transformType === 'password'">
          <el-checkbox-group v-model="password.checkbox">
            <el-checkbox :label="1">A-Z</el-checkbox>
            <el-checkbox :label="2">a-z</el-checkbox>
            <el-checkbox :label="3">0-9</el-checkbox>
            <el-checkbox :label="4">!@#$%^&*</el-checkbox>
          </el-checkbox-group>
          <el-row style="margin-top: 10px" type="flex" align="middle">
            <span style="width: 60px; font-size: 14px; margin-right: 5px">
              Length
            </span>
            <el-input-number
              v-model="password.length"
              :min="1"
              placeholder="Please enter the password length"
            />
            <el-button type="primary" @click="createPassword">Gen</el-button>
          </el-row>

          <el-input :value="password.text" style="margin-top: 10px" readonly />
        </el-row>
        <el-row v-show="transformType === 'unicode'">
          <el-input
            v-model="unicode.escape"
            type="textarea"
            :autosize="{ minRows: 2 }"
            placeholder="Please enter unescaped unicode encoding"
          />
          <el-input
            :value="unicodeUnescapeStr"
            style="margin-top: 10px"
            type="textarea"
            :autosize="{ minRows: 2 }"
            readonly
          />
        </el-row>
        <el-row v-show="transformType === 'decodeURI'">
          <el-input
            v-model="decodeURI.escape"
            type="textarea"
            :autosize="{ minRows: 2 }"
            placeholder="Please enter unescaped URI"
          />
          <el-input
            :value="decodeURI.escape ? decodeURI(decodeURI.escape) : ''"
            style="margin-top: 10px"
            type="textarea"
            :autosize="{ minRows: 2 }"
            readonly
          />
        </el-row>
        <el-row v-show="transformType === 'md5'">
          <el-input
            v-model="md5.text"
            type="textarea"
            :autosize="{ minRows: 3 }"
          />
          <el-input
            :value="hashByMD5(md5.text)"
            style="margin-top: 10px"
            readonly
          />
        </el-row>
        <el-row v-show="transformType === 'cron'">
          <el-input
            v-model="cron.expression"
            placeholder="* * * * ?"
            style="width: 450px; margin-bottom: 10px"
          />
          <el-button type="primary" @click="crontabTranslate">>></el-button>
          <el-row style="margin-left: 5px">{{ cron.chinese }}</el-row>
          <pre>
    *    *    *    *    *
    -    -    -    -    -
    |    |    |    |    |
    |    |    |    |    +----- Week (0 - 7) (0 for sunday)
    |    |    |    +---------- Month (1 - 12)
    |    |    +--------------- Day (1 - 31)
    |    +-------------------- Hour (0 - 23)
    +------------------------- Minute (0 - 59)
          </pre>
          <el-row style="padding: 0 5px">
            在以上各个字段中，还可以使用以下特殊字符：
            <p>
              星号( *
              )：代表所有可能的值，例如month字段如果是星号，则表示在满足其它字段的制约条件后每月都执行该命令操作。
            </p>
            <p>
              逗号( , )：可以用逗号隔开的值指定一个列表范围，例如，"1,2,5,7,8,9"
            </p>
            <p>
              中杠( -
              )：可以用整数之间的中杠表示一个整数范围，例如"2-6"表示"2,3,4,5,6"
            </p>
            <p>
              正斜线( /
              )：可以用正斜线指定时间的间隔频率，例如"0-23/2"表示每两小时执行一次。同时正斜线可以和星号一起使用，例如*/10，如果用在minute字段，表示每十分钟执行一次。
            </p>
          </el-row>
        </el-row>
        <el-row v-show="transformType === 'qrcode'">
          <el-input
            v-model="qrcode.text"
            type="textarea"
            :autosize="{ minRows: 2 }"
          />
          <el-row style="margin-top: 10px" type="flex" align="middle">
            <span style="width: 30px; font-size: 14px; margin-right: 10px">
              Size
            </span>
            <el-input-number v-model="qrcode.width" />
          </el-row>
          <VueQrcode
            class="text-align:center"
            :value="qrcode.text"
            :options="{ width: qrcode.width }"
          />
        </el-row>
        <el-row v-show="transformType === 'byte'" type="flex" align="middle">
          <span style="width: 40px; font-size: 14px; margin-right: 10px">
            Byte
          </span>
          <el-input
            v-model="bytes"
            style="width: 130px"
            @keyup.enter="bytesToHumanSize"
          />
          <el-select v-model="bytesUnit" style="width: 70px">
            <el-option :value="1" label="B" />
            <el-option :value="1 * 1024" label="KB" />
            <el-option :value="1024 * 1024" label="MB" />
          </el-select>
          <el-button type="primary" @click="bytesToHumanSize">>></el-button>
          <el-input v-model="humanSize" style="width: 200px" />
        </el-row>
        <el-row v-show="transformType === 'color'">
          <el-row type="flex" align="middle">
            <span style="width: 40px; font-size: 14px; margin-right: 10px">
              HEX
            </span>
            <el-input
              v-model="cHexExchange.hex"
              style="width: 200px"
              placeholder="#FFFFFF"
              clearable
              @keyup.enter="hexToRGB"
            />
            <el-button type="primary" @click="hexToRGB">>></el-button>
            <el-input v-model="cHexExchange.rgb" style="width: 200px" />
          </el-row>
          <el-row style="margin-top: 10px" type="flex" align="middle">
            <span style="width: 40px; font-size: 14px; margin-right: 10px">
              RGB
            </span>
            <el-input
              v-model="rgbExchange.rgb"
              style="width: 200px"
              placeholder="(255,255,255)"
              clearable
              @keyup.enter="rgbToHex"
            />
            <el-button type="primary" @click="rgbToHex">>></el-button>
            <el-input v-model="rgbExchange.hex" style="width: 200px" />
          </el-row>
        </el-row>
      </el-row>
    </el-dialog>
  </div>
</template>

<script>
import VueQrcode from '@chenfengyuan/vue-qrcode'
import cronstrue from 'cronstrue/i18n'
import { humanSize } from '@/utils'
import { md5 as hashByMD5 } from '@/utils/md5'
import { defineComponent } from 'vue'
import useDateTransform from './composables/useDateTransform'
import useRandomPWD from './composables/useRandomPWD'
import useUnicode from './composables/useUnicode'
import useRGBTransform from './composables/useRGBTransform'

export default defineComponent({
  components: {
    VueQrcode,
  },
  setup() {
    const {
      timeExchange,
      dateExchange,
      timestamp,
      timestampToDate,
      dateToTimestamp,
    } = useDateTransform()

    const { password, createPassword } = useRandomPWD()

    const { unicode, unicodeUnescapeStr } = useUnicode()

    const { cHexExchange, rgbExchange, hexToRGB, rgbToHex } = useRGBTransform()

    return {
      timeExchange,
      dateExchange,
      timestamp,
      timestampToDate,
      dateToTimestamp,
      password,
      createPassword,
      unicode,
      unicodeUnescapeStr,
      cHexExchange,
      rgbExchange,
      hexToRGB,
      rgbToHex,
    }
  },
  data() {
    return {
      transformVisible: false,
      transformType: '',
      qrcode: {
        text: 'https://github.com/zhenorzz/goploy',
        width: 200,
      },
      decodeURI: {
        escape: '',
      },
      md5: {
        text: '',
      },
      cron: {
        expression: '',
        chinese: '',
      },
      bytes: '',
      bytesUnit: 1,
      humanSize: '',
    }
  },
  methods: {
    hashByMD5,
    showTransformDialog(type) {
      this.transformVisible = true
      this.transformType = type
    },
    bytesToHumanSize() {
      this.humanSize = humanSize(this.bytes * this.bytesUnit)
    },
    crontabTranslate() {
      try {
        this.cron.chinese = cronstrue.toString(this.cron.expression, {
          locale: 'zh_CN',
        })
      } catch (error) {
        this.$message.error(error)
      }
    },
  },
})
</script>

<style lang="scss" scoped>
@import '@/styles/mixin.scss';
.fab {
  z-index: 20;
  position: fixed;
  right: 0;
  bottom: 40px;
  width: 36px;
  border-radius: 4px 0 0 4px;
  background-color: #fff;
  box-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.18);
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
  .el-button {
    padding: 0;
    margin: 4px 0;
    color: #606266;
    min-height: 20px;
  }
  .el-button + .el-button {
    margin-left: 0px;
  }
}
</style>
