<template>
  <a-button type="primary" class="zone-add-button" disabled>新增zone</a-button>
  <a-table :columns="columns" :dataSource="dataSource" :pagination="pagination" :loading="tableLoading"
    @change="handlePageChange">

    <template #bodyCell="{ column }">
      <template v-if="column.key === 'operations'">
        <a-typography-text disabled>
          删除<a-divider type="vertical" />更新
        </a-typography-text>
      </template>
    </template>

  </a-table>
</template>

<script setup>

import request from '@/apis/request';
import { ref, reactive } from 'vue';

const columns = ref([
  {
    title: 'zone',
    dataIndex: 'zone',
    key: 'zone',
  },
  {
    title: 'DNS服务器地址',
    dataIndex: 'ns_address',
    key: 'ns_address',
  },
  {
    title: 'TSIG',
    dataIndex: 'tsig_name',
    key: 'tsig_name',
  },
  {
    title: '描述信息',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: '创建者',
    dataIndex: 'creator',
    key: 'creator',
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
  },
  {
    title: '更新时间',
    dataIndex: 'updated_at',
    key: 'updated_at',
  },
  {
    title: '操作',
    key: 'operations',
  },
]);
const dataSource = ref([]);
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  pageSizeOptions: ['10', '20', '50'],
  locale: {
    items_per_page: '条/页'
  },
  showTotal: (total) => `共 ${total} 条`
});
const tableLoading = ref(false);

async function handlePageChange(pages) {
  pagination.current = pages.current;
  pagination.pageSize = pages.pageSize;
  await queryZones();
}

async function queryZones() {
  const requestObj = {
    page: pagination.current,
    page_size: pagination.pageSize
  };
  try {
    tableLoading.value = true;
    const rsp = await request.get('/api/v1/dns/zones', requestObj);
    dataSource.value = rsp.data.data;
    tableLoading.value = false;
  } catch (error) {
    console.log(error);
    tableLoading.value = false;
  }
}

queryZones();
// TODO: zone del, zone add, zone update
</script>

<style scoped>
.zone-add-button {
  float: right;
  margin-bottom: 10px;
}
</style>
