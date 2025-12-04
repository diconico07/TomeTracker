<script setup>
import BookCard from '@/components/BookCard.vue'
import { onMounted, ref } from 'vue'
import axios from 'axios'

const books = ref({})
const loading = ref(true)

onMounted(async () => {
  try {
    const response = await axios.get(`/api/planning`)
    books.value = response.data
    loading.value = false
  } catch (error) {
    console.error('Error fetching series:', error)
  }
})
</script>

<template>
  <template v-if="loading">
    <v-skeleton-loader type="card"></v-skeleton-loader>
  </template>
  <template v-else>
    <v-container>
      <v-row dense>
        <v-col cols="12" v-for="item in books" :key="item.isbn">
          <BookCard :title="item.title" :cover="item.cover" :tome_number="item.tome_number" :released_at="item.released_at" :isbn="item.isbn" :owned="item.owned"></BookCard>
        </v-col>
      </v-row>
    </v-container>
  </template>
</template>
