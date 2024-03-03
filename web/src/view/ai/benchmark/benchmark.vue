<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
            type="primary"
            icon="plus"
            @click="openDialog"
        >新增
        </el-button>
      </div>
      <el-table
          ref="multipleTable"
          :data="tableData"
          style="width: 100%"
          tooltip-effect="dark"
          row-key="ID"
      >
        <el-table-column
            type="selection"
            width="55"
        />
        <el-table-column
            align="left"
            label="公众号名称"
            prop="accountName"
            width="300"
        >
        </el-table-column>
        <el-table-column
            align="left"
            label="公众号ID"
            prop="accountId"
            width="100"
        />
        <el-table-column
            align="left"
            label="主题"
            prop="topic"
            width="80"
        />
        <el-table-column
            align="left"
            label="创建时间"
            width="180"
        >
          <template #default="scope">
            <span>{{ formatDate(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="操作"
            min-width="160"
        >
          <template #default="scope">
            <el-popover
                v-model="scope.row.visible"
                placement="top"
                width="160"
            >
              <p>确定要删除吗？</p>
              <div style="text-align: right; margin-top: 8px;">
                <el-button
                    type="primary"
                    link
                    @click="scope.row.visible = false"
                >取消
                </el-button>
                <el-button
                    type="primary"
                    @click="deleteWechatBenchmarkAccount(scope.row)"
                >确定
                </el-button>
              </div>
              <template #reference>
                <el-button
                    type="primary"
                    link
                    icon="delete"
                    @click="scope.row.visible = true"
                >删除
                </el-button>
              </template>
            </el-popover>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
        />
      </div>
    </div>


    <el-dialog
        v-model="dialogFormVisible"
        :before-close="closeDialog"
        title="对标账号"
    >
      <el-scrollbar height="500px">
        <el-form
            :model="form"
            label-position="right"
            label-width="90px"
        >
          <el-form-item label="微信公众号全称">
            <el-input
                v-model="form.accountName"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="微信公众号Id">
            <el-input
                v-model="form.accountId"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="主题">
            <el-input
                v-model="form.topic"
                autocomplete="off"
            />
          </el-form-item>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button
              type="primary"
              @click="enterDialog"
          >确 定
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {createBenchmark, deleteBenchmarkAccount, getBenchmarkAccountList} from '@/api/benchmarkAccount'
import {ref} from 'vue'
import {ElMessage} from 'element-plus'
import {formatDate} from '@/utils/format'
import {createPortal, updatePortal} from "@/api/portal";

defineOptions({
  name: 'Benchmark'
})


const form = ref({
  accountName: '',
  accountId: '',
  topic: '',
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})


// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async () => {
  const table = await getBenchmarkAccountList({page: page.value, pageSize: pageSize.value, ...searchInfo.value,})
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

const dialogFormVisible = ref(false)
const type = ref('')


const deleteWechatBenchmarkAccount = async (row) => {
  row.visible = false
  const res = await deleteBenchmarkAccount({ID: row.ID})
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}


const enterDialog = async () => {
  let res
  switch (type.value) {
    case 'create':
      res = await createBenchmark(form.value)
      break
      // case 'update':
      //   res = await updatePortal(form.value)
      //   break
    default:
      res = await createBenchmark(form.value)
      break
  }

  if (res.code === 0) {
    closeDialog()
    getTableData()
  }
}


const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}


const closeDialog = () => {
  dialogFormVisible.value = false
  form.value = {
    accountName: '',
    accountId: '',
    topic: '',
  }
}

</script>


<style></style>
