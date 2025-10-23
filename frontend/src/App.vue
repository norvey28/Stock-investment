<script setup lang="ts">
import { onMounted } from 'vue';
import { DataTable, Column, Button, DatePicker, Select, MultiSelect, InputText, Toast } from 'primevue';
import { useToast } from 'primevue/usetoast';
import { useItemsStore } from './stores/items';
import { storeToRefs } from 'pinia';

const store = useItemsStore();
const toast = useToast();

onMounted(async () => {
  try {
    await store.fetchItems();
    console.log('Items fetched:', store.items);
    toast.add({ severity: 'success', summary: 'Éxito', detail: 'Datos cargados correctamente', life: 2000 });
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: String(e), life: 5000 });
  }
  // initialize/reset filters from store
  store.resetFilters();
});

function updateItems() {
  store.syncItems()
    .then(() => {
      console.log('Actualización completada. Items en el componente:', items.value.length);
      toast.add({ severity: 'success', summary: 'Éxito', detail: 'Datos actualizados correctamente', life: 2000 });
    })
    .catch((e) => toast.add({ severity: 'error', summary: 'Error', detail: String(e), life: 5000 }));
}

// Expose store state/getters to the template (Pinia refs are unwrapped in template
const items = storeToRefs(store).items;
const listBroker = store.listBroker;
const listRaiting = store.listRating;
</script>


<template>

  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Recomendaciones de Analistas</title>
    <Toast />
  </head>

  <div class="p-4 md:p-6 lg:p-8 max-w-[2000px] mx-auto">
    <h1 class="text-2xl md:text-3xl lg:text-4xl font-bold mb-6 text-center text-gray-800">
      Panel de Recomendaciones de Analistas
    </h1>

    <div class="card shadow-lg rounded-2xl p-4 mb-8 bg-white">
      <DataTable :value="items || []" paginator :rows="10" dataKey="ticker" filterDisplay="row"
        class="text-sm overflow-x-auto" removableSort showGridlines stripedRows :loading="store.loading"
        v-model:filters="store.filters" responsiveLayout="scroll">
        <template #loading>
          <div class="flex flex-column align-items-center justify-content-center ">
            <span class="mt-2 text-4xl">Cargando datos...</span>
          </div>
        </template>
        <template #header>
          <div class="flex justify-between ">
            <span class="">Total de recomendaciones: {{ items.length }}</span>
            <Button label="Actualizar informacion" icon="fa fa-solid fa-rotate-right" @click="updateItems" rounded
              raised />
          </div>
        </template>
        <template #empty>
          <div class="flex flex-column align-items-center justify-content-center">
            <i class="pi pi-info-circle text-4xl"></i>
            <span class="mt-2 text-lg">No hay datos disponibles...</span>
          </div>
        </template>
        <Column field="time" filter filterField="time" header="Fecha" sortable>
          <template #body="{ data }">
            {{ new Date(data.time).toLocaleDateString() }}
          </template>
          <template #filter="{ filterModel, filterCallback }">
            <div v-if="filterModel">
              <DatePicker v-model="filterModel.value" dateFormat="dd/mm/yy" placeholder="Seleccionar fecha"
                @date-select="filterCallback" />
            </div>
          </template>
        </Column>

        <Column field="ticker" header="Símbolo" sortable filter filterPlaceholder="Buscar ticker" />

        <Column field="company" header="Empresa" sortable filter filterPlaceholder="Buscar empresa">
          <template #body="{ data }">
            {{ data.company }}
          </template>
          <template #filter="{ filterModel, filterCallback }">
            <div v-if="filterModel">
              <InputText v-model="filterModel.value" placeholder="Seleccionar empresa" :options="store.listBroker"
                @input="filterCallback" />
            </div>
          </template>
        </Column>

        <Column field="action" filter filterField="action" header="Acción" sortable>
          <template #body="{ data }">
            <span :class="{
              'text-green-600 font-semibold': data.action.includes('raised'),
              'text-red-600 font-semibold': data.action.includes('lowered'),
              'text-gray-600': data.action.includes('reiterated'),
            }">
              {{ data.action }}
            </span>
          </template>
          <template #filter="{ filterModel, filterCallback }">
            <div v-if="filterModel">
              <MultiSelect v-model="filterModel.value" placeholder="Seleccionar acción" :options="store.listAcciones"
                @change="filterCallback" />
            </div>
          </template>
        </Column>

        <Column header="Objetivo" sortable>
          <template #body="{ data }">
            <span class="font-medium text-gray-800">{{ data.target_from }}</span>
            <span class="mx-1 text-gray-400">→</span>
            <span class="font-semibold text-gray-900">{{ data.target_to }}</span>
          </template>
        </Column>

        <Column field="rating_to" filterField="rating_to" :showFilterMenu="false" header="Calificación" sortable filter
          filterPlaceholder="Filtrar calificación">
          <template #body="{ data }">
            <span :class="{
              'bg-green-100 text-green-700 px-2 py-1 rounded-full text-xs': data.rating_to === 'Buy',
              'bg-yellow-100 text-yellow-700 px-2 py-1 rounded-full text-xs': data.rating_to === 'Neutral',
              'bg-red-100 text-red-700 px-2 py-1 rounded-full text-xs': data.rating_to === 'Sell',
              'bg-blue-100 text-blue-700 px-2 py-1 rounded-full text-xs': !['Buy', 'Neutral', 'Sell'].includes(data.rating_to)
            }">
              {{ data.rating_to }}
            </span>
          </template>
          <template #filter="{ filterModel, filterCallback }">
            <div v-if="filterModel">
              <Select @change="filterCallback()" :showClear="true" v-model="filterModel.value"
                placeholder="Seleccionar Rating" :options="listRaiting" />
            </div>
          </template>
        </Column>

        <Column field="brokerage" header="Corredor" filterField="brokerage" sortable filter
          filterPlaceholder="Buscar corredor">
          <template #body="{ data }">
            {{ data.brokerage || '—' }}
          </template>
          <template #filter="{ filterModel, filterCallback }">
            <div v-if="filterModel">
              <InputText v-model="filterModel.value" placeholder="Buscar" :options="listBroker"
                @input="filterCallback" />
            </div>
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- Tablas de Recomendaciones Específicas -->
    <div class="mt-8 bg-gray-50 p-4 md:p-6 lg:p-8 rounded-3xl">
      <h2 class="text-2xl md:text-3xl font-bold mb-8 text-gray-800 text-center">Análisis por Tipo de Recomendación</h2>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Tabla de Recomendaciones de Compra -->
        <div
          class="card shadow-lg rounded-2xl p-4 md:p-6 bg-gradient-to-br from-green-50 to-white border border-green-100">
          <div class="flex items-center mb-4">
            <i class="pi pi-arrow-up text-xl md:text-2xl mr-2 text-green-600"></i>
            <h3 class="text-lg md:text-xl font-bold text-green-700">Recomendaciones de Compra</h3>
          </div>
          <DataTable
            :value="items?.filter(item => item.recommendation_score > 0).sort((a, b) => b.recommendation_score - a.recommendation_score) || []"
            paginator :rows="5" dataKey="ticker" class="text-sm" removableSort showGridlines stripedRows>
            <template #header>
              <div class="flex flex-col md:flex-row justify-between items-center gap-2">
                <div class="flex items-center">
                  <span class="text-base md:text-lg font-semibold mr-2">Total:</span>
                  <span class="bg-green-100 text-green-800 px-3 py-1 rounded-full font-bold min-w-[2.5rem] text-center">
                    {{items?.filter(item => item.recommendation_score > 0).length}}
                  </span>
                </div>
                <div class="text-xs md:text-sm text-green-600 italic">Ordenado por mayor potencial</div>
              </div>
            </template>
            <template #empty>
              <div class="flex flex-column align-items-center justify-content-center">
                <span class="mt-2 text-lg">No hay recomendaciones de compra disponibles</span>
              </div>
            </template>
            <Column field="ticker" header="Símbolo" sortable />
            <Column field="company" header="Empresa" sortable />
            <Column header="Precio Objetivo">
              <template #body="{ data }">
                <div class="flex flex-col sm:flex-row items-start sm:items-center gap-2 sm:space-x-2">
                  <div class="flex items-center space-x-2">
                    <span class="text-gray-500 text-sm md:text-base">{{ data.target_from }}</span>
                    <i class="pi pi-arrow-right text-green-500"></i>
                    <span class="font-bold text-green-700 text-sm md:text-base">{{ data.target_to }}</span>
                  </div>
                  <span class="text-xs bg-green-50 text-green-600 px-2 py-0.5 rounded-full whitespace-nowrap">
                    {{ Math.round(((data.target_to - data.target_from) / data.target_from) * 100) }}%
                  </span>
                </div>
              </template>
            </Column>
            <Column field="rating_to" header="Calificación" sortable>
              <template #body="{ data }">
                <span :class="{
                  'bg-green-100 text-green-700 px-2 py-1 rounded-full text-xs': ['Buy', 'Strong-Buy', 'Speculative Buy', 'Overweight'].includes(data.rating_to),
                  'bg-yellow-100 text-yellow-700 px-2 py-1 rounded-full text-xs': ['Outperform', 'Market Outperform', 'Sector Outperform', 'Positive'].includes(data.rating_to)
                }">
                  {{ data.rating_to }}
                </span>
              </template>
            </Column>
            <Column field="recommendation_score" header="Puntaje" sortable>
              <template #body="{ data }">
                <div class="flex items-center">
                  <div class="w-12 h-2 bg-gradient-to-r from-green-200 to-green-500 rounded-full mr-2"
                    :style="{ opacity: (data.recommendation_score / 5) + 0.2 }">
                  </div>
                  <span class="text-green-600 font-bold">+{{ data.recommendation_score }}</span>
                </div>
              </template>
            </Column>
          </DataTable>
        </div>

        <!-- Tabla de Recomendaciones de Mantener -->
        <div
          class="card shadow-lg rounded-2xl p-4 md:p-6 bg-gradient-to-br from-blue-50 to-white border border-blue-100">
          <div class="flex items-center mb-4">
            <i class="pi pi-minus text-xl md:text-2xl mr-2 text-blue-600"></i>
            <h3 class="text-lg md:text-xl font-bold text-blue-700">Recomendaciones de Mantener</h3>
          </div>
          <DataTable :value="items?.filter(item => item.recommendation_score === 0) || []" paginator :rows="5"
            dataKey="ticker" class="text-sm" removableSort showGridlines stripedRows>
            <template #header>
              <div class="flex flex-col md:flex-row justify-between items-center gap-2">
                <div class="flex items-center">
                  <span class="text-base md:text-lg font-semibold mr-2">Total:</span>
                  <span class="bg-blue-100 text-blue-800 px-3 py-1 rounded-full font-bold min-w-[2.5rem] text-center">
                    {{items?.filter(item => item.recommendation_score === 0).length}}
                  </span>
                </div>
                <div class="text-xs md:text-sm text-blue-600 italic">Recomendación neutral</div>
              </div>
            </template>
            <template #empty>
              <div class="flex flex-column align-items-center justify-content-center">
                <span class="mt-2 text-lg">No hay recomendaciones neutras disponibles</span>
              </div>
            </template>
            <Column field="ticker" header="Símbolo" sortable />
            <Column field="company" header="Empresa" sortable />
            <Column header="Precio Objetivo">
              <template #body="{ data }">
                <div class="flex items-center space-x-2">
                  <span class="text-gray-500">{{ data.target_from }}</span>
                  <i class="pi pi-arrow-right text-blue-500"></i>
                  <span class="font-bold text-gray-700">{{ data.target_to }}</span>
                  <span class="text-xs bg-blue-50 text-blue-600 px-2 py-0.5 rounded-full">
                    {{ Math.round(((data.target_to - data.target_from) / data.target_from) * 100) }}%
                  </span>
                </div>
              </template>
            </Column>
            <Column field="rating_to" header="Calificación" sortable>
              <template #body="{ data }">
                <span class="bg-blue-100 text-blue-700 px-2 py-1 rounded-full text-xs">
                  {{ data.rating_to }}
                </span>
              </template>
            </Column>
            <Column field="recommendation_score" header="Puntaje" sortable />
          </DataTable>
        </div>

        <!-- Tabla de Recomendaciones de Venta -->
        <div class="card shadow-lg rounded-2xl p-4 md:p-6 bg-gradient-to-br from-red-50 to-white border border-red-100">
          <div class="flex items-center mb-4">
            <i class="pi pi-arrow-down text-xl md:text-2xl mr-2 text-red-600"></i>
            <h3 class="text-lg md:text-xl font-bold text-red-700">Recomendaciones de Venta</h3>
          </div>
          <DataTable
            :value="items?.filter(item => item.recommendation_score < 0).sort((a, b) => a.recommendation_score - b.recommendation_score) || []"
            paginator :rows="5" dataKey="ticker" class="text-sm" removableSort showGridlines stripedRows>
            <template #header>
              <div class="flex flex-col md:flex-row justify-between items-center gap-2">
                <div class="flex items-center">
                  <span class="text-base md:text-lg font-semibold mr-2">Total:</span>
                  <span class="bg-red-100 text-red-800 px-3 py-1 rounded-full font-bold min-w-[2.5rem] text-center">
                    {{items?.filter(item => item.recommendation_score < 0).length}} </span>
                </div>
                <div class="text-xs md:text-sm text-red-600 italic">Ordenado por mayor urgencia de venta</div>
              </div>
            </template>
            <template #empty>
              <div class="flex flex-column align-items-center justify-content-center">
                <span class="mt-2 text-lg">No hay recomendaciones de venta disponibles</span>
              </div>
            </template>
            <Column field="ticker" header="Símbolo" sortable />
            <Column field="company" header="Empresa" sortable />
            <Column header="Precio Objetivo">
              <template #body="{ data }">
                <div class="flex items-center space-x-2">
                  <span class="text-gray-500">{{ data.target_from }}</span>
                  <i class="pi pi-arrow-right text-red-500"></i>
                  <span class="font-bold text-red-700">{{ data.target_to }}</span>
                  <span class="text-xs bg-red-50 text-red-600 px-2 py-0.5 rounded-full">
                    {{ Math.round(((data.target_to - data.target_from) / data.target_from) * 100) }}%
                  </span>
                </div>
              </template>
            </Column>
            <Column field="rating_to" header="Calificación" sortable>
              <template #body="{ data }">
                <span class="bg-red-100 text-red-700 px-2 py-1 rounded-full text-xs">
                  {{ data.rating_to }}
                </span>
              </template>
            </Column>
            <Column field="recommendation_score" header="Puntaje" sortable>
              <template #body="{ data }">
                <div class="flex items-center">
                  <div class="w-12 h-2 bg-gradient-to-r from-red-500 to-red-200 rounded-full mr-2"
                    :style="{ opacity: (Math.abs(data.recommendation_score) / 5) + 0.2 }">
                  </div>
                  <span class="text-red-600 font-bold">{{ data.recommendation_score }}</span>
                </div>
              </template>
            </Column>
          </DataTable>
        </div>
      </div>
    </div>
  </div>
</template>