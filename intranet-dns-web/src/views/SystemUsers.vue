<template>
  <!-- 搜索条件 -->
  <div class="user-operations">
    <a-space class="search">
      <a-input v-model:value="userNameCn" placeholder="用户中文名" allowClear></a-input>

      <a-select v-model:value="isActive" :options="boolOptions" placeholder="用户状态" allowClear>
      </a-select>

      <a-select v-model:value="sysRoles" :labelInValue="true"
        :options="Array.from(rolesMap, ([key, value]) => ({ key, value }))" placeholder="系统角色" allowClear></a-select>

      <a-button type="primary" @click="handleSearch" :loading="searchLoading">搜索</a-button>
      <a-button @click="handleReset">重置</a-button>
    </a-space>
    <a-button type="primary" @click="showAddUserForm">新增用户</a-button>
  </div>


  <a-table :columns="columns" :data-source="data" :pagination="pagination" @change="handlePageChange">

    <template #bodyCell="{ column, text, record }">
      <template v-if="column.key === 'active'">
        <template v-if="record.isEditing && record.id === editRecord.id">
          <a-select v-model:value="editRecord.active" size="small" :options="boolOptions">
          </a-select>
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

      <template v-if="column.key === 'role_ids'">
        <template v-if="record.isEditing && record.id === editRecord.id">
          <a-select v-model:value="editRecord.roles" size="small" :labelInValue="true" mode="multiple" allowClear
            :autoClearSearchValue="false" :maxTagCount="3" placeholder="绑定系统角色" :style="{ minWidth: '100px' }"
            :options="Array.from(rolesMap, ([key, value]) => ({ key, value }))">
          </a-select>
        </template>
        <template v-else>
          <a-tag color="blue" v-for="id in record.role_ids" :key="id">
            {{ rolesMap.get(id) }}
          </a-tag>
        </template>
      </template>

      <template v-if="column.key === 'email' && record.isEditing && record.id === editRecord.id">
        <a-input v-model:value="editRecord.email" size="small" allowClear></a-input>
      </template>

      <template v-if="column.key === 'operations'">
        <template v-if="record.isEditing && record.id === editRecord.id">
          <a-space>
            <a-button size="small" @click="cancelEdit">取消</a-button>
            <a-button type="primary" size="small" @click="updateUser(record)" :loading="updateLoading">确认</a-button>
          </a-space>
        </template>
        <template v-else>
          <a-popconfirm title="删除系统用户?" @confirm="deleteUser(record)" okText="确认" cancelText="取消"
            :loading="deleteLoading">
            <a href="#">删除</a>
          </a-popconfirm>
          <a-divider type="vertical" />

          <a @click="handleEdit(record)">编辑</a>
        </template>
      </template>

      <!-- 增加编辑属性 -->

    </template>

  </a-table>

  <a-modal v-model:open="openAddUserModal" title="新增用户" :footer="null">
    <a-form ref="addUserForm" name="add-user" labelAlign="left" :model="addUserFormObj" :rules="addUserFormRules"
      :label-col="{ span: 5 }" @finish="addUser" :scrollToFirstError="true">

      <a-form-item label="用户中文名" class="first-input" name="name_cn">
        <a-input v-model:value="addUserFormObj.name_cn" placeholder="请输入角色中文名" allowClear>角色中文名</a-input>
      </a-form-item>

      <a-form-item label="用户英文名" name="name">
        <a-input v-model:value="addUserFormObj.name" placeholder="请输入角色英文名" allowClear>角色英文名</a-input>
      </a-form-item>

      <a-form-item label="用户密码" name="password">
        <a-input-password v-model:value="addUserFormObj.password" placeholder="请输入至少8位数用户密码"
          allowClear>角色英文名</a-input-password>
      </a-form-item>

      <a-form-item label="用户邮箱" name="email">
        <a-input v-model:value="addUserFormObj.email" placeholder="请输入邮箱地址" allowClear>角色英文名</a-input>
      </a-form-item>

      <a-form-item label="绑定系统角色" name="roles">
        <a-select v-model:value="addUserFormObj.roles" mode="multiple" :autoClearSearchValue="false" allowClear
          placeholder="用户关联的系统角色" :maxTagCount="5" :options="Array.from(rolesMap, ([key, value]) => ({ key, value }))"
          :labelInValue="true"></a-select>
      </a-form-item>

      <a-form-item>
        <a-space class="add-user-button">
          <a-button @click="cancelAddUser">取消</a-button>
          <a-button type="primary" htmlType="submit" :disabled="addUserButton" :loading="addUserLoading">确认</a-button>
        </a-space>
      </a-form-item>
    </a-form>
  </a-modal>


</template>
<script setup>
import { getAllSysRoles } from '@/apis';
import request from '@/apis/request';
import { error } from '@/hooks/useTipsMessage';
import { message } from 'ant-design-vue';
import { ref, reactive, computed } from 'vue';

// 搜索栏
const userNameCn = ref(null);
const isActive = ref(null);
const sysRoles = ref(null);
const boolOptions = reactive([
  { value: 'true', label: '启用' },
  { value: 'false', label: '禁用' },
]);


const deleteLoading = ref(false);
async function deleteUser(record) {
  console.log(record);
  deleteLoading.value = true;
  const url = `/api/v1/users/${record.id}`;
  await request.delete(url);
  message.success(`${record.name_cn} 已删除`);
  deleteLoading.value = false;
  listUsers();
}

// 角色id-name map
const allRoles = ref([]);
async function getRoles() {
  const roles = await getAllSysRoles();
  allRoles.value = roles;
}
// map, role.id: role.name_cn
// Array.from(rolesMap, ([key, value]) => ({ key, value })), 前面解构赋值获取到map的k|v, 返回的对象使用小括号包括起来防止解析为函数体
const rolesMap = computed(() => {
  const map = new Map();
  allRoles.value.forEach(role => {
    map.set(role.id, role.name_cn);
  });
  return map;
});
// 新增角色
const openAddUserModal = ref(false);
const addUserLoading = ref(false);
const addUserFormObj = reactive({
  name: '',
  name_cn: '',
  password: '',
  email: '',
  roles: [],
});
const addUserButton = computed(() => {
  return !(addUserFormObj.name && addUserFormObj.name_cn && addUserFormObj.password && addUserFormObj.email && addUserFormObj.roles.length > 0);
});
const addUserFormRules = {
  name: [{
    required: true,
    message: '请输入用户英文名',
    whitespace: true
  }],
  name_cn: [{
    required: true,
    message: '请输入用户中文名',
    whitespace: true
  }],
  password: [{
    required: true,
    message: '请输入用户密码',
    whitespace: true
  },
  {
    min: 8,
    message: " 密码长度不能小于8位"
  }], email: [{ required: true, message: '请输入用户邮箱', whitespace: true }], roles: [{
    required:
      true, message: '请绑定系统角色',
  }],
};
function cancelAddUser() {
  openAddUserModal.value = false;
}
async function addUser() {
  try {
    addUserLoading.value = true;
    const ids = addUserFormObj.roles.map((role) => role.key);
    const requestObj = { ...addUserFormObj };
    requestObj.role_ids = ids;
    delete requestObj.roles;
    await request.post('/api/v1/users', requestObj);
    addUserLoading.value = false;
    message.success(`${addUserFormObj.name_cn} 新增成功`);
    openAddUserModal.value = false;
  } catch {
    addUserLoading.value = false;
    console.log(error);
  }
  listUsers();
}
function showAddUserForm() {
  openAddUserModal.value = true;
}


const columns = [
  {
    title: '中文名',
    dataIndex: 'name_cn',
    key: 'name_cn',
  },
  {
    title: '英文名',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '邮箱',
    dataIndex: 'email',
    key: 'email',
  },
  {
    title: '状态',
    dataIndex: 'active',
    key: 'active',
  },
  {
    title: '关联系统角色',
    dataIndex: 'role_ids',
    key: 'role_ids',
  },
  {
    title: '登录次数',
    dataIndex: 'login_times',
    key: 'login_times',
    customRender: ({ record }) => {
      return record.login_times ? record.login_times : '*';
    }
  },
  {
    title: '最近登录',
    dataIndex: 'last_login',
    key: 'last_login',
    customRender: ({ record }) => {
      return record.last_login ? record.last_login : '*';
    }
  },
  {
    title: '操作',
    dataIndex: 'operations',
    key: 'operations',
  },
];
const data = ref([]);
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

async function listUsers() {
  const requestObj = {
    page: pagination.current,
    page_size: pagination.pageSize
  };
  if (userNameCn.value) {
    requestObj.name_cn = userNameCn.value;
  }
  if (isActive.value) {
    requestObj.active = isActive.value === "true" ? true : false;
  }
  if (sysRoles.value) {
    requestObj.role_id = sysRoles.value.key;
  }
  try {
    const rsp = await request.get('/api/v1/users', requestObj);
    const pages = rsp.data.pages;
    const users = rsp.data.data;
    pagination.total = pages.total;
    data.value = users;
  } catch (error) {
    console.log(error);
  }
}

// 换页
const handlePageChange = async (pages) => {
  pagination.current = pages.current;
  pagination.pageSize = pages.pageSize;
  await listUsers();
};

const searchLoading = ref(false);
async function handleSearch() {
  searchLoading.value = true;
  await listUsers();
  searchLoading.value = false;
}

function handleReset() {
  userNameCn.value = null;
  isActive.value = null;
  sysRoles.value = null;
  pagination.current = 1;
  pagination.pageSize = 10;
  listUsers();
}

// 编辑框
const editRecord = reactive({
  id: null,
  email: null,
  active: null,
  roles: []
});

function cancelEdit() {
  editRecord.id = null;
  editRecord.email = null;
  editRecord.active = null;
  editRecord.roles = [];
}

function handleEdit(record) {
  record.isEditing = true;
  editRecord.id = record.id;
  editRecord.email = record.email;
  editRecord.active = record.active ? 'true' : 'false';
  editRecord.roles = record.role_ids.map((id) => {
    return {
      key: id,
      value: rolesMap.value.get(id)
    };
  });
}
const updateLoading = ref(false);
async function updateUser(record) {
  updateLoading.value = true;
  const requestObj = {
    id: editRecord.id,
    active: editRecord.active === 'true' ? true : false,
    email: editRecord.email,
    role_ids: editRecord.roles.map((role) => role.key)
  };

  await request.put("/api/v1/users", requestObj);
  updateLoading.value = false;
  delete record.isEditing;
  listUsers();
}

// 页面初始请求
getRoles();
listUsers();
</script>

<style scoped>
.search {
  margin-bottom: 20px;
}

.user-operations {
  display: flex;
  justify-content: space-between;
}

.first-input {
  margin-top: 24px;
}

.add-user-button {
  display: flex;
  justify-content: flex-end;
}
</style>