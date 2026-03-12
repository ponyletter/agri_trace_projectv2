<template>
  <div>
    <el-card class="page-header-card">
      <div class="page-header">
        <div>
          <h3 class="page-title">🌱 种植环节管理</h3>
          <p class="page-desc">记录作物种植、施肥、灌溉、除草等全过程信息</p>
        </div>
        <el-button type="primary" icon="Plus" @click="showDialog = true">新增种植记录</el-button>
      </div>
    </el-card>

    <el-card style="margin-top:16px">
      <el-table :data="records" stripe v-loading="loading">
        <el-table-column prop="batch_id" label="批次ID" width="80" />
        <el-table-column label="操作类型" width="100">
          <template #default="{ row }">
            <el-tag :type="nodeTypeColor(row.env_data?.operation_type)">{{ row.env_data?.operation_type || '种植' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="location" label="地理位置" />
        <el-table-column label="温度/湿度" width="120">
          <template #default="{ row }">
            {{ row.env_data?.temperature }}°C / {{ row.env_data?.humidity }}%
          </template>
        </el-table-column>
        <el-table-column label="施肥/灌溉" width="140">
          <template #default="{ row }">
            <span>{{ row.env_data?.fertilizer || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="operation_time" label="操作时间" width="160">
          <template #default="{ row }">{{ formatTime(row.operation_time) }}</template>
        </el-table-column>
        <el-table-column prop="tx_hash" label="交易哈希" width="160">
          <template #default="{ row }">
            <el-tooltip :content="row.tx_hash">
              <span class="hash-text">{{ row.tx_hash?.slice(0,16) }}...</span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="viewDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增对话框 -->
    <el-dialog v-model="showDialog" title="新增种植记录" width="600px">
      <el-form :model="form" label-width="100px" :rules="rules" ref="formRef">
        <el-form-item label="关联批次" prop="batch_id">
          <el-select v-model="form.batch_id" placeholder="选择批次" style="width:100%">
            <el-option v-for="b in batches" :key="b.id" :label="`${b.batch_no} - ${b.product_name}`" :value="b.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作类型" prop="operation_type">
          <el-select v-model="form.operation_type" style="width:100%">
            <el-option label="播种/定植" value="播种/定植" />
            <el-option label="施肥" value="施肥" />
            <el-option label="灌溉" value="灌溉" />
            <el-option label="除草" value="除草" />
            <el-option label="修剪" value="修剪" />
            <el-option label="病虫害防治" value="病虫害防治" />
            <el-option label="采摘" value="采摘" />
          </el-select>
        </el-form-item>
        <el-form-item label="地理位置" prop="location">
          <el-input v-model="form.location" placeholder="如：新疆阿克苏温宿县红旗坡农场" />
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item label="温度(°C)">
              <el-input-number v-model="form.temperature" :min="-20" :max="50" style="width:100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="湿度(%)">
              <el-input-number v-model="form.humidity" :min="0" :max="100" style="width:100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="采摘数量(kg)">
              <el-input-number v-model="form.quantity" :min="0" style="width:100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="施肥信息">
          <el-input v-model="form.fertilizer" placeholder="如：有机肥50kg/亩，磷肥20kg/亩" />
        </el-form-item>
        <el-form-item label="操作时间" prop="operation_time">
          <el-date-picker v-model="form.operation_time" type="datetime" placeholder="选择时间" style="width:100%" format="YYYY-MM-DD HH:mm:ss" value-format="YYYY-MM-DD HH:mm:ss" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" placeholder="其他说明" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确认上链</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="showDetail" title="种植记录详情" width="500px">
      <el-descriptions :column="2" border v-if="currentRecord">
        <el-descriptions-item label="批次ID">{{ currentRecord.batch_id }}</el-descriptions-item>
        <el-descriptions-item label="操作类型">{{ currentRecord.env_data?.operation_type }}</el-descriptions-item>
        <el-descriptions-item label="地理位置" :span="2">{{ currentRecord.location }}</el-descriptions-item>
        <el-descriptions-item label="温度">{{ currentRecord.env_data?.temperature }}°C</el-descriptions-item>
        <el-descriptions-item label="湿度">{{ currentRecord.env_data?.humidity }}%</el-descriptions-item>
        <el-descriptions-item label="施肥信息" :span="2">{{ currentRecord.env_data?.fertilizer || '-' }}</el-descriptions-item>
        <el-descriptions-item label="区块高度">{{ currentRecord.block_height }}</el-descriptions-item>
        <el-descriptions-item label="操作时间">{{ formatTime(currentRecord.operation_time) }}</el-descriptions-item>
        <el-descriptions-item label="交易哈希" :span="2">
          <span class="hash-text">{{ currentRecord.tx_hash }}</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listTraceRecords, addTraceRecord, listBatches } from '../../api/index'

const loading = ref(false)
const submitting = ref(false)
const showDialog = ref(false)
const showDetail = ref(false)
const records = ref<any[]>([])
const batches = ref<any[]>([])
const currentRecord = ref<any>(null)
const formRef = ref()

const form = ref({
  batch_id: null, operation_type: '施肥', location: '新疆阿克苏温宿县红旗坡农场',
  temperature: 22, humidity: 65, quantity: 0, fertilizer: '', remark: '',
  operation_time: new Date().toISOString().slice(0, 19).replace('T', ' '),
})

const rules = {
  batch_id: [{ required: true, message: '请选择批次', trigger: 'change' }],
  location: [{ required: true, message: '请输入地理位置', trigger: 'blur' }],
  operation_time: [{ required: true, message: '请选择操作时间', trigger: 'change' }],
}

function nodeTypeColor(type: string) {
  const map: Record<string, string> = { '施肥': 'success', '灌溉': 'primary', '除草': 'warning', '采摘': 'danger', '修剪': 'info' }
  return map[type] || ''
}

function formatTime(t: string) {
  return t ? t.replace('T', ' ').slice(0, 19) : '-'
}

function viewDetail(row: any) {
  currentRecord.value = row
  showDetail.value = true
}

async function loadData() {
  loading.value = true
  try {
    const [recRes, batchRes]: any[] = await Promise.all([
      listTraceRecords({ node_type: 'planting' }),
      listBatches(),
    ])
    records.value = recRes.data || []
    batches.value = batchRes.data || []
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  await formRef.value?.validate(async (valid: boolean) => {
    if (!valid) return
    submitting.value = true
    try {
      await addTraceRecord({
        batch_id: form.value.batch_id,
        node_type: 'planting',
        location: form.value.location,
        operation_time: form.value.operation_time,
        env_data: {
          operation_type: form.value.operation_type,
          temperature: form.value.temperature,
          humidity: form.value.humidity,
          quantity: form.value.quantity,
          fertilizer: form.value.fertilizer,
          remark: form.value.remark,
        },
      })
      ElMessage.success('种植记录已上链！')
      showDialog.value = false
      loadData()
    } finally {
      submitting.value = false
    }
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
</style>
