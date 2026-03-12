import request from './request'

// ==================== 认证 ====================
export const getCaptcha = () => request.get('/api/v1/auth/captcha')
export const login = (data: any) => request.post('/api/v1/auth/login', data)
export const register = (data: any) => request.post('/api/v1/auth/register', data)
export const getProfile = () => request.get('/api/v1/auth/profile')

// ==================== 仪表盘 ====================
export const getDashboardStats = () => request.get('/api/v1/dashboard/stats')

// ==================== 批次管理 ====================
export const listBatches = (params?: any) => request.get('/api/v1/batches', { params })
export const createBatch = (data: any) => request.post('/api/v1/batches', data)
export const getBatchDetail = (id: number) => request.get(`/api/v1/batches/${id}`)
export const updateBatch = (id: number, data: any) => request.put(`/api/v1/batches/${id}`, data)
export const deleteBatch = (id: number) => request.delete(`/api/v1/batches/${id}`)

// ==================== 溯源记录 ====================
export const listTraceRecords = (params?: any) => request.get('/api/v1/trace/records', { params })
export const addTraceRecord = (data: any) => request.post('/api/v1/trace/records', data)

// ==================== 电子合格证 ====================
export const getCertificate = (batchId: number) => request.get(`/api/v1/certificates/${batchId}`)
export const createCertificate = (data: any) => request.post('/api/v1/certificates', data)

// ==================== 公开溯源查询 ====================
export const queryByTraceCode = (code: string) => request.get(`/api/v1/trace/public/${code}`)

// ==================== 区块链信息 ====================
export const getBlockInfo = () => request.get('/api/v1/block/info')

// ==================== 用户管理（管理员） ====================
export const listUsers = (params?: any) => request.get('/api/v1/admin/users', { params })
export const createUser = (data: any) => request.post('/api/v1/admin/users', data)
export const updateUser = (id: number, data: any) => request.put(`/api/v1/admin/users/${id}`, data)
export const deleteUserApi = (id: number) => request.delete(`/api/v1/admin/users/${id}`)

// ==================== 文件上传 ====================
export const uploadFile = (formData: FormData) => request.post('/api/v1/upload', formData, {
  headers: { 'Content-Type': 'multipart/form-data' }
})
