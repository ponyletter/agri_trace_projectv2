<template>
  <div>
    <el-card class="page-header-card">
      <div class="page-header">
        <div>
          <h3 class="page-title">🔗 区块链浏览器</h3>
          <p class="page-desc">展示 Hyperledger Fabric 国密版底层交易日志、区块哈希与联盟链节点架构</p>
        </div>
        <el-button type="primary" icon="Refresh" @click="loadBlockInfo">刷新</el-button>
      </div>
    </el-card>

    <!-- 区块链概览卡片 -->
    <el-row :gutter="16" style="margin-top:16px">
      <el-col :span="6" v-for="card in blockCards" :key="card.label">
        <div class="block-stat-card">
          <div class="bsc-icon" :style="{ color: card.color }">{{ card.icon }}</div>
          <div class="bsc-value" :style="{ color: card.color }">{{ card.value }}</div>
          <div class="bsc-label">{{ card.label }}</div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top:16px">
      <!-- 节点架构 -->
      <el-col :span="10">
        <el-card>
          <template #header><span style="font-weight:600">🏗️ 联盟链节点架构</span></template>
          <div class="node-arch">
            <div class="org-group" v-for="org in orgNodes" :key="org.name">
              <div class="org-header" :style="{ background: org.color }">{{ org.name }}</div>
              <div class="org-peers">
                <div class="peer-item" v-for="peer in org.peers" :key="peer.name">
                  <el-badge is-dot type="success" />
                  <div class="peer-info">
                    <div class="peer-name">{{ peer.name }}</div>
                    <div class="peer-role">{{ peer.role }}</div>
                  </div>
                </div>
              </div>
            </div>
            <!-- 排序节点 -->
            <div class="orderer-group">
              <div class="orderer-header">排序服务（Raft共识）</div>
              <div class="orderer-item">
                <el-badge is-dot type="success" />
                <span>orderer.example.com:7050</span>
              </div>
            </div>
          </div>
          <el-divider />
          <div class="crypto-info">
            <div class="ci-title">🔐 国密改造说明</div>
            <el-descriptions :column="1" size="small">
              <el-descriptions-item label="签名算法">SM2（替换 ECDSA P-256）</el-descriptions-item>
              <el-descriptions-item label="哈希算法">SM3（替换 SHA-256）</el-descriptions-item>
              <el-descriptions-item label="加密算法">SM4（替换 AES-256）</el-descriptions-item>
              <el-descriptions-item label="证书格式">GM/T 0015-2012 标准</el-descriptions-item>
              <el-descriptions-item label="底层库">github.com/tjfoc/gmsm</el-descriptions-item>
            </el-descriptions>
          </div>
        </el-card>
      </el-col>

      <!-- 最新区块列表 -->
      <el-col :span="14">
        <el-card>
          <template #header>
            <div style="display:flex;justify-content:space-between;align-items:center">
              <span style="font-weight:600">📦 最新区块列表</span>
              <el-tag type="success" size="small">Channel: agrichannel</el-tag>
            </div>
          </template>
          <el-table :data="blocks" stripe size="small" @row-click="viewBlock">
            <el-table-column prop="height" label="区块高度" width="90" />
            <el-table-column label="区块哈希" width="200">
              <template #default="{ row }">
                <el-tooltip :content="row.hash">
                  <span class="hash-text">{{ row.hash?.slice(0,20) }}...</span>
                </el-tooltip>
              </template>
            </el-table-column>
            <el-table-column prop="tx_count" label="交易数" width="70" align="center" />
            <el-table-column prop="data_hash" label="数据哈希(SM3)" width="160">
              <template #default="{ row }">
                <span class="hash-text">{{ row.data_hash?.slice(0,16) }}...</span>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="时间" width="140" />
          </el-table>
        </el-card>

        <!-- 交易详情弹窗触发区 -->
        <el-card style="margin-top:12px" v-if="currentBlock">
          <template #header>
            <span style="font-weight:600">🔍 区块 #{{ currentBlock.height }} 交易详情</span>
          </template>
          <el-table :data="currentBlock.txs" size="small" stripe>
            <el-table-column prop="tx_id" label="交易ID">
              <template #default="{ row }">
                <span class="hash-text">{{ row.tx_id?.slice(0,24) }}...</span>
              </template>
            </el-table-column>
            <el-table-column prop="chaincode" label="链码" width="120" />
            <el-table-column prop="func" label="调用函数" width="140" />
            <el-table-column prop="creator" label="发起者" width="120" />
            <el-table-column prop="status" label="状态" width="80" align="center">
              <template #default="{ row }">
                <el-tag type="success" size="small">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getBlockInfo } from '../../api/index'

const currentBlock = ref<any>(null)

const blockCards = ref([
  { label: '当前区块高度', value: '2,156', icon: '📦', color: '#1677ff' },
  { label: '累计交易数', value: '15,890', icon: '📋', color: '#52c41a' },
  { label: '加密算法', value: 'SM2/SM3', icon: '🔐', color: '#722ed1' },
  { label: '共识算法', value: 'Raft', icon: '⚡', color: '#fa8c16' },
])

const orgNodes = [
  {
    name: 'Org1（种植户组织）', color: '#e6f4ff',
    peers: [
      { name: 'peer0.org1.example.com:7051', role: '背书节点' },
    ]
  },
  {
    name: 'Org2（质检机构组织）', color: '#f6ffed',
    peers: [
      { name: 'peer0.org2.example.com:9051', role: '背书节点' },
    ]
  },
  {
    name: 'Org3（物流商组织）', color: '#fff7e6',
    peers: [
      { name: 'peer0.org3.example.com:11051', role: '背书节点' },
    ]
  },
]

// 模拟区块数据
const blocks = ref([
  { height: 2156, hash: '0x8b3a4f9e2d1c5b7a6f8e9d0c1b2a3f4e5d6c7b8a9f0e1d2c3b4a5f6e7d8c9b0a', tx_count: 3, data_hash: '0xSM3_a1b2c3d4e5f6', created_at: '2025-11-15 14:30:22' },
  { height: 2155, hash: '0x7a2b3e8d1c0f5a6b4c9d8e7f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b', tx_count: 2, data_hash: '0xSM3_b2c3d4e5f6a1', created_at: '2025-11-15 14:28:10' },
  { height: 2154, hash: '0x6c1d4f7e0b9a3c5d2e8f1a4b7c0d3e6f9a2b5c8d1e4f7a0b3c6d9e2f5a8b1c4d', tx_count: 1, data_hash: '0xSM3_c3d4e5f6a1b2', created_at: '2025-11-15 14:25:55' },
  { height: 2153, hash: '0x5b0c3e6d9f2a5b8c1d4e7f0a3b6c9d2e5f8a1b4c7d0e3f6a9b2c5d8e1f4a7b0c', tx_count: 4, data_hash: '0xSM3_d4e5f6a1b2c3', created_at: '2025-11-15 14:20:33' },
  { height: 2152, hash: '0x4a9b2e5d8f1a4b7c0d3e6f9a2b5c8d1e4f7a0b3c6d9e2f5a8b1c4d7e0f3a6b9c', tx_count: 2, data_hash: '0xSM3_e5f6a1b2c3d4', created_at: '2025-11-15 14:15:18' },
])

function viewBlock(row: any) {
  currentBlock.value = {
    height: row.height,
    txs: [
      { tx_id: `0x${row.hash?.slice(2,34)}`, chaincode: 'agri-trace', func: 'AddTraceRecord', creator: 'farmer_wang@Org1', status: 'VALID' },
      { tx_id: `0x${row.hash?.slice(4,36)}`, chaincode: 'agri-trace', func: 'UpdateBatchStatus', creator: 'inspector_zhang@Org2', status: 'VALID' },
    ]
  }
}

async function loadBlockInfo() {
  try {
    const res: any = await getBlockInfo()
    const info = res.data || res
    blockCards.value[0].value = info.height?.toLocaleString() || '2,156'
    blockCards.value[1].value = info.total_tx?.toLocaleString() || '15,890'
  } catch {}
}

onMounted(loadBlockInfo)
</script>

<style scoped>
.page-header-card { margin-bottom: 0; }
.page-header { display: flex; justify-content: space-between; align-items: center; }
.page-title { margin: 0; font-size: 18px; font-weight: 600; }
.page-desc { margin: 4px 0 0; color: #888; font-size: 13px; }
.hash-text { font-family: monospace; font-size: 12px; color: #1677ff; cursor: pointer; }
.block-stat-card {
  background: #fff; border-radius: 8px; padding: 20px; text-align: center;
  box-shadow: 0 1px 4px rgba(0,0,0,.06);
}
.bsc-icon { font-size: 28px; margin-bottom: 4px; }
.bsc-value { font-size: 24px; font-weight: 700; }
.bsc-label { font-size: 13px; color: #888; margin-top: 4px; }
.node-arch { display: flex; flex-direction: column; gap: 10px; }
.org-group { border-radius: 6px; overflow: hidden; border: 1px solid #e8e8e8; }
.org-header { padding: 6px 12px; font-size: 13px; font-weight: 600; color: #333; }
.org-peers { padding: 8px 12px; display: flex; flex-direction: column; gap: 6px; }
.peer-item { display: flex; align-items: center; gap: 8px; }
.peer-info { flex: 1; }
.peer-name { font-size: 12px; color: #333; font-family: monospace; }
.peer-role { font-size: 11px; color: #888; }
.orderer-group { background: #fff7e6; border-radius: 6px; padding: 8px 12px; border: 1px solid #ffd591; }
.orderer-header { font-size: 13px; font-weight: 600; color: #fa8c16; margin-bottom: 6px; }
.orderer-item { display: flex; align-items: center; gap: 8px; font-size: 12px; font-family: monospace; }
.crypto-info { margin-top: 4px; }
.ci-title { font-size: 13px; font-weight: 600; color: #722ed1; margin-bottom: 8px; }
</style>
