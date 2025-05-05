<template>
  <a-space class="search">
    <a-input-search v-model:value.trim="searchObj.record_name" placeholder="请输入域名" allowClear @search="getProbes" />
    <a-button type="primary" @click="changeOpenAddModal">新增探测</a-button>
  </a-space>


  <a-table :pagination="pagination" :columns="columns" :dataSource="dataSource" :loading="tableLoading"
    @change="handlePageChange">
    <template #bodyCell="{ column, record, text }">
      <template v-if="column.key === 'intranet'">
        <a-tag v-if="text" color="cyan">内网域名</a-tag>
        <a-tag v-else color="blue">公网域名</a-tag>
      </template>

      <template v-if="column.key === 'succeed'">
        <a-tag v-if="text === 0" color="green">探测成功</a-tag>
        <a-tag v-if="text === 1" color="pink">查询失败</a-tag>
        <a-tag v-if="text === 2" color="red">探测失败</a-tag>
      </template>

      <template v-if="column.key === 'expect_answer' && record.id === updateRequestObj.id">
        <a-tooltip>
          <template #title>{{ updateRequestObj.expect_answer.replaceAll(',', '\n') }}</template>
          <a-input v-model:value.trim="updateRequestObj.expect_answer">
          </a-input>
        </a-tooltip>
      </template>

      <template v-if="column.key === 'operations'">
        <template v-if="record.id === updateRequestObj.id">
          <a-space>
            <a-button @click="handleCancel(record)" size="small">取消</a-button>
            <a-button type="primary" @click="handleUpdate" size="small">确认</a-button>
          </a-space>
        </template>
        <template v-else>
          <a-popconfirm title="删除定时任务?" @confirm="deleteProbe(record.id)" okText="确认" cancelText="取消">
            <a>删除</a>
          </a-popconfirm>
          <a-divider type="vertical" />
          <a @click="handleEdit(record)">编辑</a>
        </template>
      </template>
    </template>
  </a-table>

  <a-modal v-model:open="showAddModal" title="新增探测任务" :footer="null">

    <a-form ref="formRef" :rules="rules" labelAlign="right" :label-col="{ span: 4 }" :model="addRequestObj"
      @finish="handleAddProbe">
      <a-form-item label="探测域名" name="record_name" class="first-input">
        <a-input allowClear v-model:value.trim="addRequestObj.record_name"></a-input>
      </a-form-item>
      <a-form-item label="域名区域" name="zone">
        <a-select :options="zoneOptions" v-model:value="addRequestObj.zone"></a-select>
      </a-form-item>
      <a-form-item allowClear label="预期结果" name="expect_answer">
        <a-input allowClear placeholder="请输入域名预期解析结果, 英文逗号分隔"
          v-model:value.trim="addRequestObj.expect_answer"></a-input>
      </a-form-item>
      <a-form-item label="域名属性" name="intranet">
        <a-select :options="domainOptions" v-model:value="addRequestObj.intranet"></a-select>
      </a-form-item>

      <a-form-item>
        <a-space class="add-probe-button">
          <a-button @click="handleResetForm">重置</a-button>
          <a-button type="primary" htmlType="submit">确认</a-button>
        </a-space>
      </a-form-item>
    </a-form>

  </a-modal>

</template>

<script setup>
import request from '@/apis/request';
import { message } from 'ant-design-vue';
import { reactive, ref, h, computed, watch, toRaw } from 'vue';
import { intranetZones } from '@/apis/const';

// edit
const updateRequestObj = reactive({
  id: null,
  expect_answer: null,
});
const handleCancel = () => {
  Object.keys(updateRequestObj).forEach(key => {
    updateRequestObj[key] = null;
  });
};
function handleEdit(record) {
  updateRequestObj.id = record.id;
  updateRequestObj.expect_answer = record.expect_answer.join(',');
  console.log(toRaw(updateRequestObj));
};
const handleUpdate = async () => {
  updateRequestObj.expect_answer = updateRequestObj.expect_answer.split(',');
  await request.put('/api/v1/dns/probes', updateRequestObj);
  handleCancel();
  getProbes();
  message.success('更新成功');
};

// add
const zoneOptions = computed(() => {
  const zOptions = intranetZones.map((zone) => ({
    label: zone,
    value: zone
  }));
  zOptions.unshift({
    label: ".",
    value: "."
  });
  return zOptions;
});
const domainOptions = [
  {
    value: false,
    label: '公网域名'
  },
  {
    value: true,
    label: '内网域名'
  },
];
const showAddModal = ref(false);
const changeOpenAddModal = () => {
  showAddModal.value = !showAddModal.value;
};
const formRef = ref();
const handleResetForm = () => {
  formRef.value.resetFields();
};
const addRequestObj = reactive({
  record_name: null,
  zone: ".",
  expect_answer: null, // arr
  intranet: false
});
watch(() => addRequestObj.zone, () => {
  addRequestObj.intranet = intranetZones.includes(addRequestObj.zone) ? true : false;
});
const rules = {
  record_name: [{
    required: true,
    message: '请输入要探测的域名',
  }],
  zone: [{
    required: true,
    message: '域名所属域',
  }],
  expect_answer: [{
    required: true,
    message: '请输入域名预期解析结果, 英文逗号分隔',
  }],
  intranet: [{
    required: true,
    message: '请选择域名类型(公网域名 或 内网域名)',
  }]
};

const handleAddProbe = async () => {
  addRequestObj.expect_answer = addRequestObj.expect_answer.split(',');
  await request.post("/api/v1/dns/probes", addRequestObj);
  changeOpenAddModal();
  message.success('添加成功');
  getProbes();
};

// delete
const deleteProbe = async (id) => {
  await request.delete(`/api/v1/dns/probes/${id}`);
  message.success("删除成功");
  await getProbes();
};

// get
const columns = ref([
  {
    title: '探测域名',
    dataIndex: 'record_name',
    key: 'record_name',
  },
  {
    title: 'zone',
    dataIndex: 'zone',
    key: 'zone',
  },
  {
    title: '域名属性',
    dataIndex: 'intranet',
    key: 'intranet'
  },
  {
    title: '预期结果',
    dataIndex: 'expect_answer',
    key: 'expect_answer',
    customRender: ({ text }) => {
      return h(
        'div',
        {
          style: {
            whiteSpace: 'pre-wrap',
            maxWidth: '250px'
          }
        },
        text.join('\n')
      );
    }
  },
  {
    title: '探测状态',
    dataIndex: 'succeed',
    key: 'succeed',
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

const searchObj = reactive({
  record_name: null,
  page: 1,
  page_size: 10
});
const tableLoading = ref(false);
const getProbes = async () => {
  tableLoading.value = true;
  const rsp = await request.get('/api/v1/dns/probes', searchObj);
  dataSource.value = rsp.data.data;
  pagination.total = rsp.data.pages.total;
  tableLoading.value = false;
};
const handlePageChange = async (pages) => {
  pagination.current = pages.current;
  pagination.pageSize = pages.pageSize;
  searchObj.page = pagination.current;
  searchObj.page_size = pagination.pageSize;
  await getProbes();
};

getProbes();
</script>


<style scoped>
.search {
  margin-bottom: 10px;
  display: flex;
  justify-content: space-between;
}

.first-input {
  margin-top: 24px;
}

.add-probe-button {
  display: flex;
  justify-content: flex-end;
}
</style>
