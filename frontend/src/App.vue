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

  <div class="p-6 ">
    <h1 class="text-2xl font-semibold mb-4 ">Recomendaciones de Analistas</h1>

    <div class="card shadow-md rounded-2xl  p-4">
      <DataTable :value="items || []" paginator :rows="10" dataKey="ticker" filterDisplay="row" class="text-sm " removableSort
        showGridlines stripedRows :loading="store.loading" v-model:filters="store.filters">
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
  </div>
</template>