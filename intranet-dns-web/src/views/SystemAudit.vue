<template>
  <!-- 搜索条件 -->
  <a-space class="log-search">
    <a-input placeholder="用户英文名" v-model:value.trim="searchLogOBJ.user_name"></a-input>
    <a-input placeholder="请求唯一id" v-model:value.trim="searchLogOBJ.request_id"></a-input>
    <a-input placeholder="客户端ip" v-model:value.trim="searchLogOBJ.client_ip"></a-input>
    <a-range-picker :defaultValue="searchTimeArr" v-model:value="searchTimeArr" allowClear style="width: 340px"
      :showTime="showTimeConfig" :locale="locale" :format="timeFormat" :valueFormat="timeFormat"
      :presets="rangePresets"></a-range-picker>
    <a-button type="primary" @click="queryLogs" :loading="searchLoading">搜索</a-button>
    <a-button @click="resetSearchRequest">重置</a-button>
  </a-space>

  <!-- 展示数据 -->
  <a-table :loading="searchLoading" :columns="columns" :dataSource="dataSource" :pagination="pagination"
    @change="handlePageChange" :rowKey="dataSourceKey" @expand="handleExpand" :expandedRowKeys="expandedRowKeys"
    tableLayout="fixed">

    <template #bodyCell="{ column, record }">
      <template v-if="column.key === 'user_name'">
        {{ record?.user_name ? record.user_name : '-' }}
      </template>

      <template v-if="column.key === 'method'">
        <template v-if="record.method === 'DELETE'">
          <a-tag color="error" :bordered="false">{{ record.method }}</a-tag>
        </template>
        <template v-else-if="record.method === 'GET'">
          <a-tag color="success" :bordered="false">{{ record.method }}</a-tag>
        </template>
        <template v-else-if="record.method === 'PUT'">
          <a-tag color="processing" :bordered="false">{{ record.method }}</a-tag>
        </template>
        <template v-else-if="record.method === 'POST'">
          <a-tag color="lime" :bordered="false">{{ record.method }}</a-tag>
        </template>
        <template v-else>
          <a-tag color="gold" :bordered="false">{{ record.method }}</a-tag>
        </template>
      </template>

      <template v-if="column.key === 'time_cost'">
        <template v-if="record.time_cost <= 200">
          <a-tag color="success" :bordered="false">{{ record.time_cost }}</a-tag>
        </template>
        <template v-else-if="record.time_cost > 200 && record.time_cost < 500">
          <a-tag color="warning" :bordered="false">{{ record.time_cost }}</a-tag>
        </template>
        <template v-else>
          <a-tag color="error" :bordered="false">{{ record.time_cost }}</a-tag>
        </template>
      </template>

      <!-- <template v-if="column.key === 'detail'">
        <a @click.prevent="showRequestDetail(record)">查看详情</a>
      </template> -->
    </template>

    <template #expandedRowRender="{ record }">
      <div class="code_expand">
        <div class="code_container">
          <a-tag color="blue">请求内容</a-tag>
          <JsonView :jsonStr="formatBodyStr(record)[0]"></JsonView>
        </div>
        <div class="code_container">
          <a-tag color="blue">返回内容</a-tag>
          <JsonView :jsonStr="formatBodyStr(record)[1]"></JsonView>
        </div>
      </div>
    </template>

  </a-table>
</template>

<script setup>
import { computed, reactive, ref } from 'vue';
import dayjs from 'dayjs';
import locale from 'ant-design-vue/es/date-picker/locale/zh_CN';
import request from '@/apis/request';
import JsonView from '@/components/JsonView.vue';

// 请求栏
dayjs.locale('zh-cn');
const timeFormat = "YYYY-MM-DD HH:mm:ss";
const searchLoading = ref(false);
const searchTimeArr = ref([dayjs().add(-7, 'd').format(timeFormat), dayjs().format(timeFormat)]);
// 日期时间选择器
const showTimeConfig = reactive({
  minuteStep: 10,
  secondStep: 30,
});
// 预设时间
const rangePresets = reactive([
  {
    label: '最近一天',
    value: [dayjs().add(-1, 'd'), dayjs()],
  },
  {
    label: '最近三天',
    value: [dayjs().add(-3, 'd'), dayjs()],
  },
  {
    label: '最近一周',
    value: [dayjs().add(-7, 'd'), dayjs()],
  }
]);

// 表格数据
const columns = [
  {
    title: "用户",
    dataIndex: 'user_name',
    key: "user_name"
  },
  {
    title: "请求唯一id",
    dataIndex: 'request_id',
    key: "request_id",
    width: 350
  },
  {
    title: "客户端ip",
    dataIndex: 'client_ip',
    key: "client_ip",
  },
  {
    title: "请求方法",
    dataIndex: 'method',
    key: "method"
  },
  {
    title: "请求接口",
    dataIndex: 'url',
    key: "url",
    width: 300
  },
  {
    title: "请求耗时(ms)",
    dataIndex: 'time_cost',
    key: "time_cost"
  },
  {
    title: "创建时间",
    dataIndex: 'created_at',
    key: "created_at"
  },
];
const dataSource = ref([]);
const dataSourceKey = ref('id');
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
// 搜索框绑定的搜索对象
const searchLogOBJ = reactive({
  user_name: null,
  request_id: null,
  client_ip: null,

  start_time: computed(() => searchTimeArr.value[0]),
  end_time: computed(() => searchTimeArr.value[1]),
  page: computed(() => pagination.current),
  page_size: computed(() => pagination.pageSize),
});
// 重置搜索栏, page为1, page_size为10, 其他为null
function resetSearchRequest() {
  searchTimeArr.value.fill(null);
  pagination.current = 1;
  pagination.pageSize = 10;

  searchLogOBJ.user_name = null;
  searchLogOBJ.request_id = null;
  searchLogOBJ.client_ip = null;
  queryLogs();
}

// 页数变化时请求
async function handlePageChange(pages) {
  pagination.current = pages.current;
  pagination.pageSize = pages.pageSize;
  queryLogs();
}

// 审计日志后端请求
async function queryLogs() {
  try {
    searchLoading.value = true;
    const response = await request.get("/api/v1/audit_logs", searchLogOBJ);
    const result = response.data;
    dataSource.value = result.data;
    pagination.total = result.pages.total;
    searchLoading.value = false;
  } catch (error) {
    searchLoading.value = false;
    console.log(error);
  }
}

// 点击查看详情
const expandedRowKeys = computed(() => [expandRequestDetail.value.id]);

const expandRequestDetail = ref({
  id: null,
  request_body: null,
  response_body: null,
  expanded: false
});

function handleExpand(expanded, record) {
  if (!expanded) {
    expandRequestDetail.value.id = null;
    return;
  }
  const formattedArr = formatBodyStr(record);
  expandRequestDetail.value.id = record.id;
  expandRequestDetail.value.request_body = formattedArr[0];
  expandRequestDetail.value.response_body = formattedArr[1];
  expandRequestDetail.value.expanded = expanded;
}

function formatBodyStr(record) {
  let reqJson = record.request_body ? record.request_body : "null";
  let rspJson = record.response_body ? record.response_body : "null";
  let reqJsonLines, rspJsonLines;
  if (reqJson) {
    reqJson = JSON.stringify(JSON.parse(reqJson), null, 2);
    reqJsonLines = reqJson.split('\n').length;
  }
  if (rspJson) {
    rspJson = JSON.stringify(JSON.parse(rspJson), null, 2);
    rspJsonLines = rspJson.split('\n').length;
  }

  const maxLines = Math.max(reqJsonLines, rspJsonLines);
  for (let i = 0; i <= maxLines - reqJsonLines; i++) {
    reqJson += '\n';
  }
  for (let i = 0; i <= maxLines - rspJsonLines; i++) {
    rspJson += '\n';
  }
  return [reqJson, rspJson];
}

queryLogs();
</script>

<style scoped>
.log-search {
  margin-bottom: 10px;
}

.code_expand {
  display: flex;
  max-height: 300px;
}

.code_container {
  width: 50%;
  padding: 0 5px;
  overflow: auto;
}
</style>
