<template>
  <a-layout>
    <a-layout-content class="login-content">
      <a-card hoverable class="login-card">
        <template #title>
          <h3>欢迎登录DNS内网管控系统</h3>
        </template>
        <a-form class="login-form" name="login-form" labelAlign="left" :model="userLoginData" :label-col="{ span: 4 }"
          @finish="handleUserLogin">
          <a-form-item label="用户名" name="name" :rules="userNameRules">
            <a-input v-model:value="userLoginData.name" allowClear>
              <template #prefix>
                <UserOutlined></UserOutlined>
              </template>
            </a-input>
          </a-form-item>

          <a-form-item label="密码" name="password" :rules="userPasswordRules">
            <a-input-password v-model:value="userLoginData.password" allowClear>
              <template #prefix>
                <LockOutlined></LockOutlined>
              </template>
            </a-input-password>
          </a-form-item>

          <a-form-item>
            <a-button :disabled="disabled" type="primary" html-type="submit" class="login-button">
              登录
            </a-button>
          </a-form-item>
        </a-form>
      </a-card>
    </a-layout-content>
    <LayoutFooter></LayoutFooter>
  </a-layout>


</template>

<script setup>
import { LockOutlined, UserOutlined } from '@ant-design/icons-vue';
import { reactive, computed } from 'vue';
import LayoutFooter from '@/components/LayoutFooter.vue';
import { message } from 'ant-design-vue';
import { useRouter } from 'vue-router';
import { localStoreUserDataKey } from '@/apis';
import request from '@/apis/request';

const userLoginData = reactive({
  name: '',
  password: '',
});

const userNameRules = reactive(
  [{ required: true, message: '请输入用户名' }]
);

const userPasswordRules = reactive(
  [{ required: true, message: '请输入密码, 忘记密码联系管理员' }]
);

const disabled = computed(() => {
  return !(userLoginData.name && userLoginData.password);
});

const router = useRouter();

const handleUserLogin = async () => {
  try {
    const rsp = await request.post("/api/v1/users/login", userLoginData);
    localStorage.setItem(localStoreUserDataKey, JSON.stringify({
      name: userLoginData.name,
      jwt_token: rsp.data.data.jwt_token,
    }));
    message.success('登录成功');
    router.push({ name: "DnsQuery" });
  } catch (error) {
    message.error("登录失败");
    console.log(error);
  }
};
</script>

<style scoped>
.login-content {
  background-color: white;
  background-image: url(../assets/background.svg);
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  height: calc(100vh - 30px);
}

.login-card {
  animation: fadeIn 1s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }

  to {
    opacity: 1;
  }
}

.login-card h3 {
  text-align: center;
}

.login-form {
  width: 400px;
  min-width: 300px;
}

.login-button {
  width: 100%;
}
</style>
