import { defineStore } from "pinia";
import { ref, reactive, h, computed, toRaw } from "vue";
import {
  SearchOutlined,
  FileSearchOutlined,
  AuditOutlined,
  AppstoreOutlined,
  ClockCircleOutlined,
  TeamOutlined,
  UserSwitchOutlined,
  ApiOutlined,
  GlobalOutlined,
  AreaChartOutlined,
} from "@ant-design/icons-vue";
import { routesMap } from "@/router";

export const layoutStore = defineStore("layoutStore", () => {
  // 默认选择导航项
  const selectedKeys = ref(["1"]);
  const collapsed = ref(false);

  const siderNavs = reactive([
    {
      key: "域名查询",
      icon: () => h(SearchOutlined),
      label: "域名查询",
      title: "域名查询",
      routePath: routesMap["域名查询"],
    },
    {
      key: "解析记录",
      icon: () => h(FileSearchOutlined),
      label: "解析记录",
      title: "解析记录",
      routePath: routesMap["解析记录"],
    },
    {
      key: "区域管理",
      icon: () => h(GlobalOutlined),
      label: "区域管理",
      title: "区域管理",
      routePath: routesMap["区域管理"],
    },
    {
      key: "域名拨测",
      icon: () => h(ClockCircleOutlined),
      label: "域名拨测",
      title: "域名拨测",
      routePath: routesMap["域名拨测"],
    },
    {
      key: "日志图表",
      icon: () => h(AreaChartOutlined),
      label: "日志图表",
      title: "日志图表",
      routePath: routesMap["日志图表"],
    },
    {
      key: "系统管理",
      icon: () => h(AppstoreOutlined),
      label: "系统管理",
      title: "系统管理",
      children: [
        {
          key: "审计日志",
          icon: () => h(AuditOutlined),
          label: "审计日志",
          title: "审计日志",
          routePath: routesMap["审计日志"],
        },
        {
          key: "系统用户",
          icon: () => h(TeamOutlined),
          label: "系统用户",
          title: "系统用户",
          routePath: routesMap["系统用户"],
        },
        {
          key: "系统角色",
          icon: () => h(UserSwitchOutlined),
          label: "系统角色",
          title: "系统角色",
          routePath: routesMap["系统角色"],
        },
        {
          key: "接口管理",
          icon: () => h(ApiOutlined),
          label: "接口管理",
          title: "接口管理",
          routePath: routesMap["接口管理"],
        },
      ],
    },
  ]);

  const openkeys = computed(() => {
    const sub1Keys = ["审计日志", "系统用户", "系统角色", "接口管理"];
    const selectedKeysArray = toRaw(selectedKeys.value);
    const contained = sub1Keys.some((key) => selectedKeysArray.includes(key));
    const opened = !collapsed.value && contained;
    return opened ? ["系统管理"] : [];
  });

  const selecttedLabel = computed(() => {
    function findLabel(key, items) {
      for (const item of items) {
        if (item.key === key) {
          return item.label;
        }
        if (item.children) {
          const result = findLabel(key, item.children);
          if (result) {
            return result;
          }
        }
      }
      return null;
    }
    return findLabel(selectedKeys.value[0], siderNavs);
  });

  function switchCollapsed() {
    collapsed.value = !collapsed.value;
  }

  function getSelectedKey(routePath) {
    function findKey(routePath, items) {
      for (const item of items) {
        if (item.routePath === routePath) {
          return item.key;
        }
        if (item.children) {
          const result = findKey(routePath, item.children);
          if (result) {
            return result;
          }
        }
      }
      return null;
    }
    return findKey(routePath, siderNavs);
  }

  return {
    siderNavs,
    selectedKeys,
    selecttedLabel,
    openkeys,
    collapsed,
    switchCollapsed,
    getSelectedKey,
  };
});
