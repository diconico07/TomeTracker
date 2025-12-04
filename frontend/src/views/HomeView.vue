<script setup>
import SeriesCard from '@/components/SeriesCard.vue'
import { onMounted, ref } from 'vue'
import { useDisplay } from 'vuetify'
import axios from 'axios'

const { mobile } = useDisplay()
const series = ref([])
const loading = ref(true)

async function fetchSeries() {
  try {
    const response = await axios.get('./api/series')
    series.value = response.data
    loading.value = false
  } catch (error) {
    console.error('Error fetching series:', error)
  }
}

const newSeriesUrl = ref('')
const isAddingSeries = ref(false)

async function addSeries(isActive) {
  isAddingSeries.value = true

  try {
    await axios.post('/api/series', { url: newSeriesUrl.value, name: newSeriesUrl.value})
    newSeriesUrl.value = ''
    isActive.value = false
    await fetchSeries()
  } catch (error) {
    console.error('Error adding series:', error)
  } finally {
    isAddingSeries.value = false
  }
}


onMounted(async () => {
  fetchSeries()
})

</script>

<template>
  <v-container>
    <v-row dense>
      <template v-if="loading">
        <v-col cols="12" v-for="n in 3" :key="n">
          <v-skeleton-loader type="list-item-avatar-three-line"></v-skeleton-loader>
        </v-col>
      </template>
      <template v-else>
        <v-col v-for="item in series" cols="12" :key="item.id">
          <SeriesCard :total="item.volumes_total" :current="item.volumes_owned" :title="item.name" :author="item.author" :cover="item.cover" :id="item.id" />
        </v-col>
      </template>
    </v-row>
  </v-container>
  <v-dialog max-width="500">
    <template v-slot:activator="{ props: activatorProps }">
      <v-fab color="primary" icon="mdi-plus" :class="'mr-5' + ' ' + (mobile? 'mb-15' : 'mt-n4')"
        location="bottom end"
        app v-bind="activatorProps"></v-fab>
    </template>

    <template v-slot:default="{ isActive }">
      <v-card title="Add new series">
        <v-card-text>
          <v-text-field
            label="URL*"
            required
            v-model="newSeriesUrl"
          ></v-text-field>
        </v-card-text>

        <v-card-actions>

          <v-btn
            text="Close"
            @click="isActive.value = false"
          ></v-btn>
          <v-btn
            text="Add"
            color="primary"
            :loading="isAddingSeries"
            @click="addSeries(isActive)"
          ></v-btn>
        </v-card-actions>
      </v-card>
    </template>
  </v-dialog>

</template>
