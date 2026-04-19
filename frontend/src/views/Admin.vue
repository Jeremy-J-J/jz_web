<template>
  <div class="admin">
    <header class="header">
      <h1>管理后台</h1>
      <div class="header-actions">
        <span>{{ admin?.username }}</span>
        <button @click="handleLogout">退出登录</button>
      </div>
    </header>
    <div class="container">
      <div class="tabs">
        <button :class="['tab', { active: activeTab === 'resources' }]" @click="activeTab = 'resources'">资源管理</button>
        <button :class="['tab', { active: activeTab === 'categories' }]" @click="activeTab = 'categories'">类目管理</button>
        <button :class="['tab', { active: activeTab === 'stats' }]" @click="activeTab = 'stats'">访问统计</button>
      </div>

      <!-- Resource Management -->
      <div v-if="activeTab === 'resources'" class="resource-panel">
        <div class="toolbar">
          <h2>资源管理</h2>
          <router-link to="/admin/resource/new" class="add-btn">+ 新增资源</router-link>
        </div>
        <div v-if="loading" class="loading">加载中...</div>
        <div v-else-if="resources.length === 0" class="empty">暂无资源，点击上方按钮添加</div>
        <table v-else class="resource-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>封面</th>
              <th>标题</th>
              <th>类目</th>
              <th>链接</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="r in resources" :key="r.id">
              <td>{{ r.id }}</td>
              <td class="cover-cell">
                <img v-if="r.cover" :src="r.cover" alt="">
                <span v-else class="no-cover">无</span>
              </td>
              <td>{{ r.title }}</td>
              <td>{{ getCategoryName(r.category_id) }}</td>
              <td class="link-cell"><a :href="r.link" target="_blank">{{ r.link }}</a></td>
              <td>
                <span :class="['status', r.status === 1 ? 'active' : 'inactive']">
                  {{ r.status === 1 ? '上架' : '下架' }}
                </span>
              </td>
              <td class="actions">
                <button @click="toggleStatus(r.id)" class="btn-toggle">{{ r.status === 1 ? '下架' : '上架' }}</button>
                <router-link :to="`/admin/resource/${r.id}/edit`" class="btn-edit">编辑</router-link>
                <button @click="deleteResource(r.id)" class="btn-delete">删除</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Category Management -->
      <div v-if="activeTab === 'categories'">
        <CategoryManager />
      </div>

      <!-- Stats Panel -->
      <div v-if="activeTab === 'stats'" class="stats-panel">
        <h2>访问统计</h2>
        <div class="stats-grid">
          <div class="stats-card">
            <div class="stats-label">访问总数</div>
            <div class="stats-value">{{ stats.total }}</div>
          </div>
          <div class="stats-card">
            <div class="stats-label">今日访问</div>
            <div class="stats-value">{{ stats.today }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { authApi, publicApi } from '../api'
import CategoryManager from './CategoryManager.vue'

const router = useRouter()
const route = useRoute()
const resources = ref([])
const categories = ref([])
const loading = ref(true)
const activeTab = ref('resources')
const filterCategoryId = ref(null)
const admin = JSON.parse(localStorage.getItem('admin') || '{}')
const stats = ref({ total: 0, today: 0 })

const loadResources = async () => {
  try {
    const catId = filterCategoryId.value
    const resRes = catId
      ? await authApi.getAllResources({ category_id: catId })
      : await authApi.getAllResources()
    const catRes = await authApi.getCategories()
    resources.value = resRes.data || []
    categories.value = catRes.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const getCategoryName = (categoryId) => {
  const cat = categories.value.find(c => c.id === categoryId)
  return cat ? cat.name : '未分类'
}

const toggleStatus = async (id) => {
  try {
    await authApi.toggleStatus(id)
    loadResources()
  } catch (e) {
    alert('操作失败')
  }
}

const deleteResource = async (id) => {
  if (!confirm('确定要删除这个资源吗？')) return
  try {
    await authApi.deleteResource(id)
    loadResources()
  } catch (e) {
    alert('删除失败')
  }
}

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('admin')
  router.push('/admin/login')
}

const switchToResources = (categoryId) => {
  filterCategoryId.value = categoryId
  activeTab.value = 'resources'
}

const loadStats = async () => {
  try {
    const res = await publicApi.getStats()
    stats.value = res || { total: 0, today: 0 }
  } catch (e) {
    console.error(e)
  }
}

watch(() => route.query, (query) => {
  if (query.tab === 'resources') {
    activeTab.value = 'resources'
    if (query.category_id) {
      filterCategoryId.value = parseInt(query.category_id)
    } else {
      filterCategoryId.value = null
    }
    loadResources()
  }
}, { immediate: true })

watch(() => activeTab.value, (tab) => {
  if (tab === 'stats') {
    loadStats()
  }
})

onMounted(() => {
  if (route.query.tab === 'resources') {
    activeTab.value = 'resources'
    if (route.query.category_id) {
      filterCategoryId.value = parseInt(route.query.category_id)
    }
  }
  loadResources()
})
</script>

<style scoped>
.admin { min-height: 100vh; background: #f5f5f5; }
.header { background: #fff; padding: 16px 40px; display: flex; justify-content: space-between; align-items: center; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.header h1 { font-size: 20px; color: #333; }
.header-actions { display: flex; align-items: center; gap: 16px; }
.header-actions span { color: #666; }
.header-actions button { padding: 8px 16px; background: #ff4d4f; color: #fff; border: none; border-radius: 4px; cursor: pointer; }
.container { max-width: 1200px; margin: 24px auto; padding: 0 20px; }
.tabs { display: flex; gap: 0; margin-bottom: 20px; background: #fff; border-radius: 8px; overflow: hidden; }
.tab { padding: 12px 24px; border: none; background: #fafafa; cursor: pointer; font-size: 14px; color: #666; border-bottom: 2px solid transparent; }
.tab:hover { color: #1890ff; }
.tab.active { color: #1890ff; border-bottom-color: #1890ff; background: #fff; }
.resource-panel { background: #fff; border-radius: 8px; padding: 24px; }
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.toolbar h2 { font-size: 18px; color: #333; }
.add-btn { padding: 10px 20px; background: #1890ff; color: #fff; text-decoration: none; border-radius: 4px; font-size: 14px; }
.loading, .empty { text-align: center; padding: 60px; color: #999; }
.resource-table { width: 100%; }
.resource-table th, .resource-table td { padding: 12px 16px; text-align: left; border-bottom: 1px solid #f0f0f0; }
.resource-table th { background: #fafafa; font-weight: 500; color: #333; }
.resource-table td { color: #666; font-size: 14px; }
.cover-cell img { width: 60px; height: 40px; object-fit: cover; border-radius: 4px; }
.no-cover { color: #999; }
.link-cell a { color: #1890ff; text-decoration: none; max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; display: block; }
.status { padding: 4px 8px; border-radius: 4px; font-size: 12px; }
.status.active { background: #f6ffed; color: #52c41a; }
.status.inactive { background: #fff1f0; color: #ff4d4f; }
.actions { display: flex; gap: 8px; }
.actions button, .actions .btn-edit { padding: 6px 12px; border: none; border-radius: 4px; font-size: 12px; cursor: pointer; text-decoration: none; }
.btn-toggle { background: #faad14; color: #fff; }
.btn-edit { background: #1890ff; color: #fff; }
.btn-delete { background: #ff4d4f; color: #fff; }
.stats-panel { background: #fff; border-radius: 8px; padding: 24px; }
.stats-panel h2 { font-size: 18px; color: #333; margin-bottom: 20px; }
.stats-grid { display: flex; gap: 24px; }
.stats-card { flex: 1; background: #f8f8f8; border-radius: 8px; padding: 24px; text-align: center; }
.stats-label { font-size: 14px; color: #666; margin-bottom: 8px; }
.stats-value { font-size: 32px; color: #1890ff; font-weight: bold; }
</style>