<template>
  <el-row v-show="modelValue === 'cron'">
    <el-row type="flex" style="width: 100%">
      <el-input
        v-model="cron.expression"
        placeholder="* * * * ?"
        style="flex: 1"
      />
      <el-button type="primary" @click="crontabTranslate">>></el-button>
      <el-row style="margin-left: 5px">{{ cron.chinese }}</el-row>
    </el-row>
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
      <p>逗号( , )：可以用逗号隔开的值指定一个列表范围，例如，"1,2,5,7,8,9"</p>
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
</template>

<script lang="ts" setup>
import cronstrue from 'cronstrue/i18n'
import { reactive } from 'vue'
defineProps({
  modelValue: {
    type: String,
    default: '',
  },
})
const cron = reactive({
  expression: '',
  chinese: '',
})

const crontabTranslate = () => {
  try {
    cron.chinese = cronstrue.toString(cron.expression, {
      locale: 'zh_CN',
    })
  } catch (error) {
    ElMessage.error(error)
  }
}
</script>
