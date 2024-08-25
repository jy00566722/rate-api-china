import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
// import './style.css'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
const app = createApp(App)
const pinia = createPinia()
for (const i in ElementPlusIconsVue) {
    app.component(i, ElementPlusIconsVue[i])
}
app.use(pinia).use(ElementPlus).use(router).mount('#app')