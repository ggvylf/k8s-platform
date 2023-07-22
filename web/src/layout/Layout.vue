<template>
    <div class="common-layout">
        <!-- 全局container -->
        <el-container style="height: 100vh;">
            <!-- 侧边栏 -->
            <el-aside class="aside" :width="asideWidth">
                <!-- 侧边栏固钉 z-index显示优先级 -->
                <el-affix class="aside-affix" :z-index="1200">
                    <div class="aside-logo" >
                        <!-- 平台logo图片 -->
                        <el-image class="logo-image" :src="logo" />
                        <!-- 折叠侧边栏，折叠后不显示名称 -->
                        <span :class="[isCollapse ? 'is-collapse' : '']">
                            <!-- 平台显示名称 -->
                            <span class="logo-name" >k8s-platform</span>
                        </span>
                    </div>
                </el-affix>
                <!-- 侧边栏菜单 -->
                <!-- 菜单和路由做了关联 菜单折叠做了关联 -->
                <el-menu class="aside-menu"
                    router
                    :default-active="$route.path"
                    :collapse="isCollapse"
                    background-color="#131b27"
                    text-color="#bfcbd9"
                    active-text-color="#20a0ff">
                    <!-- 遍历路由生成菜单 -->
                    <div v-for="menu in routers" :key="menu">
                        <!-- 只有一个子路由的情况，例如摘要，主路由的路径就是子路由 生效的是子路由 主路由写啥都行 -->
                        <el-menu-item class="aside-menu-item" v-if="menu.children && menu.children.length == 1" :index="menu.children[0].path">
                        <!-- 引入图标     -->
                        <el-icon><component :is="menu.children[0].icon" /></el-icon>
                        <template #title>
                            {{menu.children[0].name}}
                        </template>
                        </el-menu-item>
                        <!-- 子路由有多个 显示用的主路由的信息 -->
                        <!-- el-sub-menu折叠后title不会消失，需要自行处理  el-menu-item会消失-->
                        <el-sub-menu class="aside-submenu" v-else-if="menu.children" :index="menu.path">
                            <template #title>
                                <el-icon><component :is="menu.icon" /></el-icon>
                                <span :class="[isCollapse ? 'is-collapse' : '']">{{menu.name}}</span>
                            </template>
                            <el-menu-item class="aside-menu-childitem" v-for="child in menu.children" :key="child" :index="child.path">
                                <el-icon><component :is="child.icon" /></el-icon>
                                <template #title>
                                    {{child.name}}
                                </template>
                            </el-menu-item>
                        </el-sub-menu>
                    </div>
                </el-menu>
            </el-aside>
            <!-- 页面container 3个部分 header main footer -->
            <el-container>
                <el-header class="header" >
                    <el-row :gutter="20">
                        <el-col :span="1">
                            <!-- 折叠按钮 -->
                            <div class="header-collapse" @click="onCollapse">
                                <el-icon><component :is="isCollapse ? 'expand':'fold'" /></el-icon>
                            </div>
                        </el-col>
                        <el-col :span="10" >
                            <!-- 面包屑 这里用/作为分割符 -->
                            <div class="header-breadcrumb">
                                <el-breadcrumb separator="/" v-if="this.$route.matched[0].path != '/main'">
                                    <!-- 根路径写死 -->
                                    <el-breadcrumb-item :to="{ path: '/' }">工作台</el-breadcrumb-item>
                                    <template v-for="(matched,m) in this.$route.matched" :key="m">
                                        <el-breadcrumb-item v-if="matched.name != undefined" >
                                        {{ matched.name }}
                                        </el-breadcrumb-item>
                                    </template>
                                </el-breadcrumb>
                                <el-breadcrumb separator="/" v-else>
                                    <el-breadcrumb-item>工作台</el-breadcrumb-item>
                                </el-breadcrumb> 
                            </div>
                        </el-col>
                        <!-- 用户信息相关 -->
                        <el-col class="header-menu" :span="13">
                            <el-dropdown>
                                <div class="header-dropdown">
                                    <el-image class="avator-image" :src="avator" />
                                    <span>{{ username }}</span>
                                </div>
                                <template #dropdown>
                                    <el-dropdown-menu>
                                        <el-dropdown-item icon="el-icon-switch-button" @click="logout()">退出</el-dropdown-item>
                                        <el-dropdown-item icon="el-icon-unlock">修改密码</el-dropdown-item>
                                    </el-dropdown-menu>
                                </template>
                            </el-dropdown>
                        </el-col>
                    </el-row>
                </el-header>
                <!-- main相关 -->
                <el-main class="main">
                    <!-- 嵌套路由视图 -->
                    <router-view></router-view>
                </el-main>
                <!-- footer相关 -->
                <el-footer class="footer">
                    <el-icon style="width:2em;top:3px;font-size:18px"><place/></el-icon>
                    <a class="footer el-icon-place">2023 DevOps </a>
                </el-footer>
                <!-- 返回顶部 -->
                <el-backtop target=".el-main"></el-backtop>
            </el-container>
        </el-container>
    </div>
</template>

<script>
import {useRouter} from 'vue-router'
export default {
    // 前置操作
    // 从routes中获取全部的router规则
    beforeMount() {
        this.routers = useRouter().options.routes
    },
    // 数据
    data() {
        return {
            avator: require('@/assets/avator/avator.png'),
            logo: require('@/assets/k8s/k8s-metrics.png'),
            isCollapse: false,
            asideWidth: '220px',
            routers: [],
        }
    },
    computed: {
        // 获取用户名
        username() {
            let username = localStorage.getItem('username');
            return username ? username : '未知';
        },
    },
    methods: {
        // 折叠操作
        onCollapse() {
            if (this.isCollapse) {
                this.asideWidth = '220px'
                this.isCollapse = false
            } else {
                this.isCollapse = true
                this.asideWidth = '64px'
            }
        },
        // 登出
        logout() {
            localStorage.removeItem('username');
            localStorage.removeItem('token');
            this.$router.push('/login');
        }
    },

}
</script>

<style scoped>
/* 侧边栏 */
    .aside {
        transition: all .5s;
        background-color: #131b27;
    }
    /* 固钉 */
    .aside-logo {
        background-color: #131b27;
        height: 60px;
        color: white;
        cursor: pointer;
    }
    .logo-image {
        width: 40px;
        height: 40px;
        top: 12px;
        padding-left: 12px;
    }
    .logo-name {
        font-size: 20px;
        font-weight: bold;
        padding: 10px;
    }
    /* 滚动条 */
    .aside::-webkit-scrollbar {
        display: none;
    }
    /* 修整边框 */
    .aside-affix {
        border-bottom-width: 0;
    }
    .aside-menu {
        border-right-width: 0;
    }
    /* 菜单栏的颜色 */
    /* 选中 */
    .aside-menu-item.is-active {
        background-color: #1f2a3a;
    }
    .aside-menu-item {
        padding-left: 20px !important;
    }
    .aside-menu-childitem {
        padding-left: 20px !important;
    }
    .aside-menu-childitem.is-active {
        background-color: #1f2a3a;
    }
    /* 悬停 */
    .aside-menu-childitem:hover {
        background-color: #142c4e;
    }

    /* header相关 */
    .header {
        /* 优先级 */
        z-index:1200;
        line-height: 60px;
        font-size: 24px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, .12),0 0 6px rgba(0, 0, 0, .04)
    }
    .header-collapse {
        cursor: pointer;
    }
    .header-breadcrumb {
        padding-top: 0.9em;
    }
    .header-menu {
        text-align: right;
    }
    .is-collapse {
        display: none;
    }
    .header-dropdown {
        line-height: 60px;
        cursor: pointer;
    }
    .avator-image {
        top: 12px;
        width: 40px;
        height: 40px;
        border-radius: 50%;
        margin-right: 8px;
    }
    .main {
        padding: 10px;
    }
    .footer {
        z-index: 1200;
        color: rgb(187, 184, 184);
        font-size: 14px;
        text-align: center;
        line-height: 60px;
    }
</style>