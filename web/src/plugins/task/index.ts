import type { Plugin, PluginMenuConfig, PluginRouteConfig } from '../types'
import { pluginManager } from '../manager'

/**
 * Task 任务中心插件
 * 提供任务作业、任务模板和Ansible任务管理功能
 */
class TaskPlugin implements Plugin {
  name = 'task'
  description = '任务中心插件，提供任务作业、任务模板和Ansible任务管理功能'
  version = '1.0.0'
  author = 'OpsHub Team'

  /**
   * 安装插件
   */
  async install() {
    console.log('Task 插件安装中...')
  }

  /**
   * 卸载插件
   */
  async uninstall() {
    console.log('Task 插件卸载中...')
  }

  /**
   * 获取插件菜单配置
   */
  getMenus(): PluginMenuConfig[] {
    const parentPath = '/task'

    return [
      {
        name: '任务中心',
        path: parentPath,
        icon: 'Tickets',
        sort: 90,
        hidden: false,
        parentPath: '',
      },
      {
        name: '任务作业',
        path: '/task/jobs',
        icon: 'List',
        sort: 1,
        hidden: false,
        parentPath: parentPath,
      },
      {
        name: '任务模板',
        path: '/task/templates',
        icon: 'Document',
        sort: 2,
        hidden: false,
        parentPath: parentPath,
      },
      {
        name: 'Ansible任务',
        path: '/task/ansible',
        icon: 'Setting',
        sort: 3,
        hidden: false,
        parentPath: parentPath,
      },
    ]
  }

  /**
   * 获取插件路由配置
   */
  getRoutes(): PluginRouteConfig[] {
    return [
      {
        path: '/task',
        name: 'Task',
        component: () => import('@/views/task/Index.vue'),
        meta: { title: '任务中心' },
        children: [
          {
            path: 'jobs',
            name: 'TaskJobs',
            component: () => import('@/views/task/Jobs.vue'),
            meta: { title: '任务作业' },
          },
          {
            path: 'templates',
            name: 'TaskTemplates',
            component: () => import('@/views/task/Templates.vue'),
            meta: { title: '任务模板' },
          },
          {
            path: 'ansible',
            name: 'TaskAnsible',
            component: () => import('@/views/task/Ansible.vue'),
            meta: { title: 'Ansible任务' },
          },
        ],
      },
    ]
  }
}

// 创建并注册插件实例
const plugin = new TaskPlugin()
console.log('[Task Plugin] 正在注册插件...')
pluginManager.register(plugin)
console.log('[Task Plugin] 插件注册完成')

export default plugin
