<template>
  <div class="category-manager">
    <div class="toolbar">
      <h2>类目管理</h2>
      <button @click="showAddModal = true" class="add-btn">+ 新增类目</button>
    </div>

    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="categories.length === 0" class="empty">暂无类目</div>
    <table v-else class="category-table">
      <thead>
        <tr>
          <th>排序</th>
          <th>名称</th>
          <th>描述</th>
          <th>资源数</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="cat in categories" :key="cat.id">
          <td>{{ cat.sort }}</td>
          <td>{{ cat.name }}</td>
          <td>{{ cat.description || '-' }}</td>
          <td><a href="javascript:void(0)" @click="goToResources(cat.id)" class="resource-count-link">{{ getResourceCount(cat.id) }}</a></td>
          <td class="actions">
            <button @click="editCategory(cat)" class="btn-edit">编辑</button>
            <button @click="deleteCategory(cat.id)" class="btn-delete" :disabled="getResourceCount(cat.id) > 0">删除</button>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- Add/Edit Modal -->
    <div v-if="showAddModal || showEditModal" class="modal">
      <div class="modal-content">
        <h3>{{ showEditModal ? '编辑类目' : '新增类目' }}</h3>
        <div class="form-group">
          <label>名称 *</label>
          <input v-model="form.name" type="text" placeholder="请输入类目名称" required>
        </div>
        <div class="form-group">
          <label>描述</label>
          <input v-model="form.description" type="text" placeholder="请输入类目描述">
        </div>
        <div class="form-group">
          <label>排序（数字越小越靠前）</label>
          <input v-model.number="form.sort" type="number" placeholder="0">
        </div>
        <div class="modal-actions">
          <button @click="closeModal" class="btn-cancel">取消</button>
          <button @click="submitForm" class="btn-submit">{{ showEditModal ? '保存' : '创建' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '../api'

const router = useRouter()

const categories = ref([])
const resources = ref([])
const loading = ref(true)
const showAddModal = ref(false)
const showEditModal = ref(false)
const editingId = ref(null)

const form = ref({
  name: '',
  description: '',
  sort: 0
})

const loadData = async () => {
  try {
    const [catRes, resRes] = await Promise.all([
      authApi.getCategories(),
      authApi.getAllResources()
    ])
    categories.value = catRes.data || []
    resources.value = resRes.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const getResourceCount = (categoryId) => {
  return resources.value.filter(r => r.category_id === categoryId).length
}

const goToResources = (categoryId) => {
  router.push({ path: '/admin', query: { tab: 'resources', category_id: categoryId } })
}

const editCategory = (cat) => {
  editingId.value = cat.id
  form.value = { name: cat.name, description: cat.description, sort: cat.sort }
  showEditModal.value = true
}

const deleteCategory = async (id) => {
  if (!confirm('确定要删除这个类目吗？')) return
  try {
    await authApi.deleteCategory(id)
    await loadData()
  } catch (e) {
    alert('删除失败')
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  editingId.value = null
  form.value = { name: '', description: '', sort: 0 }
}

const submitForm = async () => {
  if (!form.value.name.trim()) {
    alert('请输入类目名称')
    return
  }

  try {
    if (showEditModal.value) {
      await authApi.updateCategory(editingId.value, form.value)
    } else {
      await authApi.createCategory(form.value)
    }
    closeModal()
    await loadData()
  } catch (e) {
    alert(showEditModal.value ? '更新失败' : '创建失败')
  }
}

onMounted(loadData)
</script>

<style scoped>
.category-manager { background: #fff; border-radius: 8px; padding: 24px; }
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.toolbar h2 { font-size: 18px; color: #333; }
.add-btn { padding: 10px 20px; background: #1890ff; color: #fff; border: none; border-radius: 4px; cursor: pointer; }
.loading, .empty { text-align: center; padding: 40px; color: #999; }
.category-table { width: 100%; border-collapse: collapse; }
.category-table th, .category-table td { padding: 12px; text-align: left; border-bottom: 1px solid #f0f0f0; }
.category-table th { font-weight: 500; color: #333; }
.actions { display: flex; gap: 8px; }
.actions button { padding: 6px 12px; border: none; border-radius: 4px; cursor: pointer; font-size: 12px; }
.btn-edit { background: #1890ff; color: #fff; }
.btn-delete { background: #ff4d4f; color: #fff; }
.btn-delete:disabled { opacity: 0.5; cursor: not-allowed; }
.resource-count-link { color: #1890ff; text-decoration: none; cursor: pointer; }
.resource-count-link:hover { text-decoration: underline; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal-content { background: #fff; padding: 24px; border-radius: 8px; width: 400px; }
.modal-content h3 { margin-bottom: 20px; }
.form-group { margin-bottom: 16px; }
.form-group label { display: block; margin-bottom: 6px; font-size: 14px; color: #333; }
.form-group input { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 12px; margin-top: 24px; }
.btn-cancel, .btn-submit { padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; }
.btn-cancel { background: #f0f0f0; color: #666; }
.btn-submit { background: #1890ff; color: #fff; }
</style>// TEST_BUILD_1776551035278322173
