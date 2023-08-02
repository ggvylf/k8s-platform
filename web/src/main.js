import { createApp } from 'vue'

// vue
import App from './App.vue'

// elementplus
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

// 图标视图
import * as ELIcons from '@element-plus/icons-vue'

// 路由配置
import router from './router'

// codeimrror
import { GlobalCmComponent } from 'codemirror-editor-vue3';
// 外观
// import 'codemirror/theme/idea.css'
// 引入支持语言，这里只用到了yaml
// import 'codemirror/mode/yaml/yaml.js'


// createApp(App).mount('#app')
// 创建vue实例
const app=createApp(App)

// 图标注册为全局组件
for(let iconName in ELIcons) {
    app.component(iconName,ELIcons[iconName])

}

// 使用elementplus
app.use(ElementPlus)

// 使用codemirror
app.use(GlobalCmComponent, { componentName: "codemirror" });

// 使用路由
app.use(router)
// 挂载
app.mount('#app')