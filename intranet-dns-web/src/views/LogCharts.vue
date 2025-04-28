<template>
    <div ref="chartRef" style="width: 600px; height: 400px;"></div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import * as echarts from 'echarts';

const chartRef = ref(null);

onMounted(() => {
    const chart = echarts.init(chartRef.value);

    const option = {
        xAxis: {
            type: 'category',
            data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
        },
        yAxis: {
            type: 'value'
        },
        series: [
            {
                data: [820, 932, 901, 934, 1290, 1330, 1320],
                type: 'line'
            }
        ]
    };

    chart.setOption(option);
});
</script>