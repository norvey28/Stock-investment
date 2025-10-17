import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'
import type { FilterMatchMode } from '@primevue/core/api'

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
  global: { value: string | null; matchMode: FilterMatchMode }
  action: { value: string[] | null; matchMode: FilterMatchMode }
  brokerage: { value: string[] | null; matchMode: FilterMatchMode }
  rating_to: { value: string[] | null; matchMode: FilterMatchMode }
  time: { value: Date | null; matchMode: FilterMatchMode }
}

export const useItemsStore = defineStore('items', () => {
  // Estado
  const items = ref<Item[]>([])
  const loading = ref(false)
  const filters = ref<Filters>({
    global: { value: null, matchMode: 'contains' },
    action: { value: null, matchMode: 'in' },
    brokerage: { value: null, matchMode: 'in' },
    rating_to: { value: null, matchMode: 'in' },
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
      await axios.put('/items')
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
      global: { value: null, matchMode: 'contains' },
      action: { value: null, matchMode: 'in' },
      brokerage: { value: null, matchMode: 'in' },
      rating_to: { value: null, matchMode: 'in' },
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