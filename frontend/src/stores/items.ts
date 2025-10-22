import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

// URL por defecto para las solicitudes HTTP
axios.defaults.baseURL = 'http://localhost:8081/api/v1'

// Interface
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
  time: Date | string
  created_at: string
  recommendation_score: number
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
    time: { value: null, matchMode: 'dateAfter' }
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
      if (Array.isArray(response.data)) {
        items.value = [...response.data]
        for (const item of items.value) {
          // Convertir las cadenas de fecha a objetos Date
          item.time = new Date(item.time)
          // Manejar valores nulos o indefinidos en rating_to
          if (item.rating_to === null || item.rating_to === undefined || item.rating_to === '') {
            item.rating_to = 'N/A'
          }
          //Calcular score de recomendacion
          item.recommendation_score = calculateRecommendationScore(item)
        }
      } else {
        console.error('La respuesta del servidor no es un array:', response.data)
        items.value = []
      }
      console.log('Items cargados:', items.value.length)
    } catch (error) {
      console.error('Error cargando items:', error)
      items.value = []
      throw error
    } finally {
      loading.value = false
    }
  }

  //Actualiza los datos del backend con la informacion de la API
  async function syncItems() {
    loading.value = true
    try {
      console.log('Iniciando sincronización...')
      const syncResponse = await axios.put('/items/')
      console.log('Sincronización completada:', syncResponse.status)
      console.log('Recargando items...')
      await fetchItems() // Recargar items después de sincronizar
      console.log('Items recargados. Nuevo total:', items.value.length)
    } catch (error) {
      console.error('Error sincronizando items:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // Actualiza los filtros de búsqueda
  function updateFilters(newFilters: Partial<Filters>) {
    filters.value = { ...filters.value, ...newFilters }
  }

  // Resetea los filtros a sus valores iniciales
  function resetFilters() {
    filters.value = {
      action: { value: null, matchMode: 'in' },
      brokerage: { value: null, matchMode: 'contains' },
      rating_to: { value: null, matchMode: 'equals' },
      company: { value: null, matchMode: 'contains' },
      time: { value: null, matchMode: 'dateAfter' }
    }
  }

  function calculateRecommendationScore(item: Item): number {
    let score = 0
    // Asignar puntos según rating_to
    switch (item.rating_to) {
      case 'Buy':
      case 'Strong-Buy':
      case 'Speculative Buy':
      case 'Overweight':
        score += 2
        break
      case 'Outperform':
      case 'Market Outperform':
      case 'Sector Outperform':
      case 'Positive':
        score += 1
        break
      case 'Hold':
      case 'Neutral':
      case 'In-Line':
      case 'Equal Weight':
      case 'Market Perform':
      case 'Sector Perform':
        score += 0
        break
      case 'Underweight':
      case 'Underperform':
      case 'Reduce':
      case 'Cautious':
        score -= 1
        break
      case 'Sell':
        score -= 2
      default:
        score += 0
    }

    // Asignar puntos según action
    switch (item.action) {
      case 'upgraded by':
        score += 2
        break
      case 'target raised by':
        score += 1
        break
      case 'reiterated by':
      case 'initiated by':
      case 'target set by':
        score += 0
        break
      case 'target lowered by':
        score -= 1
        break
      case 'downgraded by':
        score -= 2
        break
    }
    return score;
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
    resetFilters,
    calculateRecommendationScore,
  }
})