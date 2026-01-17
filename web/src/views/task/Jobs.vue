<template>
  <div class="jobs-container">
    <!-- 页面标题和操作按钮 -->
    <div class="page-header">
      <div class="page-title-group">
        <div class="page-title-icon">
          <el-icon><Tickets /></el-icon>
        </div>
        <div>
          <h2 class="page-title">任务作业</h2>
          <p class="page-subtitle">管理任务作业，支持任务创建、编辑与执行</p>
        </div>
      </div>
      <div class="header-actions">
        <el-button class="black-button" @click="handleAdd">
          <el-icon style="margin-right: 6px;"><Plus /></el-icon>
          新增任务
        </el-button>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <div class="search-inputs">
        <el-input
          v-model="searchForm.keyword"
          placeholder="搜索任务名称..."
          clearable
          class="search-input"
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon class="search-icon"><Search /></el-icon>
          </template>
        </el-input>

        <el-select
          v-model="searchForm.taskType"
          placeholder="任务类型"
          clearable
          class="search-input"
          @change="handleSearch"
        >
          <el-option label="脚本执行" value="script" />
          <el-option label="文件传输" value="file" />
          <el-option label="系统命令" value="command" />
        </el-select>

        <el-select
          v-model="searchForm.status"
          placeholder="执行状态"
          clearable
          class="search-input"
          @change="handleSearch"
        >
          <el-option label="等待中" value="pending" />
          <el-option label="执行中" value="running" />
          <el-option label="成功" value="success" />
          <el-option label="失败" value="failed" />
        </el-select>
      </div>

      <div class="search-actions">
        <el-button class="reset-btn" @click="handleReset">
          <el-icon style="margin-right: 4px;"><RefreshLeft /></el-icon>
          重置
        </el-button>
      </div>
    </div>

    <!-- 表格容器 -->
    <div class="table-wrapper">
      <el-table
        :data="jobList"
        v-loading="loading"
        class="modern-table"
        :header-cell-style="{ background: '#fafbfc', color: '#606266', fontWeight: '600' }"
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />

        <el-table-column label="任务名称" prop="name" min-width="180">
          <template #default="{ row }">
            <div class="job-name-cell">
              <el-icon class="job-icon"><List /></el-icon>
              <span class="job-name">{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="任务类型" prop="taskType" width="120">
          <template #default="{ row }">
            <el-tag :type="getTaskTypeColor(row.taskType)" effect="plain">
              {{ getTaskTypeLabel(row.taskType) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="状态" prop="status" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" effect="dark">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="targetHosts" label="目标主机" min-width="150">
          <template #default="{ row }">
            <span class="description-text">{{ row.targetHosts || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="executeTime" label="执行时间" width="180" />

        <el-table-column label="操作" width="150" fixed="right" align="center">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-tooltip content="查看" placement="top">
                <el-button link class="action-btn action-view" @click="handleView(row)">
                  <el-icon><View /></el-icon>
                </el-button>
              </el-tooltip>
              <el-tooltip content="编辑" placement="top">
                <el-button link class="action-btn action-edit" @click="handleEdit(row)">
                  <el-icon><Edit /></el-icon>
                </el-button>
              </el-tooltip>
              <el-tooltip content="删除" placement="top">
                <el-button link class="action-btn action-delete" @click="handleDelete(row)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          layout="total, prev, pager, next, jumper"
          @current-change="loadJobs"
        />
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="60%"
      class="job-edit-dialog responsive-dialog"
      :close-on-click-modal="false"
      @close="handleDialogClose"
    >
      <el-form :model="jobForm" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="jobForm.name" placeholder="请输入任务名称" />
        </el-form-item>

        <el-form-item label="任务类型" prop="taskType">
          <el-select v-model="jobForm.taskType" placeholder="请选择任务类型" style="width: 100%">
            <el-option label="脚本执行" value="script" />
            <el-option label="文件传输" value="file" />
            <el-option label="系统命令" value="command" />
          </el-select>
        </el-form-item>

        <el-form-item label="任务模板" prop="templateId">
          <el-select v-model="jobForm.templateId" placeholder="请选择任务模板" clearable style="width: 100%">
            <el-option
              v-for="template in templateList"
              :key="template.id"
              :label="template.name"
              :value="template.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="目标主机" prop="targetHosts">
          <el-input
            v-model="jobForm.targetHosts"
            type="textarea"
            :rows="3"
            placeholder="请输入目标主机，多个主机用逗号分隔"
          />
        </el-form-item>

        <el-form-item label="任务参数" prop="parameters">
          <el-input
            v-model="jobForm.parameters"
            type="textarea"
            :rows="4"
            placeholder="请输入任务参数（JSON格式）"
          />
        </el-form-item>

        <el-form-item label="执行时间" prop="executeTime">
          <el-date-picker
            v-model="jobForm.executeTime"
            type="datetime"
            placeholder="选择执行时间"
            style="width: 100%"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button class="black-button" @click="handleSubmit" :loading="submitting">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 查看详情对话框 -->
    <el-dialog
      v-model="viewDialogVisible"
      title="任务详情"
      width="60%"
      class="job-view-dialog responsive-dialog"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="任务ID">{{ currentJob.id }}</el-descriptions-item>
        <el-descriptions-item label="任务名称">{{ currentJob.name }}</el-descriptions-item>
        <el-descriptions-item label="任务类型">
          <el-tag :type="getTaskTypeColor(currentJob.taskType)" effect="plain">
            {{ getTaskTypeLabel(currentJob.taskType) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentJob.status)" effect="dark">
            {{ getStatusLabel(currentJob.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="目标主机" :span="2">{{ currentJob.targetHosts || '-' }}</el-descriptions-item>
        <el-descriptions-item label="任务参数" :span="2">
          <pre style="margin: 0; white-space: pre-wrap;">{{ currentJob.parameters || '-' }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="执行时间">{{ currentJob.executeTime || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentJob.createdAt }}</el-descriptions-item>
        <el-descriptions-item label="执行结果" :span="2" v-if="currentJob.result">
          <pre style="margin: 0; white-space: pre-wrap;">{{ currentJob.result }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="错误信息" :span="2" v-if="currentJob.errorMessage">
          <span style="color: #f56c6c;">{{ currentJob.errorMessage }}</span>
        </el-descriptions-item>
      </el-descriptions>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="viewDialogVisible = false">关闭</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  Plus,
  Edit,
  Delete,
  Search,
  RefreshLeft,
  Tickets,
  List,
  View
} from '@element-plus/icons-vue'
import { getJobTaskList, createJobTask, updateJobTask, deleteJobTask, getAllJobTemplates } from '@/api/task'

// 加载状态
const loading = ref(false)
const submitting = ref(false)

// 对话框状态
const dialogVisible = ref(false)
const viewDialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)

// 表单引用
const formRef = ref<FormInstance>()

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 任务列表
const jobList = ref<any[]>([])

// 模板列表
const templateList = ref<any[]>([])

// 当前查看的任务
const currentJob = ref<any>({})

// 搜索表单
const searchForm = reactive({
  keyword: '',
  taskType: '',
  status: ''
})

// 任务表单
const jobForm = reactive({
  id: 0,
  name: '',
  templateId: undefined as number | undefined,
  taskType: '',
  targetHosts: '',
  parameters: '',
  executeTime: ''
})

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: '请输入任务名称', trigger: 'blur' },
    { min: 2, max: 100, message: '任务名称长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  taskType: [
    { required: true, message: '请选择任务类型', trigger: 'change' }
  ]
}

// 获取任务类型颜色
const getTaskTypeColor = (type: string) => {
  const colorMap: Record<string, string> = {
    script: 'success',
    file: 'warning',
    command: 'info'
  }
  return colorMap[type] || 'info'
}

// 获取任务类型标签
const getTaskTypeLabel = (type: string) => {
  const labelMap: Record<string, string> = {
    script: '脚本执行',
    file: '文件传输',
    command: '系统命令'
  }
  return labelMap[type] || type
}

// 获取状态类型
const getStatusType = (status: string) => {
  const typeMap: Record<string, string> = {
    pending: 'info',
    running: 'warning',
    success: 'success',
    failed: 'danger'
  }
  return typeMap[status] || 'info'
}

// 获取状态标签
const getStatusLabel = (status: string) => {
  const labelMap: Record<string, string> = {
    pending: '等待中',
    running: '执行中',
    success: '成功',
    failed: '失败'
  }
  return labelMap[status] || status
}

// 搜索处理
const handleSearch = () => {
  pagination.page = 1
  loadJobs()
}

// 重置搜索
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.taskType = ''
  searchForm.status = ''
  loadJobs()
}

// 加载任务列表
const loadJobs = async () => {
  loading.value = true
  try {
    const res = await getJobTaskList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: searchForm.keyword,
      taskType: searchForm.taskType,
      status: searchForm.status
    })
    jobList.value = res.list || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('获取任务列表失败:', error)
    ElMessage.error('获取任务列表失败')
  } finally {
    loading.value = false
  }
}

// 加载模板列表
const loadTemplates = async () => {
  try {
    const res = await getAllJobTemplates()
    templateList.value = res || []
  } catch (error) {
    console.error('获取模板列表失败:', error)
  }
}

// 重置表单
const resetForm = () => {
  jobForm.id = 0
  jobForm.name = ''
  jobForm.templateId = undefined
  jobForm.taskType = ''
  jobForm.targetHosts = ''
  jobForm.parameters = ''
  jobForm.executeTime = ''
  formRef.value?.clearValidate()
}

// 新增任务
const handleAdd = () => {
  resetForm()
  dialogTitle.value = '新增任务'
  isEdit.value = false
  dialogVisible.value = true
}

// 编辑任务
const handleEdit = (row: any) => {
  Object.assign(jobForm, {
    id: row.id,
    name: row.name,
    templateId: row.templateId,
    taskType: row.taskType,
    targetHosts: row.targetHosts || '',
    parameters: row.parameters || '',
    executeTime: row.executeTime || ''
  })
  dialogTitle.value = '编辑任务'
  isEdit.value = true
  dialogVisible.value = true
}

// 查看任务
const handleView = (row: any) => {
  currentJob.value = { ...row }
  viewDialogVisible.value = true
}

// 删除任务
const handleDelete = async (row: any) => {
  ElMessageBox.confirm(`确定要删除任务"${row.name}"吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteJobTask(row.id)
      ElMessage.success('删除成功')
      loadJobs()
    } catch (error: any) {
      ElMessage.error(error.message || '删除失败')
    }
  }).catch(() => {})
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      const data = { ...jobForm }

      if (isEdit.value) {
        await updateJobTask(data.id, data)
        ElMessage.success('更新成功')
      } else {
        await createJobTask(data)
        ElMessage.success('创建成功')
      }

      dialogVisible.value = false
      loadJobs()
    } catch (error: any) {
      ElMessage.error(error.message || (isEdit.value ? '更新失败' : '创建失败'))
    } finally {
      submitting.value = false
    }
  })
}

// 对话框关闭事件
const handleDialogClose = () => {
  formRef.value?.resetFields()
  resetForm()
}

onMounted(() => {
  loadJobs()
  loadTemplates()
})
</script>

<style scoped>
.jobs-container {
  padding: 0;
  background-color: transparent;
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
  padding: 16px 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.page-title-group {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.page-title-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #000 0%, #1a1a1a 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #d4af37;
  font-size: 22px;
  flex-shrink: 0;
  border: 1px solid #d4af37;
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  line-height: 1.3;
}

.page-subtitle {
  margin: 4px 0 0 0;
  font-size: 13px;
  color: #909399;
  line-height: 1.4;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

/* 搜索栏 */
.search-bar {
  margin-bottom: 12px;
  padding: 12px 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.search-inputs {
  display: flex;
  gap: 12px;
  flex: 1;
}

.search-input {
  width: 280px;
}

.search-actions {
  display: flex;
  gap: 10px;
}

.reset-btn {
  background: #f5f7fa;
  border-color: #dcdfe6;
  color: #606266;
}

.reset-btn:hover {
  background: #e6e8eb;
  border-color: #c0c4cc;
}

/* 搜索框样式 */
.search-bar :deep(.el-input__wrapper) {
  border-radius: 8px;
  border: 1px solid #dcdfe6;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  background-color: #fff;
}

.search-bar :deep(.el-input__wrapper:hover) {
  border-color: #d4af37;
  box-shadow: 0 2px 8px rgba(212, 175, 55, 0.15);
}

.search-bar :deep(.el-input__wrapper.is-focus) {
  border-color: #d4af37;
  box-shadow: 0 2px 12px rgba(212, 175, 55, 0.25);
}

.search-icon {
  color: #d4af37;
}

/* 表格容器 */
.table-wrapper {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.modern-table {
  width: 100%;
}

.modern-table :deep(.el-table__body-wrapper) {
  border-radius: 0 0 12px 12px;
}

.modern-table :deep(.el-table__row) {
  transition: background-color 0.2s ease;
}

.modern-table :deep(.el-table__row:hover) {
  background-color: #f8fafc !important;
}

/* 任务名称单元格 */
.job-name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.job-icon {
  font-size: 18px;
  color: #409eff;
  flex-shrink: 0;
}

.job-name {
  font-weight: 500;
}

.description-text {
  color: #606266;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: 8px;
  align-items: center;
  justify-content: center;
}

.action-btn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.action-btn :deep(.el-icon) {
  font-size: 16px;
}

.action-btn:hover {
  transform: scale(1.1);
}

.action-view:hover {
  background-color: #e8f4ff;
  color: #409eff;
}

.action-edit:hover {
  background-color: #e8f4ff;
  color: #409eff;
}

.action-delete:hover {
  background-color: #fee;
  color: #f56c6c;
}

.black-button {
  background-color: #000000 !important;
  color: #ffffff !important;
  border-color: #000000 !important;
  border-radius: 8px;
  padding: 10px 20px;
  font-weight: 500;
}

.black-button:hover {
  background-color: #333333 !important;
  border-color: #333333 !important;
}

/* 分页器 */
.pagination-container {
  padding: 12px 16px;
  display: flex;
  justify-content: flex-end;
  border-top: 1px solid #f0f0f0;
}

/* 对话框样式 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.job-edit-dialog),
:deep(.job-view-dialog) {
  border-radius: 12px;
}

:deep(.job-edit-dialog .el-dialog__header),
:deep(.job-view-dialog .el-dialog__header) {
  padding: 20px 24px 16px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.job-edit-dialog .el-dialog__body),
:deep(.job-view-dialog .el-dialog__body) {
  padding: 24px;
}

:deep(.job-edit-dialog .el-dialog__footer),
:deep(.job-view-dialog .el-dialog__footer) {
  padding: 16px 24px;
  border-top: 1px solid #f0f0f0;
}

/* 标签样式 */
:deep(.el-tag) {
  border-radius: 6px;
  padding: 4px 10px;
  font-weight: 500;
}

/* 输入框样式 */
:deep(.el-input__wrapper),
:deep(.el-textarea__inner) {
  border-radius: 8px;
}

:deep(.el-select .el-input__wrapper) {
  border-radius: 8px;
}

/* 响应式对话框 */
:deep(.responsive-dialog) {
  max-width: 900px;
  min-width: 500px;
}

@media (max-width: 768px) {
  :deep(.responsive-dialog .el-dialog) {
    width: 95% !important;
    max-width: none;
    min-width: auto;
  }
}
</style>
