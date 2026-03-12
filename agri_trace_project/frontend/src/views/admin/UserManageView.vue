<template>
  <div>
    <el-card class="page-header-card">
      <div class="page-header">
        <div>
          <h3 class="page-title">👥 用户与权限管理</h3>
          <p class="page-desc">管理系统用户账号、角色权限分配</p>
        </div>
        <el-button type="primary" icon="Plus" @click="openCreate">新增用户</el-button>
      </div>
    </el-card>

    <el-card style="margin-top:16px">
      <div style="margin-bottom:16px;display:flex;gap:12px">
        <el-input v-model="searchKey" placeholder="搜索用户名/姓名" clearable style="width:220px" @input="loadData" />
        <el-select v-model="filterRole" placeholder="筛选角色" clearable style="width:160px" @change="loadData">
          <el-option v-for="r in roleOptions" :key="r.value" :label="r.label" :value="r.value" />
        </el-select>
      </div>
      <el-table :data="users" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="real_name" label="真实姓名" width="100" />
        <el-table-column label="角色" width="120">
          <template #default="{ row }">
            <el-tag :type="roleColor(row.role)">{{ roleLabel(row.role) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-switch v-model="row.is_active" @change="toggleStatus(row)" />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="160">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="editUser(row)">编辑</el-button>
            <el-popconfirm title="确认删除该用户？" @confirm="deleteUser(row.id)">
              <template #reference>
                <el-button size="small" type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        style="margin-top:16px;justify-content:flex-end;display:flex"
        @change="loadData"
      />
    </el-card>

    <el-dialog v-model="showDialog" :title="editMode ? '编辑用户' : '新增用户'" width="500px">
      <el-form :model="form" label-width="90px" :rules="rules" ref="formRef">
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="form.username" :disabled="editMode" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="真实姓名" prop="real_name">
              <el-input v-model="form.real_name" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="密码" prop="password" v-if="!editMode">
          <el-input v-model="form.password" type="password" show-password placeholder="至少6位" />
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="角色" prop="role">
              <el-select v-model="form.role" style="width:100%">
                <el-option v-for="r in roleOptions" :key="r.value" :label="r.label" :value="r.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号">
              <el-input v-model="form.phone" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="邮箱">
          <el-input v-model="form.email" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listUsers, createUser, updateUser, deleteUserApi } from '../../api/index'

const loading = ref(false)
const submitting = ref(false)
const showDialog = ref(false)
const editMode = ref(false)
const users = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const searchKey = ref('')
const filterRole = ref('')
const formRef = ref()

const form = ref({ id: 0, username: '', real_name: '', password: '', role: 'farmer', phone: '', email: '' })

const roleOptions = [
  { value: 'admin', label: '🔑 系统管理员' },
  { value: 'farmer', label: '🌾 种植户' },
  { value: 'inspector', label: '🔬 质检员' },
  { value: 'transporter', label: '🚛 物流商' },
  { value: 'retailer', label: '🏪 销售商' },
]

const roleLabels: Record<string, string> = { admin: '系统管理员', farmer: '种植户', inspector: '质检员', transporter: '物流商', retailer: '销售商' }
const roleColors: Record<string, string> = { admin: 'danger', farmer: 'success', inspector: 'primary', transporter: 'warning', retailer: 'info' }

function roleLabel(r: string) { return roleLabels[r] || r }
function roleColor(r: string): any { return roleColors[r] || '' }
function formatTime(t: string) { return t ? t.replace('T', ' ').slice(0, 16) : '-' }

const rules = {
  username: [{ required: true, min: 3, max: 20, message: '用户名3-20位', trigger: 'blur' }],
  real_name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  password: [{ required: true, min: 6, message: '密码至少6位', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
}

function openCreate() {
  editMode.value = false
  form.value = { id: 0, username: '', real_name: '', password: '', role: 'farmer', phone: '', email: '' }
  showDialog.value = true
}

function editUser(row: any) {
  editMode.value = true
  form.value = { ...row, password: '' }
  showDialog.value = true
}

async function toggleStatus(row: any) {
  try {
    await updateUser(row.id, { is_active: row.is_active })
    ElMessage.success('状态已更新')
  } catch {}
}

async function deleteUser(id: number) {
  try {
    await deleteUserApi(id)
    ElMessage.success('删除成功')
    loadData()
  } catch {}
}

async function loadData() {
  loading.value = true
  try {
    const res: any = await listUsers({ page: page.value, page_size: pageSize.value, keyword: searchKey.value, role: filterRole.value })
    users.value = res.data?.list || res.data || []
    total.value = res.data?.total || users.value.length
  } catch {
    // fallback mock data
    users.value = [
      { id: 1, username: 'admin', real_name: '系统管理员', role: 'admin', phone: '13800000000', email: 'admin@aksu.com', is_active: true, created_at: '2025-01-01' },
      { id: 2, username: 'farmer_wang', real_name: '王建国', role: 'farmer', phone: '13912345678', email: 'wang@aksu.com', is_active: true, created_at: '2025-01-02' },
      { id: 3, username: 'inspector_zhang', real_name: '张质检', role: 'inspector', phone: '13987654321', email: 'zhang@aksu.com', is_active: true, created_at: '2025-01-03' },
      { id: 4, username: 'transporter_sun', real_name: '孙师傅', role: 'transporter', phone: '13611112222', email: 'sun@aksu.com', is_active: true, created_at: '2025-01-04' },
      { id: 5, username: 'retailer_zhao', real_name: '赵店长', role: 'retailer', phone: '13733334444', email: 'zhao@aksu.com', is_active: true, created_at: '2025-01-05' },
    ]
    total.value = users.value.length
  } finally { loading.value = false }
}

async function handleSubmit() {
  await formRef.value?.validate(async (valid: boolean) => {
    if (!valid) return
    submitting.value = true
    try {
      if (editMode.value) {
        await updateUser(form.value.id, form.value)
        ElMessage.success('更新成功')
      } else {
        await createUser(form.value)
        ElMessage.success('创建成功')
      }
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
</style>
