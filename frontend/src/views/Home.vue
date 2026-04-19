<template>
  <div class="home">
    <header class="header">
      <h1>极智视界数据集分享</h1>
    </header>

    <div class="nav-bar">
      <div class="nav-container">
        <div class="category-nav">
          <button
            :class="['cat-btn', { active: selectedCategory === 0 }]"
            @click="selectCategory(0)"
          >全部</button>
          <button
            v-for="cat in categories"
            :key="cat.id"
            :class="['cat-btn', { active: selectedCategory === cat.id }]"
            @click="selectCategory(cat.id)"
          >{{ cat.name }}</button>
        </div>
        <div class="search-box">
          <input
            v-model="searchKeyword"
            type="text"
            placeholder="搜索资源..."
            @keyup.enter="doSearch"
          >
          <button @click="doSearch">搜索</button>
        </div>
      </div>
    </div>

    <div class="container">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="resources.length === 0" class="empty">
        {{ searchKeyword ? '未找到匹配的资源' : '暂无资源' }}
      </div>
      <div v-else class="resource-grid">
        <div v-for="resource in resources" :key="resource.id" class="resource-card" @click="goDetail(resource.id)">
          <div class="card-cover">
            <img v-if="resource.cover" :src="resource.cover" :alt="resource.title">
            <div v-else class="cover-placeholder">暂无封面</div>
          </div>
          <div class="card-content">
            <h3>{{ resource.title }}</h3>
            <p>{{ resource.description || '暂无描述' }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { publicApi } from '../api'

const resources = ref([])
const categories = ref([])
const loading = ref(true)
const selectedCategory = ref(0)
const searchKeyword = ref('')

const loadCategories = async () => {
  try {
    const res = await publicApi.getCategories()
    categories.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

const loadResources = async () => {
  loading.value = true
  try {
    const params = {}
    if (selectedCategory.value > 0) {
      params.category_id = selectedCategory.value
    }
    if (searchKeyword.value.trim()) {
      params.keyword = searchKeyword.value.trim()
    }
    const res = await publicApi.getResources(params)
    resources.value = res.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const selectCategory = (id) => {
  selectedCategory.value = id
  loadResources()
}

const doSearch = () => {
  loadResources()
}

const goDetail = (id) => {
  window.location.href = `/resource/${id}`
}

onMounted(async () => {
  await loadCategories()
  loadResources()
})
</script>

<style scoped>
.home { min-height: 100vh; background: #f5f5f5; }
.header { background: #fff; padding: 20px 40px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.header h1 { font-size: 24px; color: #333; text-align: center; }
.nav-bar { background: #fff; border-bottom: 1px solid #eee; position: sticky; top: 0; z-index: 100; }
.nav-container { max-width: 1200px; margin: 0 auto; padding: 12px 20px; display: flex; justify-content: space-between; align-items: center; gap: 20px; }
.category-nav { display: flex; gap: 8px; flex-wrap: wrap; flex: 1; }
.cat-btn { padding: 8px 16px; border: 1px solid #ddd; background: #fff; border-radius: 4px; cursor: pointer; font-size: 14px; transition: all 0.2s; }
.cat-btn:hover { border-color: #1890ff; color: #1890ff; }
.cat-btn.active { background: #1890ff; border-color: #1890ff; color: #fff; }
.search-box { display: flex; gap: 8px; }
.search-box input { padding: 8px 12px; border: 1px solid #ddd; border-radius: 4px; width: 200px; font-size: 14px; }
.search-box button { padding: 8px 16px; background: #1890ff; color: #fff; border: none; border-radius: 4px; cursor: pointer; font-size: 14px; }
.container { max-width: 1200px; margin: 40px auto; padding: 0 20px; }
.loading, .empty { text-align: center; padding: 60px; color: #999; }
.resource-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 24px; }
.resource-card { background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08); cursor: pointer; transition: transform 0.2s, box-shadow 0.2s; }
.resource-card:hover { transform: translateY(-4px); box-shadow: 0 8px 24px rgba(0,0,0,0.12); }
.card-cover { height: 180px; overflow: hidden; background: #f0f0f0; }
.card-cover img { width: 100%; height: 100%; object-fit: cover; }
.cover-placeholder { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; color: #999; }
.card-content { padding: 16px; }
.card-content h3 { font-size: 16px; margin-bottom: 8px; color: #333; }
.card-content p { font-size: 14px; color: #666; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
