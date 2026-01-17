import type { Plugin, PluginMenuConfig, PluginRouteConfig } from '../types'
import { pluginManager } from '../manager'
import DomainMonitor from './components/DomainMonitor.vue'
import TroubleHandling from './components/TroubleHandling.vue'

/**
 * 监控中心插件
 * 提供域名监控、故障处理等功能
 */
class MonitorPlugin implements Plugin {
  name = 'monitor'
  prettyName = '监控中心'
  description = '监控中心插件，提供域名监控、故障处理等功能'
  version = '1.0.0'
  author = 'OpsHub Team'

  /**
   * 安装插件
   */
  async install() {
    console.log('[Monitor Plugin] 插件安装中...')
    // 初始化操作
  }

  /**
   * 卸载插件
   */
  async uninstall() {
    console.log('[Monitor Plugin] 插件卸载中...')
    // 清理资源
  }

  /**
   * 获取插件菜单配置
   */
  getMenus(): PluginMenuConfig[] {
    const parentPath = '/monitor'

    return [
      {
        name: '监控中心',
        path: parentPath,
        icon: 'Monitor',
        sort: 20,
        hidden: false,
        parentPath: '',
      },
      {
        name: '域名监控',
        path: '/monitor/domain',
        icon: 'Monitor',
        sort: 1,
        hidden: false,
        parentPath: parentPath,
      },
      {
        name: '故障处理',
        path: '/monitor/trouble',
        icon: 'Warning',
        sort: 2,
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
        path: '/monitor',
        name: 'Monitor',
        component: () => import('./components/DomainMonitor.vue'),
        redirect: '/monitor/domain',
        meta: { title: '监控中心' },
      },
      {
        path: '/monitor/domain',
        name: 'DomainMonitor',
        component: () => import('./components/DomainMonitor.vue'),
        meta: { title: '域名监控' },
      },
      {
        path: '/monitor/trouble',
        name: 'TroubleHandling',
        component: () => import('./components/TroubleHandling.vue'),
        meta: { title: '故障处理' },
      },
    ]
  }
}

// 创建并注册插件实例
const plugin = new MonitorPlugin()
console.log('[Monitor Plugin] 正在注册插件...')
pluginManager.register(plugin)
console.log('[Monitor Plugin] 插件注册完成')

export default plugin
