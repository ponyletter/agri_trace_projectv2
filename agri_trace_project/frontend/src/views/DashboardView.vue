<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <el-row :gutter="16" class="stat-cards">
      <el-col :span="6" v-for="card in statCards" :key="card.key">
        <div class="stat-card" :style="{ borderTop: `4px solid ${card.color}` }">
          <div class="stat-icon" :style="{ background: card.color + '20', color: card.color }">
            <el-icon :size="28"><component :is="card.icon" /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ card.value }}</div>
            <div class="stat-label">{{ card.label }}</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 图表区 -->
    <el-row :gutter="16" style="margin-top:16px">
      <el-col :span="16">
        <el-card class="chart-card">
          <template #header>
            <span class="card-title">📈 近15日溯源查询趋势</span>
          </template>
          <div ref="queryChartRef" class="chart-box"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="chart-card">
          <template #header>
            <span class="card-title">🔗 区块链节点状态</span>
          </template>
          <div class="node-list">
            <div v-for="node in blockNodes" :key="node.name" class="node-item">
              <el-badge is-dot type="success" />
              <span class="node-name">{{ node.name }}</span>
              <el-tag size="small" type="success">在线</el-tag>
            </div>
          </div>
          <el-divider />
          <div class="block-info">
            <div class="bi-item"><span>当前区块高度</span><strong>2,156</strong></div>
            <div class="bi-item"><span>共识算法</span><strong>Raft</strong></div>
            <div class="bi-item"><span>加密算法</span><strong>国密SM2/SM3</strong></div>
            <div class="bi-item"><span>累计交易数</span><strong>15,890</strong></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top:16px">
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header><span class="card-title">📦 批次状态分布</span></template>
          <div ref="statusChartRef" class="chart-box" style="height:220px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header><span class="card-title">💰 近15日交易金额趋势（万元）</span></template>
          <div ref="amountChartRef" class="chart-box" style="height:220px"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 操作日志 -->
    <el-card style="margin-top:16px">
      <template #header><span class="card-title">📋 最新操作日志</span></template>
      <el-table :data="logs" size="small" stripe>
        <el-table-column prop="username" label="操作人" width="120" />
        <el-table-column prop="action" label="操作" width="120" />
        <el-table-column prop="module" label="模块" width="100" />
        <el-table-column prop="detail" label="详情" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="160" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts'
import { getDashboardStats } from '../api/index'
import { Box, User, View, Connection } from '@element-plus/icons-vue'

const queryChartRef = ref<HTMLElement>()
const statusChartRef = ref<HTMLElement>()
const amountChartRef = ref<HTMLElement>()

const statCards = ref([
  { key: 'batches', label: '累计批次数', value: '5', color: '#1677ff', icon: Box },
  { key: 'users', label: '注册用户数', value: '10', color: '#52c41a', icon: User },
  { key: 'visitors', label: '今日访客', value: '534', color: '#fa8c16', icon: View },
  { key: 'tx', label: '区块链交易数', value: '15,890', color: '#722ed1', icon: Connection },
])

const blockNodes = [
  { name: 'peer0.org1（种植户节点）' },
  { name: 'peer0.org2（质检机构节点）' },
  { name: 'peer0.org3（物流商节点）' },
  { name: 'orderer（排序节点）' },
]

const logs = ref([
  { username: 'farmer_wang', action: '创建批次', module: 'batch', detail: '创建批次 BATCH-APPLE-001', status: 'success', created_at: '2025-11-15 09:12:34' },
  { username: 'inspector_zhang', action: '添加质检记录', module: 'inspection', detail: '为批次1出具质检报告', status: 'success', created_at: '2025-11-14 14:30:22' },
  { username: 'transporter_sun', action: '添加物流记录', module: 'logistics', detail: '批次1开始运输，目的地上海', status: 'success', created_at: '2025-11-13 08:05:11' },
  { username: 'admin', action: '生成合格证', module: 'cert', detail: '为批次1生成电子合格证', status: 'success', created_at: '2025-11-12 16:22:45' },
  { username: 'retailer_zhao', action: '查询溯源', module: 'trace', detail: '消费者扫码查询 AKS2025100001', status: 'success', created_at: '2025-11-11 11:08:30' },
])

// 默认图表数据（后端无数据时使用）
const defaultDates = ['11-01','11-02','11-03','11-04','11-05','11-06','11-07','11-08','11-09','11-10','11-11','11-12','11-13','11-14','11-15']
const defaultQueries = [45,62,38,75,88,52,94,67,83,71,95,108,89,76,112]
const defaultAmounts = [12.5,18.3,9.8,22.1,31.4,15.6,28.9,19.7,25.3,21.8,34.5,41.2,29.6,23.4,38.7]

onMounted(async () => {
  let dates = defaultDates
  let queries = defaultQueries
  let amounts = defaultAmounts

  try {
    const res: any = await getDashboardStats()
    const statsData = res.data || []
    if (statsData.length > 0) {
      dates = statsData.map((s: any) => s.stat_date?.slice(5))
      queries = statsData.map((s: any) => s.total_queries)
      amounts = statsData.map((s: any) => s.total_amount)
    }
  } catch (e) {
    // 使用默认数据
  }

  // 查询趋势图
  if (queryChartRef.value) {
    const chart = echarts.init(queryChartRef.value)
    chart.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: 40, right: 20, top: 20, bottom: 30 },
      xAxis: { type: 'category', data: dates, axisLabel: { fontSize: 11 } },
      yAxis: { type: 'value', name: '次' },
      series: [{
        name: '溯源查询次数',
        type: 'line',
        data: queries,
        smooth: true,
        areaStyle: { color: 'rgba(22,119,255,0.1)' },
        lineStyle: { color: '#1677ff' },
        itemStyle: { color: '#1677ff' },
      }]
    })
  }

  // 交易金额趋势
  if (amountChartRef.value) {
    const chart2 = echarts.init(amountChartRef.value)
    chart2.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: 50, right: 20, top: 20, bottom: 30 },
      xAxis: { type: 'category', data: dates, axisLabel: { fontSize: 11 } },
      yAxis: { type: 'value', name: '万元' },
      series: [{
        name: '交易金额(万元)',
        type: 'bar',
        data: amounts,
        itemStyle: { color: '#52c41a', borderRadius: [4,4,0,0] },
      }]
    })
  }

  // 批次状态饼图
  if (statusChartRef.value) {
    const chart3 = echarts.init(statusChartRef.value)
    chart3.setOption({
      tooltip: { trigger: 'item' },
      legend: { bottom: 0 },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        data: [
          { value: 2, name: '已上架', itemStyle: { color: '#52c41a' } },
          { value: 2, name: '运输中', itemStyle: { color: '#1677ff' } },
          { value: 1, name: '种植中', itemStyle: { color: '#fa8c16' } },
        ],
        label: { show: true, formatter: '{b}: {c}批' }
      }]
    })
  }
})
</script>

<style scoped>
.dashboard { padding: 0; }
.stat-cards { margin-bottom: 0; }
.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 1px 4px rgba(0,0,0,.06);
}
.stat-icon {
  width: 56px; height: 56px; border-radius: 12px;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.stat-value { font-size: 28px; font-weight: 700; color: #1a1a1a; line-height: 1; }
.stat-label { font-size: 13px; color: #888; margin-top: 4px; }
.chart-card { border-radius: 8px; }
.card-title { font-size: 14px; font-weight: 600; color: #333; }
.chart-box { height: 280px; }
.node-list { display: flex; flex-direction: column; gap: 10px; }
.node-item { display: flex; align-items: center; gap: 8px; }
.node-name { flex: 1; font-size: 13px; color: #555; }
.block-info { display: flex; flex-direction: column; gap: 8px; }
.bi-item { display: flex; justify-content: space-between; font-size: 13px; color: #666; }
.bi-item strong { color: #333; }
</style>
