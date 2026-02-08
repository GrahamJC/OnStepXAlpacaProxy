<script setup>
import { useOnStepXStore } from '../stores/onstepx'
import { storeToRefs } from 'pinia'

const store = useOnStepXStore()
const { firmwareVersion, comPort, connectionStatus, isConnected, proxyVersion, position } = storeToRefs(store)
</script>

<template>
  <header class="main-header glass-panel">
      <div class="header-content">
          <h1>OnStepX Proxy</h1>
          <span class="subtitle">Alpaca Driver & Controller <span id="proxy-version-display">{{ proxyVersion }}</span></span>
      </div>
      <div class="header-badges">
          <div class="status-badge" id="com-port-badge">
              <span class="value">{{ comPort }}</span>
          </div>
          <div class="status-badge" id="firmware-badge">
              <span class="label">FW</span>
              <span class="value">{{ firmwareVersion }}</span>
          </div>
          <div class="connection-pill" id="connection-status-pill">
              <span id="connection-indicator" :class="{ connected: isConnected, disconnected: !isConnected }"></span>
              <span id="connection-text">{{ connectionStatus }}</span>
          </div>
      </div>
      <div>
        RA: {{ position.ra }}
        Dec: {{ position.dec }}
        Alt: {{ position.alt }}
        Az: {{ position.az }}
      </div>
  </header>
</template>

<style scoped>
/* Scoped styles can be added here if needed, but we rely on global style.css for now */
</style>
