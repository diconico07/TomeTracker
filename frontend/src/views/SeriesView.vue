<script setup>
import BookCard from '@/components/BookCard.vue'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const series = ref({})
const title = ref('')
const loading = ref(true)

onMounted(async () => {
  try {
    const response = await axios.get(`/api/series/${route.params.id}`)
    series.value = response.data
    title.value = series.value.name
    loading.value = false
  } catch (error) {
    console.error('Error fetching series:', error)
  }
})
</script>

<template>
  <v-app-bar scroll-behavior="elevate">
    <template v-slot:prepend>
      <v-app-bar-nav-icon icon="mdi-chevron-left" to="/"></v-app-bar-nav-icon>
    </template>

    <v-app-bar-title>{{ title }}</v-app-bar-title>
  </v-app-bar>
  <template v-if="loading">
    <v-skeleton-loader type="card"></v-skeleton-loader>
  </template>
  <template v-else>
    <v-container>
      <v-row dense>
        <v-col cols="12" v-for="item in series.edges.books" :key="item.isbn">
          <BookCard :title="item.title" :cover="item.cover" :tome_number="item.tome_number" :released_at="item.released_at" :isbn="item.isbn" :owned="item.owned"></BookCard>
        </v-col>
      </v-row>
    </v-container>
  </template>
</template>
