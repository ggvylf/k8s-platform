[TOC]
# readme
k8s-platform
k8s管理平台demo

后端：
gin client-go

前端：
vue

# 代码参考
https://github.com/dqzboy/DKube

https://github.com/dnsjia/luban



# 目录结构
```shell
├── config # 配置参数
├── controller # api接口 router配置
├── dao # db交互crud，model中定义的struct的操作
├── db # 连接db相关
├── docs # 文档
├── middle # 中间件
├── model # 数据结构体定义
├── service  #业务逻辑 跟k8s交互 跟db交互
├── utils # 其他工具
└── web  # 前端

```

# 后端流程
## 开发流程
model --> dao -->  service -->  controller -->  router

## web请求处理过程
router --> controller --> service --> dao

# 后端TODO
按照开发进度
- [x] 路由初始化
- [x] clientset初始化
- [x] 数组的排序 过滤 分页
  - [x] 排序 通过sort.Sort()实现
  - [x] 过滤 
  - [x] 分页
- [ ] 工作负载workload 
  - [x] pod
    - [x] podList
    - [x] 获取pod信息
    - [x] 删除pod
    - [x] 更新pod
    - [x] 各个ns下的pod数量
  - [x] containers
    - [x] 获取pod中container.Name的list
    - [x] 获取container的log
  - [x] pod的gin route
  - [ ] deploment
    - [x] deploment list
    - [x] 更新deployment
    - [x] 创建deployment
    - [x] 重启deployment
    - [x] 各个ns下的deployment数量
  - [ ] deployment的gin route
  - [ ] daemonset
  - [ ] statefuset
- [ ] 集群资源
  - [ ] node
  - [ ] ns
  - [ ] pv
- [ ] 网络
  - [ ] service
    - [x] createservice 
  - [ ] ingress
    - [x] createingress
- [ ] 存储
  - [ ] configmap
  - [ ] secret
  - [ ] pvc
- [ ] workflow
  - [x] gorm db连接初始化 
  - [x] model
  - [x] dao
  - [x] service
  - [x] controller
  - [x] router
- [ ] middleware
  - [x] cors
  - [x] jwt
- [x] web终端

# 前端流程
## 开发流程
/route/index.js --> src/views/xx.vue --> html+css 布局 小组件 js 动态数据

## 请求处理流程

index.html --> App.vue --> route/index.js --> src/views/xx.vue

## 代码目录
```shell
.
├── public
└── src
    ├── assets
    │   ├── avator
    │   ├── img
    │   └── k8s
    ├── components
    │   └── icons
    ├── layout # 布局
    ├── router # 路由
    ├── utils # 工具 拦截器
    └── views # 视图
        ├── common #公共视图和后端接口
        ├── damonset
        ├── deployment
        ├── home
        └── pod

```
## npm命令相关
```shell
# 初始化vue项目
npm init vue@latest

# 依赖包
npm install element-plus vue-router nprogess  axios codemirror-editor-vue3 js-yaml json2yaml xterm xterm-addon-fit echarts


# 测试
npm run dev

# 构建
npm build
```


# 前端TODO
- [x] 初始化vue项目 npm init vue@latest  npm install
- [x] 初始化vite.config.js
- [x] 初始化main.js 
- [x] 初始化App.vue
- [x] 创建views router layout utils目录
- [x] 创建路由 route/index.js vue-router
- [x] 安装相关插件 npm install element-plus
- [x] 启动项目 npm run dev
- [x] 进度条 Nprogress
- [x] axios http 拦截器
- [x] 错误页处理
- [x] 根路径跳转/home
- [x] 静态资源引用
- [x] 整体框架布局
- [x] 后端api接口 Config.js
- [x] workload页面布局
  - [ ] deployment 
    - [x] deployment的header1
      - [x] 下拉列表
      - [x] 刷新
    - [x] deployment的header2
      - [x] 创建按钮 弹出隐藏抽屉
      - [x] 搜索框
      - [x] 搜索按钮  inline
    - [x] deployment的隐藏抽屉
      - [x] 抽屉的结构 
      - [x] 隐藏抽屉的气泡弹出框
      - [x] 隐藏抽屉的表单校验
      - [x] 隐藏抽屉的表单提交
    - [ ] 数据表格
      - [x] getDeployment() 函数
      - [x] 表格结构
      - [x] label处理
      - [x] 分页操作
      - [x] 操作按钮
        - [x]  npm install codemirror-editor-vue3 js-yaml json2yaml 
        - [x] 全局加载codemirror
        - [x] yaml详细信息
        - [x] 调整replicas
        - [x] 删除deployment
        - [x] 重启depoyment
  - [ ] pod
    - [x] 拷贝deployment.vue 替换变量
    - [x] 删减功能 创建pod 扩缩容 重启
    - [x] 修改pod表格字段，增加pod需要的内容
    - [x] 折叠框展开和点击podname展开
    - [x] container日志
    - [ ] webshell npm install xterm xterm-addon-fit 
      - [x] socket
      - [x] xtrem 
  - [ ] daemonset
    - [x] 增加路由规则 
    - [x] 拷贝deployment代码
    - [x] 修改容器组的计算参数，跟deployment不同
  - [ ] statusful
    - [x] 增加路由规则 
    - [x] 拷贝daemonset代码
    - [x] 修改容器组的计算参数，跟deployment相同
- [ ] 集群管理
  - [x] 修改路由规则
  - [x] 创建文件和目录
  - [x] node
  - [x] namespace
  - [x] pv
- [ ] 负载均衡
  - [x] 修改路由规则
  - [x] 创建文件和目录
  - [x] service 修改表格字段 
  - [x] 抽屉显示换行问题 label-width="100px" label-width="auto"
  - [x] ingress
- [ ] 存储和配置
  - [x]] pvc
  - [x] configmap
  - [x] secret
- [ ] 工作概览 /home
  - [ ] 仪表盘  
    - [x] 布局 使用折叠面板
    - [x] 安装echarts组件 npm install echarts
    - [x] 圆形进度条使用el-progress
    - [x] echarts柱状图