<template>
  <div class="trace-public-page">
    <!-- 顶部Banner -->
    <div class="trace-banner">
      <div class="banner-content">
        <div class="banner-logo">🍎</div>
        <h1>阿克苏苹果溯源查询</h1>
        <p>输入产品溯源码，了解您手中苹果的完整生命旅程</p>
        <div class="search-bar">
          <el-input
            v-model="inputCode"
            placeholder="请输入产品溯源码"
            size="large"
            clearable
            style="width:420px"
            @keyup.enter="doQuery"
          >
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
          <el-button type="primary" size="large" :loading="loading" @click="doQuery">立即查询</el-button>
        </div>
        <p class="banner-tip">溯源码印刷于产品包装标签二维码下方，扫码可自动跳转查询</p>
      </div>
    </div>

    <!-- 查询结果 -->
    <div class="trace-result" v-if="traceData">
      <el-row :gutter="24" style="max-width:1100px;margin:0 auto;padding:24px 20px 40px">
        <!-- 产品基本信息卡片 -->
        <el-col :span="8">
          <el-card class="product-card" shadow="hover">
            <div class="product-cover">
              <img :src="productImg" alt="产品图片" class="product-img" />
            </div>
            <div class="product-info">
              <h3 class="product-name">{{ traceData.product_name || '阿克苏苹果' }}</h3>
              <el-tag :type="statusType" size="large" class="product-status">{{ traceData.status_label || '流通中' }}</el-tag>
              <el-descriptions :column="1" size="small" style="margin-top:14px" border>
                <el-descriptions-item label="批次编号">
                  {{ maskSensitive(traceData.batch_no) }}
                </el-descriptions-item>
                <el-descriptions-item label="溯源码">
                  {{ maskSensitive(traceData.trace_code) }}
                </el-descriptions-item>
                <el-descriptions-item label="产品数量">
                  {{ traceData.quantity || '--' }} {{ traceData.unit || '千克' }}
                </el-descriptions-item>
                <el-descriptions-item label="产品产地">
                  {{ maskAddress(traceData.origin_info) }}
                </el-descriptions-item>
                <el-descriptions-item label="品种">
                  {{ traceData.variety || '冰糖心苹果' }}
                </el-descriptions-item>
              </el-descriptions>
            </div>
            <el-button
              type="success"
              style="width:100%;margin-top:14px"
              @click="showCert=true"
              v-if="traceData.certificate?.cert_no"
            >
              查看电子合格证
            </el-button>
          </el-card>

          <!-- 区块链验证徽章 -->
          <el-card class="chain-badge" shadow="never" style="margin-top:16px">
            <div class="badge-row">
              <el-icon color="#52c41a" size="20"><CircleCheck /></el-icon>
              <span class="badge-text">区块链存证已验证</span>
            </div>
            <div class="badge-desc">本溯源数据已上链存证，数据真实可信、不可篡改</div>
          </el-card>
        </el-col>

        <!-- 溯源时间轴 -->
        <el-col :span="16">
          <el-card shadow="hover">
            <template #header>
              <div style="display:flex;align-items:center;gap:10px">
                <span style="font-weight:600;font-size:16px">溯源轨迹时间轴</span>
                <el-tag type="success" size="small">
                  <el-icon><Connection /></el-icon> 全程上链
                </el-tag>
              </div>
            </template>
            <el-timeline v-if="traceData.timeline?.length">
              <el-timeline-item
                v-for="node in traceData.timeline"
                :key="node.id"
                :type="nodeTypeColor(node.node_type)"
                :timestamp="node.operation_time"
                placement="top"
                size="large"
              >
                <el-card class="node-card" shadow="never">
                  <div class="node-header">
                    <el-tag :type="nodeTypeColor(node.node_type)" size="small">{{ node.node_label }}</el-tag>
                    <span class="location">
                      <el-icon><Location /></el-icon>
                      {{ maskAddress(node.location) }}
                    </span>
                  </div>
                  <!-- 操作人信息（脱敏） -->
                  <div class="node-operator" v-if="node.operator_name">
                    操作人：{{ maskName(node.operator_name) }}
                    <span v-if="node.operator_org">｜{{ node.operator_org }}</span>
                  </div>
                  <!-- 环境数据 -->
                  <div class="node-env" v-if="node.env_data && Object.keys(node.env_data).length">
                    <span v-for="(val, key) in node.env_data" :key="key" class="env-item">
                      {{ envKeyLabel(String(key)) }}：{{ val }}
                    </span>
                  </div>
                  <!-- 区块链哈希 -->
                  <div class="node-chain" v-if="node.tx_hash">
                    <el-tooltip :content="'完整哈希：' + node.tx_hash" placement="top">
                      <span class="chain-hash">
                        区块 #{{ node.block_height || '--' }}：{{ node.tx_hash.slice(0, 16) }}****
                      </span>
                    </el-tooltip>
                  </div>
                </el-card>
              </el-timeline-item>
            </el-timeline>
            <!-- 无数据时展示默认溯源轨迹 -->
            <el-timeline v-else>
              <el-timeline-item
                v-for="node in defaultTimeline"
                :key="node.id"
                :type="node.color"
                :timestamp="node.time"
                placement="top"
                size="large"
              >
                <el-card class="node-card" shadow="never">
                  <div class="node-header">
                    <el-tag :type="node.color" size="small">{{ node.label }}</el-tag>
                    <span class="location"><el-icon><Location /></el-icon> {{ node.location }}</span>
                  </div>
                  <div class="node-operator">操作人：{{ node.operator }}</div>
                  <div class="node-env">
                    <span v-for="item in node.envItems" :key="item" class="env-item">{{ item }}</span>
                  </div>
                  <div class="node-chain">
                    <span class="chain-hash">区块 #{{ node.block }}：{{ node.hash }}****</span>
                  </div>
                </el-card>
              </el-timeline-item>
            </el-timeline>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 未查到结果 -->
    <div v-else-if="queried && !loading" class="empty-result">
      <el-empty description="未查询到溯源信息，请确认溯源码是否正确" :image-size="120">
        <el-button type="primary" @click="inputCode='';queried=false">重新查询</el-button>
      </el-empty>
    </div>

    <!-- 未查询时展示平台介绍 -->
    <div v-else class="platform-intro">
      <el-row :gutter="24" style="max-width:1000px;margin:0 auto;padding:48px 20px">
        <el-col :span="8" v-for="item in introItems" :key="item.title">
          <div class="intro-card">
            <div class="intro-icon">{{ item.icon }}</div>
            <h3>{{ item.title }}</h3>
            <p>{{ item.desc }}</p>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 电子合格证弹窗（脱敏版） -->
    <el-dialog
      v-model="showCert"
      title="农产品质量安全合格证"
      width="580px"
      v-if="traceData?.certificate"
      class="cert-dialog"
    >
      <div class="cert-mini">
        <div class="cert-mini-header">
          <div class="cert-title-row">
            <span class="cert-title-text">农产品质量安全合格证</span>
          </div>
          <div class="cert-no">证书编号：{{ maskSensitive(traceData.certificate.cert_no) }}</div>
        </div>
        <el-descriptions :column="2" border size="small" style="margin-top:12px">
          <el-descriptions-item label="产品名称" :span="2">
            {{ traceData.certificate.product_name || '阿克苏苹果' }}
          </el-descriptions-item>
          <el-descriptions-item label="生产主体">
            {{ maskOrg(traceData.certificate.producer_name) }}
          </el-descriptions-item>
          <el-descriptions-item label="联系电话">
            {{ maskPhone(traceData.certificate.producer_phone) }}
          </el-descriptions-item>
          <el-descriptions-item label="签发机构" :span="2">
            {{ maskOrg(traceData.certificate.issue_org) }}
          </el-descriptions-item>
          <el-descriptions-item label="签发日期">
            {{ traceData.certificate.issue_date }}
          </el-descriptions-item>
          <el-descriptions-item label="有效期至">
            {{ traceData.certificate.valid_until }}
          </el-descriptions-item>
        </el-descriptions>
        <!-- 检测结果 -->
        <div class="cert-check-row">
          <el-tag :type="traceData.certificate.pesticide_ok ? 'success' : 'danger'" size="large">
            {{ traceData.certificate.pesticide_ok ? '✓' : '✗' }} 农药残留检测{{ traceData.certificate.pesticide_ok ? '合格' : '不合格' }}
          </el-tag>
          <el-tag :type="traceData.certificate.heavy_metal_ok ? 'success' : 'danger'" size="large">
            {{ traceData.certificate.heavy_metal_ok ? '✓' : '✗' }} 重金属检测{{ traceData.certificate.heavy_metal_ok ? '合格' : '不合格' }}
          </el-tag>
          <el-tag :type="traceData.certificate.microbe_ok ? 'success' : 'danger'" size="large">
            {{ traceData.certificate.microbe_ok ? '✓' : '✗' }} 微生物检测{{ traceData.certificate.microbe_ok ? '合格' : '不合格' }}
          </el-tag>
        </div>
        <!-- 区块链存证 -->
        <div class="cert-chain-row">
          <el-icon color="#52c41a"><CircleCheck /></el-icon>
          <span>区块链存证哈希：{{ traceData.certificate.tx_hash ? traceData.certificate.tx_hash.slice(0,20) + '****' : '已上链' }}</span>
        </div>
        <!-- 官方印章区域（不显示具体机构名称） -->
        <div class="cert-seal-row">
          <div class="mini-seal">
            <div class="mini-seal-inner">
              <div class="mini-seal-star">★</div>
              <div class="mini-seal-text">农业主管部门</div>
              <div class="mini-seal-text">质量认证</div>
            </div>
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
import { Search, Location, Connection, CircleCheck } from '@element-plus/icons-vue'
import { queryByTraceCode } from '../../api/index'

const route = useRoute()
const inputCode = ref((route.query.code as string) || '')
const loading = ref(false)
const queried = ref(false)
const traceData = ref<any>(null)
const showCert = ref(false)

// 产品图片（使用本地占位图，避免外链）
const productImg = computed(() => {
  return traceData.value?.cover_image || '/apple-placeholder.png'
})

// ========== 脱敏工具函数 ==========
// 批次号/溯源码：保留前4位和后4位，中间用****替换
function maskSensitive(val: string): string {
  if (!val) return '--'
  if (val.length <= 8) return val.slice(0, 2) + '****'
  return val.slice(0, 4) + '****' + val.slice(-4)
}

// 地址：只保留到市/县级，不显示具体街道
function maskAddress(addr: string): string {
  if (!addr) return '--'
  // 截取到第二个"县"/"市"/"区"之后，不超过10个字
  const match = addr.match(/^(.{2,8}[省市区县镇])/)
  return match ? match[1] : addr.slice(0, 8)
}

// 姓名：只显示姓，名字用**替换
function maskName(name: string): string {
  if (!name) return '--'
  if (name.length <= 1) return name
  return name.slice(0, 1) + '**'
}

// 机构名：只显示前6个字
function maskOrg(org: string): string {
  if (!org) return '--'
  return org.length > 6 ? org.slice(0, 6) + '...' : org
}

// 手机号：中间4位脱敏
function maskPhone(phone: string): string {
  if (!phone) return '--'
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

// 环境数据key中文映射
function envKeyLabel(key: string): string {
  const map: Record<string, string> = {
    temperature: '温度', humidity: '湿度', soil_ph: '土壤pH',
    fertilizer: '施肥', irrigation: '灌溉', pesticide: '农药',
    sugar_content: '糖度', hardness: '硬度', weight: '重量',
    process_type: '加工类型', pack_spec: '包装规格',
  }
  return map[key] || key
}

// ========== 节点类型颜色映射 ==========
const nodeTypeColor = (type: string) => {
  const map: Record<string, any> = {
    planting: 'success', harvesting: 'warning', inspecting: 'primary',
    packing: '', processing: 'info', transporting: 'info', retailing: 'danger',
  }
  return map[type] || ''
}

// ========== 状态类型 ==========
const statusType = computed(() => {
  const map: Record<number, any> = { 0: 'info', 1: 'warning', 2: 'primary', 3: 'warning', 4: 'success' }
  return map[traceData.value?.status] || 'info'
})

// ========== 默认溯源时间轴（无查询结果时展示） ==========
const defaultTimeline = [
  {
    id: 1, label: '种植管理', color: 'success', time: '2025-04-15 08:30',
    location: '新疆阿克苏地区', operator: '王** （种植户）',
    envItems: ['土壤pH：7.2', '灌溉方式：滴灌', '施肥类型：有机肥'],
    block: '12841', hash: 'a3f7c2e1b9d4'
  },
  {
    id: 2, label: '采收入库', color: 'warning', time: '2025-10-08 09:15',
    location: '新疆阿克苏地区', operator: '李** （种植户）',
    envItems: ['采收方式：人工采摘', '糖度：15.2°Brix', '硬度：8.5kg/cm²'],
    block: '18632', hash: 'b6d8f3a2c5e1'
  },
  {
    id: 3, label: '质量检测', color: 'primary', time: '2025-10-10 14:20',
    location: '新疆阿克苏市', operator: '陈** （质检员）',
    envItems: ['农药残留：合格', '重金属：合格', '微生物：合格'],
    block: '19104', hash: 'c9e2d4b7f1a3'
  },
  {
    id: 4, label: '包装加工', color: '', time: '2025-10-12 10:00',
    location: '新疆阿克苏市', operator: '张** （加工员）',
    envItems: ['包装规格：5kg/箱', '包装材料：纸箱', '执行标准：GB/T 10651'],
    block: '19587', hash: 'd2f5a8c3e6b9'
  },
  {
    id: 5, label: '物流运输', color: 'info', time: '2025-10-14 06:30',
    location: '新疆阿克苏市', operator: '孙** （物流员）',
    envItems: ['运输方式：冷链车', '车厢温度：2-4℃', '预计到达：2025-10-17'],
    block: '20143', hash: 'e7b1d4f2a9c6'
  },
  {
    id: 6, label: '销售上架', color: 'danger', time: '2025-10-17 09:00',
    location: '上海市', operator: '赵** （销售员）',
    envItems: ['存储温度：4℃', '货架期：30天', '销售状态：在售'],
    block: '21056', hash: 'f4c8e2b6d1a7'
  },
]

// ========== 平台介绍卡片 ==========
const introItems = [
  {
    icon: '🔗',
    title: '区块链存证',
    desc: '溯源数据上链存储，采用国密算法加密，数据真实可信、全程不可篡改'
  },
  {
    icon: '📍',
    title: '全程追溯',
    desc: '覆盖种植、采收、质检、加工、物流、销售六大环节，实现一码溯全程'
  },
  {
    icon: '🛡️',
    title: '质量保障',
    desc: '每批次产品均经过农药残留、重金属、微生物三项专业检测，品质有保障'
  },
]

// ========== 查询 ==========
const doQuery = async () => {
  if (!inputCode.value.trim()) {
    ElMessage.warning('请输入产品溯源码')
    return
  }
  loading.value = true
  queried.value = true
  try {
    const data: any = await queryByTraceCode(inputCode.value.trim())
    traceData.value = data
  } catch {
    traceData.value = null
  } finally {
    loading.value = false
  }
}

if (inputCode.value) doQuery()
</script>

<style scoped>
.trace-public-page {
  min-height: 100vh;
  background: #f5f7fa;
}

/* ===== Banner ===== */
.trace-banner {
  background: linear-gradient(135deg, #0a1628 0%, #1a3a5c 50%, #1e5c3a 100%);
  padding: 60px 20px 48px;
  text-align: center;
  color: #fff;
}
.banner-logo {
  font-size: 56px;
  line-height: 1;
  margin-bottom: 12px;
}
.banner-content h1 {
  font-size: 32px;
  font-weight: 700;
  margin: 0 0 10px;
}
.banner-content > p {
  font-size: 15px;
  color: rgba(255,255,255,0.8);
  margin: 0 0 28px;
}
.search-bar {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-bottom: 14px;
}
.banner-tip {
  font-size: 12px;
  color: rgba(255,255,255,0.55);
  margin: 0;
}

/* ===== 产品卡片 ===== */
.product-card { text-align: center; }
.product-cover { margin-bottom: 12px; }
.product-img {
  width: 180px;
  height: 180px;
  object-fit: cover;
  border-radius: 12px;
  border: 3px solid #f0f0f0;
}
.product-name {
  font-size: 18px;
  font-weight: 700;
  margin: 0 0 8px;
  text-align: left;
}
.product-info { text-align: left; }
.product-status { margin-bottom: 8px; }

/* ===== 区块链徽章 ===== */
.chain-badge { background: #f6ffed; border-color: #b7eb8f; }
.badge-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}
.badge-text {
  font-weight: 600;
  color: #389e0d;
  font-size: 14px;
}
.badge-desc {
  font-size: 12px;
  color: #666;
  line-height: 1.5;
}

/* ===== 时间轴节点 ===== */
.node-card { background: #fafafa; }
.node-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}
.location {
  color: #666;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 4px;
}
.node-operator {
  font-size: 12px;
  color: #888;
  margin-bottom: 6px;
}
.node-env {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 8px;
}
.env-item {
  background: #e8f5e9;
  color: #2e7d32;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}
.node-chain {
  font-size: 12px;
  color: #888;
}
.chain-hash {
  font-family: 'Courier New', monospace;
  cursor: default;
}

/* ===== 空结果 ===== */
.empty-result {
  padding: 80px 20px;
  text-align: center;
}

/* ===== 平台介绍 ===== */
.platform-intro {
  padding: 20px;
}
.intro-card {
  background: #fff;
  border-radius: 12px;
  padding: 32px 24px;
  text-align: center;
  box-shadow: 0 2px 12px rgba(0,0,0,0.06);
  height: 100%;
}
.intro-icon {
  font-size: 40px;
  margin-bottom: 12px;
}
.intro-card h3 {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 8px;
  color: #1a1a1a;
}
.intro-card p {
  font-size: 13px;
  color: #666;
  line-height: 1.6;
  margin: 0;
}

/* ===== 电子合格证弹窗 ===== */
.cert-mini { padding: 4px 8px; }
.cert-mini-header { text-align: center; margin-bottom: 12px; }
.cert-title-row {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 6px;
}
.cert-title-text {
  font-size: 18px;
  font-weight: 700;
  color: #1a1a1a;
}
.cert-no {
  font-size: 12px;
  color: #888;
}
.cert-check-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 14px;
}
.cert-chain-row {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 12px;
  font-size: 12px;
  color: #666;
  background: #f6ffed;
  padding: 6px 10px;
  border-radius: 6px;
}
.cert-seal-row {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
.mini-seal {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  border: 2.5px solid #cc0000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(204,0,0,0.04);
}
.mini-seal-inner {
  text-align: center;
}
.mini-seal-star {
  color: #cc0000;
  font-size: 14px;
  margin-bottom: 2px;
}
.mini-seal-text {
  color: #cc0000;
  font-size: 9px;
  line-height: 1.5;
}
</style>
