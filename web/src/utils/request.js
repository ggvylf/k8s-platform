import axios from 'axios'


// 创建axios实例
const httpClient = axios.create({
    // 相应状态码验证，不在范围的不给相应
    validateStatus:{
        return status >=200 && status < 504
    },
    // 超时，单位是ms
    timeout: 10000
})

// 请求重试相关
httpClient.defaults.retry=3
httpClient.defaults.retryDelay=1000
httpClient.defaults.shouldRetry=true


// 请求拦截器
httpClient.interceptors.request.use(
    config => {
        // 添加header
        config.headers['Content-Type']='appliaction/json'
        config.headers['Accept-Language']='zh-CN'
        config.headers['authorization']=localStorage.getItem("token")

        // 处理post请求
        if (config.method =='post') {
            if (!config.data) {
                config.data={}
            }
        }
        return config
    },
    err => {
      return  Promise.reject(err)
    }
)

// 响应拦截器
httpClient.interceptors.response.use(
    response =>{
        // 处理状态码
        if (response.status !==200) {
            return Promise.reject(response.data)
        } else {
            return response.data
        }

    },
    err => {
        return Promise.reject(err)
    }
)


// 导出httpClient
export default httpClient