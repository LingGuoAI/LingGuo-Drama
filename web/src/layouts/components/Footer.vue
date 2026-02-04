<template>
  <div class="footer">
    Copyright © {{ new Date().getFullYear() }} 灵果短剧. All Rights Reserved
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'Footer'
})
</script>

<style scoped>
.footer {
  color: #999;
  line-height: 20px;
  text-align: center;
  padding: 16px 0;
}
</style>