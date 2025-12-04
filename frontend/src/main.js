import { createApp } from 'vue'
// Vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import 'unfonts.css'
import '@mdi/font/css/materialdesignicons.css' // Ensure you are using css-loader

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(router)

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'system',
  },
})

app.use(vuetify)


app.mount('#app')
