<template>
  <v-card class="mx-auto" :title="title">
    <template v-slot:subtitle>
      <v-chip density="compact" prepend-icon="mdi-calendar">{{ date.format(released_at, 'fullDate') }}</v-chip>
      <v-chip density="compact" class="ml-1" prepend-icon="mdi-book" v-if="mobile">{{ props.tome_number }}</v-chip>
    </template>
    <template v-slot:prepend>
      <v-img
        class="mx-auto"
        size="125"
        :src="cover || missing"
        height="125"
        width="125"
      ></v-img>
    </template>

    <template v-slot:append v-if="!mobile">
      <v-btn
        :key="`info-${owned}-wide`"
        :color="owned ? 'success' : 'primary'"
        :prepend-icon="owned ? 'mdi-check' : 'mdi-plus'"
        :text="owned ? 'Owned' : 'Add'"
        @click="handleButtonClick"
        :loading="isLoading"
      ></v-btn>
    </template>
    <template v-if="mobile">
      <v-card-actions>
        <v-btn
        :key="`info-${owned}`"
        :color="owned ? 'success' : 'primary'"
        :prepend-icon="owned ? 'mdi-check' : 'mdi-plus'"
        :text="owned ? 'Owned' : 'Add'"
        @click="handleButtonClick"
        :loading="isLoading"
      ></v-btn>
      </v-card-actions>
    </template>
  </v-card>
  <v-dialog v-model="isDialogOpen" max-width="500">
    <v-card>
      <v-card-title>Confirmation Needed</v-card-title>
      <v-card-text> Are you sure you want to remove that book from your collection ? </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="grey" variant="text" @click="isDialogOpen = false">Cancel</v-btn>
        <v-btn color="primary" @click="confirmDialogAction">Confirm</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import missing from '@/assets/missing.webp'
const props = defineProps(['title', 'owned', 'cover', 'isbn', 'released_at', 'tome_number'])
import axios from 'axios'
import { useDate } from 'vuetify'
import { useDisplay } from 'vuetify'
import { ref } from 'vue'

// --- State ---
// The variable that decides the path
const isDialogOpen = ref(false)
const isLoading = ref(false)
const owned = ref(props.owned)
const date = useDate()
const { mobile } = useDisplay()

// --- Logic ---

// 1. Main Handler
function handleButtonClick() {
  if (owned.value) {
    // Path A: Open Dialog
    isDialogOpen.value = true
  } else {
    // Path B: Direct API Call
    toggleOwned()
  }
}

// 2. Dialog Confirmation Handler
function confirmDialogAction() {
  isDialogOpen.value = false // Close dialog
  toggleOwned() // Proceed to API
}

// 3. The Shared API Logic
async function toggleOwned() {
  isLoading.value = true
  console.log('Calling API to toggle Owned to ' + !owned.value + '...')

  try {
    await axios.patch(`/api/books/${props.isbn}`, {owned: !owned.value})
    owned.value = !owned.value
  } catch (error) {
    console.error('Error changing state:', error)
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
/* Add any specific styles for your card here if needed */
</style>
