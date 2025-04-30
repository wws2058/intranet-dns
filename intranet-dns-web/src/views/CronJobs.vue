<template>
    <div class="search">
        <a-space>
            <a-input allowClear placeholder="任务名称" v-model:value.trim="searchObj.name"></a-input>
            <a-input allowClear placeholder="创建者" v-model:value.trim="searchObj.creator"></a-input>
            <a-select :options="taskTypeOptions" v-model:value="searchObj.task_type"></a-select>
            <a-select :options="inUseOptions" v-model:value="searchObj.started"></a-select>
            <a-select :options="runningStatusOPtions" v-model:value="searchObj.last_succeed"></a-select>
            <a-button type="primary" @click="handleSearch" :loading="tableLoading">搜索</a-button>
            <a-button @click="resetSearchObj">重置</a-button>
        </a-space>
        <a-button type="primary" @click="changeOpenAddModal">新增任务</a-button>
    </div>

    <a-table :columns="columns" :dataSource="dataSource" :pagination="pagination" :loading="tableLoading"
        @change="handlePageChange">

        <template #bodyCell="{ column, record, text }">
            <template v-if="column.key === 'name' && updateRequestObj.id == record.id">
                <a-input v-model:value.trim="updateRequestObj.name" size="small"></a-input>
            </template>
            <template v-if="column.key === 'spec' && updateRequestObj.id == record.id">
                <a-input v-model:value.trim="updateRequestObj.spec" size="small"></a-input>
            </template>
            <template v-if="column.key === 'description' && updateRequestObj.id == record.id">
                <a-input v-model:value.trim="updateRequestObj.description" size="small"></a-input>
            </template>

            <template v-if="column.key === 'started'">
                <template v-if="record.id == updateRequestObj.id">
                    <a-switch v-model:checked="updateRequestObj.started" checkedChildren="开启" unCheckedChildren="关闭" />
                </template>
                <template v-else>
                    <a-switch @click="handleSwitch(record)" :checked="text" checkedChildren="开启"
                        unCheckedChildren="关闭" />
                </template>
            </template>

            <template v-if="column.key === 'last_succeed'">
                <a-tag color="green" v-if="text">成功</a-tag>
                <a-tag color="red" v-if="!text">失败</a-tag>
            </template>

            <template v-if="column.key === 'task_type'">
                {{taskTypeOptions.find(item => item.value === text).label}}
            </template>

            <template v-if="column.key === 'task_args'">
                <template v-if="updateRequestObj.id == record.id && record.task_type === 'function'">
                    <a-select :options="internalFunctionOptions" v-model:value="updateRequestObj.args"></a-select>
                </template>
                <template v-else-if="updateRequestObj.id == record.id && record.task_type === 'http'">
                    <a-input v-model:value.trim="updateRequestObj.args" size="small"></a-input>
                </template>
                <template v-else>
                    <a-tooltip :title="Object.entries(text)
                        .map(([key, value]) => `${key}: ${value}`)
                        .join(', ')" :color="'green'"><a>查看详情</a></a-tooltip>
                </template>
            </template>

            <template v-if="column.key === 'history'">
                <a @click="changeshowHistoryModal(record)">查看详情</a>
            </template>

            <template v-if="column.key === 'operations'">
                <template v-if="updateRequestObj.id == record.id">
                    <a-space>
                        <a-button @click="handleCancel(record)" size="small">取消</a-button>
                        <a-button type="primary" @click="handleUpdate" size="small">确认</a-button>
                    </a-space>
                </template>
                <template v-else>
                    <a-popconfirm title="删除定时任务?" @confirm="handleDeleteJob(record)" okText="确认" cancelText="取消">
                        <a>删除</a>
                    </a-popconfirm>
                    <a-divider type="vertical" />
                    <a @click="handleEdit(record)">编辑</a>
                </template>
            </template>

        </template>

    </a-table>


    <a-modal v-model:open="showAddModal" title="新增定时任务" :footer="null">

        <a-form ref="formRef" :rules="rules" labelAlign="right" :label-col="{ span: 4 }" :model="addRequestObj"
            @finish="handleAddJob">
            <a-form-item label="任务名称" name="name" class="first-input">
                <a-input allowClear v-model:value.trim="addRequestObj.name"></a-input>
            </a-form-item>
            <a-form-item label="任务描述" name="description">
                <a-input allowClear placeholder="简要描述任务" v-model:value.trim="addRequestObj.description"></a-input>
            </a-form-item>
            <a-form-item allowClear label="调度周期" name="spec">
                <a-input placeholder="分时日月周" v-model:value.trim="addRequestObj.spec"></a-input>
            </a-form-item>
            <a-form-item label="任务类型" name="task_type">
                <a-select :options="taskTypeOptions" v-model:value="addRequestObj.task_type"></a-select>
            </a-form-item>

            <a-form-item label="任务参数" name="args">
                <template v-if="addRequestObj.task_type === 'function'">
                    <a-select :options="internalFunctionOptions" v-model:value="addRequestObj.args"></a-select>
                </template>
                <template v-else-if="addRequestObj.task_type === 'http'">
                    <a-input placeholder="http get方法调用的url" v-model:value="addRequestObj.args"></a-input>
                </template>
            </a-form-item>


            <a-form-item>
                <a-space class="add-job-button">
                    <a-button @click="handleResetForm">重置</a-button>
                    <a-button type="primary" htmlType="submit">确认</a-button>
                </a-space>
            </a-form-item>
        </a-form>

    </a-modal>

    <a-modal v-model:open="showHistoryModal" title="定时任务执行历史" :footer="null" width="640px">
        <a-table :columns="historyColumns" :dataSource="historyData">

            <template #bodyCell="{ column, text, record }">
                <template v-if="column.key === 'succeed'">
                    <a-tag color="green" v-if="text">成功</a-tag>
                    <a-tooltip :title="record.error" :color="'red'"><a-tag color="red"
                            v-if="!text">失败</a-tag></a-tooltip>
                </template>
            </template>

        </a-table>
    </a-modal>


</template>

<script setup>
import { reactive, ref, watch } from 'vue';
import request from '@/apis/request';
import { isValidCron } from 'cron-validator';
import { message } from 'ant-design-vue';

// show history modal
const historyColumns = [
    {
        title: 'uid',
        dataIndex: 'uid',
        key: 'uid',
    },
    {
        title: '执行状态',
        dataIndex: 'succeed',
        key: 'succeed',
    },
    {
        title: '调用时间',
        dataIndex: 'call_at',
        key: 'call_at',
    },
];
const historyData = ref([]);
const showHistoryModal = ref(false);
const changeshowHistoryModal = (record) => {
    showHistoryModal.value = !showHistoryModal.value;
    historyData.value = record.history;
};


// edit job
const updateRequestObj = reactive({
    id: null,
    name: null,
    spec: null,
    started: null,
    description: null,
    task_type: null,
    args: null,
});
const handleCancel = () => {
    Object.keys(updateRequestObj).forEach(key => {
        updateRequestObj[key] = null;
    });
};
const handleEdit = (record) => {
    updateRequestObj.id = record.id;
    updateRequestObj.name = record.name;
    updateRequestObj.spec = record.spec;
    updateRequestObj.started = record.started;
    updateRequestObj.description = record.description;
    updateRequestObj.task_type = record.task_type;
    updateRequestObj.args = record.task_type === "http" ? record.task_args.url : record.task_args.function_name;
};
const handleUpdate = async () => {
    // 检查cron表达式
    if (!isValidCron(updateRequestObj.spec)) {
        message.error("spec cronjob表达式不合法");
        return;
    }

    const taskType = updateRequestObj.task_type;
    if (taskType === "function") {
        updateRequestObj.task_args = {
            function_name: updateRequestObj.args
        };
    }
    if (taskType === "http") {
        updateRequestObj.task_args = {
            url: updateRequestObj.args
        };
    }
    await request.put('/api/v1/cronjobs', updateRequestObj);
    message.success("更新成功");
    handleCancel();
    queryCronjobs(searchObj);
};
const handleSwitch = async (record) => {
    const requestObj = {
        id: record.id,
        started: !record.started
    };
    await request.put('/api/v1/cronjobs', requestObj);
    queryCronjobs(searchObj);
};

// del job
const handleDeleteJob = async (record) => {
    await request.delete(`/api/v1/cronjobs/${record.id}`);
    message.success(`${record.name} 删除成功`);
    queryCronjobs(searchObj);
};

// add job
const formRef = ref();
const handleResetForm = () => {
    formRef.value.resetFields();
};
const addRequestObj = reactive({
    name: null,
    description: null,
    spec: null,
    task_type: "function",
    args: null,
});
watch(() => addRequestObj.task_type, () => {
    addRequestObj.task_args = null;
});
const showAddModal = ref(false);
const changeOpenAddModal = () => {
    showAddModal.value = !showAddModal.value;
};
const rules = {
    name: [{
        required: true,
        message: '请输入任务名称',
    }],
    description: [{
        required: true,
        message: '请输入任务描述',
    }],
    spec: [{
        required: true,
        message: '请输入cron表达式',
    },
    {
        validator: (_, value) => {
            if (value && !isValidCron(value)) {
                return Promise.reject('请输入有效的cron表达式');
            }
            return Promise.resolve();
        }
    }],
    task_type: [{
        required: true,
        message: '请选择任务类型',
    }],
    args: [{
        required: true,
        message: '任务参数必填',
    }],
};
const internalFunctionOptions = ref([]);

const getInternalFunctionOptions = async () => {
    const rsp = await request.get("/api/v1/cronjobs/functions");
    internalFunctionOptions.value = rsp.data.data.map((functionName) => ({
        label: functionName,
        value: functionName
    }));
};

const handleAddJob = async () => {
    const taskType = addRequestObj.task_type;
    if (taskType === "function") {
        addRequestObj.task_args = {
            function_name: addRequestObj.args
        };
    }
    if (taskType === "http") {
        addRequestObj.task_args = {
            url: addRequestObj.args
        };
    }
    await request.post("/api/v1/cronjobs", addRequestObj);
    changeOpenAddModal();
    message.success(`${addRequestObj.name} 任务添加成功`);
    queryCronjobs(searchObj);
};


// search 
const searchObj = reactive(
    {
        name: null,
        creator: null,
        task_type: null,
        started: null,
        last_succeed: null,
        page: 1,
        page_size: 10
    }
);

const resetSearchObj = () => {
    searchObj.name = null;
    searchObj.creator = null;
    searchObj.task_type = null;
    searchObj.started = null;
    searchObj.last_succeed = null;
    searchObj.page = 1;
    searchObj.page_size = 10;

    pagination.current = 1;
    pagination.pageSize = 10;

    queryCronjobs(searchObj);
};
const handleSearch = () => {
    queryCronjobs(searchObj);
};

const taskTypeOptions = ref([
    {
        label: "任务类型",
        value: null,
        disabled: true
    },
    {
        label: "http调用",
        value: "http"
    },
    {
        label: "内置函数",
        value: "function"
    },
    {
        label: "shell脚本",
        value: "shell",
        disabled: true
    },
]);
const inUseOptions = ref([
    {
        label: "启用状态",
        value: null,
        disabled: true
    },
    {
        label: "启动",
        value: "true"
    },
    {
        label: "停止",
        value: "false"
    },
]);
const runningStatusOPtions = ref([
    {
        label: "运行状态",
        value: null,
        disabled: true
    },
    {
        label: "成功",
        value: "true"
    },
    {
        label: "失败",
        value: "false"
    },
]);

const columns = ref([
    {
        title: '任务名称',
        dataIndex: 'name',
        key: 'name',
    },
    {
        title: '调度时间',
        dataIndex: 'spec',
        key: 'spec',
    },
    {
        title: '创建者',
        dataIndex: 'creator',
        key: 'creator',
    },
    {
        title: '描述信息',
        dataIndex: 'description',
        key: 'description',
    },
    {
        title: '启用状态',
        dataIndex: 'started',
        key: 'started',
    },
    {
        title: '运行状态',
        dataIndex: 'last_succeed',
        key: 'last_succeed',
    },
    {
        title: '任务类型',
        dataIndex: 'task_type',
        key: 'task_type',
    },
    {
        title: '任务参数',
        dataIndex: 'task_args',
        key: 'task_args',
    },
    {
        title: '历史记录',
        dataIndex: 'history',
        key: 'history',
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
    await queryCronjobs(searchObj);
};

async function queryCronjobs(searchObj) {
    try {
        tableLoading.value = true;
        const rsp = await request.get('/api/v1/cronjobs', searchObj);
        dataSource.value = rsp.data.data;
        pagination.total = rsp.data.pages.total;
        tableLoading.value = false;
    } catch (error) {
        console.log(error);
        tableLoading.value = false;
    }
}

getInternalFunctionOptions();
queryCronjobs(searchObj)

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

.add-job-button {
    display: flex;
    justify-content: flex-end;
}
</style>