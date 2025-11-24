<script setup>
// import { reactive } from 'vue'

// const state = reactive({
//   count: 0,
// })

// const setCount = () => {
//   state.count = 0
// }

import { ref } from 'vue'
const state = ref(0)

const setCount = () => {
  // 脚本区域修改reef产生的响应式对象数据必须通过.value属性
  state.value = 0
}
</script>

<template>
  <div>
    <button @click="state++">count is: {{ state }}</button>
    <button @click="setCount">{{ state }}</button>
  </div>
</template>
