<template>
  <div class="query">

    <a-input v-model:value.trim="queryDomain" placeholder="请输入要查询的域名" allowClear @pressEnter="handleQuery" size="large"
      :class="queryInputClass"></a-input>

    <a-table :columns="columns" :dataSource="dataSource" :pagination="false" class="dns-result"
      :loading="dnsQueryLoading" v-show="showTable"></a-table>

  </div>
</template>

<script setup>
import request from '@/apis/request';
import { intranetZones } from '@/apis/const';
import { computed, ref, watch, h } from 'vue';

const queryDomain = ref('');
const queryInputClass = ref(['search-input']);
watch(queryDomain, () => {
  if (!queryDomain.value) {
    dataSource.value = [];
    columns.value = [];
    showTable.value = false;
    queryInputClass.value = ['search-input'];
  }
});
const intranetQueryObj = computed(() => {
  let domain = queryDomain.value;
  if (!queryDomain.value.endsWith('.')) {
    domain += ".";
  }
  return {
    isIntranet: intranetZones.some((zone) => domain.endsWith(zone)),
    zone: intranetZones.find((zone) => domain.endsWith(zone)),
    domain: domain
  };
});

const dnsQueryLoading = ref(false);
async function isIntranetQuery() {
  const requestObj = {
    domain: queryDomain.value,
    zone: intranetQueryObj.value.zone
  };
  dnsQueryLoading.value = true;
  const rsp = await request.get('/api/v1/dns/rrs', requestObj);
  const rrs = rsp.data.data.sort((a, b) => {
    if (a.record_type < b.record_type) {
      return -1;
    }
    if (a.record_type > b.record_type) {
      return 1;
    }
    if (a.record_content < b.record_content) {
      return -1;
    }
    if (a.record_content > b.record_content) {
      return 1;
    }
    return 0;
  });
  dnsQueryLoading.value = false;
  dataSource.value = rrs;
}

// edns查询119.29.29.29
async function ednsQuery() {
  try {
    const requestObj = {
      domain: queryDomain.value,
    };
    dnsQueryLoading.value = true;
    const rsp = await request.get('/api/v1/dns/edns', requestObj);
    const result = rsp.data.data;
    const ispRRs = [];
    result.forEach(data => {
      const tmpObj = {
        isp: data.isp,
        client_ip: data.client_ip,
        type_cname: [],
        type_a: [],
        type_aaaa: []
      };
      data.dns_rrs.forEach(rr => {
        switch (rr.record_type) {
          case 'A':
            tmpObj.type_a.push(rr.record_content);
            break;
          case 'AAAA':
            tmpObj.type_aaaa.push(rr.record_content);
            break;
          case 'CNAME':
            tmpObj.type_cname.push(rr.record_content);
            break;
          default:
            break;
        }
      });
      tmpObj.type_a.sort();
      tmpObj.type_aaaa.sort();
      tmpObj.type_cname.sort();
      ispRRs.push(tmpObj);
    });

    ispRRs.sort((a, b) => {
      if (a.isp < b.isp) {
        return -1;
      }
      if (a.isp > b.isp) {
        return 1;
      }
      if (a.client_ip < b.client_ip) {
        return -1;
      }
      if (a.client_ip > b.client_ip) {
        return 1;
      }
      return 0;
    });
    dataSource.value = ispRRs;
    dnsQueryLoading.value = false;
  } catch (error) {
    dnsQueryLoading.value = false;
    console.log(error);
  }
}

// 内网结果的zone
const columns = ref([]);
const dataSource = ref([]);
const showTable = ref(false);
const intranetColumn = [
  {
    title: '内网域名',
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
    title: '记录值',
    dataIndex: 'record_content',
    key: 'record_content',
  },
];
const publicColumn = [
  {
    title: '省份运营商',
    dataIndex: 'isp',
    key: 'isp',
  },
  {
    title: '客户端地址',
    dataIndex: 'client_ip',
    key: 'client_ip',
  },
  {
    title: 'CNAME记录',
    dataIndex: 'type_cname',
    key: 'type_cname',
    customRender: ({ text }) => {
      return h(
        'div',
        {
          style: {
            whiteSpace: 'pre-wrap'
          }
        },
        text.join('\n')
      );
    }
  },
  {
    title: 'A记录',
    dataIndex: 'type_a',
    key: 'type_a',
    customRender: ({ text }) => {
      return h(
        'div',
        {
          style: {
            whiteSpace: 'pre-wrap'
          }
        },
        text.join('\n')
      );
    }
  },
  {
    title: 'AAAA记录',
    dataIndex: 'type_aaaa',
    key: 'type_aaaa',
    customRender: ({ text }) => {
      return h(
        'div',
        {
          style: {
            whiteSpace: 'pre-wrap'
          }
        },
        text.join('\n')
      );
    }
  }
];

const handleQuery = async () => {
  queryInputClass.value = ['search-input-after'];
  if (intranetQueryObj.value.isIntranet) {
    columns.value = intranetColumn;
    showTable.value = true;
    isIntranetQuery();
  } else {
    columns.value = publicColumn;
    showTable.value = true;
    ednsQuery();
  }
};
</script>

<style scoped>
.query {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  min-height: 100%;
}

.search-input {
  width: 600px;
  min-width: 400px;
  margin-bottom: 15%;
}

.search-input-after {
  width: 600px;
  min-width: 400px;
  margin-bottom: 10px;
  opacity: 0.5;
}

.dns-result {
  min-width: 600px;
}
</style>
