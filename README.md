[TOC]
# readme
k8s-platform
k8s管理平台demo

后端：
gin client-go

前端：
vue

# 目录结构
```shell
├── config # 配置参数
├── controller # api
├── dao # db交互
├── db # db相关
├── docs # 文档
├── middle # 中间件
├── modle # struct定义
├── service #业务逻辑
├── utils # 其他工具
└── web  # 前端

```


# 后端TODO
按照开发进度
- [x] 路由初始化
- [x] clientset初始化
- [x] 数组的排序 过滤 分页
  - [x] 排序 通过sort.Sort()实现
  - [x] 过滤 
  - [x] 分页
- [ ] 工作负载workload 
  - [ ] pod
    - [x] podList
    - [x] 获取pod信息
    - [x] 删除pod
    - [x] 更新pod
    - [x] 各个ns下的pod数量
  - [ ] containers
    - [x] 获取pod中container.Name的list
    - [x] 获取container的log
  - [ ] pod的gin route
  - [ ] deploment
    - [x] deploment list
    - [x] 更新deployment
    - [ ] 创建deployment
  - [ ] 
- [ ] cluster操作
- [ ] 





# 前端TODO
- [ ] 

