<script setup>
import { h, computed, reactive, ref, toRefs, onMounted, onUnmounted } from 'vue'
import {
    SearchUser
} from '../../wailsjs/go/main/App'
import { LogInfo } from '../../wailsjs/runtime/runtime'
import { ElMessageBox, ElLoading, ElMessage } from 'element-plus'

const tableData = reactive({
    data: []
})
function tableRowClassName({ row, rowIndex }) {
    if (rowIndex % 2 == 0) {
        return 'success-row'
    }
    return ''
}

const is_running = ref(false)
const searchForm = reactive({
    keyword: ""
})
const onSubmit = () => {
    if (is_running.value) {
        ElMessageBox({
            title: '警告',
            message: '系统运行中...'
        })
        return
    }
    const loading = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)'
    })
    is_running.value = true

    SearchUser(searchForm.keyword).then((callback_data) => {
        loading.close()
        is_running.value = false

        console.log("result:", callback_data)
        let user_item = JSON.parse(callback_data['user_item'])
        tableData.data = user_item
    }).catch(e => {
        console.error(e)
        loading.close()
    })
}

const tableHeight = ref(300)
function onResize() {
    tableHeight.value = window.innerHeight - 70
    console.log(window.innerHeight - 20)
}
onMounted(function () {
    window.addEventListener('resize', onResize)
    onResize()
})
onUnmounted(() => {
    window.removeEventListener('resize', onResize)
})
</script>

<template>
    <el-affix :offset="0" id="header_form">
        <el-form :inline="true" :model="searchForm" class="demo-form-inline">
            <el-form-item label="关键词">
                <el-input v-model="searchForm.keyword" placeholder="请输入要搜索的关键词" />
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="onSubmit" :disabled="is_running">{{ !is_running ? '搜索' : '搜索中'
                }}</el-button>
            </el-form-item>
        </el-form>
    </el-affix>

    <el-table :data="tableData.data" :style="{ 'width': '100%', 'height': tableHeight + 'px' }"
        :row-class-name="tableRowClassName">
        <el-table-column prop="Nickname" label="昵称" width="120" />
        <el-table-column prop="Unique_id" label="抖音号" width="100">
            <template #default="scope">
                {{ scope.row.Unique_id != "" ? scope.row.Unique_id : scope.row.Short_id }}
            </template>
        </el-table-column>
        <el-table-column prop="Follower_count" label="粉丝数" width="100" />
        <el-table-column prop="Custom_verify" label="认证信息" width="120" />
        <el-table-column prop="Signature" label="个性签名" />
    </el-table>
</template>

<style>
.el-table .warning-row {
    --el-table-tr-bg-color: var(--el-color-warning-light-9);
}

.el-table .success-row {
    --el-table-tr-bg-color: var(--el-color-success-light-9);
}
</style>
