<template>
  <template v-if="!isFullscreen">
    <a-layout class="layout">
      <!-- 导航栏 -->
      <LayoutSider></LayoutSider>
      <a-layout>
        <!-- 页眉 -->
        <LayoutHeader></LayoutHeader>
        <!-- 路由信息展示区 -->
        <LayoutMain></LayoutMain>
        <!-- 页脚 -->
        <LayoutFooter></LayoutFooter>
      </a-layout>
    </a-layout>
  </template>
  <template v-else>
    <!-- 全局显示, 如用户登录页面 -->
    <router-view></router-view>
  </template>
</template>

<script setup>
import LayoutHeader from './components/LayoutHeader.vue';
import LayoutMain from './components/LayoutMain.vue';
import LayoutFooter from './components/LayoutFooter.vue';
import LayoutSider from './components/LayoutSider.vue';
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();

// 全屏页面如用户登录
const isFullscreen = computed(() => {
  return route.meta.fullScreen ? route.meta.fullScreen : false;
});

</script>

<style>
body {
  margin: 0;
}

.layout {
  height: 100vh;
}
</style>
