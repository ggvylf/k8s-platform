// 导入路由模式
import {createRouter,createWebHistory} from 'vue-router'

// 进度条
import Nprogress from 'nprogress'
import 'nprogress/nprogress.css'

// 布局
import Layout from '@/layout/Layout.vue'

// 路由规则
const routers = [

    // 访问/跳转到/home
    {
        path:'/',
        redirect:'/home',
    },


    // 主页面
    {
        path:'/home',
        component: Layout,
        // component:() => import('@/layout/Layout.vue'),
        meta:{
            title:"k8s",
            requireAuth: false,
        },
        children: [
            {
                // 视图
                path:'/home',
                name: "概要",
                // 视图组件
                component:() => import('@/views/home/Home.vue'),
                // 图标
                icon: 'odometer',
        
                // meta信息
                meta:{
                    title:"概要",
                    requireAuth: false,
                }
                
            },
        ]
        
    },


    // 工作负载
    {
        path:'/workload',
        name: '工作负载',
        component: Layout,
        icon:'menu',
        meta:{
            title:"工作负载",
            requireAuth: false,
        },
        children: [
            {
                path:'/workload/deployment',
                name: 'Deployment',
                icon: 'el-icon-s-data',
                meta:{
                    title:"deployment",
                    requireAuth: false,
                },
                component:() => import('@/views/deployment/Deployment.vue'),
                        
            },
            {
                path:'/workload/pod',
                name: 'Pod',
                icon: 'el-icon-document-add',
                meta:{
                    title:"pod",
                    requireAuth: false,
                },
                component:() => import('@/views/pod/Pod.vue'),
                        
            },
        ]
        
    },
    



    // 错误页
    {
        path:'/404',
        component:() => import('@/views/common/404.vue'),
        meta:{
            title:"404 page",
            requireAuth: false,
        }
        
    },
    {
        path:'/403',
        component:() => import('@/views/common/403.vue'),
        meta:{
            title:"403 page",
            requireAuth: false,
        }
        
    },
    // 其他页跳转到404
    {
        path:"/:pathMatch(.*)",
        redirect: '/404',
        
    },


]
// 创建路由实例
const router = createRouter({
    history: createWebHistory(),
    // 注意参数名，别写错
    routes:routers
})

// 进度条配置
// 进度条递增值
Nprogress.inc(0.2)
Nprogress.configure({
    // 动画效果
    easing: 'ease',
    // 速度
    speed: 600,
    // 进度环
    showSpinner: false,
})

// 路由守卫
router.beforeEach((to,from,next) => {
    // 启动进度条
    Nprogress.start()
    // 设置头部
    if (to.meta.title) {
        document.title=to.meta.title
    } else {
        document.title="k8s"
    }


    // 放行
    next()
})

router.afterEach((to,from,next) => {
    Nprogress.done()
})


// 导出路由实例
export default router