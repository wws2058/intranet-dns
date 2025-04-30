<template>

  <!-- 搜索栏 -->

  <div class="dns-search">
    <a-space>
      <a-input v-model:value.trim="searchObj.record_name" placeholder="请输入要查询的域名" allowClear></a-input>
      <a-input v-model:value.trim="searchObj.record_content" placeholder="请输入要查询的解析记录" allowClear></a-input>
      <a-select v-model:value="searchObj.zone" :options="zoneOptions"></a-select>
      <a-select v-model:value="searchObj.record_type" :options="typeOptions"></a-select>
      <a-button type="primary" @click="handleSearchRecords" :loading="tableLoading">搜索</a-button>
      <a-button @click="resetSearchObj">重置</a-button>
    </a-space>
    <a-button type="primary" @click="changeOpenAddModal">新增记录</a-button>
  </div>


  <!-- 表格栏 -->
  <a-table :columns="columns" :dataSource="dataSource" :pagination="pagination" :loading="tableLoading"
    @change="handlePageChange">

    <template #bodyCell="{ column, record }">

      <template v-if="column.key === 'operations'">

        <template v-if="record.id === editRecordObj.id">
          <a-space>
            <a-button size="small" @click="handleCancelEdit">取消</a-button>
            <a-button type="primary" size="small" @click="handleUpdateRecord">确认</a-button>
          </a-space>
        </template>
        <template v-else>
          <template v-if="['SOA', 'NS'].includes(record.record_type)">
            <a-typography-text disabled>
              删除<a-divider type="vertical" />更新
            </a-typography-text>
          </template>
          <template v-else>
            <a-popconfirm title="删除dns记录?" @confirm="deleteRecord(record)" okText="确认" cancelText="取消" placement="left">
              <a href="#">删除</a>
            </a-popconfirm>
            <a-divider type="vertical" />
            <a @click="handleEdit(record)">编辑</a>
          </template>
        </template>

      </template>

      <template v-if="column.key === 'record_content' && record.id === editRecordObj.id">
        <a-input size="small" v-model:value="editRecordObj.record_content"></a-input>
      </template>

      <template v-if="column.key === 'record_ttl' && record.id === editRecordObj.id">
        <a-input size="small" v-model:value.number="editRecordObj.record_ttl"></a-input>
      </template>

    </template>

  </a-table>

  <!-- 交互栏 -->
  <a-modal v-model:open="showAddModal" title="新增dns解析记录" :footer="null">
    <a-form ref="formRef" :rules="addDnsRules" labelAlign="left" :model="addDnsObj" :label-col="{ span: 5 }"
      :scrollToFirstError="true" @finish="handletAddDns">

      <a-form-item label="新增域名" name="record_name" class="first-input">
        <a-input v-model:value.trim="addDnsObj.record_name" allowClear placeholder="my.test.com"></a-input>
      </a-form-item>

      <a-form-item label="域名TTL" name="record_ttl">
        <a-input v-model:value.number="addDnsObj.record_ttl" allowClear placeholder="60"></a-input>
      </a-form-item>

      <a-form-item label="域名所属域" name="zone">
        <a-select v-model:value="addDnsObj.zone" :options="zoneOptions" allowClear placeholder="test.com."></a-select>
      </a-form-item>

      <a-form-item label="解析类型" name="record_type">
        <a-select v-model:value="addDnsObj.record_type" :options="typeOptions" allowClear placeholder="A"></a-select>
      </a-form-item>

      <a-form-item label="解析内容" name="record_content">
        <a-input v-model:value.trim="addDnsObj.record_content" allowClear
          placeholder="单个1.1.1.1, 多个英文逗号分隔如1.1.1.1,2.2.2.2"></a-input>
      </a-form-item>

      <a-form-item>
        <a-space class="add-dns-button">
          <a-button @click="handleResetAddDns">重置</a-button>
          <a-button type="primary" htmlType="submit">确认</a-button>
        </a-space>
      </a-form-item>

    </a-form>
  </a-modal>


</template>

<script setup>
import { dnsTypes, intranetZones } from '@/apis/const';
import request from '@/apis/request';
import { message } from 'ant-design-vue';
import { computed, reactive, ref } from 'vue';

const searchObj = reactive({
  record_name: null,
  zone: null,
  record_type: null,
  record_content: null,
  page: 1,
  page_size: 10,
});
const resetSearchObj = () => {
  searchObj.record_name = null;
  searchObj.zone = null;
  searchObj.record_type = null;
  searchObj.record_content = null;
  searchObj.page = 1;
  searchObj.page_size = 10;

  pagination.current = 1;
  pagination.pageSize = 10;

  queryRecords(searchObj);
};
const handleSearchRecords = () => {
  queryRecords(searchObj);
};

const zoneOptions = computed(() => {
  const zOptions = intranetZones.map((zone) => ({
    label: zone,
    value: zone
  }));
  zOptions.unshift({
    label: 'zones',
    value: null,
    disabled: true
  });
  return zOptions;
});
const typeOptions = computed(() => {
  const tOptions = dnsTypes.map((type) => ({
    label: type,
    value: type
  }));
  tOptions.unshift({
    label: 'types',
    value: null,
    disabled: true
  });
  return tOptions;
});


const columns = ref([
  {
    title: '域名',
    dataIndex: 'record_name',
    key: 'record_name',
  },
  {
    title: '所属域',
    dataIndex: 'zone',
    key: 'zone',
  },
  {
    title: '记录类型',
    dataIndex: 'record_type',
    key: 'record_type',
  },
  {
    title: 'TTL',
    dataIndex: 'record_ttl',
    key: 'record_ttl',
  },
  {
    title: '解析记录',
    dataIndex: 'record_content',
    key: 'record_content',
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
    dataIndex: 'operations',
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

const handlePageChange = async (pages) => {
  pagination.current = pages.current;
  pagination.pageSize = pages.pageSize;
  searchObj.page = pagination.current;
  searchObj.page_size = pagination.pageSize;
  await queryRecords(searchObj);
};

async function queryRecords(searchObj) {
  try {
    tableLoading.value = true;
    const rsp = await request.get('/api/v1/dns/records', searchObj);
    dataSource.value = rsp.data.data;
    pagination.total = rsp.data.pages.total;
    tableLoading.value = false;
  } catch (error) {
    console.log(error);
    tableLoading.value = false;
  }
}
const deleteRecord = async (record) => {
  await request.delete("/api/v1/dns/records", {
    id: record.id
  });
  message.success("记录已删除");
  queryRecords(searchObj);
};


const editRecordObj = reactive({
  id: 0,
  record_name: null,
  record_ttl: null,
  record_content: null
});
const handleEdit = (record) => {
  editRecordObj.id = record.id;
  editRecordObj.record_name = record.record_name;
  editRecordObj.record_ttl = record.record_ttl;
  editRecordObj.record_content = record.record_content;
};
const handleCancelEdit = () => {
  editRecordObj.id = 0;
  editRecordObj.record_name = null;
  editRecordObj.record_content = null;
  editRecordObj.record_ttl = null;
};
const handleUpdateRecord = async () => {
  await request.put("/api/v1/dns/records", editRecordObj);
  message.success("更新成功");
  await queryRecords(searchObj);
  handleCancelEdit();
};

const formRef = ref();
const showAddModal = ref(false);
const changeOpenAddModal = () => {
  showAddModal.value = !showAddModal.value;
};
const addDnsObj = reactive({
  record_name: null,
  record_ttl: 60,
  zone: 'test.com.',
  record_type: 'A',
  record_content: null
});
const addDnsRules = {
  record_name: [{
    required: true,
    message: '请输入dns域名',
  }],
  record_ttl: [{
    required: true,
    message: '请输入域名ttl',
  }],
  zone: [{
    required: true,
    message: '请选择域名所属域',
  }],
  record_type: [{
    required: true,
    message: '请选择域名解析类型',
  }],
  record_content: [{
    required: true,
    message: '请输入域名解析内容',
  }]
};
const handletAddDns = async () => {
  await request.post('/api/v1/dns/records', addDnsObj);
  message.success("添加成功");
  changeOpenAddModal();
  queryRecords(searchObj);
};
const handleResetAddDns = () => {
  formRef.value.resetFields();
};


queryRecords(searchObj);
</script>

<style scoped>
.dns-search {
  margin-bottom: 10px;
  display: flex;
  justify-content: space-between;
}

.first-input {
  margin-top: 24px;
}

.add-dns-button {
  display: flex;
  justify-content: flex-end;
}
</style>
