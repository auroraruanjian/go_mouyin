<script setup>
import { h, computed, reactive, ref, toRefs, onMounted, onUnmounted, provide } from 'vue'
import {
    Greet,
    DrawData,
    GetCode,
    CheckQrcode,
    LogRedirect,
    GetQueryUser
} from '../../wailsjs/go/main/App'
import { LogInfo } from '../../wailsjs/runtime/runtime'
import { ElMessageBox, ElLoading, ElMessage } from 'element-plus'
import {
    Document,
    Menu as IconMenu,
    Location,
    Setting,
    Search
} from '@element-plus/icons-vue'
import SearchPage from "./Search.vue"
import User from "./User.vue"
import Follow from "./Follow.vue"
import Aweme from "./Aweme.vue"

const menu = reactive({
    currentMenu: "search"
})
const handleOpen = (key, keyPath) => {
    console.log(key, keyPath)
}
const handleClose = (key, keyPath) => {
    console.log(key, keyPath)
}
const handleCheckMenu = (key, keyPath) => {
    if (menu.currentMenu == key) return
    if (key == "login") {
        showQrcode()
    } else {
        menu.currentMenu = key
    }
    console.log(key, keyPath)
}

const userInfo = reactive({
    session_id: '',
    login_info: {},
    info: {},
    follow: []
})
function setUserInfo(key, value) {
    console.log(key, value)
    userInfo[key] = value
    localStorage.setItem('userInfo', JSON.stringify(userInfo))
}
function resetUserInfo() {
    userInfo = {
        session_id: '',
        login_info: {},
        info: {}
    }
    localStorage.removeItem("userInfo")
}
provide("userInfo", userInfo)

const dialogLogin = reactive({
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

    dialogLogin.qrcode = ''
    dialogLogin.qrcode_token = ''
    dialogLogin.show_info = ''

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
            dialogLogin.qrcode =
                `data:image/png;base64,` + result_json['data']['qrcode']
            dialogLogin.qrcode_token = result_json['data']['token']
            dialogLogin.visible = true

            if (dialogLogin.check_interval != null) clearInterval(dialogLogin.check_interval)
            dialogLogin.check_interval = setInterval(function () {
                checkQrcode()
            }, 2000)
        })
        .catch((e) => {
            ElMessage.error("GetCode:" + e)
            console.error(e)
            loading.close()
        })
}
function closeQrcode() {
    clearInterval(dialogLogin.check_interval)
    console.log("弹出框关闭")
}
function checkQrcode() {
    if (dialogLogin.qrcode_token == '') {
        ElMessage.error('未获取到Token')
        clearInterval(dialogLogin.check_interval)
        return
    }
    CheckQrcode(dialogLogin.qrcode_token)
        .then((callback_data) => {
            let result = callback_data['json']
            if (result['error'] != null) {
                ElMessage.error('程序异常：' + result['error'])
                return
            }
            console.log(result)
            console.log(dialogLogin.qrcode_token)

            let result_json = JSON.parse(result)
            if (
                result_json &&
                (typeof result_json['error'] != 'undefined' ||
                    (typeof result_json['error_code'] != 'undefined' &&
                        result_json['error_code'] != 0))
            ) {
                ElMessage.error('接口请求错误：' + result)
                clearInterval(dialogLogin.check_interval)
                return
            }

            let status = parseInt(result_json['data']['status'])
            switch (status) {
                case 5:
                    dialogLogin.qrcode =
                        `data:image/png;base64,` + result_json['data']['qrcode']
                    dialogLogin.qrcode_token = result_json['data']['token']
                    console.log('刷新Token')
                    break
                case 3:
                    clearInterval(dialogLogin.check_interval)
                    dialogLogin.visible = false
                    userInfo.session_id = callback_data['sessionid']

                    const loading = ElLoading.service({
                        lock: true,
                        text: 'Loading',
                        background: 'rgba(0, 0, 0, 0.7)'
                    })
                    LogRedirect(unescape(result_json['data']['redirect_url']))
                        .then((redirect_data) => {
                            loading.close()
                            console.log(redirect_data)
                            switch (parseInt(redirect_data['status'])) {
                                // 1:获取sessionid成功，但获取用户信息失败
                                case 1:
                                // 2:获取用户信息成功，存储用户信息
                                case 2:
                                    setUserInfo("session_id", redirect_data['sessionid'])

                                    let redirect_json = {}
                                    try {
                                        redirect_json = JSON.parse(redirect_data['json'])
                                    } catch (e) { }

                                    if (parseInt(redirect_data['status']) == 1 ||
                                        typeof redirect_json['user_uid'] == "undefined" ||
                                        redirect_json['user_uid'] == null
                                    ) {
                                        // 数据解析错误，调用 GetQueryUser 函数,重试
                                        ElMessageBox.confirm(
                                            "状态码：" + redirect_data['status'] + ",返回：" + redirect_data['json'],
                                            'Warning',//error
                                            {
                                                confirmButtonText: '重试',
                                                cancelButtonText: '取消',
                                                type: 'warning',
                                            }
                                        ).then(() => {
                                            // 调用 GetQueryUser 函数,重试
                                            apiGetQueryUser()
                                        }).catch(() => {
                                        })
                                    } else {
                                        setUserInfo("login_info", redirect_json)
                                        ElMessage('登录成功')
                                    }
                                    break;
                                // 0:获取session id失败
                                default:
                                    resetUserInfo()
                                    ElMessage.error('未获取sessionid，登录失败:', redirect_data['error'])
                            }
                        })
                        .catch((e) => {
                            console.log(e)
                            loading.close()
                        })
                    console.log(redirect)
                    break
                case 2:
                    console.log('已扫码')
                    dialogLogin.show_info = '已扫码'
                    break
                case 1:
                default:
                    console.log(result_json['data']['status'])
            }
        })
        .catch((e) => {
            ElMessage.error("CheckQrcode:" + e)
            console.error(e)
        })
}
function apiGetQueryUser() {
    const loading = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)'
    })

    GetQueryUser(userInfo.session_id).then(res => {
        loading.close()
        console.log(res)
        if (typeof res['error'] != "undefined" && res['error'] != null) {
            ElMessage.error('获取用户信息，失败:', redirect_data['error'])
            return
        }

        let redirect_json = {}
        try {
            redirect_json = JSON.parse(res['json'])
        } catch (e) { }

        if (typeof redirect_json['user_uid'] == "undefined" ||
            redirect_json['user_uid'] == null) {
            // 数据解析错误，调用 GetQueryUser 函数,重试
            ElMessageBox.confirm(
                "状态码：" + redirect_data['status'] + ",返回：" + redirect_data['json'],
                'Warning',//error
                {
                    confirmButtonText: '重试',
                    cancelButtonText: '取消',
                    type: 'warning',
                }
            ).then(() => {
                // 调用 GetQueryUser 函数,重试
                apiGetQueryUser()
            }).catch(() => {
                resetUserInfo()
            })
        } else {
            setUserInfo("login_info", redirect_json)
            ElMessage('登录成功')
        }
    }).catch(e => {
        loading.close()
        ElMessage.error("apiGetQueryUser:" + e)
        console.error(e)
    })
}

const menu_height = ref(500)
function onResize() {
    menu_height.value = window.innerHeight - 20
}
onMounted(function () {
    window.addEventListener('resize', onResize)
    onResize()

    try {
        let local_userInfo = JSON.parse(localStorage.getItem("userInfo"))
        if (local_userInfo != null) {
            userInfo.session_id = local_userInfo.session_id
            userInfo.login_info = local_userInfo.login_info
            userInfo.info = local_userInfo.info
            userInfo.follow = local_userInfo.follow
        }
    } catch (e) {
        console.log(e)
    }
})
onUnmounted(() => {
    closeQrcode()
    window.removeEventListener('resize', onResize)
})
</script>

<template>
    <el-row class="tac" :gutter="10">
        <el-col :span="5">
            <el-menu :style="{ height: menu_height + 'px' }" default-active="search" class="el-menu-vertical-demo"
                @open="handleOpen" @close="handleClose" @select="handleCheckMenu">
                <el-menu-item index="search">
                    <el-icon>
                        <location />
                    </el-icon>
                    <span>搜索用户</span>
                </el-menu-item>
                <el-menu-item index="user" :disabled="userInfo.session_id == ''">
                    <el-icon><icon-menu /></el-icon>
                    <span>我的</span>
                </el-menu-item>
                <el-menu-item index="follow" :disabled="userInfo.session_id == ''">
                    <el-icon>
                        <document />
                    </el-icon>
                    <span>关注用户</span>
                </el-menu-item>
                <el-menu-item index="aweme" :disabled="userInfo.session_id == ''">
                    <el-icon>
                        <setting />
                    </el-icon>
                    <span>关注用户视频</span>
                </el-menu-item>
                <el-menu-item index="log">
                    <el-icon>
                        <setting />
                    </el-icon>
                    <span>运行日志</span>
                </el-menu-item>
                <el-menu-item index="login" v-if="userInfo.session_id == ''">
                    <el-icon>
                        <setting />
                    </el-icon>
                    <span>登录</span>
                </el-menu-item>
            </el-menu>
        </el-col>
        <el-col :span="19">
            <SearchPage v-if="menu.currentMenu == 'search'"></SearchPage>
            <User v-if="menu.currentMenu == 'user'" @setUserInfo="setUserInfo"></User>
            <Follow v-if="menu.currentMenu == 'follow'"></Follow>
            <Aweme v-if="menu.currentMenu == 'aweme'"></Aweme>
        </el-col>
    </el-row>

    <el-dialog v-model="dialogLogin.visible" title="登录" @close="closeQrcode">
        <el-image style="width: 200px; height: 200px" :src="dialogLogin.qrcode" fit="fill" />
        <p style="color: red" v-if="dialogLogin.show_info.length > 0">
            {{ dialogLogin.show_info }}
        </p>
    </el-dialog>
</template>

<style scoped></style>