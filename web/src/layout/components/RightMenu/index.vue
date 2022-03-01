<template>
  <div class="fab">
    <el-popover placement="left-end" trigger="click" popper-class="fab-popper">
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
          @click="showTransformDialog('json')"
        >
          JSON Pretty
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
        <TheDatetransform v-model="transformType" />
        <TheJSONPretty v-model="transformType" />
        <TheRandomPWD v-model="transformType" />
        <TheUnicode v-model="transformType" />
        <el-row v-show="transformType === 'decodeURI'">
          <el-input
            v-model="uri.escape"
            type="textarea"
            :autosize="{ minRows: 2 }"
            placeholder="Please enter unescaped URI"
          />
          <el-input
            :value="uri.escape ? decodeURI(uri.escape) : ''"
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
        <TheCronstrue v-model="transformType" />
        <TheByteTransform v-model="transformType" />
        <TheRGBTransform v-model="transformType" />
      </el-row>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import VueQrcode from '@chenfengyuan/vue-qrcode'
import { md5 as hashByMD5 } from '@/utils/md5'
import TheDatetransform from './TheDatetransform.vue'
import TheJSONPretty from './TheJSONPretty.vue'
import TheRandomPWD from './TheRandomPWD.vue'
import TheUnicode from './TheUnicode.vue'
import TheCronstrue from './TheCronstrue.vue'
import TheByteTransform from './TheByteTransform.vue'
import TheRGBTransform from './TheRGBTransform.vue'
import { ref, reactive } from 'vue'

const transformVisible = ref(false)
const transformType = ref('')
const qrcode = reactive({
  text: 'https://github.com/zhenorzz/goploy',
  width: 200,
})
const uri = reactive({
  escape: '',
})
const md5 = reactive({
  text: '',
})
function showTransformDialog(type: string) {
  transformVisible.value = true
  transformType.value = type
}
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
