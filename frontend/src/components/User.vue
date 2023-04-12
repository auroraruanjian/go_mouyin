<script setup>
import { reactive, ref, onMounted, onUnmounted, inject, defineEmits } from 'vue'
import {
    GetRelation
} from '../../wailsjs/go/main/App'
import { LogInfo } from '../../wailsjs/runtime/runtime'
import { ElMessageBox, ElLoading, ElMessage } from 'element-plus'

const emit = defineEmits(['setUserInfo'])

const userInfo = inject('userInfo')
const followings = ref([])
const myInfo = reactive({})

function apiGetRelation() {
    const loading = ElLoading.service({
        lock: true,
        text: 'Loading',
        background: 'rgba(0, 0, 0, 0.7)'
    })
    GetRelation(userInfo.session_id).then((callback_data) => {
        loading.close()

        console.log("result:", callback_data)
        let relation_json = JSON.parse(callback_data['json'])
        //console.log(relation_json)

        if (relation_json['status_code'] == 0) {
            //followings = relation_json['followings']
            for (let key in relation_json['followings']) {
                //console.log(relation_json['followings'][key]['uid'], userInfo['login_info']['user_uid'])
                if (relation_json['followings'][key]['uid'] == userInfo['login_info']['user_uid']) {
                    myInfo["uid"] = relation_json['followings'][key]['uid']
                    myInfo["avatar_small"] = relation_json['followings'][key]['avatar_small']['url_list'][0]
                    myInfo["custom_verify"] = relation_json['followings'][key]['custom_verify']
                    myInfo["nickname"] = relation_json['followings'][key]['nickname']
                    myInfo["sec_uid"] = relation_json['followings'][key]['sec_uid']
                    myInfo["signature"] = relation_json['followings'][key]['signature']
                    myInfo["short_id"] = relation_json['followings'][key]['short_id']
                    myInfo["unique_id"] = relation_json['followings'][key]['unique_id']
                    myInfo["signature"] = relation_json['followings'][key]['signature']
                } else {
                    followings.value.push({
                        "uid": relation_json['followings'][key]['uid'],
                        "avatar_small": relation_json['followings'][key]['avatar_small']['url_list'][0],
                        "custom_verify": relation_json['followings'][key]['custom_verify'],
                        "nickname": relation_json['followings'][key]['nickname'],
                        "sec_uid": relation_json['followings'][key]['sec_uid'],
                        "signature": relation_json['followings'][key]['signature'],
                        "short_id": relation_json['followings'][key]['short_id'],
                        "unique_id": relation_json['followings'][key]['unique_id'],
                        "signature": relation_json['followings'][key]['signature'],
                    })
                }
            }
            emit("setUserInfo", 'info', myInfo)
            emit("setUserInfo", 'follow', followings)
        } else {
            if (relation_json['status_code'] == 0) {
                ElMessage.error(callback_data)
            }
        }
    }).catch(e => {
        console.error(e)
        loading.close()
    })
}

const boxHeight = ref(500)
function onResize() {
    boxHeight.value = window.innerHeight - 20
}
onMounted(function () {
    window.addEventListener('resize', onResize)
    onResize()
    apiGetRelation()
})
onUnmounted(() => {
    window.removeEventListener('resize', onResize)
})
</script>

<template>
    <div style="overflow-y: scroll;overflow-x: hidden;padding:2px" :style="{ height: boxHeight + 'px' }">
        <el-card :body-style="{ padding: '10px' }" v-if="myInfo['uid'] != null">
            <el-row :gutter="5">
                <el-col :span="4">
                    <img :src="myInfo.avatar_small" class="image" width="100" height="100" />
                </el-col>
                <el-col :span="20" style="text-align: left;">
                    <div style="padding: 3px; font-size: 12px;">
                        <div>昵称：{{ myInfo.nickname }}</div>
                        <div>个性签名：{{ myInfo.signature }}</div>
                    </div>
                </el-col>
            </el-row>
        </el-card>
        <el-divider />
        <el-row :gutter="5">
            <el-col v-for="(item, key) in followings" :key="key" :span="4" style="margin-bottom:5px">
                <el-card :body-style="{ padding: '10px' }">
                    <img :src="item.avatar_small" class="image" width="100" height="100" />
                    <div style="padding: 3px; font-size: 12px; height: 30px; overflow: hidden;display: block;">
                        <div>{{ item.nickname }}</div>
                        <div>{{ item.signature }}</div>
                    </div>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<style></style>