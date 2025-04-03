<template>
  <a-button type="primary" class="add-button">新增角色</a-button>
  <a-table :columns="columns" :dataSource="data">
    <!-- 自定义列 -->
    <template #bodyCell="{ column }">
      <template v-if="column.key === 'operations'">
        <span>
          <a>删除</a>
          <a-divider type="vertical" />
          <a>更新</a>
        </span>
      </template>
      <template v-else-if="column.key === 'apis'">
        <a>查看详情</a>
      </template>
    </template>
  </a-table>
</template>
<script setup>

const columns = [
  {
    title: "ID",
    dataIndex: 'id'
  },
  {
    title: "角色名(英)",
    dataIndex: 'name'
  },
  {
    title: "角色名(中)",
    dataIndex: 'name_cn'
  },
  {
    title: "创建时间",
    dataIndex: 'created_at'
  },
  {
    title: "更新时间",
    dataIndex: 'updated_at'
  },
  {
    title: "api详情",
    dataIndex: 'api_ids',
    key: "apis"
  },
  {
    title: '操作',
    key: 'operations'
  },
];
const data = [
  {
    "id": 18,
    "created_at": "2025-04-02 14:04:15",
    "updated_at": "2025-04-02 14:04:15",
    "name": "super_admin",
    "name_cn": "超级管理员"
  },
  {
    "id": 2,
    "created_at": "2024-12-09 15:55:15",
    "updated_at": "2024-12-09 15:55:15",
    "name": "dev",
    "name_cn": "开发",
    "api_ids": [
      2,
      3
    ]
  }
];
</script>

<style scope>
.add-button {
  float: right;
  margin-bottom: 10px;
}
</style>
