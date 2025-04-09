<template>
  <!-- 新增角色api -->
  <a-button type="primary" class="add-button" @click="showAddSysRoleFormModal">新增角色</a-button>

  <a-modal v-model:open="openAddRoleModal" title="新增系统角色" :footer="null">
    <a-form ref="addSysRoleForm" name="add-role-form" :rules="addSysRoleFormRules" labelAlign="left"
      :model="addSysRoleFormObj" :label-col="{ span: 5 }" @finish="addSysRoleRequest" :scrollToFirstError="true">
      <a-form-item label="角色中文名" name="name_cn" class="first-input">
        <a-input v-model:value="addSysRoleFormObj.name_cn" placeholder="请输入角色中文名" allowClear>角色中文名</a-input>
      </a-form-item>

      <a-form-item label=" 角色英文名" name="name">
        <a-input v-model:value="addSysRoleFormObj.name" placeholder="请输入角色英文名" allowClear>角色英文名</a-input>
      </a-form-item>

      <a-form-item label="api访问权限" name="apis">
        <a-select v-model:value="addSysRoleFormObj.apis" mode="multiple" :autoClearSearchValue="false" allowClear
          :labelInValue="true" :options="allApis.map(api => ({ value: api.detail, key: api.id }))"
          placeholder="绑定访问的api" :maxTagCount="5"></a-select>
      </a-form-item>

      <a-form-item>
        <a-space class="add-role-button">
          <a-button htmlType="reset" @click="resetAddSysRoleForm">重置</a-button>
          <a-button type="primary" htmlType="sumit" :disabled="addSysRoleButtonDisable"
            :loading="addSysButtonRoleLoading">确认</a-button>
        </a-space>
      </a-form-item>
    </a-form>
  </a-modal>

  <!-- 查询角色api -->
  <a-table :loading="loading" :columns="columns" :dataSource="dataSource" :pagination="pagination"
    @change="handlePageChange">
    <!-- 特殊列 -->
    <template #bodyCell="{ record, column }">
      <template v-if="column.key === 'operations'">
        <span>
          <a-popconfirm title="删除系统角色?" @confirm="deleteOneRole(record)" okText="确认" cancelText="取消">
            <a href="#">删除</a>
          </a-popconfirm>
          <a-divider type="vertical" />
          <a @click="showUpdateSysRoleFormModal(record)">更新</a>
        </span>
      </template>
      <template v-else-if="column.key === 'apis'">
        <a @click="showRoleApis(record)">查看详情</a>
      </template>
    </template>
  </a-table>

  <!-- 角色详情 -->
  <a-modal v-model:open="openSysRoleDetailModal" :footer="null" :afterClose="closeSysRoleDetailsModal" :closable="false"
    width="400px">
    <a-input v-model:value.lazy="searchPath" placeholder="搜索api path" allowClear class="detail-search" />
    <div class="api-details">
      <ul class="api-details-list">
        <li v-for="api in filterSysRoleApiDetails" :key="api.id">{{ api.method }} - {{ api.path }} - {{ api.description
        }}
        </li>
      </ul>
    </div>

  </a-modal>

  <!-- 更新角色 -->
  <a-modal v-model:open="openUpdateRoleModal" title="更新系统角色" :footer="null">
    <a-form ref="updateSysRoleForm" name="update-role-form" :rules="addSysRoleFormRules" labelAlign="left"
      :model="updateSysRoleObj" :label-col="{ span: 5 }" @finish="updateSysRoleRequest" :scrollToFirstError="true">
      <a-form-item label="角色中文名" name="name_cn" class="first-input">
        <a-input v-model:value="updateSysRoleObj.name_cn" placeholder="请输入角色中文名" allowClear>角色中文名</a-input>
      </a-form-item>

      <a-form-item label=" 角色英文名" name="name">
        <a-input v-model:value="updateSysRoleObj.name" placeholder="请输入角色英文名" allowClear>角色英文名</a-input>
      </a-form-item>

      <a-form-item label="api访问权限" name="apis">
        <a-select v-model:value="updateSysRoleObj.apis" mode="multiple" :autoClearSearchValue="false" allowClear
          :labelInValue="true" :options="allApis.map(api => ({ value: api.detail, key: api.id }))"
          placeholder="绑定访问的api" :maxTagCount="5"></a-select>
      </a-form-item>

      <a-form-item>
        <a-space class="add-role-button">
          <a-button htmlType="reset" @click="resetUpdateSysRoleForm">重置</a-button>
          <a-button type="primary" htmlType="sumit" :disabled="updateSysRoleButtonDisable"
            :loading="addSysButtonRoleLoading">确认</a-button>
        </a-space>
      </a-form-item>
    </a-form>
  </a-modal>

</template>
<script setup>
import { deleteSysRole, getRoleApis, getSysRoles, getAllSysApis, newSysRoles, updateSysRole } from '@/apis';
import { message } from 'ant-design-vue';
import { reactive, ref, computed, toRaw } from 'vue';

// 定义列名
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

// 定义分页器
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  pageSizeOptions: ['10', '20', '30'],
  locale: {
    items_per_page: '条/页'
  },
  showTotal: (total) => `共 ${total} 条`
});

// 定义表格数据
const dataSource = ref([]);

// 加载项
const loading = ref(false);

// 调用接口
const getRoles = async () => {
  loading.value = true;
  try {
    const response = await getSysRoles({ page: pagination.current, page_size: pagination.pageSize });
    const data = response?.data;
    const pages = response?.pages;
    if (data !== undefined) {
      dataSource.value = data;
    }
    if (pages !== undefined) {
      pagination.total = pages.total;
    }
  } catch (error) {
    console.log(error);
  } finally {
    loading.value = false;
  }
};

// page回调
const handlePageChange = async (pages) => {
  pagination.current = pages.current;
  pagination.pageSize = pages.pageSize;
  await getRoles();
};

// 加载页面时获取初始数据
getRoles();

// 删除角色
const deleteOneRole = async (record) => {
  await deleteSysRole(record.id);
  message.success(`${record.name_cn} 已删除`);
  getRoles();
};


// 新增角色
const openAddRoleModal = ref(false);
const showAddSysRoleFormModal = () => {
  openAddRoleModal.value = true;
  getAllApis();
};
const addSysRoleFormObj = reactive({
  name: '',
  name_cn: '',
  apis: [],
});
const addSysRoleFormRules = {
  name: [{
    required: true,
    message: '请输入角色英文名',
  }],
  name_cn: [{
    required: true,
    message: '请输入角色中文名',
  }],
  apis: [{
    required: true,
    message: '请绑定具有权限的接口',
  }]
};
const addSysRoleButtonDisable = computed(() => {
  return !(addSysRoleFormObj.name && addSysRoleFormObj.name_cn && addSysRoleFormObj.apis.length > 0);
});
const addSysButtonRoleLoading = ref(false);
const addSysRoleRequest = async () => {
  try {
    const request = {
      name: addSysRoleFormObj.name,
      name_cn: addSysRoleFormObj.name_cn,
      api_ids: addSysRoleFormObj.apis.map(api => {
        return api.key;
      })
    };
    addSysButtonRoleLoading.value = true;
    await newSysRoles(request);
    addSysButtonRoleLoading.value = false;
    openAddRoleModal.value = false;
    getRoles();
    message.success(`${request.name_cn} 添加成功`);
    resetAddSysRoleForm();
  } catch (error) {
    console.log(error);
    addSysButtonRoleLoading.value = false;
  }

};
const addSysRoleForm = ref();
function resetAddSysRoleForm() {
  addSysRoleForm.value.resetFields();
}

// {detail,id}
const allApis = ref([]);
// 获取所有的api
async function getAllApis() {
  allApis.value = await getAllSysApis();
  allApis.value = allApis.value.map((api) => {
    return {
      detail: `${api.method} - ${api.path} - ${api.description}`,
      id: api.id
    };
  });
  allApis.value.sort((a, b) => {
    if (a.detail < b.detail) {
      return -1;
    }
    if (a.detail > b.detail) {
      return 1;
    }
    return 0;
  });
}

const openSysRoleDetailModal = ref(false);
const sysRoleApiDetails = ref([]);
const searchPath = ref('');
const filterSysRoleApiDetails = computed(() => {
  const filterApis = sysRoleApiDetails.value.filter((api) => {
    return api.path.includes(searchPath.value);
  });
  return filterApis;
});


// 查看角色详情
async function showRoleApis(record) {
  openSysRoleDetailModal.value = true;
  const response = await getRoleApis(record.id);
  const accessibleApis = response.data.accessible_apis;
  accessibleApis.sort((a, b) => {
    if (a.method < b.method) {
      return -1;
    }
    if (a.method > b.method) {
      return 1;
    }
    return 0;
  });
  sysRoleApiDetails.value = accessibleApis;
}
const closeSysRoleDetailsModal = () => {
  sysRoleApiDetails.value = [];
};

// 更新角色
const openUpdateRoleModal = ref(false);
const updateSysRoleForm = ref();
const updateSysRoleObj = reactive({
  id: 0,
  name: '',
  name_cn: '',
  apis: [],
});
let originStoreApis = [];
const showUpdateSysRoleFormModal = async (record) => {
  updateSysRoleObj.name = record.name;
  updateSysRoleObj.name_cn = record.name_cn;
  updateSysRoleObj.id = record.id;
  await getAllApis();
  allApis.value.forEach(api => {
    if (record.name === "super_admin") {
      updateSysRoleObj.apis.push(
        {
          value: api.detail,
          key: api.id
        }
      );
      return;
    }
    if (record.api_ids.includes(api.id) && !updateSysRoleObj.apis.some(existApi => existApi.key === api.id)) {
      updateSysRoleObj.apis.push(
        {
          value: api.detail,
          key: api.id
        }
      );
    }
  });
  openUpdateRoleModal.value = true;
  originStoreApis = updateSysRoleObj.apis;
  console.log("showUpdateSysRoleFormModal", toRaw(updateSysRoleObj.apis));
};
function resetUpdateSysRoleForm() {
  updateSysRoleForm.value.resetFields();
  updateSysRoleObj.apis = originStoreApis;
}
const updateSysRoleButtonDisable = computed(() => {
  return !(updateSysRoleObj.name && updateSysRoleObj.name_cn && updateSysRoleObj.apis.length > 0);
});
const updateSysRoleRequest = async () => {
  try {
    const request = {
      id: updateSysRoleObj.id,
      name: updateSysRoleObj.name,
      name_cn: updateSysRoleObj.name_cn,
      api_ids: updateSysRoleObj.apis.map(api => {
        return api.key;
      })
    };
    addSysButtonRoleLoading.value = true;
    await updateSysRole(request);
    addSysButtonRoleLoading.value = false;
    openUpdateRoleModal.value = false;
    getRoles();
    message.success(`${request.name_cn} 更新成功`);
  } catch (error) {
    console.log(error);
    addSysButtonRoleLoading.value = false;
  };
};
</script>

<style scope>
.add-button {
  float: right;
  margin-bottom: 10px;
}

.first-input {
  margin-top: 24px;
}

.add-role-button {
  display: flex;
  justify-content: flex-end;
}

.detail-search {
  margin-bottom: 10px;
}

.api-details {
  max-height: 50vh;
  overflow: auto;
  white-space: nowrap;
}

.api-details-list {
  list-style-type: none;
  padding: 0;
  margin: 0;
}
</style>
