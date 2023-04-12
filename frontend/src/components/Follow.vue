<script setup>
import { h, computed, reactive, ref, toRefs, onMounted, onUnmounted, inject } from 'vue'
import {
    GetFollow,
    GetAweme
} from '../../wailsjs/go/main/App'
import { LogInfo } from '../../wailsjs/runtime/runtime'
import { ElMessageBox, ElLoading, ElMessage } from 'element-plus'

const userInfo = inject('userInfo')

const multipleTable = ref()
const tableData = reactive({
    data: []
})
function tableRowClassName({ row, rowIndex }) {
    if (rowIndex % 2 == 0) {
        return 'success-row'
    }
    return ''
}
const handleRowClick = (row, column, event) => {
    multipleTable.value.toggleRowSelection(row);
}
const handleSelectionChange = (val) => {
    multipleTable.value = val
}
const toggleAll = (choose) => {
    if (choose) {
        multipleTable.value.toggleAllSelection()
    } else {
        multipleTable.value.clearSelection()
    }
}
/*
const toggleSelection = (rows) => {
    if (rows) {
        rows.forEach((row) => {
            // TODO: improvement typing when refactor table
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-expect-error
            multipleTable.value.toggleRowSelection(row, undefined)
        })
    } else {
        multipleTableRef.value.clearSelection()
    }
}
*/
const drawAweme = () => {
    let rows = multipleTable.value.getSelectionRows()
    for (let x in rows) {
        GetAweme(rows[x]['Sec_uid'], "0").then((callback_data) => {
            console.log("result:", callback_data)
            let ameme_list = JSON.parse(callback_data['ameme_list'])
            console.debug(ameme_list)
        }).catch(e => {
            console.error(e)
        })
    }
}

const is_running = ref(false)
const searchForm = reactive({
    max_time: "0",
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

    GetFollow(userInfo.session_id, searchForm.sec_uid, searchForm.max_time).then((callback_data) => {
        loading.close()
        is_running.value = false

        console.log("result:", callback_data)
        let user_item = JSON.parse(callback_data['user_item'])
        console.debug(user_item)
        tableData.data = user_item
    }).catch(e => {
        console.error(e)
        loading.close()
    })
}

const tableHeight = ref(300)
function onResize() {
    tableHeight.value = window.innerHeight - 130
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

    <el-table :data="tableData.data" ref="multipleTable" :style="{ 'width': '100%', 'height': tableHeight + 'px' }"
        @row-click="handleRowClick" :row-class-name="tableRowClassName">
        <el-table-column type="selection" width="55" />
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
    <div style="margin-top: 20px;text-align: left;">
        <el-button @click="toggleAll(true)">全选</el-button>
        <el-button @click="toggleAll(false)">反选</el-button>
        <el-button type="danger" @click="drawAweme">抓取</el-button>
    </div>
</template>

<style>
.el-table .warning-row {
    --el-table-tr-bg-color: var(--el-color-warning-light-9);
}

.el-table .success-row {
    --el-table-tr-bg-color: var(--el-color-success-light-9);
}
</style>
