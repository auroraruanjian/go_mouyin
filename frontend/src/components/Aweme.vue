<script setup>
import { h, computed, reactive, ref, toRefs, onMounted, onUnmounted, inject } from 'vue'
import {
    GetAweme
} from '../../wailsjs/go/main/App'
import { LogInfo, BrowserOpenURL } from '../../wailsjs/runtime/runtime'
import { ElMessageBox, ElLoading, ElMessage } from 'element-plus'

const userInfo = inject('userInfo')

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
    max_cursor: "0",
    sec_uid: ""
})
const showSelect = computed(() => {
    let follow = userInfo['follow'].concat([])
    follow.unshift(userInfo['info'])
    return follow
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

    GetAweme(searchForm.sec_uid, searchForm.max_cursor).then((callback_data) => {
        loading.close()
        setTimeout(function () {
            is_running.value = false
        }, 300)

        console.log("result:", callback_data)
        let ameme_list = JSON.parse(callback_data['ameme_list'])
        console.debug(ameme_list)
        tableData.data = ameme_list
    }).catch(e => {
        console.error(e)
        loading.close()
    })
}
function openBrowser(url) {
    BrowserOpenURL(url)
}

const tableHeight = ref(300)
function onResize() {
    tableHeight.value = window.innerHeight - 70
    console.log(window.innerHeight - 20)
}
onMounted(function () {
    searchForm.sec_uid = userInfo.info.sec_uid
    onSubmit()

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
            <el-form-item label="用户">
                <el-select v-model="searchForm.sec_uid" class="m-2" placeholder="">
                    <el-option v-for="item in showSelect" :key="item.sec_uid" :label="item.nickname"
                        :value="item.sec_uid" />
                </el-select>
            </el-form-item>
            <el-form-item>
                <el-button type="warning" @click="searchForm.sec_uid = userInfo.info.sec_uid">重置</el-button>
                <el-button type="primary" @click="onSubmit" :disabled="is_running">{{ !is_running ? '查询' :
                    '查询中' }}</el-button>
            </el-form-item>
        </el-form>
    </el-affix>

    <el-table :data="tableData.data" :style="{ 'width': '100%', 'height': tableHeight + 'px' }"
        :row-class-name="tableRowClassName">
        <el-table-column prop="Aweme_id" label="ID" width="120" />
        <el-table-column prop="Ip_attribution" label="IP地址" width="100" />
        <el-table-column prop="Address_info" label="地址" width="100"></el-table-column>
        <el-table-column prop="Desc" label="描述">
            <template #default="scope">
                <a style="color: rgb(90, 122, 208); cursor: pointer;" @click="openBrowser(scope.row.Share_url)">{{
                    scope.row.Desc }}</a>
            </template>
        </el-table-column>
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
