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


# 目录结构
```shell
├── config # 配置参数
├── controller # api router接口
├── dao # db交互相关，crud model中定义的struc的操作
├── db # db相关
├── docs # 文档
├── middle # 中间件
├── model # 数据结构体定义
├── service  #业务逻辑 跟k8s交互 跟db交互
├── utils # 其他工具
└── web  # 前端

```
# 开发流程
model --> dao -->  service -->  controller -->  router

# 请求处理过程
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
  - [ ] 数据处理
  - [ ] client-go crud 
  - [ ] mysql crud
  - [ ] 缓存crud 例如redis
- [ ] model
  - [x] workflow 
- [ ] dao
  - [x] workflow  
- [ ] service
  - [x] workflow  
- [ ] middleware
- [ ] 




# 前端TODO
- [ ] 

