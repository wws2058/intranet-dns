<template>
  <a-layout-content class="content">
    <router-view></router-view>
  </a-layout-content>
</template>

<script setup>
import { useRouter, useRoute } from 'vue-router';
import { layoutStore } from '@/store/layout';
import { watch } from 'vue';
import { routesMap } from '@/router';

const store = layoutStore();
const router = useRouter();
const route = useRoute();

watch(() => store.selecttedLabel, () => {
  const routePath = routesMap[store.selecttedLabel];
  if (routePath) {
    router.push(routePath);
  } else {
    router.push("/");
  }
});

watch(() => route.path, () => {
  store.selectedKeys = [store.getSelectedKey(route.path)];
})

</script>

<style scoped>
.content {
  background-color: white;
  margin: 8px 8px;
  padding: 24px;
}
</style>
