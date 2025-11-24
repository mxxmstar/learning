<!-- <script >
export default {
    setup() {
      console.log('setup', this)
      const message = 'this is message'
      const logMessage = () => {
        console.log(message)
      }
      return {
        message,
        logMessage,
      }
    },
    beforeCreate() {
      console.log('beforeCreate')
    },
  }
</script> -->
<script setup>
  const message = 'this is message'
  const logMessage = () => {
    console.log(message)
  }
</script>
<template>
  <div>
    <h1>Hello Vue!</h1>
  </div>
  <div>
    {{ message }}
    <button @click="logMessage">log message</button>
  </div>
</template>
