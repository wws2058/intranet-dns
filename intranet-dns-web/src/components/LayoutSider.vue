<template>
  <a-layout-sider v-model:collapsed="store.collapsed" collapsible :trigger="null">
    <div class="sider-logo">
      <DnsLogo></DnsLogo>
      <span v-if="!store.collapsed">DNS内网管控系统</span>
    </div>
    <a-menu v-model:selectedKeys="store.selectedKeys" mode="inline" theme="dark" :items="store.siderNavs"
      :openKeys="store.openkeys"></a-menu>
  </a-layout-sider>
</template>

<script setup>
import { layoutStore } from '@/store/layout';
const store = layoutStore();
</script>

<style scoped>
.sider-logo {
  display: flex;
  height: 32px;
  margin: 16px;
  text-align: left;
  font-size: em;
  color: white;
  line-height: 32px;
  justify-content: center;
}
</style>
