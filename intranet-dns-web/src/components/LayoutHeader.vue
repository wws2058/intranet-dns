<template>
  <a-layout-header class="header">
    <div class="header-left">
      <MenuUnfoldOutlined v-if="store.collapsed" class="collapsed-icon" @click="store.switchCollapsed" />
      <MenuFoldOutlined v-else class="collapsed-icon" @click="store.switchCollapsed" />
      <span class="select-label">{{ store.selecttedLabel }}</span>
    </div>
    <div class="header-center">
      <div class="roll-msg-container">
        <span class="roll-msg">内部数据, 注意信息安全!!!</span>
      </div>
    </div>
    <div class="header-right">
      <a-dropdown arrow>
        <a class="ant-dropdown-link" @click.prevent>
          <UserOutlined></UserOutlined> {{ username }}
        </a>
        <template #overlay>
          <a-menu>
            <a-menu-item @click="logout">
              <span>退出登录</span>
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
  </a-layout-header>
</template>

<script setup>
import { MenuFoldOutlined, MenuUnfoldOutlined, UserOutlined } from '@ant-design/icons-vue';
import { layoutStore } from '@/store/layout';
import { localStoreUserDataKey } from '@/apis';
import { computed } from 'vue';
import { useRouter } from 'vue-router';

const store = layoutStore();

const userdataObj = JSON.parse(localStorage.getItem(localStoreUserDataKey));
const username = computed(() => {
  return userdataObj?.name ?? "anybody";
});

// useRouter依赖于Vue的响应式系统和组件实例, 只能在Vue组件的setup函数或者生命周期钩子中使用
// 如果这条语言写在logout函数中那么会导致报错, userRouter在非Vue组件的上下文中被调用, 无法正确获取到Vue组件实例
const router = useRouter();
function logout() {
  localStorage.removeItem(localStoreUserDataKey);
  router.push({ name: "UserLogin" });
}

// console.log(logout());
</script>

<style scoped>
/* 右侧布局头部 */
.header {
  background-color: white;
  padding: 0;
  height: 64px;
  font-size: 20px;
  overflow: hidden;
  display: flex;
  align-items: center;
}

.select-label {
  color: gray;
}

.header-center {
  flex-grow: 1;
  flex-shrink: 0;
  text-align: center;
}

.header-left,
.header-right {
  flex-basis: 150px;
  flex-shrink: 0;
}

.header-right {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  text-align: right;
  padding-right: 16px;
}

/* 折叠触发器 */
.collapsed-icon {
  padding: 0 10px 0 24px;
  cursor: pointer;
  transition: color 0.3s;
}

.collapsed-icon:hover {
  color: #1890ff;
}

.roll-msg-container {
  margin: 0 auto;
  overflow: hidden;
  width: 300px;
}

.roll-msg {
  display: inline-block;
  width: 300px;
  white-space: nowrap;
  animation: rollText 8s linear infinite normal;
  color: red;
  /* opacity: 0.5; */
}

@keyframes rollText {
  0% {
    transform: translateX(100%);
  }

  100% {
    transform: translateX(-100%);
  }
}
</style>
