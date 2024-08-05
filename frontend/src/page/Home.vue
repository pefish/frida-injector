<script lang="ts" setup>
import { computed, reactive, Ref, ref } from 'vue'
import { ListDevices, ListProcessesOfDevice, AttachProcess } from "../../wailsjs/go/main/App.js"
import { main } from "../../wailsjs/go/models.js"
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';

const router = useRouter()
const dialogVisible: Ref<boolean> = ref(false)
const devices: Ref<string[]> = ref([])
const form = reactive({
    selectedDevice: '',
})
const searchProcessName = ref("")
const processes: Ref<main.ProcessInfo[]> = ref([])
const selectedDeviceName = ref("")
const selectedProcess: Ref<main.ProcessInfo | null> = ref(null)

const filterProcesses: Ref<main.ProcessInfo[]> = computed(() =>
    processes.value.filter(
        (data) =>
            !searchProcessName.value ||
            data.name.toLowerCase().includes(searchProcessName.value.toLowerCase())
    )
)

async function attachOnClick() {
    devices.value = await ListDevices()
    dialogVisible.value = true
    if (devices.value.length > 0) {
        form.selectedDevice = devices.value[0]
        await deviceSelectOnChange(form.selectedDevice)
    }
}

async function deviceSelectOnChange(deviceName: string) {
    processes.value = []
    processes.value = await ListProcessesOfDevice(deviceName)
    selectedDeviceName.value = deviceName
}

async function processRowOnClick(row: main.ProcessInfo) {
    selectedProcess.value = row
}

function tableRowClassName(
    {
        row,
        rowIndex,
    }: {
        row: main.ProcessInfo
        rowIndex: number
    }
): string {
    if (selectedProcess.value !== null && row.id === selectedProcess.value.id) {
        return "success-row"
    }
    return ""
}

async function attachConfirmOnClick() {
    if (!selectedDeviceName.value) {
        ElMessage({
            message: "Select device first.",
            type: "error",
        })
        return
    }
    if (selectedProcess.value == null) {
        ElMessage({
            message: "Select process first.",
            type: "error",
        })
        return
    }
    console.log("selectedDeviceName.value", selectedDeviceName.value)
    await AttachProcess(selectedDeviceName.value, selectedProcess.value.id)
    dialogVisible.value = false
    router.push({
        name: "inject",
        query: {
            processName: selectedProcess.value.name,
            processPath: selectedProcess.value.path,
        }
    })
}

</script>

<template>
    <div class="flex justify-center items-center px-20 py-6 h-full">
        <el-button class="ml-2" type="primary" @click="attachOnClick">Attach</el-button>
    </div>
    <el-dialog v-model="dialogVisible" title="Attach" width="900">
        <el-form :model="form" label-width="auto">
            <el-form-item label="Device">
                <el-select v-model="form.selectedDevice" @change="deviceSelectOnChange">
                    <el-option v-for="d in devices" :label="d" :value="d" />
                </el-select>
            </el-form-item>
            <el-form-item label="Process">
                <el-table @row-click="processRowOnClick" :row-class-name="tableRowClassName" :data="filterProcesses"
                    :default-sort="{ prop: 'name', order: 'ascending' }" style="width: 100%" max-height="350">
                    <el-table-column fixed width="250">
                        <template #header>
                            <div class="flex">
                                <span>Name</span>
                                <el-input v-model="searchProcessName" size="small" placeholder="Search"
                                    class="w-[150px] ml-4" />
                            </div>
                        </template>
                        <template #default="{ row }">
                            <span>{{ row.name }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column fixed prop="id" label="PID" width="100" />
                    <el-table-column prop="parent_id" label="PPID" width="100" />
                    <el-table-column sortable prop="user" label="User" width="100" />
                    <el-table-column prop="path" label="Path" width="300" />
                </el-table>
            </el-form-item>
        </el-form>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="dialogVisible = false">Cancel</el-button>
                <el-button type="primary" @click="attachConfirmOnClick">
                    Confirm
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>

<style>
.el-table .cell {
    overflow-wrap: anywhere;
}

.el-table__row {
    cursor: pointer;
}

.el-table .success-row {
    background: darkgray;
}
</style>
