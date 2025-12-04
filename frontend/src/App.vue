<script setup>
import { RouterView } from 'vue-router'
import { useDisplay } from 'vuetify'
import logo from '@/assets/logo.svg'

// 1. Destructure 'mobile' from useDisplay to detect screen size
// 'mobile' returns true if the screen is 'xs' (phones) or 'sm' (small tablets)
const { mobile, lgAndDown } = useDisplay()

const items = [
  { title: 'Collection', value: 'home', icon: 'mdi-bookshelf', path: '/' },
  { title: 'Planning', value: 'planning', icon: 'mdi-calendar', path: '/planning' },
  { title: 'Missing Books', value: 'missing', icon: 'mdi-crosshairs-question', path: '/missing' },
]
</script>

<template>
    <v-app id="tometracker">
    <v-navigation-drawer
    v-if="!mobile"
        :expand-on-hover="lgAndDown"
        permanent
        :rail="lgAndDown"
      >
        <v-list>
          <v-list-item
            :prepend-avatar="logo"
            title="TomeTracker"
          ></v-list-item>
        </v-list>

        <v-divider></v-divider>

        <v-list density="compact" nav>
          <v-list-item
          v-for="item in items"
          :value="item.value"
          :prepend-icon="item.icon"
          :title="item.title"
          :to="item.path"
          :key="item.title"
        ></v-list-item>
        </v-list>
      </v-navigation-drawer>

    <v-main>
      <RouterView />
    </v-main>

    <v-bottom-navigation
      v-if="mobile"
      color="primary"
      app 
      grow 
    >
      <v-btn 
        v-for="item in items" 
        :value="item.value" 
        :to="item.path"
        :key="item.title"
      >
        <v-icon>{{ item.icon }}</v-icon>
        <span>{{ item.title }}</span>
      </v-btn>
    </v-bottom-navigation>
  </v-app>
</template>