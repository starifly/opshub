<template>
  <div class="data-container">
    <!-- 页面标题和操作按钮 -->
    <div class="page-header">
      <div class="page-title-group">
        <div class="page-title-icon">
          <el-icon><DataLine /></el-icon>
        </div>
        <div>
          <h2 class="page-title">数据管理</h2>
          <p class="page-subtitle">管理监控数据和配置信息</p>
        </div>
      </div>
      <div class="header-actions">
        <el-button @click="loadData">
          <el-icon style="margin-right: 6px;"><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 数据统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-icon" style="background: #ecf5ff;">
          <el-icon style="color: #409eff;"><Monitor /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">主机总数</div>
          <div class="stat-value">156</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #e8f5e9;">
          <el-icon style="color: #67c23a;"><Collection /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">业务分组</div>
          <div class="stat-value">12</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #fdf6ec;">
          <el-icon style="color: #e6a23c;"><DataLine /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">监控指标</div>
          <div class="stat-value">2,450</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: #fef0f0;">
          <el-icon style="color: #f56c6c;"><Warning /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">告警数量</div>
          <div class="stat-value">8</div>
        </div>
      </div>
    </div>

    <!-- 表格容器 -->
    <div class="table-wrapper">
      <el-table
        :data="tableData"
        v-loading="loading"
        class="modern-table"
        :header-cell-style="{ background: '#fafbfc', color: '#606266', fontWeight: '600' }"
      >
        <el-table-column label="数据源" prop="source" min-width="200">
          <template #default="{ row }">
            <span style="display: inline-flex; align-items: center;">
              <el-icon style="color: #e6a23c; margin-right: 8px;"><DataLine /></el-icon>
              {{ row.source }}
            </span>
          </template>
        </el-table-column>

        <el-table-column label="数据类型" prop="type" min-width="150" />

        <el-table-column label="大小" prop="size" width="120" align="right" />

        <el-table-column label="记录数" prop="records" width="120" align="right">
          <template #default="{ row }">
            <el-tag type="primary">{{ row.records.toLocaleString() }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column label="更新时间" prop="updateTime" min-width="180" />

        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" effect="dark">
              {{ row.status === 'active' ? '活跃' : '归档' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" fixed="right" align="center">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-tooltip content="查看" placement="top">
                <el-button link class="action-btn action-view" @click="handleView(row)">
                  <el-icon><View /></el-icon>
                </el-button>
              </el-tooltip>
              <el-tooltip content="清理" placement="top">
                <el-button link class="action-btn action-clean" @click="handleClean(row)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Monitor,
  Collection,
  DataLine,
  Warning,
  View,
  Delete,
  Refresh
} from '@element-plus/icons-vue'

const loading = ref(false)
const tableData = ref([])

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const loadData = async () => {
  loading.value = true
  // TODO: 调用API获取数据
  setTimeout(() => {
    tableData.value = [
      { source: '主机监控数据', type: '时序数据', size: '2.5 GB', records: 1250000, status: 'active', updateTime: '2024-01-15 10:00:00' },
      { source: '容器指标数据', type: '时序数据', size: '1.8 GB', records: 890000, status: 'active', updateTime: '2024-01-15 10:05:00' },
      { source: '告警历史记录', type: '日志数据', size: '450 MB', records: 310000, status: 'active', updateTime: '2024-01-15 09:30:00' },
      { source: '操作审计日志', type: '日志数据', size: '280 MB', records: 156000, status: 'active', updateTime: '2024-01-15 09:15:00' },
      { source: '历史配置备份', type: '文档数据', size: '120 MB', records: 4500, status: 'archived', updateTime: '2024-01-14 18:00:00' }
    ]
    pagination.total = 5
    loading.value = false
  }, 500)
}

const handleView = (row: any) => {
  ElMessage.info('查看数据功能开发中...')
}

const handleClean = (row: any) => {
  ElMessage.info('清理数据功能开发中...')
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.data-container {
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

/* 统计卡片 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 12px;
  margin-bottom: 12px;
}

.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 24px;
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

.modern-table :deep(.el-table__body-wrapper) {
  border-radius: 0 0 12px 12px;
}

.modern-table :deep(.el-table__row) {
  transition: background-color 0.2s ease;
}

.modern-table :deep(.el-table__row:hover) {
  background-color: #f8fafc !important;
}

.pagination-wrapper {
  padding: 16px;
  display: flex;
  justify-content: flex-end;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: 8px;
  align-items: center;
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

.action-clean:hover {
  background-color: #fef0f0;
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
</style>
