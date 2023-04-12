<script setup>
import { h, computed, reactive, ref, toRefs, onMounted, onUnmounted } from 'vue'
import {
    Greet,
    DrawData,
    GetCode,
    CheckQrcode,
    LogRedirect
} from '../../wailsjs/go/main/App'
import { LogInfo } from '../../wailsjs/runtime/runtime'
import { ElMessageBox, ElLoading, ElMessage } from 'element-plus'

const userInfo = reactive({
    session_id: '',
    info: {}
})

const formInline = reactive({
    url: '',
    resultText: false
})

const message_list = ref([])
const showInfomation = () => {
    let message_list_dom = []
    for (let value of message_list.value) {
        message_list_dom.push(h('li', null, value))
    }
    ElMessageBox({
        title: '运行状态',
        message:
            message_list_dom.length > 0
                ? h(
                      'ul',
                      {
                          style: {
                              maxHeight: '200px',
                              overflowX: 'hidden',
                              overflowY: 'auto',
                              listStyle: 'none',
                              padding: '10px'
                          }
                      },
                      message_list_dom
                  )
                : ''
    })
}

const table_width = ref(300)
function onResize() {
    table_width.value =
        window.innerHeight -
        document.getElementById('header_form').clientHeight -
        20
}
onMounted(function () {
    window.addEventListener('resize', onResize)
    onResize()
})
onUnmounted(() => {
    closeQrcode()
    window.removeEventListener('resize', onResize)
})

const tableData = reactive([
    {
        date: '2016-05-03',
        name: 'Tom',
        address: 'No. 189, Grove St, Los Angeles'
    }
])

function tableRowClassName({ row, rowIndex }) {
    if (rowIndex % 2 == 0) {
        return 'success-row'
    }
    return ''
}

const is_running = ref(false)
const onSubmit = () => {
    if (is_running.value) {
        ElMessageBox({
            title: '警告',
            message: '系统运行中...'
        })
        return
    }
    is_running.value = true

    DrawData(formInline.url).then((result) => {
        formInline.resultText = result
        console.log(result, '11111111')
        is_running.value = false
    })
}

function greet() {
    Greet(formInline.url).then((result) => {
        formInline.resultText = result
    })
}

const dialogTable = reactive({
    visible: false,
    qrcode: '',
    qrcode_token: '',
    check_interval: null,
    show_info: ''
})
function showQrcode() {
    const loading = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)'
    })

    dialogTable.qrcode = ''
    dialogTable.qrcode_token = ''
    dialogTable.show_info = ''

    GetCode()
        .then((callback_data) => {
            console.log(callback_data)
            loading.close()

            let result = callback_data['json']
            if (result['error'] != null) {
                ElMessage.error('程序异常：' + result['error'])
                return
            }

            let result_json = JSON.parse(result)

            if (
                result_json &&
                (typeof result_json['error'] != 'undefined' ||
                    (typeof result_json['error_code'] != 'undefined' &&
                        result_json['error_code'] != 0))
            ) {
                ElMessage.error('接口请求错误：' + result)
                return
            }
            dialogTable.qrcode =
                `data:image/png;base64,` + result_json['data']['qrcode']
            dialogTable.qrcode_token = result_json['data']['token']
            dialogTable.visible = true
            dialogTable.check_interval = setInterval(function () {
                checkQrcode()
            }, 2000)
        })
        .catch((e) => {
            ElMessage.error(e)
            loading.close()
        })
}
function closeQrcode() {
    clearInterval(dialogTable.check_interval)
}
function checkQrcode() {
    if (dialogTable.qrcode_token == '') {
        ElMessage.error('未获取到Token')
        clearInterval(dialogTable.check_interval)
        return
    }
    CheckQrcode(dialogTable.qrcode_token)
        .then((callback_data) => {
            let result = callback_data['json']
            if (result['error'] != null) {
                ElMessage.error('程序异常：' + result['error'])
                return
            }
            console.log(result)
            console.log(dialogTable.qrcode_token)

            let result_json = JSON.parse(result)
            if (
                result_json &&
                (typeof result_json['error'] != 'undefined' ||
                    (typeof result_json['error_code'] != 'undefined' &&
                        result_json['error_code'] != 0))
            ) {
                ElMessage.error('接口请求错误：' + result)
                clearInterval(dialogTable.check_interval)
                return
            }

            let status = parseInt(result_json['data']['status'])
            switch (status) {
                case 5:
                    dialogTable.qrcode =
                        `data:image/png;base64,` + result_json['data']['qrcode']
                    dialogTable.qrcode_token = result_json['data']['token']
                    console.log('刷新Token')
                    break
                case 3:
                    clearInterval(dialogTable.check_interval)
                    dialogTable.visible = false
                    userInfo.session_id = callback_data['sessionid']

                    const loading = ElLoading.service({
                        lock: true,
                        text: 'Loading',
                        background: 'rgba(0, 0, 0, 0.7)'
                    })
                    LogRedirect(unescape(result_json['data']['redirect_url']))
                        .then((redirect_data) => {
                            loading.close()
                            if (redirect_data['sessionid'] != '') {
                                userInfo.session_id = redirect_data['sessionid']
                                ElMessage('登录成功')
                            } else {
                                userInfo.session_id = ''
                                ElMessage.error('未获取sessionid，登录失败')
                            }
                            console.log(redirect_data)
                        })
                        .catch((e) => {
                            console.log(e)
                            loading.close()
                        })
                    console.log(redirect)
                    break
                case 2:
                    console.log('已扫码')
                    dialogTable.show_info = '已扫码'
                    break
                case 1:
                default:
                    console.log(result_json['data']['status'])
            }
        })
        .catch((e) => {
            ElMessage.error(e)
        })
}
</script>

<template>
    <main>
        <el-dialog
            v-model="dialogTable.visible"
            title="登录"
            @close="closeQrcode"
        >
            <el-image
                style="width: 200px; height: 200px"
                :src="dialogTable.qrcode"
                fit="fill"
            />
            <p style="color: red" v-if="dialogTable.show_info.length > 0">
                {{ dialogTable.show_info }}
            </p>
        </el-dialog>

        <el-affix :offset="0" id="header_form">
            <el-form
                :inline="true"
                :model="formInline"
                class="demo-form-inline"
            >
                <el-form-item label="网址" style="width: 50%">
                    <el-input v-model="formInline.url" placeholder="URL地址" />
                </el-form-item>
                <el-form-item>
                    <el-button
                        type="primary"
                        @click="onSubmit"
                        :disabled="is_running"
                        >{{ !is_running ? '抓取' : '抓取中' }}</el-button
                    >
                    <el-button type="warning" @click="showInfomation"
                        >运行状态</el-button
                    >
                    <el-button
                        type="info"
                        v-if="userInfo.session_id != ''"
                        @click="showQrcode"
                        >登录</el-button
                    >
                </el-form-item>
            </el-form>
        </el-affix>

        <el-table
            :data="tableData"
            style="width: 100%"
            :height="table_width"
            :row-class-name="tableRowClassName"
        >
            <el-table-column prop="date" label="Date" width="180" />
            <el-table-column prop="name" label="Name" width="180" />
            <el-table-column prop="address" label="Address" />
        </el-table>
    </main>
</template>

<style >
.el-table .warning-row {
    --el-table-tr-bg-color: var(--el-color-warning-light-9);
}
.el-table .success-row {
    --el-table-tr-bg-color: var(--el-color-success-light-9);
}
</style>
