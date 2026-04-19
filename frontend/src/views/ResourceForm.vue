<template>
  <div class="form-page">
    <header class="header">
      <button class="back-btn" @click="$router.back()">&larr; 返回</button>
      <h1>{{ isEdit ? '编辑资源' : '新增资源' }}</h1>
    </header>
    <div class="container">
      <form @submit.prevent="handleSubmit" class="resource-form">
        <div class="form-group">
          <label>类目</label>
          <select v-model="form.category_id">
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>
        <div class="form-group">
          <label>标题 *</label>
          <input v-model="form.title" type="text" placeholder="请输入资源标题" required>
        </div>
        <div class="form-group">
          <label>封面图片</label>
          <div class="cover-upload">
            <input type="file" accept="image/*" @change="handleCoverChange" ref="coverInput">
            <div v-if="form.cover" class="cover-preview">
              <img :src="form.cover" alt="封面预览">
              <button type="button" @click="removeCover">删除</button>
            </div>
            <div v-else class="upload-placeholder" @click="$refs.coverInput.click()">
              点击上传封面图片
            </div>
          </div>
        </div>
        <div class="form-group">
          <label>跳转链接 *</label>
          <input v-model="form.link" type="url" placeholder="请输入外部跳转链接（http://或https://开头）" required>
        </div>
        <div class="form-group">
          <label>描述</label>
          <textarea v-model="form.description" placeholder="请输入资源描述" rows="5"></textarea>
        </div>
        <div class="form-group">
          <label>状态</label>
          <select v-model="form.status">
            <option :value="1">上架</option>
            <option :value="0">下架</option>
          </select>
        </div>
        <div class="form-actions">
          <button type="button" @click="$router.back()" class="btn-cancel">取消</button>
          <button type="submit" :disabled="loading" class="btn-submit">{{ loading ? '保存中...' : '保存' }}</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { authApi, publicApi } from '../api'

const router = useRouter()
const route = useRoute()
const isEdit = computed(() => !!route.params.id)
const loading = ref(false)
const coverInput = ref(null)
const categories = ref([])

const form = reactive({
  category_id: 1,
  title: '',
  cover: '',
  link: '',
  description: '',
  status: 1
})

const loadCategories = async () => {
  try {
    const res = await authApi.getCategories()
    categories.value = res.data || []
    if (categories.value.length > 0 && form.category_id === 1) {
      form.category_id = categories.value[0].id
    }
  } catch (e) {
    console.error(e)
  }
}

const loadResource = async () => {
  if (!isEdit.value) return
  try {
    const res = await publicApi.getResource(route.params.id)
    const r = res.data
    form.category_id = r.category_id || 1
    form.title = r.title
    form.cover = r.cover
    form.link = r.link
    form.description = r.description
    form.status = r.status
  } catch (e) {
    alert('加载资源失败')
    router.back()
  }
}

const handleCoverChange = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  try {
    const res = await authApi.uploadCover(file)
    form.cover = res.url
  } catch (e) {
    alert('上传封面失败')
  }
}

const removeCover = () => {
  form.cover = ''
}

const handleSubmit = async () => {
  loading.value = true
  try {
    if (isEdit.value) {
      await authApi.updateResource(route.params.id, form)
    } else {
      await authApi.createResource(form)
    }
    router.push('/admin')
  } catch (e) {
    alert(isEdit.value ? '更新失败' : '创建失败')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadCategories()
  loadResource()
})
</script>

<style scoped>
.form-page { min-height: 100vh; background: #f5f5f5; }
.header { background: #fff; padding: 16px 40px; display: flex; align-items: center; gap: 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.back-btn { background: none; border: none; color: #666; cursor: pointer; font-size: 16px; }
.back-btn:hover { color: #1890ff; }
.header h1 { font-size: 20px; color: #333; }
.container { max-width: 600px; margin: 40px auto; padding: 0 20px; }
.resource-form { background: #fff; padding: 32px; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.08); }
.form-group { margin-bottom: 24px; }
.form-group label { display: block; font-size: 14px; color: #333; margin-bottom: 8px; font-weight: 500; }
.form-group input[type="text"],
.form-group input[type="url"],
.form-group textarea,
.form-group select { width: 100%; padding: 12px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px; }
.form-group textarea { resize: vertical; }
.form-group input:focus, .form-group textarea:focus, .form-group select:focus { outline: none; border-color: #1890ff; }
.cover-upload { border: 2px dashed #ddd; border-radius: 8px; padding: 20px; text-align: center; }
.cover-upload input[type="file"] { display: none; }
.upload-placeholder { color: #999; cursor: pointer; padding: 20px; }
.cover-preview { position: relative; display: inline-block; }
.cover-preview img { max-width: 200px; max-height: 150px; border-radius: 4px; }
.cover-preview button { position: absolute; top: -8px; right: -8px; background: #ff4d4f; color: #fff; border: none; border-radius: 50%; width: 24px; height: 24px; cursor: pointer; font-size: 12px; }
.form-actions { display: flex; justify-content: flex-end; gap: 12px; margin-top: 32px; }
.btn-cancel, .btn-submit { padding: 12px 24px; border: none; border-radius: 4px; font-size: 14px; cursor: pointer; }
.btn-cancel { background: #f0f0f0; color: #666; }
.btn-submit { background: #1890ff; color: #fff; }
.btn-submit:disabled { opacity: 0.6; cursor: not-allowed; }
</style>