<template>
  <div>
    <el-card class="page-header-card">
      <div class="page-header">
        <div>
          <h3 class="page-title">🔬 质检管理</h3>
          <p class="page-desc">记录农残、重金属、微生物等检测结果，生成电子质检报告</p>
        </div>
        <el-button type="primary" icon="Plus" @click="showDialog = true">新增质检记录</el-button>
      </div>
    </el-card>

    <el-card style="margin-top:16px">
      <el-table :data="records" stripe v-loading="loading">
        <el-table-column prop="batch_id" label="批次ID" width="80" />
        <el-table-column label="检测机构" width="180">
          <template #default="{ row }">{{ row.env_data?.inspect_org || '-' }}</template>
        </el-table-column>
        <el-table-column label="农残检测" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.env_data?.pesticide_ok ? 'success' : 'danger'" size="small">
              <el-icon><component :is="row.env_data?.pesticide_ok ? 'Check' : 'Close'" /></el-icon>
              {{ row.env_data?.pesticide_ok ? '合格' : '不合格' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="重金属检测" width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="row.env_data?.heavy_metal_ok ? 'success' : 'danger'" size="small">
              {{ row.env_data?.heavy_metal_ok ? '合格' : '不合格' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="微生物检测" width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="row.env_data?.microbe_ok ? 'success' : 'danger'" size="small">
              {{ row.env_data?.microbe_ok ? '合格' : '不合格' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="综合结论" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="isAllPass(row) ? 'success' : 'danger'">
              {{ isAllPass(row) ? '✅ 通过' : '❌ 不合格' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="operation_time" label="检测时间" width="160">
          <template #default="{ row }">{{ formatTime(row.operation_time) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="80">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="viewDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showDialog" title="新增质检记录" width="600px">
      <el-form :model="form" label-width="110px" :rules="rules" ref="formRef">
        <el-form-item label="关联批次" prop="batch_id">
          <el-select v-model="form.batch_id" placeholder="选择批次" style="width:100%">
            <el-option v-for="b in batches" :key="b.id" :label="`${b.batch_no} - ${b.product_name}`" :value="b.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="检测机构" prop="inspect_org">
          <el-input v-model="form.inspect_org" placeholder="如：阿克苏市农业农村局质检中心" />
        </el-form-item>
        <el-form-item label="检测地点" prop="location">
          <el-input v-model="form.location" placeholder="检测地点" />
        </el-form-item>
        <el-form-item label="农残检测">
          <el-radio-group v-model="form.pesticide_ok">
            <el-radio :value="true">合格</el-radio>
            <el-radio :value="false">不合格</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="重金属检测">
          <el-radio-group v-model="form.heavy_metal_ok">
            <el-radio :value="true">合格</el-radio>
            <el-radio :value="false">不合格</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="微生物检测">
          <el-radio-group v-model="form.microbe_ok">
            <el-radio :value="true">合格</el-radio>
            <el-radio :value="false">不合格</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="检测时间" prop="operation_time">
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

    <el-dialog v-model="showDetail" title="质检记录详情" width="500px">
      <el-descriptions :column="2" border v-if="currentRecord">
        <el-descriptions-item label="批次ID">{{ currentRecord.batch_id }}</el-descriptions-item>
        <el-descriptions-item label="检测机构">{{ currentRecord.env_data?.inspect_org }}</el-descriptions-item>
        <el-descriptions-item label="检测地点" :span="2">{{ currentRecord.location }}</el-descriptions-item>
        <el-descriptions-item label="农残检测">
          <el-tag :type="currentRecord.env_data?.pesticide_ok ? 'success' : 'danger'" size="small">{{ currentRecord.env_data?.pesticide_ok ? '合格' : '不合格' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="重金属检测">
          <el-tag :type="currentRecord.env_data?.heavy_metal_ok ? 'success' : 'danger'" size="small">{{ currentRecord.env_data?.heavy_metal_ok ? '合格' : '不合格' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="微生物检测">
          <el-tag :type="currentRecord.env_data?.microbe_ok ? 'success' : 'danger'" size="small">{{ currentRecord.env_data?.microbe_ok ? '合格' : '不合格' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="综合结论">
          <el-tag :type="isAllPass(currentRecord) ? 'success' : 'danger'">{{ isAllPass(currentRecord) ? '通过' : '不合格' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="区块高度">{{ currentRecord.block_height }}</el-descriptions-item>
        <el-descriptions-item label="检测时间">{{ formatTime(currentRecord.operation_time) }}</el-descriptions-item>
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
  batch_id: null, inspect_org: '阿克苏市农业农村局质检中心',
  location: '阿克苏市农业农村局质检中心实验室',
  pesticide_ok: true, heavy_metal_ok: true, microbe_ok: true, remark: '',
  operation_time: new Date().toISOString().slice(0, 19).replace('T', ' '),
})

const rules = {
  batch_id: [{ required: true, message: '请选择批次', trigger: 'change' }],
  inspect_org: [{ required: true, message: '请输入检测机构', trigger: 'blur' }],
  location: [{ required: true, message: '请输入检测地点', trigger: 'blur' }],
  operation_time: [{ required: true, message: '请选择检测时间', trigger: 'change' }],
}

function isAllPass(row: any) {
  return row.env_data?.pesticide_ok && row.env_data?.heavy_metal_ok && row.env_data?.microbe_ok
}
function formatTime(t: string) { return t ? t.replace('T', ' ').slice(0, 19) : '-' }
function viewDetail(row: any) { currentRecord.value = row; showDetail.value = true }

async function loadData() {
  loading.value = true
  try {
    const [recRes, batchRes]: any[] = await Promise.all([
      listTraceRecords({ node_type: 'inspecting' }),
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
        batch_id: form.value.batch_id, node_type: 'inspecting',
        location: form.value.location, operation_time: form.value.operation_time,
        env_data: {
          inspect_org: form.value.inspect_org, pesticide_ok: form.value.pesticide_ok,
          heavy_metal_ok: form.value.heavy_metal_ok, microbe_ok: form.value.microbe_ok, remark: form.value.remark,
        },
      })
      ElMessage.success('质检记录已上链！')
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
