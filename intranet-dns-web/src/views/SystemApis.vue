<template>
  <!-- 搜索条件 -->
  <a-space class="api-search">
    <a-input v-model:value="apiPath" placeholder="api路径" allowClear></a-input>

    <a-select v-model:value="isActive" :options="boolOptions" placeholder="接口状态" allowClear>
    </a-select>

    <a-select v-model:value="apiMethod" :options="apiMethodOptions" placeholder="接口方法" allowClear></a-select>

    <a-button type="primary" @click="handleSearch">搜索</a-button>
    <a-button @click="handleReset">重置</a-button>
  </a-space>

  <!-- 表格 -->
  <a-table :loading="pageloading" :columns="columns" :dataSource="dataSource" :pagination="pagination"
    @change="handlePageChange">
    <!-- 特殊列 -->
    <template #bodyCell="{ column, text, record }">
      <template v-if="column.key === 'operations'">
        <template v-if="record.isEditing && record.id === bindSelectObj.id">
          <a-space>
            <a-button @click="handleCancel(record)" size="small">取消</a-button>
            <a-button type="primary" @click="handleConfirm(record)" size="small" :loading="buttonLoading">确认</a-button>
          </a-space>
        </template>
        <template v-else>
          <a @click="handleEdit(record)">编辑</a>
        </template>
      </template>

      <template v-if="['audit', 'active'].includes(column.key)">
        <template v-if="record.isEditing && record.id === bindSelectObj.id">
          <template v-if="column.key === 'active'">
            <a-select :options="boolOptions" size="small" v-model:value="bindSelectObj.active">
            </a-select>
          </template>
          <template v-if="column.key === 'audit'">
            <a-select :options="boolOptions" size="small" v-model:value="bindSelectObj.audit">
            </a-select>
          </template>
        </template>
        <template v-else>
          <template v-if="text">
            <a-tag color="green">启用</a-tag>
          </template>
          <template v-else>
            <a-tag color="red">禁用</a-tag>
          </template>
        </template>
      </template>

    </template>
  </a-table>
</template>


<script setup>
import { message } from 'ant-design-vue';
import { reactive, ref } from 'vue';
import request from '@/apis/request';

// 搜索栏
const apiMethod = ref(null);
const apiPath = ref(null);
const isActive = ref(null);

const boolOptions = reactive([
  { value: 'true', label: '启用' },
  { value: 'false', label: '禁用' },
]);

const apiMethodOptions = reactive([
  { value: 'POST', label: 'POST' },
  { value: 'PUT', label: 'PUT' },
  { value: 'GET', label: 'GET' },
  { value: 'DELETE', label: 'DELETE' },
]);


function handleSearch() {
  const params = { page: pagination.current, page_size: pagination.pageSize };
  if (apiMethod.value !== null && apiMethod.value !== undefined) {
    params["method"] = apiMethod.value;
  }
  if (isActive.value !== null && isActive.value !== undefined) {
    params["active"] = isActive.value === "true" ? true : false;
  }
  if (apiPath.value !== null && apiPath.value !== '') {
    params["path"] = apiPath.value;
  }
  getApis(params);
}

function handleReset() {
  apiMethod.value = null;
  apiPath.value = null;
  isActive.value = null;
  pagination.current = 1;
  pagination.pageSize = 10;
  getApis({ page: pagination.current, page_size: pagination.pageSize });
}


// 定义列名
const columns = [
  {
    title: "ID",
    dataIndex: 'id'
  },
  {
    title: "方法",
    dataIndex: 'method'
  },
  {
    title: "路径",
    dataIndex: 'path'
  },
  {
    title: "接口描述",
    dataIndex: 'description'
  },
  {
    title: "接口状态",
    dataIndex: 'active',
    key: 'active'
  },
  {
    title: "日志审计",
    dataIndex: 'audit',
    key: 'audit'
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
    title: '操作',
    key: 'operations'
  },
];

// 定义表格数据
const dataSource = ref([]);

// 加载项
const pageloading = ref(false);

// 定义分页器
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

// 调用接口, 获取数据更新pagination
const getApis = async (params) => {
  pageloading.value = true;
  try {
    const response = await request.get("/api/v1/apis", params);
    const data = response?.data?.data;
    const pages = response?.data?.pages;
    if (data !== undefined) {
      dataSource.value = data;
    }
    if (pages !== undefined) {
      pagination.total = pages.total;
    }
  } catch (error) {
    console.log(error);
  } finally {
    pageloading.value = false;
  }
};

// page回调
function handlePageChange(pages) {
  pagination.current = pages.current;
  pagination.pageSize = pages.pageSize;
  handleSearch();
};

// 更新判断渲染
const buttonLoading = ref(false);
const bindSelectObj = reactive({
  id: 0,
  active: '',
  audit: '',
});

function handleEdit(record) {
  bindSelectObj.active = record.active === true ? "true" : "false";
  bindSelectObj.audit = record.audit === true ? "true" : "false";
  bindSelectObj.id = record.id;
  record.isEditing = true;
}

function handleCancel(record) {
  delete record.isEditing;
  bindSelectObj.id = 0;
  bindSelectObj.active = '';
  bindSelectObj.audit = '';
}

const handleConfirm = async (record) => {
  buttonLoading.value = true;
  buttonLoading.value = false;
  const params = {
    id: bindSelectObj.id,
    active: bindSelectObj.active === "true" ? true : false,
    audit: bindSelectObj.audit === "true" ? true : false,
  };
  await request.put("/api/v1/apis", params);
  handleCancel(record);
  await getApis({ page: pagination.current, page_size: pagination.pageSize });
  message.success("更新成功");
  buttonLoading.value = false;
};

// 加载页面时获取初始数据
getApis({ page: pagination.current, page_size: pagination.pageSize });
</script>

<style scoped>
.api-search {
  margin-bottom: 10px;
}
</style>
