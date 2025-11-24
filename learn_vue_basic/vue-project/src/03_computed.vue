<script setup>
import { ref } from 'vue'
import { computed } from 'vue'

const list = ref([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
const computedList = computed(() => {
  return list.value.filter(item => item > 5)
})

setTimeout(() => {
  list.value.push(11, 12)
}, 3000)
</script>

<template>
  <div>
    原始响应数组- {{ list }}
  </div>
  <div>
    计算属性- {{ computedList }}
  </div>
</template>
