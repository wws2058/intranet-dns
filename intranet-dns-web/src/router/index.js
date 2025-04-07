import { createRouter, createWebHashHistory } from "vue-router";
import { message } from "ant-design-vue";
import { layoutStore } from "@/store/layout";

import DnsRecords from "@/views/DnsRecords.vue";
import DnsProbes from "@/views/DnsProbes.vue";
import SystemAudit from "@/views/SystemAudit.vue";
import SystemUsers from "@/views/SystemUsers.vue";
import SystemRoles from "@/views/SystemRoles.vue";
import UserLogin from "@/views/UserLogin.vue";
import DnsQuery from "@/views/DnsQuery.vue";
import { localStoreUserDataKey, isTokenValid } from "@/apis";
import SystemApis from "@/views/SystemApis.vue";

// name string:path string
const routesMap = {
  用户登录: "/login",
  域名查询: "/dns_query",
  解析管理: "/dns_records",
  域名拨测: "/dns_probes",
  审计日志: "/system_audit",
  系统用户: "/system_users",
  系统角色: "/system_roles",
  接口管理: "/system_apis",
};

const routes = [
  {
    name: "Home",
    path: "/",
    meta: {
      title: "首页",
    },
    redirect: { name: "DnsQuery" },
  },
  {
    name: "UserLogin",
    path: routesMap["用户登录"],
    component: UserLogin,
    meta: {
      title: "用户登录",
      fullScreen: true, // 用户登录全屏
    },
  },
  {
    name: "DnsQuery",
    path: routesMap["域名查询"],
    component: DnsQuery,
    meta: {
      title: "域名查询",
    },
  },
  {
    name: "DnsRecords",
    path: routesMap["解析管理"],
    component: DnsRecords,
    meta: {
      title: "解析管理",
    },
  },
  {
    name: "DnsProbes",
    path: routesMap["域名拨测"],
    component: DnsProbes,
    meta: {
      title: "域名拨测",
    },
  },
  {
    name: "SystemAudit",
    path: routesMap["审计日志"],
    component: SystemAudit,
    meta: {
      title: "审计日志",
    },
  },
  {
    name: "SystemUsers",
    path: routesMap["系统用户"],
    component: SystemUsers,
    meta: {
      title: "域名查询",
    },
  },
  {
    name: "SystemRoles",
    path: routesMap["系统角色"],
    component: SystemRoles,
    meta: {
      title: "系统角色",
    },
  },
  {
    name: "SystemApis",
    path: routesMap["接口管理"],
    component: SystemApis,
    meta: {
      title: "接口管理",
    },
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

// 前置路由守卫, 使用Vue4-router写法
router.beforeEach((to) => {
  const userdataObj = JSON.parse(localStorage.getItem(localStoreUserDataKey));
  if (to.name === "UserLogin") {
    if (isTokenValid(userdataObj?.jwt_token)) {
      return { name: "DnsQuery" };
    }
    return true;
  }

  if (!isTokenValid(userdataObj?.jwt_token)) {
    message.warn("token不存在或过期, 请先登录");
    return { name: "UserLogin" };
  }
  return true;
});

// 后置路由守卫, 设置页面标题和key
router.afterEach((to) => {
  if (to.meta.title) {
    document.title = to.meta.title;
  }
  const store = layoutStore();
  store.selectedKeys = [store.getSelectedKey(to.path)];
});

export default router;
export { routesMap };
