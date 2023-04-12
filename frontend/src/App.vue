<script>
</script>
<script setup>
import HelloWorld from './components/HelloWorld.vue'
import Main from './components/Main.vue'

import { onMounted, onUnmounted } from 'vue'
import { LogInfo, EventsOn, EventsOff } from '../wailsjs/runtime/runtime'
import crypto from 'crypto'

function msToken(length) {
    const characters =
        'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    const randomBytes = crypto.randomBytes(length)
    return Array.from(
        randomBytes,
        (byte) => characters[byte % characters.length]
    ).join('')
}

onMounted(function () {
    EventsOn('douyin_signer', async function (args) {
        console.log('前端事件运行', args)
        LogInfo('前端事件运行' + args.length)
    })
})
onUnmounted(() => {
    window.removeEventListener('resize', onResize)
    EventsOff('douyin_signer')
})
</script>
<template>
    <Main />
</template>

<style>
</style>
