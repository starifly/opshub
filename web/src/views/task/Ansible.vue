<template>
  <div class="ansible-container">
    <!-- 页面标题和操作按钮 -->
    <div class="page-header">
      <div class="page-title-group">
        <div class="page-title-icon">
          <el-icon><Setting /></el-icon>
        </div>
        <div>
          <h2 class="page-title">Ansible任务</h2>
          <p class="page-subtitle">管理Ansible自动化任务，支持Playbook执行与管理</p>
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
          v-model="searchForm.status"
          placeholder="任务状态"
          clearable
          class="search-input"
          @change="handleSearch"
        >
          <el-option label="等待中" value="pending" />
          <el-option label="运行中" value="running" />
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
        :data="ansibleList"
        v-loading="loading"
        class="modern-table"
        :header-cell-style="{ background: '#fafbfc', color: '#606266', fontWeight: '600' }"
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />

        <el-table-column label="任务名称" prop="name" min-width="180">
          <template #default="{ row }">
            <div class="ansible-name-cell">
              <el-icon class="ansible-icon"><Operation /></el-icon>
              <span class="ansible-name">{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="Playbook路径" prop="playbookPath" min-width="200">
          <template #default="{ row }">
            <span class="description-text">{{ row.playbookPath || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column label="并发数" prop="fork" width="100" align="center">
          <template #default="{ row }">
            <el-tag type="info" effect="plain">{{ row.fork }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column label="超时(秒)" prop="timeout" width="100" align="center">
          <template #default="{ row }">
            <span class="description-text">{{ row.timeout }}</span>
          </template>
        </el-table-column>

        <el-table-column label="状态" prop="status" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" effect="dark">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="lastRunTime" label="最后运行时间" width="180" />

        <el-table-column label="操作" width="180" fixed="right" align="center">
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
          @current-change="loadAnsibleTasks"
        />
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="70%"
      class="ansible-edit-dialog responsive-dialog"
      :close-on-click-modal="false"
      @close="handleDialogClose"
    >
      <el-form :model="ansibleForm" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="ansibleForm.name" placeholder="请输入任务名称" />
        </el-form-item>

        <el-form-item label="Playbook类型" prop="playbookType">
          <el-radio-group v-model="playbookType">
            <el-radio label="path">文件路径</el-radio>
            <el-radio label="content">直接内容</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item
          v-if="playbookType === 'path'"
          label="Playbook路径"
          prop="playbookPath"
        >
          <el-input
            v-model="ansibleForm.playbookPath"
            placeholder="请输入Playbook文件路径，如：/etc/ansible/playbooks/site.yml"
          />
        </el-form-item>

        <el-form-item
          v-if="playbookType === 'content'"
          label="Playbook内容"
          prop="playbookContent"
        >
          <el-input
            v-model="ansibleForm.playbookContent"
            type="textarea"
            :rows="10"
            placeholder="请输入Playbook内容（YAML格式）"
            style="font-family: monospace;"
          />
        </el-form-item>

        <el-form-item label="Inventory" prop="inventory">
          <el-input
            v-model="ansibleForm.inventory"
            type="textarea"
            :rows="4"
            placeholder="请输入Inventory内容（INI或YAML格式）"
          />
        </el-form-item>

        <el-form-item label="额外变量" prop="extraVars">
          <el-input
            v-model="ansibleForm.extraVars"
            type="textarea"
            :rows="3"
            placeholder="请输入额外变量（JSON或YAML格式）"
          />
        </el-form-item>

        <el-form-item label="Tags" prop="tags">
          <el-input
            v-model="ansibleForm.tags"
            placeholder="请输入Tags，多个标签用逗号分隔"
          />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="并发数" prop="fork">
              <el-input-number v-model="ansibleForm.fork" :min="1" :max="100" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="超时时间(秒)" prop="timeout">
              <el-input-number v-model="ansibleForm.timeout" :min="1" :max="7200" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="详细程度" prop="verbose">
              <el-select v-model="ansibleForm.verbose" placeholder="请选择" style="width: 100%">
                <el-option label="正常" value="v" />
                <el-option label="详细" value="vv" />
                <el-option label="更详细" value="vvv" />
                <el-option label="调试" value="vvvv" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
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
      width="70%"
      class="ansible-view-dialog responsive-dialog"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="任务ID">{{ currentAnsible.id }}</el-descriptions-item>
        <el-descriptions-item label="任务名称">{{ currentAnsible.name }}</el-descriptions-item>
        <el-descriptions-item label="Playbook路径" :span="2">
          {{ currentAnsible.playbookPath || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="并发数">{{ currentAnsible.fork }}</el-descriptions-item>
        <el-descriptions-item label="超时时间">{{ currentAnsible.timeout }} 秒</el-descriptions-item>
        <el-descriptions-item label="详细程度">{{ currentAnsible.verbose }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentAnsible.status)" effect="dark">
            {{ getStatusLabel(currentAnsible.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="最后运行时间" :span="2">
          {{ currentAnsible.lastRunTime || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="最后运行结果" :span="2" v-if="currentAnsible.lastRunResult">
          <pre style="margin: 0; white-space: pre-wrap; max-height: 200px; overflow-y: auto;">{{ currentAnsible.lastRunResult }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="Inventory" :span="2">
          <pre style="margin: 0; white-space: pre-wrap; max-height: 200px; overflow-y: auto;">{{ currentAnsible.inventory || '-' }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="额外变量" :span="2">
          <pre style="margin: 0; white-space: pre-wrap;">{{ currentAnsible.extraVars || '-' }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="Tags" :span="2">{{ currentAnsible.tags || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentAnsible.createdAt }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ currentAnsible.updatedAt }}</el-descriptions-item>
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
  Setting,
  Operation,
  View
} from '@element-plus/icons-vue'
import { getAnsibleTaskList, createAnsibleTask, updateAnsibleTask, deleteAnsibleTask } from '@/api/task'

// 加载状态
const loading = ref(false)
const submitting = ref(false)

// 对话框状态
const dialogVisible = ref(false)
const viewDialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)

// Playbook类型
const playbookType = ref('path')

// 表单引用
const formRef = ref<FormInstance>()

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// Ansible任务列表
const ansibleList = ref<any[]>([])

// 当前查看的任务
const currentAnsible = ref<any>({})

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: ''
})

// Ansible任务表单
const ansibleForm = reactive({
  id: 0,
  name: '',
  playbookContent: '',
  playbookPath: '',
  inventory: '',
  extraVars: '',
  tags: '',
  fork: 5,
  timeout: 300,
  verbose: 'v'
})

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: '请输入任务名称', trigger: 'blur' },
    { min: 2, max: 100, message: '任务名称长度在 2 到 100 个字符', trigger: 'blur' }
  ]
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
    running: '运行中',
    success: '成功',
    failed: '失败'
  }
  return labelMap[status] || status
}

// 搜索处理
const handleSearch = () => {
  pagination.page = 1
  loadAnsibleTasks()
}

// 重置搜索
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  loadAnsibleTasks()
}

// 加载Ansible任务列表
const loadAnsibleTasks = async () => {
  loading.value = true
  try {
    const res = await getAnsibleTaskList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: searchForm.keyword,
      status: searchForm.status
    })
    ansibleList.value = res.list || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('获取Ansible任务列表失败:', error)
    ElMessage.error('获取Ansible任务列表失败')
  } finally {
    loading.value = false
  }
}

// 重置表单
const resetForm = () => {
  ansibleForm.id = 0
  ansibleForm.name = ''
  ansibleForm.playbookContent = ''
  ansibleForm.playbookPath = ''
  ansibleForm.inventory = ''
  ansibleForm.extraVars = ''
  ansibleForm.tags = ''
  ansibleForm.fork = 5
  ansibleForm.timeout = 300
  ansibleForm.verbose = 'v'
  playbookType.value = 'path'
  formRef.value?.clearValidate()
}

// 新增任务
const handleAdd = () => {
  resetForm()
  dialogTitle.value = '新增Ansible任务'
  isEdit.value = false
  dialogVisible.value = true
}

// 编辑任务
const handleEdit = (row: any) => {
  Object.assign(ansibleForm, {
    id: row.id,
    name: row.name,
    playbookContent: row.playbookContent || '',
    playbookPath: row.playbookPath || '',
    inventory: row.inventory || '',
    extraVars: row.extraVars || '',
    tags: row.tags || '',
    fork: row.fork,
    timeout: row.timeout,
    verbose: row.verbose
  })
  // 设置playbook类型
  if (row.playbookContent) {
    playbookType.value = 'content'
  } else {
    playbookType.value = 'path'
  }
  dialogTitle.value = '编辑Ansible任务'
  isEdit.value = true
  dialogVisible.value = true
}

// 查看任务
const handleView = (row: any) => {
  currentAnsible.value = { ...row }
  viewDialogVisible.value = true
}

// 删除任务
const handleDelete = async (row: any) => {
  ElMessageBox.confirm(`确定要删除Ansible任务"${row.name}"吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteAnsibleTask(row.id)
      ElMessage.success('删除成功')
      loadAnsibleTasks()
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
      const data: any = { ...ansibleForm }

      // 根据playbook类型清理数据
      if (playbookType.value === 'path') {
        data.playbookContent = ''
      } else {
        data.playbookPath = ''
      }

      if (isEdit.value) {
        await updateAnsibleTask(data.id, data)
        ElMessage.success('更新成功')
      } else {
        await createAnsibleTask(data)
        ElMessage.success('创建成功')
      }

      dialogVisible.value = false
      loadAnsibleTasks()
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
  loadAnsibleTasks()
})
</script>

<style scoped>
.ansible-container {
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

/* Ansible任务名称单元格 */
.ansible-name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.ansible-icon {
  font-size: 18px;
  color: #e6a23c;
  flex-shrink: 0;
}

.ansible-name {
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

:deep(.ansible-edit-dialog),
:deep(.ansible-view-dialog) {
  border-radius: 12px;
}

:deep(.ansible-edit-dialog .el-dialog__header),
:deep(.ansible-view-dialog .el-dialog__header) {
  padding: 20px 24px 16px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.ansible-edit-dialog .el-dialog__body),
:deep(.ansible-view-dialog .el-dialog__body) {
  padding: 24px;
}

:deep(.ansible-edit-dialog .el-dialog__footer),
:deep(.ansible-view-dialog .el-dialog__footer) {
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
  max-width: 1100px;
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
