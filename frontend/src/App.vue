<script setup lang="ts">
import axios from 'axios'
import { onMounted, ref } from 'vue';
import { DataTable, Column, Button, DatePicker, MultiSelect } from 'primevue';
import {FilterMatchMode} from '@primevue/core/api';
axios.defaults.baseURL = 'http://localhost:8080/api/v1';

const items = ref<any[]>([]);
const loading = ref(false);
const listAcciones = ref<any[]>([]);
const listBroker = ref<any[]>([]);
const listRaiting = ref<any[]>([]);
const filters = ref([]);

onMounted(() => {
  getItems();
  init();
});

function getItems() {
  loading.value = true;
  axios.get('/items')
    .then(response => {
      console.log('Respuesta del backend:', response.data);
      items.value = response.data;
      listAcciones.value = [...new Set(items.value.map((item: { action: any; }) => item.action))];
      listBroker.value = [...new Set(items.value.map((item: { brokerage: any; }) => item.brokerage))];
      listRaiting.value = [...new Set(items.value.map((item: { rating_to: any; }) => item.rating_to))];
    }).catch(error => {
      console.error('Error al conectar con el backend:', error);
    }).finally(() => {
      loading.value = false;
    });
}

function init() {
  filters.value = {
    global: { value: null, matchMode: FilterMatchMode.CONTAINS },
    action: { value: null, matchMode: FilterMatchMode.IN },
    brokerage: { value: null, matchMode: FilterMatchMode.IN },
    rating_to: { value: null, matchMode: FilterMatchMode.IN },
    time: { value: null, matchMode: FilterMatchMode.DATE_IS }
  };
}
</script>


<template>

  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Recomendaciones de Analistas</title>
  </head>

  <div class="p-6 ">
    <h1 class="text-2xl font-semibold mb-4 ">Recomendaciones de Analistas</h1>

    <div class="card shadow-md rounded-2xl  p-4">
      <DataTable :value="items" paginator :rows="10" dataKey="ticker" filterDisplay="row" class="text-sm "
        removableSort showGridlines stripedRows :loading="loading"  v-model:filters="filters" >
        <template #loading>
          <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
        </template>
        <template #header>
          <div class="flex justify-between ">
            <div class="">
              <span class="">Total de recomendaciones: {{ items.length }}</span>
            </div>
            <div class="">
              <Button label="Actualizar informacion" icon="fa fa-solid fa-rotate-right" rounded raised />
            </div>
          </div>
        </template>
        <Column field="time" filterField="time" header="Fecha" sortable>
          <template #body="{ data }">
            {{ new Date(data.time).toLocaleDateString() }}
          </template>
          <template #filter="filterModel">
            <DatePicker v-model="filterModel.value" dateFormat="yy-mm-dd" placeholder="Seleccionar fecha" />
          </template>
        </Column>

        <Column field="ticker" header="Ticker" sortable filter filterPlaceholder="Buscar ticker" />

        <Column field="company" header="Empresa" sortable filter filterPlaceholder="Buscar empresa" />

        <Column field="action" filterField="action" header="Acción" sortable >
          <template #body="{ data }">
            <span :class="{
              'text-green-600 font-semibold': data.action.includes('raised'),
              'text-red-600 font-semibold': data.action.includes('lowered'),
              'text-gray-600': data.action.includes('reiterated'),
            }">
              {{ data.action }}
            </span>
          </template>
          <template #filter="filterModel, filterCallback">
            <MultiSelect v-model="filterModel.value" placeholder="Seleccionar acción" :options="listAcciones" @change="filterCallback()" />
          </template>
        </Column>

        <Column header="Target" sortable>
          <template #body="{ data }">
            <span class="font-medium text-gray-800">{{ data.target_from }}</span>
            <span class="mx-1 text-gray-400">→</span>
            <span class="font-semibold text-gray-900">{{ data.target_to }}</span>
          </template>
        </Column>

        <Column field="rating_to" header="Rating" sortable filter filterPlaceholder="Filtrar rating">
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
          <template #filter="filterModel, filterCallback">
            <MultiSelect v-model="filterModel.value" placeholder="Seleccionar acción" :options="listRaiting"  @click="filterCallback()" />
          </template>
        </Column>

        <Column field="brokerage" header="Broker" sortable filter filterPlaceholder="Buscar broker">
          <template #body="{ data }">
            {{ data.brokerage || '—' }}
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>