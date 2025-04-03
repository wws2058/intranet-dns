<template>
  <div>
    <h3>欢迎来到首页</h3>
    <p>这是基于Vue2 + Vuex3 + VueRouter3制作的练手项目, 帮助熟悉基本的Vue语法. </p>
  </div>
</template>
<script>
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'RootPage'
});
</script>
<style scoped>
div {
  width: 400px;
}
</style>
