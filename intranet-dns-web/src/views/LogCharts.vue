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
            <!-- 折线图 -->

            <div style="height: 30px;"></div>

            <div class="sub-title">实时</div>
            <!-- 饼状图 -->
        </a-tab-pane>
    </a-tabs>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { PieChart } from "echarts/charts";
import {
    TitleComponent,
    TooltipComponent,
    LegendComponent,
} from "echarts/components";
import VChart from "vue-echarts";
import { debounce } from 'lodash';

use([
    CanvasRenderer,
    PieChart,
    TitleComponent,
    TooltipComponent,
    LegendComponent
]);

const activeKey = ref('gather');
function getRandomNum() {
    const num = (Math.random() + 1) * 1000;
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

function getCharObject(num) {
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
        series: [
            {
                type: "pie",
                radius: "80%",
                center: ["50%", "50%"],
                data: getCharObject(num),
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

function handleResize() {
    location.reload();
}
const debouncedHandleResize = debounce(handleResize, 50);

onMounted(() => {
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