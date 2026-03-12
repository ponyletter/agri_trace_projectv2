<template>
  <div class="trace-public-page">
    <!-- 顶部Banner -->
    <div class="trace-banner">
      <div class="banner-content">
        <h1>🍎 阿克苏苹果溯源查询</h1>
        <p>输入溯源码，了解您手中苹果的完整生命旅程</p>
        <div class="search-bar">
          <el-input
            v-model="inputCode"
            placeholder="请输入溯源码，如：AKS2025100001"
            size="large"
            clearable
            style="width:400px"
            @keyup.enter="doQuery"
          >
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
          <el-button type="primary" size="large" :loading="loading" @click="doQuery">立即查询</el-button>
        </div>
        <div class="demo-codes">
          <span>示例溯源码：</span>
          <el-tag v-for="code in demoCodes" :key="code" @click="inputCode=code;doQuery()" style="cursor:pointer;margin:0 4px">{{ code }}</el-tag>
        </div>
      </div>
    </div>

    <!-- 查询结果 -->
    <div class="trace-result" v-if="traceData">
      <el-row :gutter="24" style="max-width:1100px;margin:0 auto;padding:24px 20px 0">
        <!-- 产品基本信息 -->
        <el-col :span="8">
          <el-card class="product-card">
            <div class="product-cover">
              <img :src="traceData.cover_image || defaultAppleImg" alt="产品图片" class="product-img" />
            </div>
            <div class="product-info">
              <h3 class="product-name">{{ traceData.product_name }}</h3>
              <el-tag :type="statusType" size="large" class="product-status">{{ traceData.status_label }}</el-tag>
              <el-descriptions :column="1" size="small" style="margin-top:12px">
                <el-descriptions-item label="批次号">{{ traceData.batch_no }}</el-descriptions-item>
                <el-descriptions-item label="溯源码">{{ traceData.trace_code }}</el-descriptions-item>
                <el-descriptions-item label="数量">{{ traceData.quantity }} {{ traceData.unit }}</el-descriptions-item>
                <el-descriptions-item label="产地">{{ traceData.origin_info }}</el-descriptions-item>
              </el-descriptions>
            </div>
            <el-button type="success" style="width:100%;margin-top:12px" @click="showCert=true" v-if="traceData.certificate?.cert_no">
              📜 查看电子合格证
            </el-button>
          </el-card>
        </el-col>

        <!-- 溯源时间轴 -->
        <el-col :span="16">
          <el-card>
            <template #header>
              <span style="font-weight:600;font-size:16px">📍 溯源轨迹时间轴</span>
              <el-tag type="success" style="margin-left:12px" size="small">
                <el-icon><Connection /></el-icon> 区块链已验证
              </el-tag>
            </template>
            <el-timeline>
              <el-timeline-item
                v-for="node in traceData.timeline"
                :key="node.id"
                :type="nodeTypeColor(node.node_type)"
                :timestamp="node.operation_time"
                placement="top"
              >
                <el-card class="node-card" shadow="never">
                  <div class="node-header">
                    <el-tag :type="nodeTypeColor(node.node_type)" size="small">{{ node.node_label }}</el-tag>
                    <span class="location"><el-icon><Location /></el-icon> {{ node.location }}</span>
                  </div>
                  <div class="node-env" v-if="node.env_data">
                    <span v-for="(val, key) in node.env_data" :key="key" class="env-item">
                      {{ key }}: {{ val }}
                    </span>
                  </div>
                  <div class="node-chain">
                    <el-tooltip :content="node.tx_hash" placement="top">
                      <span class="chain-hash">🔗 区块 #{{ node.block_height }}: {{ node.tx_hash?.slice(0, 20) }}...</span>
                    </el-tooltip>
                  </div>
                  <div v-if="node.ipfs_files?.length" class="ipfs-files">
                    <a v-for="f in node.ipfs_files" :key="f.cid" :href="f.url" target="_blank" class="ipfs-link">
                      📎 {{ f.file_name }}
                    </a>
                  </div>
                </el-card>
              </el-timeline-item>
            </el-timeline>
            <el-empty v-if="!traceData.timeline?.length" description="暂无溯源记录" />
          </el-card>
        </el-col>
      </el-row>
    </div>

    <el-empty v-else-if="queried && !loading" description="未找到溯源信息，请检查溯源码" style="margin-top:60px" />

    <!-- 电子合格证弹窗 -->
    <el-dialog v-model="showCert" title="电子合格证" width="600px" v-if="traceData?.certificate">
      <div class="cert-mini">
        <div class="cert-mini-header">
          <h3>农产品质量安全合格证</h3>
          <div class="cert-mini-no">证书编号：{{ traceData.certificate.cert_no }}</div>
        </div>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="产品名称" :span="2">{{ traceData.certificate.product_name }}</el-descriptions-item>
          <el-descriptions-item label="生产主体">{{ traceData.certificate.producer_name }}</el-descriptions-item>
          <el-descriptions-item label="联系电话">{{ traceData.certificate.producer_phone }}</el-descriptions-item>
          <el-descriptions-item label="签发机构" :span="2">{{ traceData.certificate.issue_org }}</el-descriptions-item>
          <el-descriptions-item label="签发日期">{{ traceData.certificate.issue_date }}</el-descriptions-item>
          <el-descriptions-item label="有效期至">{{ traceData.certificate.valid_until }}</el-descriptions-item>
        </el-descriptions>
        <div style="margin-top:12px;display:flex;flex-wrap:wrap;gap:8px">
          <el-tag :type="traceData.certificate.pesticide_ok ? 'success' : 'danger'">
            {{ traceData.certificate.pesticide_ok ? '✅' : '❌' }} 农药残留：{{ traceData.certificate.pesticide_ok ? '合格' : '不合格' }}
          </el-tag>
          <el-tag :type="traceData.certificate.heavy_metal_ok ? 'success' : 'danger'">
            {{ traceData.certificate.heavy_metal_ok ? '✅' : '❌' }} 重金属：{{ traceData.certificate.heavy_metal_ok ? '合格' : '不合格' }}
          </el-tag>
          <el-tag :type="traceData.certificate.microbe_ok ? 'success' : 'danger'">
            {{ traceData.certificate.microbe_ok ? '✅' : '❌' }} 微生物：{{ traceData.certificate.microbe_ok ? '合格' : '不合格' }}
          </el-tag>
        </div>
        <div class="cert-seal-row">
          <div class="mini-seal">
            <div class="mini-seal-text">新疆阿克苏<br/>农业农村局<br/>质量认证</div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, Location, Connection } from '@element-plus/icons-vue'
import { queryByTraceCode } from '../../api/index'

const route = useRoute()
const inputCode = ref((route.query.code as string) || '')
const loading = ref(false)
const queried = ref(false)
const traceData = ref<any>(null)
const showCert = ref(false)

const demoCodes = ['AKS2025100001', 'AKS2025100002', 'AKS2025100003']
const defaultAppleImg = 'https://images.unsplash.com/photo-1560806887-1e4cd0b6cbd6?w=300&h=300&fit=crop'

const nodeTypeColor = (type: string) => {
  const map: Record<string, any> = {
    planting: 'success', harvesting: 'warning', inspecting: 'primary',
    packing: '', processing: 'info', transporting: 'info', retailing: 'danger',
  }
  return map[type] || ''
}

const statusType = computed(() => {
  const map: Record<number, any> = { 0: 'info', 1: 'warning', 2: 'primary', 3: 'warning', 4: 'success' }
  return map[traceData.value?.status] || 'info'
})

const doQuery = async () => {
  if (!inputCode.value.trim()) { ElMessage.warning('请输入溯源码'); return }
  loading.value = true
  queried.value = true
  try {
    const res: any = await queryByTraceCode(inputCode.value.trim())
    traceData.value = res.data || res
  } catch {
    traceData.value = null
  } finally {
    loading.value = false
  }
}

if (inputCode.value) doQuery()
</script>

<style scoped>
.trace-public-page { min-height: 100vh; background: #f5f7fa; }
.trace-banner {
  background: linear-gradient(135deg, #0d1b2a 0%, #1a3a5c 50%, #2d6a4f 100%);
  padding: 60px 20px 40px; text-align: center; color: #fff;
}
.banner-content h1 { font-size: 32px; font-weight: 700; margin: 0 0 8px; }
.banner-content p { font-size: 16px; color: rgba(255,255,255,0.8); margin: 0 0 24px; }
.search-bar { display: flex; justify-content: center; gap: 12px; margin-bottom: 16px; }
.demo-codes { font-size: 13px; color: rgba(255,255,255,0.7); }
.product-card { text-align: center; }
.product-cover { margin-bottom: 12px; }
.product-img { width: 200px; height: 200px; object-fit: cover; border-radius: 8px; }
.product-name { font-size: 18px; font-weight: 700; margin: 0 0 8px; }
.product-info { text-align: left; }
.node-card { background: #fafafa; }
.node-header { display: flex; align-items: center; gap: 12px; margin-bottom: 8px; }
.location { color: #666; font-size: 13px; display: flex; align-items: center; gap: 4px; }
.node-env { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 8px; }
.env-item { background: #e8f5e9; color: #2e7d32; padding: 2px 8px; border-radius: 4px; font-size: 12px; }
.node-chain { font-size: 12px; color: #888; margin-bottom: 6px; }
.chain-hash { cursor: pointer; font-family: monospace; }
.ipfs-files { display: flex; flex-wrap: wrap; gap: 8px; }
.ipfs-link { font-size: 12px; color: #27ae60; text-decoration: none; }
.cert-mini { padding: 8px; }
.cert-mini-header { text-align: center; margin-bottom: 16px; }
.cert-mini-header h3 { margin: 0 0 4px; }
.cert-mini-no { font-size: 13px; color: #888; }
.cert-seal-row { display: flex; justify-content: flex-end; margin-top: 16px; }
.mini-seal {
  width: 80px; height: 80px; border-radius: 50%;
  border: 2px solid #cc0000; display: flex; align-items: center; justify-content: center;
  background: rgba(204,0,0,0.05);
}
.mini-seal-text { color: #cc0000; font-size: 9px; text-align: center; line-height: 1.4; }
</style>
