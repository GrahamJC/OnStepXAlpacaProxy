<script setup>
import { onMounted, ref } from 'vue'
import { useOnStepXStore } from './stores/onstepx'
import Header from './components/Header.vue'

const store = useOnStepXStore()

onMounted(() => {
    store.startPolling()
})
</script>

<template>
  <div class="container">
    <Header />
   </div>
</template>

<style scoped>
.section-spacing {
    margin-top: 1.5rem;
}
</style>
