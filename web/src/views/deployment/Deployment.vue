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
        <!-- header2 用来创建资源 -->
        <el-col :span="24">
            <div>
                456
            </div>
        </el-col>
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
            namespaceValue:"default",
            namespaceList:[],
            namespaceListUrl: common.k8sNamespaceList,
            
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
        }
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

</style>