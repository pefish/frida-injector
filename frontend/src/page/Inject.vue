<script lang="ts" setup>
import { reactive, Ref, ref } from 'vue'
import MonacoEditor from 'monaco-editor-vue3';
import { main } from "../../wailsjs/go/models.js"
import { useRoute, useRouter } from 'vue-router';
import { InjectScript, CancelScript, DetachProcess } from "../../wailsjs/go/main/App.js"
import { EventsOn, EventsOff } from '../../wailsjs/runtime'
import { ElMessage } from 'element-plus';


const router = useRouter()
const route = useRoute()

console.log("route.fullPath", route.fullPath, route.params, route.query)

const logTextareaDom: Ref<any> = ref(null)
const script = ref(`
setInterval(
    () => {
        console.log("This is script.")
    },
    1000,
)
`)
const logTextarea = ref(``)
const isInjected = ref(false)

function appendLog(level: string, msg: string) {
    logTextarea.value += `<${level.toUpperCase()}>\t${msg}\n`
    logTextareaDom.value.textarea.scrollTop = logTextareaDom.value.textarea.scrollHeight;
}

async function injectOnClick() {
    await InjectScript(script.value)
    EventsOn('log-parse-error', (data: any) => {
        appendLog("log-parse-error", data)
    })
    EventsOn('error', (data: any) => {
        appendLog(data.level, data.payload)
    })
    EventsOn('log', (data: any) => {
        appendLog(data.level, data.payload)
    })
    isInjected.value = true
    // ElMessage({
    //     message: "Injected.",
    //     type: "success",
    // })
    appendLog("INFO", "Injected.")
}

async function cancelOnClick() {
    await CancelScript()
    EventsOff('log-parse-error')
    EventsOff('error')
    EventsOff('log')
    isInjected.value = false
    // ElMessage({
    //     message: "Canceled.",
    //     type: "success",
    // })
    appendLog("INFO", "Canceled.")
}

async function detachOnClick() {
    await CancelScript()
    await DetachProcess()
    router.back()
}

</script>

<template>
    <div class="flex flex-col justify-center h-full px-4 py-2">
        <div class="px-20 py-6 flex-1">
            <el-form label-width="auto">
                <el-form-item label="Script">
                    <MonacoEditor theme="vs" :options="{
                        dragAndDrop: true,
                    }" language="javascript" :height="300" :diffEditor="false" v-model:value="script"></MonacoEditor>
                </el-form-item>
                <el-form-item>
                    <el-button class="ml-2" :disabled="isInjected" type="primary"
                        @click="injectOnClick">Inject</el-button>
                    <el-button class="ml-2" :disabled="!isInjected" type="primary"
                        @click="cancelOnClick">Cancel</el-button>
                </el-form-item>
            </el-form>
            <el-input ref="logTextareaDom" v-model="logTextarea" :autosize="{ minRows: 14, maxRows: 14 }"
                type="textarea" readonly />
        </div>

        <div class="flex items-center">
            <span class="text-red-400">Attached to <{{ route.query.processName }}></span>
            <el-button class="ml-2" type="primary" @click="detachOnClick">Detach</el-button>
        </div>
    </div>
</template>

<style></style>
