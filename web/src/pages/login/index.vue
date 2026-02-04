<template>
  <div class="login-wrapper">

    <div class="login-container">
      <div class="title-container">
        <h1 class="title">灵果短剧 - 管理系统</h1>

      </div>

      <login v-if="type === 'login'" />
      <!-- <register v-else @register-success="switchType('login')" /> -->
      <!-- <tdesign-setting /> -->
    </div>

    <footer class="copyright">Copyright @ 2021-2025 灵果短剧. All Rights Reserved</footer>
  </div>
</template>
<script lang="ts">
export default {
  name: 'LoginIndex',
};
</script>
<script setup lang="ts">
import { ref } from 'vue';

import Login from './components/Login.vue';

const type = ref('login');

</script>

<style lang="less" scoped>
@import './index.less';
</style>
