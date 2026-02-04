<template>
  <t-loading :loading="statisticsLoading" size="large">
    <t-row :gutter="[16, 16]">
      <t-col v-for="(item, index) in statisticsData" :key="item.title" :xs="6" :xl="3">
        <t-card
            :title="item.title"
            :bordered="false"
            :class="{ 'dashboard-item': true, 'dashboard-item--main-color': index == 0 }"
        >
          <div class="dashboard-item-top">
            <span :style="{ fontSize: `${resizeTime * 28}px` }">{{ item.number }}</span>
          </div>
          <div class="dashboard-item-left">
            <div
                v-if="index === 0"
                id="moneyContainer"
                class="dashboard-chart-container"
                :style="{ width: `${resizeTime * 120}px`, height: '100px', marginTop: '-24px' }"
            ></div>
            <div
                v-else-if="index === 1"
                id="refundContainer"
                class="dashboard-chart-container"
                :style="{ width: `${resizeTime * 120}px`, height: '56px', marginTop: '-24px' }"
            ></div>
            <span v-else-if="index === 2" :style="{ marginTop: `-24px` }">
              <usergroup-icon />
            </span>
            <span v-else :style="{ marginTop: '-24px' }">
              <file-icon />
            </span>
          </div>
          <template #footer>
            <div class="dashboard-item-bottom">
              <div class="dashboard-item-block">
                {{ t('pages.dashboardBase.topPanel.cardTips') }}
                <trend
                    class="dashboard-item-trend"
                    :type="item.upTrend ? 'up' : 'down'"
                    :is-reverse-color="index === 0"
                    :describe="item.upTrend || item.downTrend"
                />
              </div>
              <t-icon name="chevron-right" />
            </div>
          </template>
        </t-card>
      </t-col>
    </t-row>
  </t-loading>
</template>

<script lang="ts">
export default {
  name: 'DashboardBase',
};
</script>

<script setup lang="ts">
import { useWindowSize } from '@vueuse/core';
import { BarChart, LineChart } from 'echarts/charts';
import * as echarts from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { FileIcon, UsergroupIcon } from 'tdesign-icons-vue-next';
import { MessagePlugin } from 'tdesign-vue-next';
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue';

import {
  getStatisticsData,
  type StatisticsCard
} from '@/api/statistics'

// 导入样式
import Trend from '@/components/trend/index.vue';
import { t } from '@/locales';
import { useSettingStore } from '@/store';
import { changeChartsTheme } from '@/utils/color';

import { constructInitDashboardDataset } from '../index';

echarts.use([LineChart, BarChart, CanvasRenderer]);

const store = useSettingStore();
const resizeTime = ref(1);

// ===== 替换为API数据 =====
const statisticsLoading = ref(false);
const statisticsData = ref<StatisticsCard[]>([]);

// 获取统计数据
const fetchStatisticsData = async () => {
  try {
    statisticsLoading.value = true;
    const res = await getStatisticsData();
    if (res.code === 0) {
      statisticsData.value = res.data;
    } else {
      MessagePlugin.error(res.msg || '获取统计数据失败');
      // 如果API失败，使用默认数据
      statisticsData.value = [
        {
          title: '总收入',
          number: '¥ 0.00',
          upTrend: '0%',
          leftType: 'echarts-line',
        },
        {
          title: '总支出',
          number: '¥ 0.00',
          downTrend: '0%',
          leftType: 'echarts-bar',
        },
        {
          title: '用户总数',
          number: '0',
          upTrend: '0%',
          leftType: 'icon-usergroup',
        },
        {
          title: '订单总数',
          number: '0',
          downTrend: '0%',
          leftType: 'icon-file-paste',
        },
      ];
    }
  } catch (error) {
    console.error('获取统计数据失败:', error);
    MessagePlugin.error('获取统计数据失败');
    // 错误时使用默认数据
    statisticsData.value = [
      {
        title: '总收入',
        number: '¥ 0.00',
        upTrend: '0%',
        leftType: 'echarts-line',
      },
      {
        title: '总支出',
        number: '¥ 0.00',
        downTrend: '0%',
        leftType: 'echarts-bar',
      },
      {
        title: '用户总数',
        number: '0',
        upTrend: '0%',
        leftType: 'icon-usergroup',
      },
      {
        title: '订单总数',
        number: '0',
        downTrend: '0%',
        leftType: 'icon-file-paste',
      },
    ];
  } finally {
    statisticsLoading.value = false;
  }
};
// ===== API数据替换结束 =====

// moneyCharts
let moneyContainer: HTMLElement;
let moneyChart: echarts.ECharts;
const renderMoneyChart = () => {
  if (!moneyContainer) {
    moneyContainer = document.getElementById('moneyContainer');
  }
  moneyChart = echarts.init(moneyContainer);
  moneyChart.setOption(constructInitDashboardDataset('line'));
};

// refundCharts
let refundContainer: HTMLElement;
let refundChart: echarts.ECharts;
const renderRefundChart = () => {
  if (!refundContainer) {
    refundContainer = document.getElementById('refundContainer');
  }
  refundChart = echarts.init(refundContainer);
  refundChart.setOption(constructInitDashboardDataset('bar'));
};

const renderCharts = () => {
  renderMoneyChart();
  renderRefundChart();
};

// chartSize update
const updateContainer = () => {
  if (document.documentElement.clientWidth >= 1400 && document.documentElement.clientWidth < 1920) {
    resizeTime.value = Number((document.documentElement.clientWidth / 2080).toFixed(2));
  } else if (document.documentElement.clientWidth < 1080) {
    resizeTime.value = Number((document.documentElement.clientWidth / 1080).toFixed(2));
  } else {
    resizeTime.value = 1;
  }
  if (moneyChart) {
    moneyChart.resize({
      width: resizeTime.value * 120,
    });
  }
  if (refundChart) {
    refundChart.resize({
      width: resizeTime.value * 120,
    });
  }
};

onMounted(async () => {
  // 先获取统计数据
  await fetchStatisticsData();

  // 然后渲染图表
  await nextTick();
  renderCharts();
  updateContainer();
});

const { width, height } = useWindowSize();
watch([width, height], () => {
  updateContainer();
});

watch(
    () => store.brandTheme,
    () => {
      if (refundChart) {
        changeChartsTheme([refundChart]);
      }
    },
);

watch(
    () => store.mode,
    () => {
      [moneyChart, refundChart].forEach((item) => {
        if (item) {
          item.dispose();
        }
      });

      renderCharts();
    },
);
</script>

<style lang="less" scoped>
.dashboard-item {
  padding: var(--td-comp-paddingTB-xl) var(--td-comp-paddingLR-xxl);

  :deep(.t-card__header) {
    padding: 0;
  }

  :deep(.t-card__footer) {
    padding: 0;
  }

  :deep(.t-card__title) {
    font: var(--td-font-body-medium);
    color: var(--td-text-color-secondary);
  }

  :deep(.t-card__body) {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    flex: 1;
    position: relative;
    padding: 0;
    margin-top: var(--td-comp-margin-s);
    margin-bottom: var(--td-comp-margin-xxl);
  }

  &:hover {
    cursor: pointer;
  }

  &-top {
    display: flex;
    flex-direction: row;
    align-items: flex-start;

    > span {
      display: inline-block;
      color: var(--td-text-color-primary);
      font-size: var(--td-font-size-headline-medium);
      line-height: var(--td-line-height-headline-medium);
    }
  }

  &-bottom {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;

    > .t-icon {
      cursor: pointer;
      font-size: var(--td-comp-size-xxxs);
    }
  }

  &-block {
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--td-text-color-placeholder);
  }

  &-trend {
    margin-left: var(--td-comp-margin-s);
  }

  &-left {
    position: absolute;
    top: 0;
    right: 0;

    > span {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: var(--td-comp-size-xxxl);
      height: var(--td-comp-size-xxxl);
      background: var(--td-brand-color-light);
      border-radius: 50%;

      .t-icon {
        font-size: 24px;
        color: var(--td-brand-color);
      }
    }
  }

  // 针对第一个卡片需要反色处理
  &--main-color {
    background: var(--td-brand-color);
    color: var(--td-text-color-primary);

    :deep(.t-card__title),
    .dashboard-item-top span,
    .dashboard-item-bottom {
      color: var(--td-text-color-anti);
    }

    .dashboard-item-block {
      color: var(--td-text-color-anti);
      opacity: 0.6;
    }

    .dashboard-item-bottom {
      color: var(--td-text-color-anti);
    }
  }
}
</style>