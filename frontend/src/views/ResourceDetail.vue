<template>
  <div class="detail">
    <header class="header">
      <button class="back-btn" @click="$router.back()">&larr; 返回</button>
      <router-link to="/" class="home-link">首页</router-link>
    </header>
    <div class="container">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <div v-else-if="resource" class="resource-detail">
        <div class="detail-cover">
          <img v-if="resource.cover" :src="resource.cover" :alt="resource.title">
          <div v-else class="cover-placeholder">暂无封面</div>
        </div>
        <div class="detail-content">
          <h1>{{ resource.title }}</h1>
          <p class="description">{{ resource.description || '暂无描述' }}</p>
          <a :href="resource.link" target="_blank" class="jump-btn">访问资源 &rarr;</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { publicApi } from '../api'

const route = useRoute()
const resource = ref(null)
const loading = ref(true)
const error = ref('')

const loadResource = async () => {
  try {
    const res = await publicApi.getResource(route.params.id)
    resource.value = res.data
  } catch (e) {
    error.value = '资源不存在或已下架'
  } finally {
    loading.value = false
  }
}

onMounted(loadResource)
</script>

<style scoped>
.detail { min-height: 100vh; background: #f5f5f5; }
.header { background: #fff; padding: 16px 40px; display: flex; align-items: center; gap: 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.back-btn { background: none; border: none; color: #666; cursor: pointer; font-size: 16px; }
.back-btn:hover { color: #1890ff; }
.home-link { color: #666; text-decoration: none; margin-left: auto; }
.home-link:hover { color: #1890ff; }
.container { max-width: 900px; margin: 40px auto; padding: 0 20px; }
.loading, .error { text-align: center; padding: 60px; color: #999; }
.error { color: #ff4d4f; }
.resource-detail { background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08); }
.detail-cover { height: 400px; background: #f0f0f0; overflow: hidden; }
.detail-cover img { width: 100%; height: 100%; object-fit: contain; background: #000; }
.cover-placeholder { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; color: #999; }
.detail-content { padding: 32px; }
.detail-content h1 { font-size: 28px; margin-bottom: 16px; color: #333; }
.description { font-size: 16px; color: #666; line-height: 1.8; margin-bottom: 32px; white-space: pre-wrap; }
.jump-btn { display: inline-block; padding: 14px 40px; background: #1890ff; color: #fff; text-decoration: none; border-radius: 4px; font-size: 16px; transition: background 0.2s; }
.jump-btn:hover { background: #40a9ff; }
</style>