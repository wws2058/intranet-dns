<template>
    <a-tabs v-model:activeKey="activeKey">
        <a-tab-pane key="gather" tab="汇总分析">
            <span class="sub-title">统计</span>
            <div class="flex-container">
                <template v-for="i in 6" :key="i">
                    <div class="card">
                        <span class="span font">统计项{{ i }}</span>
                        <div class="separator"></div>
                        <span class="span">{{ getRandomNum() }}</span>
                    </div>
                </template>
            </div>

            <div style="height: 30px;"></div>

            <span class="sub-title">分布</span>
            <div class="flex-container">
                <template v-for="i in 4" :key="i">
                    <v-chart class="pie" :option="getPieOptions(i, i * 3)" />
                </template>
            </div>
        </a-tab-pane>

        <a-tab-pane key="trend" tab="趋势分析">
            <div class="sub-title">趋势</div>
            <v-chart class="line-bar" :option="getLineOptions(20)"></v-chart>
            <div style="height: 10px;"></div>
            <div class="sub-title">实时</div>
            <v-chart class="line-bar" :option="getBarOptions(20)"></v-chart>
        </a-tab-pane>
    </a-tabs>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { PieChart, LineChart, BarChart } from "echarts/charts";
import {
    TitleComponent,
    TooltipComponent,
    LegendComponent,
    ToolboxComponent,
    GridComponent,
} from "echarts/components";
import VChart from "vue-echarts";
import { debounce } from 'lodash';
import dayjs from 'dayjs';

const activeKey = ref('trend');

use([
    CanvasRenderer,
    PieChart,
    LineChart,
    TitleComponent,
    TooltipComponent,
    ToolboxComponent,
    LegendComponent,
    GridComponent,
    BarChart
]);

function getPastNDayArray(n) {
    const dateArray = [];
    const currentDate = dayjs();
    for (let i = 0; i < n; i++) {
        const pastDate = currentDate.subtract(i, 'day');
        const formattedDate = pastDate.format('YYYY-MM-DD');
        dateArray.push(formattedDate);
    }
    dateArray.reverse();
    return dateArray;
}

function getRandomNumArray(num) {
    const numArray = [];
    for (let i = 0; i < num; i++) {
        const num = (Math.random() + 1) * (Math.random() * 1000);
        numArray.push(parseFloat(num.toFixed(2)));
    }
    return numArray;
}

function getRandomNum() {
    const num = (Math.random() + 1) * (Math.random() * 1000);
    return parseFloat(num.toFixed(2));
}

function getCharArr(num) {
    const letterArray = [];
    const startCharCode = 'a'.charCodeAt(0);
    for (let i = 0; i < num; i++) {
        letterArray.push(String.fromCharCode(startCharCode + i));
    }
    return letterArray;
}

function getLineCharObject(num) {
    const letterArray = getCharArr(num);
    const lineObjectArray = [];
    letterArray.forEach((letter) => {
        lineObjectArray.push(
            {
                name: letter,
                type: 'line',
                smooth: true,
                data: getRandomNumArray(num)
            }
        );
    });
    return lineObjectArray;
}

function getBarCharObject(num) {
    const codeArray = ['401', '403', '504'];
    const barObjectArray = [];
    codeArray.forEach((letter) => {
        barObjectArray.push(
            {
                name: letter,
                type: 'bar',
                data: getRandomNumArray(num)
            }
        );
    });
    return barObjectArray;
}

function getPieCharObject(num) {
    const letterArray = getCharArr(num);
    const letterObjectArray = [];
    letterArray.forEach((letter) => {
        letterObjectArray.push(
            {
                name: letter,
                value: getRandomNum()
            }
        );
    });
    return letterObjectArray;
}

// 饼状图
function getPieOptions(index, num) {
    return {
        title: {
            text: `分布统计项${index}`,
            left: 'center',
            top: '5%',
            textStyle: {
                fontSize: 16
            }
        },
        tooltip: {
            trigger: "item",
            formatter: "{b}: {c}/{d}%",
        },
        legend: {
            type: "scroll",
            bottom: "5%",
            data: getCharArr(num),
            itemWidth: 10,
            itemHeight: 10,
        },
        toolbox: {
            show: true,
            feature: {
                saveAsImage: {}
            }
        },
        series: [
            {
                type: "pie",
                radius: "75%",
                center: ["50%", "50%"],
                data: getPieCharObject(num),
                labelLine: {
                    show: false
                },
                label: {
                    show: false
                },
                emphasis: {
                    itemStyle: {
                        shadowBlur: 10,
                        shadowOffsetX: 0,
                        shadowColor: "rgba(0, 0, 0, 0.5)"
                    }
                }
            }
        ]
    };
}

// 折线图
function getLineOptions(num) {
    return {
        title: {
            text: '查询请求',
            left: '1%',
            textStyle: {
                fontSize: 14
            }
        },
        tooltip: {
            trigger: 'axis',
            enterable: true,
            position: function (point) {
                return [point[0] + 10, point[1] - 10];
            },
            extraCssText: 'overflow-y: auto; height: 200px;',
            formatter: function (params) {
                params.sort((a, b) => b.value - a.value);
                let tooltipStr = `${params[0].axisValue}<br/>`;
                params.forEach(item => {
                    tooltipStr += `<span style="display:inline-block;margin-right:5px;border-radius:10px;width:9px;height:9px;background-color:${item.color};"></span>`;
                    tooltipStr += `${item.seriesName}: ${item.value}<br/>`;
                });
                return tooltipStr;
            }
        },
        legend: {
            data: getCharArr(num).map((item) => item + ".host"),
            type: "scroll",
            itemWidth: 10,
            itemHeight: 10,
            width: '30%',
        },
        grid: {
            left: '1%',
            right: '4%',
            bottom: '1%',
            containLabel: true
        },
        toolbox: {
            feature: {
                saveAsImage: {}
            }
        },
        xAxis: {
            type: 'category',
            boundaryGap: false,
            data: getPastNDayArray(num),
        },
        yAxis: {
            type: 'value'
        },
        series: getLineCharObject(num).map((objectItem) => {
            return {
                ...objectItem,
                name: objectItem.name + '.host'
            };
        })
    };
}

// 直方图
function getBarOptions(num) {
    return {
        title: {
            text: '错误统计',
            left: '1%',
            textStyle: {
                fontSize: 14
            }
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            }
        },
        legend: {
            type: "scroll",
            itemWidth: 10,
            itemHeight: 10,
            width: '30%',
        },
        toolbox: {
            feature: {
                saveAsImage: {}
            }
        },
        grid: {
            left: '1%',
            right: '4%',
            bottom: '1%',
            containLabel: true
        },
        xAxis: {
            type: 'category',
            boundaryGap: true,
            data: getCharArr(num).map((item) => item + '.service')
        },
        yAxis: {
            type: 'value',
        },
        series: getBarCharObject(num)
    };
}

// 防抖resize重绘
function handleResize() {
    sessionStorage.setItem('a-tab-pane', activeKey.value);
    location.reload();
}
const debouncedHandleResize = debounce(handleResize, 50);

onMounted(() => {
    const savedKey = sessionStorage.getItem('a-tab-pane');
    if (savedKey) {
        activeKey.value = savedKey;
    }
    window.addEventListener('resize', debouncedHandleResize);
});

onUnmounted(() => {
    window.removeEventListener('resize', debouncedHandleResize);
    debouncedHandleResize.cancel();
});
</script>

<style scoped>
.sub-title {
    margin-bottom: 10px;
    padding-left: 10px;
    border-left: 3px solid blue;
    height: 30px;
    line-height: 30px;
    font-weight: bold;
    display: block;
}

.flex-container {
    display: flex;
    justify-content: space-between;
    margin-bottom: 10px;
    gap: 10px;
    width: 100%;
}

.card {
    border: 1px solid #e8e8e8;
    border-radius: 5px;
    box-sizing: border-box;
    flex-grow: 1;
}

.pie {
    box-sizing: border-box;
    flex-grow: 1;
    height: 400px;
    /* min-width: 400px; */
    background-color: rgb(248, 248, 248);
}

.line-bar {
    height: 285px;
    /* background-color: rgba(163, 11, 11, 0.289); */
    background-color: rgba(233, 230, 230, 0.15);
}

.separator {
    border-top: 1px solid #e8e8e8;
}

.font {
    font-weight: bold;
}

.span {
    line-height: 40px;
    height: 40px;
    display: block;
    text-align: center;
}
</style>