<template>
    <div class="daemonset">
        <el-row>
            <!-- header1 用来选择ns 和刷新 -->
            <el-col :span="24">
                <div>
                    <el-card class="daemonset-head-card" shadow="never" :body-style="{padding:'10'}">
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
                            <!-- 刷新的时候会重新请求表格数据 -->
                            <el-col :span="2" :offset=16>
                                <div>
                                    <el-button style="border-radius:2px" icon="Refresh" plain @click="getDaemonsetments()">刷新</el-button>
                                </div>
                            </el-col>
                        </el-row>

                    </el-card>
                </div>
            </el-col>

            <!-- header2 用来创建资源 和提供搜索-->
            <el-col :span="24">
                <div>
                    <el-card class="daemonset-head-card" shadow="never" :body-style="{padding:'10px'}">
                        <el-row >
                            <!-- 创建资源按钮 -->
                            <el-col :span="2" >
                                <div>
                                    <el-button style="border-radius: 2px;" icon="edit" type="primary" v-loading.fullscreen.lock="fullscreenloading" @click="createDaemonsetmentDrawer=true">创建</el-button>
                                </div>
                            </el-col>
                            <!-- 搜索框和搜索按钮 -->
                            <el-col :span="6" >
                                <div >
                                    <el-input class="daemonset-head-search" clearable placeholder="请输入" v-model="searchInput"></el-input>
                                    <el-button style="border-radius: 2px;" icon="search" type="primary" plain @click="getDaemonsetments()">搜索</el-button>
                                </div>
                            </el-col>
                    </el-row>
                    </el-card>
                </div>
            </el-col>

            <!-- 数据表格 -->
            <el-col :span="24">
                <div>
                    <el-card class="daemonset-body-card" shadow="never" :body-style="{padding:'5px'}">
                        <el-table style="width: 100%;font-size: 12px;margin-bottom: 10px;"
                        :data="daemonsetList">
                            <el-table-column width="20"></el-table-column>
                            <el-table-column align="center" label="名称">
                                <template v-slot="scope">
                                    <a class="daemonset-body-daemonsetname">{{ scope.row.metedata.name }}</a>
                                </template>
                            </el-table-column>
                            <el-table-column align="center" label="标签">
                                <template v-slot="scope">
                                    <!-- 遍历label -->
                                    <div v-for="(val,key) in scope.row.metadata.labels" key="key">
                                        <!-- 气泡框 -->
                                        <el-popover
                                        placement="right"
                                        :width="200"
                                        trigger="hover"
                                        :content="key+':'+value">
                                            <template #reference>
                                                <!-- 对超过长度的label做省略处理 -->
                                                <el-tag style="margin-bottom: 5px;" type="warning">ellipsis({{ key+':'+value }})</el-tag>
                                            </template>
                                        </el-popover>
                                    </div>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="容器组">
                                <template v-slot="scope">
                                    <!-- 显示当前容器和总容器个数 -->
                                    <span>{{ scope.row.status.numberAvailable>0?scope.row.status.numberAvailable:0  }} / {{ scope.row.status.desiredNumberScheduled>0?scope.row.status.desiredNumberScheduled:0 }} </span>
                                </template>
                            </el-table-column>
                            <el-table-column align=center min-width="100" label="创建时间">
                                <template v-slot="scope">
                                    <!-- 时间转换 -->
                                    <el-tag type="info">{{ timeTrans(scope.row.metadata.creationTimestamp) }} </el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column align=center label="镜像">
                                <template v-slot="scope">
                                    <!-- 遍历image列表 -->
                                    <div v-for="(val, key) in scope.row.spec.template.spec.containers" :key="key">
                                        <el-popover
                                            placement="right"
                                            :width="200"
                                            trigger="hover"
                                            :content="val.image">
                                            <template #reference>
                                                <!-- 对image名称做分割 超过长度省略处理-->
                                                <el-tag style="margin-bottom: 5px">{{ ellipsis(val.image.split('/')[2]==undefined?val.image:val.image.split('/')[2]) }}</el-tag>
                                            </template>
                                        </el-popover>
                                    </div>
                                </template>
                            </el-table-column>
                            <!-- 资源操作按钮 -->
                            <el-table-column align="center" label="操作" width="400" >
                                <template v-slot="scope">
                                    <el-button size="small" style="border-radius:2px;" icon="Edit" type="primary" plain @click="getDaemonsetmentDetail(scope)">YAML</el-button>
                                    <el-button size="small" style="border-radius:2px;" icon="Plus" type="primary" @click="handleScale(scope)">扩缩</el-button>
                                    <el-button size="small" style="border-radius:2px;" icon="RefreshLeft" type="primary" @click="handleConfirm(scope, '重启', restartDaemonsetment)">重启</el-button>
                                    <el-button size="small" style="border-radius:2px;" icon="Delete" type="danger" @click="handleConfirm(scope, '删除', delDaemonsetment)">删除</el-button>
                                </template>
                            </el-table-column>
                        </el-table>

                        <!-- 分页 -->
                        <!-- page-size和currnet-page变更时触发对应函数 -->
                        <el-pagination
                            class="daemonset-body-pagination"
                            background
                            @size-change="handleSizeChange"
                            @current-change="handleCurrentChange"
                            current-page="currentPage"
                            :page-sizes="pagesizeList"
                            page-size="pagesize"
                            layout="total, sizes, prev, pager, next, jumper"
                            :total="daemonsetTotal">
                        </el-pagination>
                    </el-card>
                </div>
            </el-col>
        </el-row>

        <!-- 创建daemonset抽屉 -->
        <el-drawer 
            v-model="createDaemonsetmentDrawer" 
            :direction="direction" 
            :before-close="handleClose">
            <!-- 抽屉标题 -->
            <template #header>
                <h4>创建Daemonsetment</h4>
            </template>
            <!-- 抽屉表单 -->
            <template #default>
                <el-row type="flex" justify="center">
                    <el-col :span="20">
                        <!-- 增加校验 -->
                        <el-form ref="createDaemonsetment" 
                        :rules="createDaemonsetmentRules" 
                        :model="createDaemonsetment">
                        <el-form-item class="daemonset-create-form" label="名称" prop="name">
                                <el-input v-model="createDaemonsetment.name"></el-input>
                        </el-form-item>
                        <el-form-item class="daemonset-create-form" label="命名空间" prop="namespace">
                            <el-select v-model="createDaemonsetment.namespace" filterable placeholder="请选择">
                                <el-option
                                v-for="(item, index) in namespaceList"
                                :key="index"
                                :label="item.metadata.name"
                                :value="item.metadata.name">
                                </el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item class="daemonset-create-form" label="副本数" prop="replicas">
                            <el-input-number v-model="createDaemonsetment.replicas" :min="1" :max="10"></el-input-number>
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
                        <el-form-item class="daemonset-create-form" label="镜像" prop="image">
                            <el-input v-model="createDaemonsetment.image"></el-input>
                        </el-form-item>
                        <el-form-item class="daemonset-create-form" label="标签" prop="label_str">
                            <el-input v-model="createDaemonsetment.label_str" placeholder="示例: project=ms,app=gateway"></el-input>
                        </el-form-item>
                        <el-form-item class="daemonset-create-form" label="资源配额" prop="resource">
                            <!-- 常用规格 格式cpu/mem  用/做分隔符 -->
                            <el-select v-model="createDaemonsetment.resource" placeholder="请选择">
                                <el-option value="0.5/1" label="0.5C1G"></el-option>
                                <el-option value="1/2" label="1C2G"></el-option>
                                <el-option value="2/4" label="2C4G"></el-option>
                                <el-option value="4/8" label="4C8G"></el-option>
                                <el-option value="8/16" label="8C16G"></el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item class="daemonset-create-form" label="容器端口" prop="container_port">
                            <el-input v-model="createDaemonsetment.container_port" placeholder="示例: 80"></el-input>
                        </el-form-item>
                        <el-form-item class="daemonset-create-form" label="健康检查" prop="health">
                            <el-switch v-model="createDaemonsetment.health_check" />
                        </el-form-item>
                        <el-form-item class="daemonset-create-form" label="检查路径" prop="healthPath">
                            <el-input v-model="createDaemonsetment.health_path" placeholder="示例: /health"></el-input>
                        </el-form-item>
                        </el-form>
                    </el-col>
                </el-row>
            </template>
            <!-- 抽屉底部 -->
            <template #footer>
                <el-button @click="createDaemonsetmentDrawer = false">取消</el-button>
                <!-- 提交ref匹配的form表单 -->
                <el-button type="primary" @click="submitForm('createDaemonsetment')">立即创建</el-button>
            </template>
        </el-drawer>

        <!-- 资源操作弹出框 -->
        <el-dialog title="YAML信息" v-model="yamlDialog" width="45%" top="2%">
            <codemirror
                :value="contentYaml"
                border
                :options="cmOptions"
                height="500"
                style="font-size:14px;"
                @change="onChange"
            ></codemirror>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="this.yamlDialog = false">取 消</el-button>
                    <el-button type="primary" @click="updateDaemonsetment()">更 新</el-button>
                </span>
            </template>
        </el-dialog>

        <el-dialog title="副本数调整" v-model="scaleDialog" width="25%">
            <div style="text-align:center">
                <span>实例数: </span>
                <el-input-number :step="1" v-model="scaleNum" :min="0" :max="30" label="描述文字"></el-input-number>
            </div>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="scaleDialog = false">取 消</el-button>
                    <el-button type="primary" @click="scaleDaemonsetment()">更 新</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>


<script>
import common from "@/views/common/Config"
import httpClient from "@/utils/request"
// 局部引入
import yaml2obj from 'js-yaml';
import json2yaml from 'json2yaml';
export default { 
    data() {
        return {
            // ns
            namespaceValue:"default",
            namespaceList:[],
            namespaceListUrl: common.k8sNamespaceList,

            // daemonset
            createDaemonsetmentDrawer: false,
            fullscreenloading:false,
            direction:"rtl",

            // 抽屉的默认值
            createDaemonsetment:{
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
            // 创建daemonset
            createDaemonsetmentData:{
                url: common.k8sDaemonsetmentCreate,
                params:{}
            },
            // 表单校验
            createDaemonsetmentRules:{
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
            searchInput:"",

            // daemonset列表
            daemonsetList:[],
            // 表格loading动画
            apploading:false,
            daemonsetTotal:0,
            getDaemonsetmentData:{
                url:common.k8sDaemonsetmentList,
                params:{
                    filter_name:"",
                    namespace:"",
                    page:1,
                    limit:10,
                },
            },
            // 分页
            currentPage:1,
            pagesize:10,
            pagesizeList:[10,20,50],


            // 编辑器配置
            cmOptions: common.cmOptions,
            contentYaml:'',
    


            // 资源操作初始数据
            daemonsetDetail: {},
            getDaemonsetmentDetailData: {
                url: common.k8sDaemonsetmentDetail,
                params: {
                    daemonset_name: '',
                    namespace: ''
                }
            },
            yamlDialog: false,
            updateDaemonsetmentData: {
                url: common.k8sDaemonsetmentUpdate,
                params: {
                    namespace: '',
                    content: ''
                }
            },
            scaleNum: 0,
            scaleDialog: false,
            scaleDaemonsetmentData: {
                url: common.k8sDaemonsetmentScale,
                params: {
                    daemonset_name: '',
                    namespace: '',
                    scale_num: ''
                }
            },
            restartDaemonsetmentData: {
                url: common.k8sDaemonsetmentRestart,
                params: {
                    daemonset_name: '',
                    namespace: '',
                }
            },
            delDaemonsetmentData: {
                url: common.k8sDaemonsetmentDel,
                params: {
                    daemonset_name: '',
                    namespace: '',
                }
            },
            
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

        // 分页操作
        handleSizeChange(size) {
            this.pagesize = size;
            this.getDaemonsetments()
        },
        handleCurrentChange(currentPage) {
            this.currentPage = currentPage;
            this.getDaemonsetments()
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

        // 创建daemonset
        submitForm(formName){
            // 验证表单
            this.$refs[formName].validate((valid)=>{
                if (valid) {
                    this.createDaemonsetFunc()
                }else {
                    return false
                }
            })
        },

        createDaemonsetFunc(){
            // 验证label字段是否合规
            let reg = new RegExp("(^[A-Za-z]+=[A-Za-z0-9]+).*")
            if (!reg.test(this.createDaemonsetment.label_str)) {
                this.$message.warning({
                    message: "标签填写异常，请确认后重新填写"
                })
                return
            }

            // 加载loading动画
            this.fullscreenLoading = true

            // 处理cpu和mem
            let cpu, memory
            let resourceList = this.createDaemonsetment.resource.split("/")

            // 处理label，把字符串转换成map
            let label = new Map()
            let a = (this.createDaemonsetment.label_str).split(",")
            a.forEach(item => {
                let b = item.split("=")
                label[b[0]] = b[1]
            })


            // 赋值表单数据
            cpu = resourceList[0]
            memory = resourceList[1] + "Gi"
            this.createDaemonsetmentData.params = this.createDaemonsetment
            this.createDaemonsetmentData.params.container_port = parseInt(this.createDaemonsetment.container_port)
            this.createDaemonsetmentData.params.label = label
            this.createDaemonsetmentData.params.cpu = cpu
            this.createDaemonsetmentData.params.memory = memory


            // 提交表单
            httpClient.post(this.createDaemonsetmentData.url, this.createDaemonsetmentData.params)
            .then(res => {
                this.$message.success({
                message: res.msg
                })
                this.getDaemonsetments()
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })

            // 后续操作
            this.resetForm('createDaemonsetment')
            this.fullscreenLoading = false
            this.createDaemonsetmentDrawer = false
             
        },

        // 重置表单
        resetForm(formName) {
            this.$refs[formName].resetFields()
        },


        // 获取daemonset列表
        getDaemonsetments() {
            // 填充参数
            this.apploading=true
            this.getDaemonsetmentData.params.filter_name=this.searchInput
            this.getDaemonsetmentData.params.namespace=this.namespaceValue
            this.getDaemonsetmentData.params.page=this.currentPage
            this.getDaemonsetmentData.params.pagesize=this.pagesize 

            // 请求后端
            httpClient.get(this.getDaemonsetmentData.url,{params:this.getDaemonsetmentData.params})
            .then(res => {
    
                this.daemonsetList=res.data.items
                this.daemonsetTotal=res.data.total
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })

            // 请求完成关闭加载动画
            this.apploading=false
        },

        // 字符串截取,截取前15个字符+...
        // app:test,myapp:...
        ellipsis(value) {
            return value.length>15?value.substring(0,15)+'...':value
        },

        // 时间转换UTC到CTS
        timeTrans(timestamp) {
            let date = new Date(new Date(timestamp).getTime() + 8 * 3600 * 1000)
            date = date.toJSON();
            date = date.substring(0, 19).replace('T', ' ')
            return date 
        },


        // 资源操作

        // content转yaml
        transYaml(content) {
            return json2yaml.stringify(content)
        },


        // yaml转content
        transObj(content ) {
            return yaml2obj.load(content)
        },

        // 编辑器内容变化时更新contentYaml
        onChange(val) {
            this.contentYaml=val
        },

        // 查看daemonset详情，e表示行的数据
        getDaemonsetmentDetail(e) {
            this.getDaemonsetmentDetailData.params.daemonset_name = e.row.metadata.name
            this.getDaemonsetmentDetailData.params.namespace = this.namespaceValue
            httpClient.get(this.getDaemonsetmentDetailData.url, {params: this.getDaemonsetmentDetailData.params})
            .then(res => {
                this.daemonsetDetail = res.data
                // 显示到编辑器里
                this.contentYaml = this.transYaml(this.daemonsetDetail)
                // 打开编辑器弹出框
                this.yamlDialog = true
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },

        // 更新daemonset
        updateDaemonsetment() {
            // 把编辑器里的内容转换格式，赋值给content
            let content = JSON.stringify(this.transObj(this.contentYaml))
            this.updateDaemonsetmentData.params.namespace = this.namespaceValue
            this.updateDaemonsetmentData.params.content = content
            // 更新用put操作
            httpClient.put(this.updateDaemonsetmentData.url, this.updateDaemonsetmentData.params)
            .then(res => {
                this.$message.success({
                message: res.msg
                })
                // 更新成功刷新列表
                this.getDaemonsetments()
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            // 关闭弹出框
            this.yamlDialog = false
        },

        // 扩缩容
        handleScale(e) {
            // 打开弹出框
            this.scaleDialog = true
            // 赋值daemonset的详细信息
            this.daemonsetDetail = e.row
            // 把replicas赋值给当前的值进行展示
            this.scaleNum = e.row.spec.replicas
        },

        // 修改daemonset副本数
        scaleDaemonsetment() {
            this.scaleDaemonsetmentData.params.daemonset_name = this.daemonsetDetail.metadata.name
            this.scaleDaemonsetmentData.params.namespace = this.namespaceValue
            // 这里的this.scaleNum是调整以后点击按钮提交的值。
            this.scaleDaemonsetmentData.params.scale_num = this.scaleNum
            httpClient.put(this.scaleDaemonsetmentData.url, this.scaleDaemonsetmentData.params)
            .then(res => {
                this.$message.success({
                message: res.msg
                })
                this.getDaemonsetments()
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
            this.scaleDialog = false
        },

        // 重启daemonset
        restartDaemonsetment(e) {
            this.restartDaemonsetmentData.params.daemonset_name = e.row.metadata.name
            this.restartDaemonsetmentData.params.namespace = this.namespaceValue
            httpClient.put(this.restartDaemonsetmentData.url, this.restartDaemonsetmentData.params)
            .then(res => {
                this.$message.success({
                message: res.msg
                })
                this.getDaemonsetments()
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },

        // 删除daemonset
        delDaemonsetment(e) {
            this.delDaemonsetmentData.params.daemonset_name = e.row.metadata.name
            this.delDaemonsetmentData.params.namespace = this.namespaceValue
            httpClient.delete(this.delDaemonsetmentData.url, {data: this.delDaemonsetmentData.params})
            .then(res => {
                this.$message.success({
                message: res.msg
                })
                this.getDaemonsetments()
            })
            .catch(res => {
                this.$message.error({
                message: res.msg
                })
            })
        },

        // 弹出确认框
        handleConfirm(obj, operateName, fn) {
            this.confirmContent = '确认继续 ' + operateName + ' 操作吗？'

            this.$confirm(this.confirmContent,'提示',{
                confirmButtonText: '确定',
                cancelButtonText: '取消',
            })
            // 点击确认后执行对应的函数
            .then(() => {
                fn(obj)
            })
            .catch(() => {
                this.$message.info({
                    message: '已取消操作'
                })          
            })
        },

    },
    watch: {
        // namespace变更相关
        namespaceValue: {
            handler() {

                // 写入localstorage
                localStorage.setItem('namespace',this.getNamespaceValue)

                // 重置page页数，不重置看不到数据
                this.currentPage=1
                // 获取depoyment列表
                this.getDaemonsetments()
            }
        }
    },
    beforeMount() {
        // 把ns的值同步到其他页面
        if (localStorage.getItem('namespace') !== "undfined" && localStorage.getItem('namespace') !==null) {
            this.namespaceValue=localStorage.getItem('namespace')
        }
        // 获取ns列表
        this.getNamespaceList()
        // 获取depoyment列表
        this.getDaemonsetments()

    },
}


</script>

<style scoped>
    .daemonset-head-card,.daemonset-body-card {
        border-radius: 1px;
        margin-bottom: 5px;
    }
    .daemonset-head-search {
        width:160px;
        margin-right: 10px; 
    }
    .daemonset-body-daemonsetname {
        color: #4795EE;
    }
    .daemonset-body-daemonsetname:hover {
        color: rgb(84, 138, 238);
        cursor: pointer;
        font-weight: bold;
    }
    .dialog-footer {
        margin-right: 10px;
    }
</style>