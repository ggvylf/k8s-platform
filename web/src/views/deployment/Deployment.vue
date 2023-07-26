<template>
<div class="deploy">
    <el-row>
        <!-- header1 用来选择ns 和刷新 -->
        <el-col :span="24">
            <div>
                <el-card class="dfeploy-head-card" shadow="never" :body-style="{padding:'10'}">
                    <el-row>
                        <!-- 下拉列表 -->
                        <el-col :span="6">
                            <div>
                                <span>命名空间: </span>
                                <el-select v-model="namespaceValue" filterable placeholder="请选择">
                                    <el-option v-for="(item,index) in namespaceList" 
                                    :key="index" 
                                    :label="item.metadata.name" 
                                    :value="item.metadata.name">
                                    </el-option>
                                </el-select>
                            </div>
                        </el-col>
                        <!-- 刷新按钮 -->
                        <el-col :span="2" :offset=16>
                            <div>
                                <el-button style="border-radius:2px" icon="Refresh" plain>刷新</el-button>
                            </div>
                        </el-col>
                    </el-row>

                </el-card>
            </div>
        </el-col>
        <!-- header2 用来创建资源 和提供搜索-->
        <el-col :span="24">
            <div>
                <el-card class="deploy-head-card" shadow="never" :body-style="{padding:'10px'}">
                    <el-row>
                        <!-- 创建资源按钮 -->
                        <el-col :span="2">
                            <div>
                                <el-button style="border-radius: 2px;" icon="edit" type="primary" v-loading.fullscreen.lock="fullscreenloading" @click="createDeploymentDrawer=true">创建</el-button>
                            </div>
                        </el-col>
                        <!-- 搜索框和搜索按钮 -->
                        <el-col :span="6">
                            <div>
                                <el-input class="deploy-head-serach" clearable placeholder="请输入" v-model="searchinput"></el-input>
                                <el-button style="border-radius: 2px;" icon="search" type="primary" plain @click="getDeployments()">搜索</el-button>
                            </div>
                        </el-col>
                </el-row>
                </el-card>
            </div>
        </el-col>
        <!-- 抽屉 -->
        <el-drawer 
        v-model="createDeploymentDrawer" 
        :direction="direction" 
        :before-close="handleClose">
            <!-- 抽屉标题 -->
            <template #header>
                <h4>创建Deployment</h4>
            </template>
            <!-- 抽屉表单 -->
            <template #default>
                <el-row type="flex" justify="center">
                    <el-col :span="20">
                        <!-- 增加校验 -->
                        <el-form ref="createDeployment" 
                        :rules="createDeploymentRules" 
                        :model="createDeployment">
                        <el-form-item class="deploy-create-form" label="名称" prop="name">
                                <el-input v-model="createDeployment.name"></el-input>
                        </el-form-item>
                        <el-form-item class="deploy-create-form" label="命名空间" prop="namespace">
                            <el-select v-model="createDeployment.namespace" filterable placeholder="请选择">
                                <el-option
                                v-for="(item, index) in namespaceList"
                                :key="index"
                                :label="item.metadata.name"
                                :value="item.metadata.name">
                                </el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item class="deploy-create-form" label="副本数" prop="replicas">
                            <el-input-number v-model="createDeployment.replicas" :min="1" :max="10"></el-input-number>
                            <!-- 气泡弹出框 -->
                                <el-popover
                                    placement="top"
                                    :width="100"
                                    trigger="hover"
                                    content="申请副本数上限为10个">
                                    <template #referen ce>
                                        <el-icon style="width:2em;font-size:18px;color:#4795EE"><WarningFilled/></el-icon>
                                    </template>
                                </el-popover>
                        </el-form-item>
                        <el-form-item class="deploy-create-form" label="镜像" prop="image">
                            <el-input v-model="createDeployment.image"></el-input>
                        </el-form-item>
                        <el-form-item class="deploy-create-form" label="标签" prop="label_str">
                            <el-input v-model="createDeployment.label_str" placeholder="示例: project=ms,app=gateway"></el-input>
                        </el-form-item>
                        <el-form-item class="deploy-create-form" label="资源配额" prop="resource">
                            <!-- 常用规格 格式cpu/mem  用/做分隔符 -->
                            <el-select v-model="createDeployment.resource" placeholder="请选择">
                                <el-option value="0.5/1" label="0.5C1G"></el-option>
                                <el-option value="1/2" label="1C2G"></el-option>
                                <el-option value="2/4" label="2C4G"></el-option>
                                <el-option value="4/8" label="4C8G"></el-option>
                                <el-option value="8/16" label="8C16G"></el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item class="deploy-create-form" label="容器端口" prop="container_port">
                            <el-input v-model="createDeployment.container_port" placeholder="示例: 80"></el-input>
                        </el-form-item>
                        <el-form-item class="deploy-create-form" label="健康检查" prop="health">
                            <el-switch v-model="createDeployment.health_check" />
                        </el-form-item>
                        <el-form-item class="deploy-create-form" label="检查路径" prop="healthPath">
                            <el-input v-model="createDeployment.health_path" placeholder="示例: /health"></el-input>
                        </el-form-item>
                        </el-form>
                    </el-col>
                </el-row>
            </template>
            <!-- 抽屉底部 -->
            <template #footer>
                <el-button @click="createDeploymentDrawer = false">取消</el-button>
                <!-- 提交ref匹配的form表单 -->
                <el-button type="primary" @click="submitForm('createDeployment')">立即创建</el-button>
            </template>
        </el-drawer>
        <!-- 数据表格 -->
        <el-col :span="24">
            <div>
                789
            </div>
        </el-col>
    </el-row>
</div>
</template>


<script>
import common from "@/views/common/Config"
import httpClient from "@/utils/request"
export default { 
    data() {
        return {
            // ns
            namespaceValue:"default",
            namespaceList:[],
            namespaceListUrl: common.k8sNamespaceList,

            // deploy
            createDeploymentDrawer: false,
            fullscreenloading:false,
            direction:"rtl",

            // 抽屉的默认值
            createDeployment:{
                name:"",
                namespace: "default",
                replicas: 1,
                resource:"",
                health_check: false,
                health_path:'',
                label_str:"",
                label:{},
                container_port:"",
            },
            // 创建deploy
            createDeploymentData:{
                url: common.k8sDeploymentCreate,
                params:{}
            },
            // 表单校验
            createDeploymentRules:{
                name: [{
                    required: true,
                    message: '请填写名称',
                    trigger: 'change'
                }],
                image: [{
                    required: true,
                    message: '请填写镜像',
                    trigger: 'change'
                }],
                namespace: [{
                    required: true,
                    message: '请选择命名空间',
                    trigger: 'change',
                }],
                resource: [{
                    required: true,
                    message: '请选择配额',
                    trigger: 'change'
                }],
                label_str: [{
                    required: true,
                    message: '请填写标签',
                    trigger: 'change'
                }],
                container_port: [{
                    required: true,
                    message: '请填写容器端口',
                    trigger: 'change'
                }],

            },


            // search
            searchinput:"",
            
        }
    },
    methods: {
        getNamespaceList() {
            httpClient.get(this.namespaceListUrl)
            .then(res =>{
                this.namespaceList=res.data.items
            })
            .catch(res => {
                this.$message.error({
                    message: res.meg
                })
            })
        },
        // 抽屉的关闭确认
        handleClose(done) {
            this.confirm("确认关闭")
            .then(()=>{
                done();
            })
            .catch(()=>{
                
            })

        },

        // 创建deploy
        submitForm(formName){
            // 验证表单
            this.$refs[formName].validate((valid)=>{
                if (valid) {
                    this.createDeployFunc()
                }else {
                    return false
                }
            })
        },
        createDeployFunc(){
            // 验证label字段是否合规
            let reg = new RegExp("(^[A-Za-z]+=[A-Za-z0-9]+).*")
            if (!reg.test(this.createDeployment.label_str)) {
                this.$message.warning({
                    message: "标签填写异常，请确认后重新填写"
                })
                return
            }

            // 加载loading动画
            this.fullscreenLoading = true

            // 处理cpu和mem
            let cpu, memory
            let resourceList = this.createDeployment.resource.split("/")

            // 处理label，把字符串转换成map
            let label = new Map()
            let a = (this.createDeployment.label_str).split(",")
            a.forEach(item => {
                let b = item.split("=")
                label[b[0]] = b[1]
            })


            // 赋值表单数据
            cpu = resourceList[0]
            memory = resourceList[1] + "Gi"
            this.createDeploymentData.params = this.createDeployment
            this.createDeploymentData.params.container_port = parseInt(this.createDeployment.container_port)
            this.createDeploymentData.params.label = label
            this.createDeploymentData.params.cpu = cpu
            this.createDeploymentData.params.memory = memory


            // 提交表单
            httpClient.post(this.createDeploymentData.url, this.createDeploymentData.params)
            .then(res => {
                this.$message.success({
                message: res.msg
                })
                this.getDeployments()
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })

            // 后续操作
            this.resetForm('createDeployment')
            this.fullscreenLoading = false
            this.createDeploymentDrawer = false
             
        },

        // 重置表单
        resetForm(formName) {
            this.$refs[formName].resetFields()
        },

    },
    watch: {
        // 把ns写入localstorage，方便其他页面调用
        namespaceValue: {
            handler() {
                localStorage.setItem('namespace',this.getNamespaceValue)
            }
        }
    },
    beforeMount() {
        // 把ns的值同步到其他页面
        if (localStorage.getItem('namespace') !== "undfined" && localStorage.getItem('namespace') !==null) {
            this.namespaceValue=localStorage.getItem('namespace')
        }
        this.getNamespaceList()

    },
}


</script>

<style>
    .deploy-head-card,.deploy-body-card{
        border-radius: 1px;
        margin-bottom: 5px;

    }
    .deploy-head-search {
        width: 160px;
        margin-right: 10px;
    }

</style>