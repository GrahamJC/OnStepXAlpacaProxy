<script setup lang="ts">
import onstepx from './store';
import Site from './components/Site.vue'
</script>

<template>
  <h1>OnstepX</h1>
  <p class="">
    Firmware: {{ onstepx.firmwareVersion }}
  </p>
  <button @click="onstepx.firmwareVersion = 'x.xx'">Change</button>
  <Site></Site>
</template>

<style scoped></style>
