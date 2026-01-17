<template>
  <div class="templates-container">
    <!-- 页面标题和操作按钮 -->
    <div class="page-header">
      <div class="page-title-group">
        <div class="page-title-icon">
          <el-icon><Document /></el-icon>
        </div>
        <div>
          <h2 class="page-title">任务模板</h2>
          <p class="page-subtitle">管理任务模板，支持模板创建、编辑与复用</p>
        </div>
      </div>
      <div class="header-actions">
        <el-button class="black-button" @click="handleAdd">
          <el-icon style="margin-right: 6px;"><Plus /></el-icon>
          新增模板
        </el-button>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <div class="search-inputs">
        <el-input
          v-model="searchForm.keyword"
          placeholder="搜索模板名称..."
          clearable
          class="search-input"
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon class="search-icon"><Search /></el-icon>
          </template>
        </el-input>

        <el-select
          v-model="searchForm.category"
          placeholder="模板分类"
          clearable
          class="search-input"
          @change="handleSearch"
        >
          <el-option label="系统管理" value="system" />
          <el-option label="应用部署" value="deploy" />
          <el-option label="监控巡检" value="monitor" />
          <el-option label="备份恢复" value="backup" />
        </el-select>

        <el-select
          v-model="searchForm.platform"
          placeholder="适用平台"
          clearable
          class="search-input"
          @change="handleSearch"
        >
          <el-option label="Linux" value="linux" />
          <el-option label="Windows" value="windows" />
          <el-option label="通用" value="general" />
        </el-select>

        <el-select
          v-model="searchForm.status"
          placeholder="启用状态"
          clearable
          class="search-input"
          @change="handleSearch"
        >
          <el-option label="启用" :value="1" />
          <el-option label="禁用" :value="0" />
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
        :data="templateList"
        v-loading="loading"
        class="modern-table"
        :header-cell-style="{ background: '#fafbfc', color: '#606266', fontWeight: '600' }"
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />

        <el-table-column label="模板名称" prop="name" min-width="180">
          <template #default="{ row }">
            <div class="template-name-cell">
              <el-icon class="template-icon"><Tickets /></el-icon>
              <span class="template-name">{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="模板编码" prop="code" min-width="150">
          <template #header>
            <span class="header-with-icon">
              <el-icon class="header-icon header-icon-gold"><Key /></el-icon>
              模板编码
            </span>
          </template>
        </el-table-column>

        <el-table-column label="分类" prop="category" width="120">
          <template #default="{ row }">
            <el-tag :type="getCategoryType(row.category)" effect="plain">
              {{ getCategoryLabel(row.category) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="适用平台" prop="platform" width="120">
          <template #default="{ row }">
            <span class="description-text">{{ row.platform || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column label="超时时间(秒)" prop="timeout" width="120" align="center">
          <template #default="{ row }">
            <span class="description-text">{{ row.timeout }}</span>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" effect="dark">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="sort" label="排序" width="80" align="center" />

        <el-table-column prop="createdAt" label="创建时间" min-width="180" />

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
          @current-change="loadTemplates"
        />
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="70%"
      class="template-edit-dialog responsive-dialog"
      :close-on-click-modal="false"
      @close="handleDialogClose"
    >
      <el-form :model="templateForm" :rules="rules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="模板名称" prop="name">
              <el-input v-model="templateForm.name" placeholder="请输入模板名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="模板编码" prop="code">
              <el-input v-model="templateForm.code" placeholder="请输入模板编码" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="模板分类" prop="category">
              <el-select v-model="templateForm.category" placeholder="请选择分类" style="width: 100%">
                <el-option label="系统管理" value="system" />
                <el-option label="应用部署" value="deploy" />
                <el-option label="监控巡检" value="monitor" />
                <el-option label="备份恢复" value="backup" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="适用平台" prop="platform">
              <el-select v-model="templateForm.platform" placeholder="请选择平台" style="width: 100%">
                <el-option label="Linux" value="linux" />
                <el-option label="Windows" value="windows" />
                <el-option label="通用" value="general" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="超时时间" prop="timeout">
              <el-input-number v-model="templateForm.timeout" :min="1" :max="3600" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="templateForm.sort" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="描述" prop="description">
          <el-input
            v-model="templateForm.description"
            type="textarea"
            :rows="2"
            placeholder="请输入模板描述"
          />
        </el-form-item>

        <el-form-item label="模板内容" prop="content">
          <el-input
            v-model="templateForm.content"
            type="textarea"
            :rows="8"
            placeholder="请输入模板内容（脚本或命令）"
          />
        </el-form-item>

        <el-form-item label="变量定义" prop="variables">
          <el-input
            v-model="templateForm.variables"
            type="textarea"
            :rows="4"
            placeholder="请输入变量定义（JSON格式）"
          />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="templateForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
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
      title="模板详情"
      width="70%"
      class="template-view-dialog responsive-dialog"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="模板ID">{{ currentTemplate.id }}</el-descriptions-item>
        <el-descriptions-item label="模板名称">{{ currentTemplate.name }}</el-descriptions-item>
        <el-descriptions-item label="模板编码">{{ currentTemplate.code }}</el-descriptions-item>
        <el-descriptions-item label="分类">
          <el-tag :type="getCategoryType(currentTemplate.category)" effect="plain">
            {{ getCategoryLabel(currentTemplate.category) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="适用平台">{{ currentTemplate.platform || '-' }}</el-descriptions-item>
        <el-descriptions-item label="超时时间">{{ currentTemplate.timeout }} 秒</el-descriptions-item>
        <el-descriptions-item label="排序">{{ currentTemplate.sort }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentTemplate.status === 1 ? 'success' : 'danger'" effect="dark">
            {{ currentTemplate.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentTemplate.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="模板内容" :span="2">
          <pre style="margin: 0; white-space: pre-wrap; max-height: 300px; overflow-y: auto;">{{ currentTemplate.content || '-' }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="变量定义" :span="2">
          <pre style="margin: 0; white-space: pre-wrap;">{{ currentTemplate.variables || '-' }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentTemplate.createdAt }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ currentTemplate.updatedAt }}</el-descriptions-item>
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
  Document,
  Tickets,
  Key,
  View
} from '@element-plus/icons-vue'
import { getJobTemplateList, createJobTemplate, updateJobTemplate, deleteJobTemplate } from '@/api/task'

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

// 模板列表
const templateList = ref<any[]>([])

// 当前查看的模板
const currentTemplate = ref<any>({})

// 搜索表单
const searchForm = reactive({
  keyword: '',
  category: '',
  platform: '',
  status: undefined as number | undefined
})

// 模板表单
const templateForm = reactive({
  id: 0,
  name: '',
  code: '',
  description: '',
  content: '',
  variables: '',
  category: '',
  platform: '',
  timeout: 300,
  sort: 0,
  status: 1
})

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: '请输入模板名称', trigger: 'blur' },
    { min: 2, max: 100, message: '模板名称长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入模板编码', trigger: 'blur' },
    { min: 2, max: 50, message: '模板编码长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择模板分类', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请输入模板内容', trigger: 'blur' }
  ],
  timeout: [
    { required: true, message: '请输入超时时间', trigger: 'blur' }
  ]
}

// 获取分类类型
const getCategoryType = (category: string) => {
  const typeMap: Record<string, string> = {
    system: 'info',
    deploy: 'success',
    monitor: 'warning',
    backup: 'danger'
  }
  return typeMap[category] || 'info'
}

// 获取分类标签
const getCategoryLabel = (category: string) => {
  const labelMap: Record<string, string> = {
    system: '系统管理',
    deploy: '应用部署',
    monitor: '监控巡检',
    backup: '备份恢复'
  }
  return labelMap[category] || category
}

// 搜索处理
const handleSearch = () => {
  pagination.page = 1
  loadTemplates()
}

// 重置搜索
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.category = ''
  searchForm.platform = ''
  searchForm.status = undefined
  loadTemplates()
}

// 加载模板列表
const loadTemplates = async () => {
  loading.value = true
  try {
    const res = await getJobTemplateList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: searchForm.keyword,
      category: searchForm.category,
      platform: searchForm.platform,
      status: searchForm.status
    })
    templateList.value = res.list || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('获取模板列表失败:', error)
    ElMessage.error('获取模板列表失败')
  } finally {
    loading.value = false
  }
}

// 重置表单
const resetForm = () => {
  templateForm.id = 0
  templateForm.name = ''
  templateForm.code = ''
  templateForm.description = ''
  templateForm.content = ''
  templateForm.variables = ''
  templateForm.category = ''
  templateForm.platform = ''
  templateForm.timeout = 300
  templateForm.sort = 0
  templateForm.status = 1
  formRef.value?.clearValidate()
}

// 新增模板
const handleAdd = () => {
  resetForm()
  dialogTitle.value = '新增模板'
  isEdit.value = false
  dialogVisible.value = true
}

// 编辑模板
const handleEdit = (row: any) => {
  Object.assign(templateForm, {
    id: row.id,
    name: row.name,
    code: row.code,
    description: row.description || '',
    content: row.content,
    variables: row.variables || '',
    category: row.category,
    platform: row.platform || '',
    timeout: row.timeout,
    sort: row.sort || 0,
    status: row.status
  })
  dialogTitle.value = '编辑模板'
  isEdit.value = true
  dialogVisible.value = true
}

// 查看模板
const handleView = (row: any) => {
  currentTemplate.value = { ...row }
  viewDialogVisible.value = true
}

// 删除模板
const handleDelete = async (row: any) => {
  ElMessageBox.confirm(`确定要删除模板"${row.name}"吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteJobTemplate(row.id)
      ElMessage.success('删除成功')
      loadTemplates()
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
      const data = { ...templateForm }

      if (isEdit.value) {
        await updateJobTemplate(data.id, data)
        ElMessage.success('更新成功')
      } else {
        await createJobTemplate(data)
        ElMessage.success('创建成功')
      }

      dialogVisible.value = false
      loadTemplates()
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
  loadTemplates()
})
</script>

<style scoped>
.templates-container {
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
  width: 240px;
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

/* 表头图标 */
.header-with-icon {
  display: flex;
  align-items: center;
  gap: 6px;
}

.header-icon {
  font-size: 16px;
}

.header-icon-gold {
  color: #d4af37;
}

/* 模板名称单元格 */
.template-name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.template-icon {
  font-size: 18px;
  color: #67c23a;
  flex-shrink: 0;
}

.template-name {
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

:deep(.template-edit-dialog),
:deep(.template-view-dialog) {
  border-radius: 12px;
}

:deep(.template-edit-dialog .el-dialog__header),
:deep(.template-view-dialog .el-dialog__header) {
  padding: 20px 24px 16px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.template-edit-dialog .el-dialog__body),
:deep(.template-view-dialog .el-dialog__body) {
  padding: 24px;
}

:deep(.template-edit-dialog .el-dialog__footer),
:deep(.template-view-dialog .el-dialog__footer) {
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
