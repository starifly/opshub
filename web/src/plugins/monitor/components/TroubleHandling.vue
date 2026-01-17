<template>
  <div class="trouble-handling-container">
    <!-- 页面标题和操作按钮 -->
    <div class="page-header">
      <div class="page-title-group">
        <div class="page-title-icon">
          <el-icon><Warning /></el-icon>
        </div>
        <div>
          <h2 class="page-title">故障处理</h2>
          <p class="page-subtitle">管理和处理系统故障告警，跟踪故障处理进度</p>
        </div>
      </div>
      <div class="header-actions">
        <el-button class="black-button" @click="handleAdd">
          <el-icon style="margin-right: 6px;"><Plus /></el-icon>
          新建工单
        </el-button>
        <el-button @click="loadData">
          <el-icon style="margin-right: 6px;"><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <div class="search-inputs">
        <el-input
          v-model="searchForm.keyword"
          placeholder="搜索故障标题或描述..."
          clearable
          class="search-input"
        >
          <template #prefix>
            <el-icon class="search-icon"><Search /></el-icon>
          </template>
        </el-input>

        <el-select
          v-model="searchForm.level"
          placeholder="故障级别"
          clearable
          class="search-input"
        >
          <el-option label="严重" value="critical" />
          <el-option label="高" value="high" />
          <el-option label="中" value="medium" />
          <el-option label="低" value="low" />
        </el-select>

        <el-select
          v-model="searchForm.status"
          placeholder="处理状态"
          clearable
          class="search-input"
        >
          <el-option label="待处理" value="pending" />
          <el-option label="处理中" value="processing" />
          <el-option label="已解决" value="resolved" />
          <el-option label="已关闭" value="closed" />
        </el-select>
      </div>

      <div class="search-actions">
        <el-button class="reset-btn" @click="handleReset">
          <el-icon style="margin-right: 4px;"><RefreshLeft /></el-icon>
          重置
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-icon stat-icon-primary">
          <el-icon><Document /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">工单总数</div>
          <div class="stat-value">{{ stats.total }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon stat-icon-warning">
          <el-icon><Clock /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">待处理</div>
          <div class="stat-value">{{ stats.pending }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon stat-icon-info">
          <el-icon><Loading /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">处理中</div>
          <div class="stat-value">{{ stats.processing }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon stat-icon-success">
          <el-icon><CircleCheck /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">已解决</div>
          <div class="stat-value">{{ stats.resolved }}</div>
        </div>
      </div>
    </div>

    <!-- 表格容器 -->
    <div class="table-wrapper">
      <el-table
        :data="filteredData"
        v-loading="loading"
        class="modern-table"
        :header-cell-style="{ background: '#fafbfc', color: '#606266', fontWeight: '600' }"
      >
        <el-table-column label="工单编号" prop="ticketNo" width="140" />

        <el-table-column label="故障标题" prop="title" min-width="250" />

        <el-table-column label="故障级别" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.level === 'critical'" type="danger" effect="dark">严重</el-tag>
            <el-tag v-else-if="row.level === 'high'" type="danger" effect="plain">高</el-tag>
            <el-tag v-else-if="row.level === 'medium'" type="warning" effect="plain">中</el-tag>
            <el-tag v-else type="info" effect="plain">低</el-tag>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.status === 'pending'" type="warning">待处理</el-tag>
            <el-tag v-else-if="row.status === 'processing'" type="primary">处理中</el-tag>
            <el-tag v-else-if="row.status === 'resolved'" type="success">已解决</el-tag>
            <el-tag v-else type="info">已关闭</el-tag>
          </template>
        </el-table-column>

        <el-table-column label="优先级" width="100" align="center">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.priority }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column label="责任人" prop="assignee" width="120" />

        <el-table-column label="创建时间" prop="createTime" width="180" />

        <el-table-column label="截止时间" prop="deadline" width="180" />

        <el-table-column label="操作" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-tooltip content="查看详情" placement="top">
                <el-button link class="action-btn action-view" @click="handleView(row)">
                  <el-icon><View /></el-icon>
                </el-button>
              </el-tooltip>
              <el-tooltip content="处理" placement="top">
                <el-button link class="action-btn action-handle" @click="handleProcess(row)">
                  <el-icon><Tools /></el-icon>
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
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="700px"
      class="trouble-edit-dialog"
      :close-on-click-modal="false"
      @close="handleDialogClose"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="故障标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入故障标题" />
        </el-form-item>
        <el-form-item label="故障级别" prop="level">
          <el-radio-group v-model="form.level">
            <el-radio label="critical" border>严重</el-radio>
            <el-radio label="high" border>高</el-radio>
            <el-radio label="medium" border>中</el-radio>
            <el-radio label="low" border>低</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-select v-model="form.priority" placeholder="请选择优先级" style="width: 100%;">
            <el-option label="P0 - 紧急" value="P0" />
            <el-option label="P1 - 高" value="P1" />
            <el-option label="P2 - 中" value="P2" />
            <el-option label="P3 - 低" value="P3" />
          </el-select>
        </el-form-item>
        <el-form-item label="故障描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            placeholder="请详细描述故障情况"
          />
        </el-form-item>
        <el-form-item label="责任人" prop="assignee">
          <el-input v-model="form.assignee" placeholder="请输入责任人" />
        </el-form-item>
        <el-form-item label="截止时间" prop="deadline">
          <el-date-picker
            v-model="form.deadline"
            type="datetime"
            placeholder="选择截止时间"
            style="width: 100%;"
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
import {
  Plus,
  Search,
  RefreshLeft,
  Refresh,
  Edit,
  Delete,
  View,
  Warning,
  Document,
  Clock,
  Loading,
  CircleCheck,
  Tools
} from '@element-plus/icons-vue'

const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitting = ref(false)
const formRef = ref<FormInstance>()

// 搜索表单
const searchForm = reactive({
  keyword: '',
  level: '',
  status: ''
})

// 统计数据
const stats = ref({
  total: 0,
  pending: 0,
  processing: 0,
  resolved: 0
})

// 表单数据
const form = reactive({
  id: 0,
  title: '',
  level: 'medium',
  priority: 'P2',
  description: '',
  assignee: '',
  deadline: ''
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入故障标题', trigger: 'blur' }],
  level: [{ required: true, message: '请选择故障级别', trigger: 'change' }],
  priority: [{ required: true, message: '请选择优先级', trigger: 'change' }],
  description: [{ required: true, message: '请输入故障描述', trigger: 'blur' }],
  assignee: [{ required: true, message: '请输入责任人', trigger: 'blur' }],
  deadline: [{ required: true, message: '请选择截止时间', trigger: 'change' }]
}

// 模拟数据
const tableData = ref([
  {
    id: 1,
    ticketNo: 'T20250116001',
    title: '服务器CPU使用率过高告警',
    level: 'high',
    status: 'processing',
    priority: 'P1',
    assignee: '张三',
    createTime: '2025-01-16 09:30:00',
    deadline: '2025-01-16 18:00:00'
  },
  {
    id: 2,
    ticketNo: 'T20250116002',
    title: '数据库连接池耗尽',
    level: 'critical',
    status: 'pending',
    priority: 'P0',
    assignee: '李四',
    createTime: '2025-01-16 10:15:00',
    deadline: '2025-01-16 12:00:00'
  },
  {
    id: 3,
    ticketNo: 'T20250115099',
    title: 'SSL证书即将到期',
    level: 'medium',
    status: 'resolved',
    priority: 'P2',
    assignee: '王五',
    createTime: '2025-01-15 14:20:00',
    deadline: '2025-01-17 18:00:00'
  }
])

// 过滤后的数据
const filteredData = computed(() => {
  if (!searchForm.keyword && !searchForm.level && !searchForm.status) {
    return tableData.value
  }
  return tableData.value.filter(item => {
    const matchKeyword = !searchForm.keyword ||
      item.title.toLowerCase().includes(searchForm.keyword.toLowerCase())
    const matchLevel = !searchForm.level || item.level === searchForm.level
    const matchStatus = !searchForm.status || item.status === searchForm.status
    return matchKeyword && matchLevel && matchStatus
  })
})

// 重置搜索
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.level = ''
  searchForm.status = ''
}

// 加载数据
const loadData = () => {
  loading.value = true
  setTimeout(() => {
    stats.value = {
      total: tableData.value.length,
      pending: tableData.value.filter(i => i.status === 'pending').length,
      processing: tableData.value.filter(i => i.status === 'processing').length,
      resolved: tableData.value.filter(i => i.status === 'resolved').length
    }
    loading.value = false
    ElMessage.success('数据已刷新')
  }, 500)
}

// 新增
const handleAdd = () => {
  dialogTitle.value = '新建故障工单'
  Object.assign(form, {
    id: 0,
    title: '',
    level: 'medium',
    priority: 'P2',
    description: '',
    assignee: '',
    deadline: ''
  })
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = '编辑故障工单'
  Object.assign(form, row)
  dialogVisible.value = true
}

// 查看详情
const handleView = (row: any) => {
  ElMessage.info(`查看工单 ${row.ticketNo} 的详情`)
}

// 处理工单
const handleProcess = (row: any) => {
  ElMessage.success(`开始处理工单 ${row.ticketNo}`)
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除该故障工单吗？', '提示', { type: 'warning' })
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') console.error(error)
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        // TODO: 调用API
        ElMessage.success(form.id ? '更新成功' : '创建成功')
        dialogVisible.value = false
        loadData()
      } catch (error) {
        ElMessage.error('操作失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

// 对话框关闭
const handleDialogClose = () => {
  formRef.value?.resetFields()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.trouble-handling-container {
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

/* 统计卡片 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin-bottom: 12px;
}

.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}

.stat-icon-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}

.stat-icon-warning {
  background: linear-gradient(135deg, #e6a23c 0%, #d9972c 100%);
  color: #fff;
}

.stat-icon-info {
  background: linear-gradient(135deg, #409eff 0%, #3a8ee6 100%);
  color: #fff;
}

.stat-icon-success {
  background: linear-gradient(135deg, #67c23a 0%, #5daf34 100%);
  color: #fff;
}

.stat-content {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
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

.action-handle:hover {
  background-color: #e8f5e9;
  color: #67c23a;
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

/* 对话框样式 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.trouble-edit-dialog) {
  border-radius: 12px;
}

:deep(.trouble-edit-dialog .el-dialog__header) {
  padding: 20px 24px 16px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.trouble-edit-dialog .el-dialog__body) {
  padding: 24px;
}

:deep(.trouble-edit-dialog .el-dialog__footer) {
  padding: 16px 24px;
  border-top: 1px solid #f0f0f0;
}
</style>
