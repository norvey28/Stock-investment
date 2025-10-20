import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

// Default backend base URL
axios.defaults.baseURL = 'http://localhost:8081/api/v1'

// Interfaces para TypeScript
interface Item {
  id: string
  ticker: string
  target_from: string
  target_to: string
  company: string
  action: string
  brokerage: string
  rating_from: string
  rating_to: string
  time: string
  created_at: string
}

interface Filters {
  action: { value: string[] | null; matchMode: string }
  brokerage: { value: string | null; matchMode: string }
  rating_to: { value: string | null; matchMode: string }
  company: { value: string | null; matchMode: string }
  time: { value: Date | string | null; matchMode: string }
}

export const useItemsStore = defineStore('items', () => {
  // Estado
  const items = ref<Item[]>([])
  const loading = ref(false)
  const filters = ref<Filters>({
    action: { value: null, matchMode: 'in' },
    brokerage: { value: null, matchMode: 'contains' },
    rating_to: { value: null, matchMode: 'equals' },
    company: { value: null, matchMode: 'contains' },
    time: { value: null, matchMode: 'dateIs' }
  })

  // Getters computados
  const listAcciones = computed(() => [...new Set(items.value.map(item => item.action))])
  const listBroker = computed(() => [...new Set(items.value.map(item => item.brokerage))])
  const listRating = computed(() => [...new Set(items.value.map(item => item.rating_to))])
  const totalItems = computed(() => items.value.length)

  // Acciones
  async function fetchItems() {
    loading.value = true
    try {
      const response = await axios.get('/items')
      console.log('Respuesta del servidor:', JSON.stringify(response.data))
      items.value = response.data
      console.log('Items cargados:', items.value.length)
    } catch (error) {
      console.error('Error cargando items:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function syncItems() {
    loading.value = true
    try {
      await axios.put('/items/')
      await fetchItems() // Recargar items después de sincronizar
    } catch (error) {
      console.error('Error sincronizando items:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  function updateFilters(newFilters: Partial<Filters>) {
    filters.value = { ...filters.value, ...newFilters }
  }

  function resetFilters() {
    filters.value = {
      action: { value: null, matchMode: 'in' },
      brokerage: { value: null, matchMode: 'contains' },
      rating_to: { value: null, matchMode: 'equals' },
      company: { value: null, matchMode: 'contains' },
      time: { value: null, matchMode: 'dateIs' }
    }
  }

  return {
    // Estado
    items,
    loading,
    filters,
    // Getters
    listAcciones,
    listBroker,
    listRating,
    totalItems,
    // Acciones
    fetchItems,
    syncItems,
    updateFilters,
    resetFilters
  }
})