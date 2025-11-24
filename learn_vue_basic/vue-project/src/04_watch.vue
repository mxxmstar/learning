<script setup>
import { watch } from 'vue'
import { ref } from 'vue'
// const count = ref(0)
// const setCount = () => {
//   count.value++
// }
// // 调用watch侦听
// watch(count, (newCount, oldCount) => {
//   console.log(newCount)
//   console.log(oldCount)
// })

const count = ref(0)
const setCount = () => {
  count.value++
}
const name = ref('mx')
const setName = () => {
  name.value = 'mxxm'
}
watch([count, name], ([newCount, newName], [oldCount, oldName]) => {
  console.log(newCount)
  console.log(oldCount)
  console.log(newName)
  console.log(oldName)
}, {
  immediate: true // 立即执行回调
})

const obj = ref({
  a: 1,
  b: 2
})
const setObj = () => {
  obj.value.a++
}
const setObjB = () => {
  obj.value.b++
}
watch(obj, (newObj, oldObj) => {
  console.log('obj')
  console.log(newObj)
  console.log(oldObj)
}, {
  deep: true  // 改变对象的属性值，也会触发回调
})
// 精确监听
// 写两个回调
watch(() => obj.value.a, (newA, oldA) => {
  console.log('a')
  console.log(newA)
  console.log(oldA)
})
</script>

<template>
  <div>
    <button @click="setCount">count is: {{ count }}</button>
    <button @click="setName">name is: {{ name }}</button>
  </div>
  <div>
    <button @click="setObj">obj.a is: {{ obj.a }}</button>
    <button @click="setObjB">obj.b is: {{ obj.b }}</button>
  </div>
</template>

