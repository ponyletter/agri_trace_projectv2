<template>
  <div>
    <el-card class="page-header-card">
      <div class="page-header">
        <div>
          <h3 class="page-title">🚛 物流追踪定位</h3>
          <p class="page-desc">记录产品流转部门、时间及地理位置，可视化物流轨迹</p>
        </div>
        <el-button type="primary" icon="Plus" @click="showDialog = true">新增物流记录</el-button>
      </div>
    </el-card>

    <el-row :gutter="16" style="margin-top:16px">
      <!-- 物流记录列表 -->
      <el-col :span="14">
        <el-card>
          <template #header><span style="font-weight:600">📋 物流记录列表</span></template>
          <el-table :data="records" stripe v-loading="loading" @row-click="selectRecord" highlight-current-row>
            <el-table-column prop="batch_id" label="批次ID" width="80" />
            <el-table-column label="物流状态" width="100">
              <template #default="{ row }">
                <el-tag type="primary">{{ row.env_data?.logistics_status || '运输中' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="location" label="当前位置" />
            <el-table-column label="运输公司" width="120">
              <template #default="{ row }">{{ row.env_data?.company || '-' }}</template>
            </el-table-column>
            <el-table-column label="目的地" width="120">
              <template #default="{ row }">{{ row.env_data?.destination || '-' }}</template>
            </el-table-column>
            <el-table-column prop="operation_time" label="时间" width="140">
              <template #default="{ row }">{{ formatTime(row.operation_time) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- 地图可视化 -->
      <el-col :span="10">
        <el-card>
          <template #header><span style="font-weight:600">🗺️ 物流轨迹地图</span></template>
          <div class="map-container">
            <div class="map-placeholder">
              <div class="map-route">
                <div v-for="(point, idx) in routePoints" :key="idx" class="route-point">
                  <div class="point-dot" :class="{ active: idx === selectedIdx }">
                    <span class="point-num">{{ idx + 1 }}</span>
                  </div>
                  <div class="point-info">
                    <div class="point-name">{{ point.name }}</div>
                    <div class="point-time">{{ point.time }}</div>
                    <div class="point-coords">{{ point.lat }}, {{ point.lng }}</div>
                  </div>
                  <div v-if="idx < routePoints.length - 1" class="route-line"></div>
                </div>
              </div>
              <div class="map-tip">
                <el-icon><Location /></el-icon>
                <span>物流轨迹示意图（实际部署可接入高德地图API）</span>
              </div>
            </div>
          </div>
        </el-card>

        <!-- 当前选中记录详情 -->
        <el-card style="margin-top:12px" v-if="currentRecord">
          <template #header><span style="font-weight:600">📍 当前节点详情</span></template>
          <el-descriptions :column="1" size="small">
            <el-descriptions-item label="当前位置">{{ currentRecord.location }}</el-descriptions-item>
            <el-descriptions-item label="运输公司">{{ currentRecord.env_data?.company }}</el-descriptions-item>
            <el-descriptions-item label="车牌号">{{ currentRecord.env_data?.vehicle_no }}</el-descriptions-item>
            <el-descriptions-item label="目的地">{{ currentRecord.env_data?.destination }}</el-descriptions-item>
            <el-descriptions-item label="温度">{{ currentRecord.env_data?.temperature }}°C</el-descriptions-item>
            <el-descriptions-item label="交易哈希">
              <span class="hash-text">{{ currentRecord.tx_hash?.slice(0,20) }}...</span>
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="showDialog" title="新增物流记录" width="600px">
      <el-form :model="form" label-width="110px" :rules="rules" ref="formRef">
        <el-form-item label="关联批次" prop="batch_id">
          <el-select v-model="form.batch_id" placeholder="选择批次" style="width:100%">
            <el-option v-for="b in batches" :key="b.id" :label="`${b.batch_no} - ${b.product_name}`" :value="b.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="物流状态" prop="logistics_status">
          <el-select v-model="form.logistics_status" style="width:100%">
            <el-option label="已发货" value="已发货" />
            <el-option label="运输中" value="运输中" />
            <el-option label="中转站" value="中转站" />
            <el-option label="已到达" value="已到达" />
            <el-option label="已签收" value="已签收" />
          </el-select>
        </el-form-item>
        <el-form-item label="当前位置" prop="location">
          <el-input v-model="form.location" placeholder="如：新疆阿克苏市物流中心" />
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="纬度" prop="lat">
              <el-input-number v-model="form.lat" :precision="6" :step="0.000001" style="width:100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="经度" prop="lng">
              <el-input-number v-model="form.lng" :precision="6" :step="0.000001" style="width:100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="运输公司">
              <el-input v-model="form.company" placeholder="物流公司名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="车牌号">
              <el-input v-model="form.vehicle_no" placeholder="如：新A12345" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="目的地">
              <el-input v-model="form.destination" placeholder="最终目的地" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="车厢温度(°C)">
              <el-input-number v-model="form.temperature" :min="-20" :max="30" style="width:100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="操作时间" prop="operation_time">
          <el-date-picker v-model="form.operation_time" type="datetime" placeholder="选择时间" style="width:100%" format="YYYY-MM-DD HH:mm:ss" value-format="YYYY-MM-DD HH:mm:ss" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确认上链</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Location } from '@element-plus/icons-vue'
import { listTraceRecords, addTraceRecord, listBatches } from '../../api/index'

const loading = ref(false)
const submitting = ref(false)
const showDialog = ref(false)
const records = ref<any[]>([])
const batches = ref<any[]>([])
const currentRecord = ref<any>(null)
const selectedIdx = ref(0)
const formRef = ref()

const form = ref({
  batch_id: null, logistics_status: '运输中', location: '新疆阿克苏市物流中心',
  lat: 41.168400, lng: 80.260480, company: '新疆顺丰冷链物流',
  vehicle_no: '新A12345', destination: '上海盒马鲜生仓储中心', temperature: 4,
  operation_time: new Date().toISOString().slice(0, 19).replace('T', ' '),
})

const rules = {
  batch_id: [{ required: true, message: '请选择批次', trigger: 'change' }],
  location: [{ required: true, message: '请输入位置', trigger: 'blur' }],
  operation_time: [{ required: true, message: '请选择时间', trigger: 'change' }],
}

// 模拟路由点（实际从records生成）
const routePoints = computed(() => {
  if (records.value.length > 0) {
    return records.value.slice(0, 5).map((r, i) => ({
      name: r.location,
      time: formatTime(r.operation_time),
      lat: r.env_data?.lat || (41.168 + i * 0.5).toFixed(4),
      lng: r.env_data?.lng || (80.260 + i * 2.1).toFixed(4),
    }))
  }
  return [
    { name: '新疆阿克苏市物流中心', time: '2025-11-10 08:00', lat: '41.1684', lng: '80.2605' },
    { name: '乌鲁木齐中转仓', time: '2025-11-11 14:30', lat: '43.8256', lng: '87.6168' },
    { name: '兰州分拨中心', time: '2025-11-12 22:00', lat: '36.0611', lng: '103.8343' },
    { name: '郑州转运中心', time: '2025-11-13 16:00', lat: '34.7466', lng: '113.6254' },
    { name: '上海盒马鲜生仓储中心', time: '2025-11-14 10:00', lat: '31.2304', lng: '121.4737' },
  ]
})

function formatTime(t: string) { return t ? t.replace('T', ' ').slice(0, 16) : '-' }

function selectRecord(row: any) {
  currentRecord.value = row
  selectedIdx.value = records.value.indexOf(row)
}

async function loadData() {
  loading.value = true
  try {
    const [recRes, batchRes]: any[] = await Promise.all([
      listTraceRecords({ node_type: 'transporting' }),
      listBatches(),
    ])
    records.value = recRes.data || []
    batches.value = batchRes.data || []
    if (records.value.length > 0) currentRecord.value = records.value[0]
  } finally { loading.value = false }
}

async function handleSubmit() {
  await formRef.value?.validate(async (valid: boolean) => {
    if (!valid) return
    submitting.value = true
    try {
      await addTraceRecord({
        batch_id: form.value.batch_id, node_type: 'transporting',
        location: form.value.location, operation_time: form.value.operation_time,
        env_data: {
          logistics_status: form.value.logistics_status, lat: form.value.lat, lng: form.value.lng,
          company: form.value.company, vehicle_no: form.value.vehicle_no,
          destination: form.value.destination, temperature: form.value.temperature,
        },
      })
      ElMessage.success('物流记录已上链！')
      showDialog.value = false
      loadData()
    } finally { submitting.value = false }
  })
}

onMounted(loadData)
</script>

<style scoped>
.page-header-card { margin-bottom: 0; }
.page-header { display: flex; justify-content: space-between; align-items: center; }
.page-title { margin: 0; font-size: 18px; font-weight: 600; }
.page-desc { margin: 4px 0 0; color: #888; font-size: 13px; }
.hash-text { font-family: monospace; font-size: 12px; color: #1677ff; }
.map-container { padding: 8px 0; }
.map-placeholder { background: #f0f7ff; border-radius: 8px; padding: 16px; min-height: 240px; }
.map-route { display: flex; flex-direction: column; gap: 0; }
.route-point { display: flex; align-items: flex-start; gap: 10px; position: relative; padding-bottom: 8px; }
.point-dot {
  width: 28px; height: 28px; border-radius: 50%;
  background: #d9e8ff; border: 2px solid #1677ff;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.point-dot.active { background: #1677ff; }
.point-dot.active .point-num { color: #fff; }
.point-num { font-size: 12px; font-weight: 700; color: #1677ff; }
.point-info { flex: 1; }
.point-name { font-size: 13px; font-weight: 600; color: #333; }
.point-time { font-size: 11px; color: #888; }
.point-coords { font-size: 11px; color: #aaa; font-family: monospace; }
.route-line {
  position: absolute; left: 13px; top: 28px;
  width: 2px; height: 20px; background: #1677ff; opacity: 0.3;
}
.map-tip { display: flex; align-items: center; gap: 4px; font-size: 11px; color: #999; margin-top: 8px; }
</style>
