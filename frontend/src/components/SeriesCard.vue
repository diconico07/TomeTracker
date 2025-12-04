<template>
  <v-card class="mx-auto" :title="title" :subtitle="author" :to="'/series/' + id ">
    <template v-slot:prepend>
      <v-img
        class="mx-auto"
        size="125"
        :src="cover || missing"
        height="125"
        width="125"
        cover
      ></v-img>
    </template>
    <v-card-text>
      <v-sheet class="d-flex align-center mx-auto">
        <v-progress-linear
          :max="total"
          min="0"
          :model-value="props.current"
          rounded
          height="20"
        >
        <strong>{{ current }}/{{ total }}</strong>
      </v-progress-linear>
      </v-sheet>
    </v-card-text>
  </v-card>
</template>

<script setup>
import missing from '@/assets/missing.webp'
const props = defineProps(['total', 'current', 'title', 'author', 'cover', 'id'])

</script>

<style scoped>
/* Add any specific styles for your card here if needed */
</style>
