<template>
  <div>
    <el-card class="page-header-card">
      <div class="page-header">
        <div>
          <h3 class="page-title">⚙️ 加工环节管理</h3>
          <p class="page-desc">精细化记录清洗、烘干、分级、包装等每次加工操作，责任到人</p>
        </div>
        <el-button type="primary" icon="Plus" @click="showDialog = true">新增加工记录</el-button>
      </div>
    </el-card>

    <el-card style="margin-top:16px">
      <el-table :data="records" stripe v-loading="loading">
        <el-table-column prop="batch_id" label="批次ID" width="80" />
        <el-table-column label="加工工序" width="120">
          <template #default="{ row }">
            <el-tag type="warning">{{ row.env_data?.process_type || '加工' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="location" label="加工厂地址" />
        <el-table-column label="操作人" width="100">
          <template #default="{ row }">{{ row.env_data?.operator_name || '-' }}</template>
        </el-table-column>
        <el-table-column label="操作人电话" width="130">
          <template #default="{ row }">{{ row.env_data?.operator_phone || '-' }}</template>
        </el-table-column>
        <el-table-column label="投入量/产出量" width="130">
          <template #default="{ row }">
            {{ row.env_data?.input_qty || '-' }} / {{ row.env_data?.output_qty || '-' }} kg
          </template>
        </el-table-column>
        <el-table-column prop="operation_time" label="操作时间" width="160">
          <template #default="{ row }">{{ formatTime(row.operation_time) }}</template>
        </el-table-column>
        <el-table-column prop="tx_hash" label="交易哈希" width="140">
          <template #default="{ row }">
            <el-tooltip :content="row.tx_hash">
              <span class="hash-text">{{ row.tx_hash?.slice(0,14) }}...</span>
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

    <el-dialog v-model="showDialog" title="新增加工记录" width="620px">
      <el-form :model="form" label-width="110px" :rules="rules" ref="formRef">
        <el-form-item label="关联批次" prop="batch_id">
          <el-select v-model="form.batch_id" placeholder="选择批次" style="width:100%">
            <el-option v-for="b in batches" :key="b.id" :label="`${b.batch_no} - ${b.product_name}`" :value="b.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="加工工序" prop="process_type">
          <el-select v-model="form.process_type" style="width:100%">
            <el-option label="清洗" value="清洗" />
            <el-option label="分级" value="分级" />
            <el-option label="烘干" value="烘干" />
            <el-option label="打蜡" value="打蜡" />
            <el-option label="包装" value="包装" />
            <el-option label="装箱" value="装箱" />
            <el-option label="冷藏" value="冷藏" />
          </el-select>
        </el-form-item>
        <el-form-item label="加工厂地址" prop="location">
          <el-input v-model="form.location" placeholder="如：阿克苏市农产品加工园区A栋" />
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="操作人姓名" prop="operator_name">
              <el-input v-model="form.operator_name" placeholder="操作人姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="操作人电话" prop="operator_phone">
              <el-input v-model="form.operator_phone" placeholder="联系电话" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item label="投入量(kg)">
              <el-input-number v-model="form.input_qty" :min="0" style="width:100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="产出量(kg)">
              <el-input-number v-model="form.output_qty" :min="0" style="width:100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="温度(°C)">
              <el-input-number v-model="form.temperature" :min="-20" :max="50" style="width:100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="操作时间" prop="operation_time">
          <el-date-picker v-model="form.operation_time" type="datetime" placeholder="选择时间" style="width:100%" format="YYYY-MM-DD HH:mm:ss" value-format="YYYY-MM-DD HH:mm:ss" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确认上链</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showDetail" title="加工记录详情" width="500px">
      <el-descriptions :column="2" border v-if="currentRecord">
        <el-descriptions-item label="批次ID">{{ currentRecord.batch_id }}</el-descriptions-item>
        <el-descriptions-item label="加工工序">{{ currentRecord.env_data?.process_type }}</el-descriptions-item>
        <el-descriptions-item label="加工厂地址" :span="2">{{ currentRecord.location }}</el-descriptions-item>
        <el-descriptions-item label="操作人">{{ currentRecord.env_data?.operator_name }}</el-descriptions-item>
        <el-descriptions-item label="操作人电话">{{ currentRecord.env_data?.operator_phone }}</el-descriptions-item>
        <el-descriptions-item label="投入量">{{ currentRecord.env_data?.input_qty }} kg</el-descriptions-item>
        <el-descriptions-item label="产出量">{{ currentRecord.env_data?.output_qty }} kg</el-descriptions-item>
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
  batch_id: null, process_type: '清洗', location: '阿克苏市农产品加工园区A栋',
  operator_name: '', operator_phone: '', input_qty: 0, output_qty: 0,
  temperature: 4, remark: '',
  operation_time: new Date().toISOString().slice(0, 19).replace('T', ' '),
})

const rules = {
  batch_id: [{ required: true, message: '请选择批次', trigger: 'change' }],
  location: [{ required: true, message: '请输入地址', trigger: 'blur' }],
  operator_name: [{ required: true, message: '请输入操作人姓名', trigger: 'blur' }],
  operator_phone: [{ required: true, message: '请输入操作人电话', trigger: 'blur' }],
  operation_time: [{ required: true, message: '请选择操作时间', trigger: 'change' }],
}

function formatTime(t: string) { return t ? t.replace('T', ' ').slice(0, 19) : '-' }
function viewDetail(row: any) { currentRecord.value = row; showDetail.value = true }

async function loadData() {
  loading.value = true
  try {
    const [recRes, batchRes]: any[] = await Promise.all([
      listTraceRecords({ node_type: 'processing' }),
      listBatches(),
    ])
    records.value = recRes.data || []
    batches.value = batchRes.data || []
  } finally { loading.value = false }
}

async function handleSubmit() {
  await formRef.value?.validate(async (valid: boolean) => {
    if (!valid) return
    submitting.value = true
    try {
      await addTraceRecord({
        batch_id: form.value.batch_id, node_type: 'processing',
        location: form.value.location, operation_time: form.value.operation_time,
        env_data: {
          process_type: form.value.process_type, operator_name: form.value.operator_name,
          operator_phone: form.value.operator_phone, input_qty: form.value.input_qty,
          output_qty: form.value.output_qty, temperature: form.value.temperature, remark: form.value.remark,
        },
      })
      ElMessage.success('加工记录已上链！')
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
</style>
