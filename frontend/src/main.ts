import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import './style.css';
import { createMemoryHistory, createRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import HomePage from './page/Home.vue'
import InjectPage from './page/Inject.vue'

const routes = [
    {
        path: '/',
        name: "home",
        component: HomePage,
    },
    {
        path: '/inject',
        name: "inject",
        component: InjectPage,
    },
]

const router = createRouter({
    history: createMemoryHistory(),
    routes,
})

const app = createApp(App)

app.use(ElementPlus)
app.use(router)
app.mount('#app')

app.config.errorHandler = (err, vm, info) => {
    ElMessage({
        message: err as any,
        type: "error",
    })
};
