// 导入路由模式
import {createRouter,createWebHistory} from 'vue-router'


// 路由规则
const routers= [
    {
        // 视图
        path:'/home',
        // 视图组件
        component:() => import('@/views/home/Home.vue'),
        // 图标
        icon: 'odometer',

        // meta信息
        meta:{
            title:"概要",
            requireAuth: false,
        },
        
    },
]
// 创建路由实例
const router = createRouter({
    history:createWebHistory(),
    routers
})

// 导出路由实例
export default router