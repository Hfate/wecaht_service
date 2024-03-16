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
        <el-button
            type="primary"
            icon="plus"
            @click="openWxDialog"
        >更新token
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
            label="初始爬取文章数量"
            prop="initNum"
            width="100"
        />
        <el-table-column
            align="left"
            label="已爬取文章数量"
            prop="finishNum"
            width="100"
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
            <el-button
                type="primary"
                link
                icon="edit"
                @click="updateWechatBenchmark(scope.row)"
            >变更
            </el-button>
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
      <el-scrollbar height="600px">
        <el-form
            :model="form"
            label-position="right"
            label-width="180px"
        >
          <el-form-item label="微信公众号全称">
            <el-input
                v-model="form.accountName"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="主题">
            <el-select
                v-model="form.topic"
                class="w-56"
            >
              <el-option
                  v-for="item in topicArr"
                  :value="item"
                  :label="item"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="该公众号任意微信文章链接">
            <el-input
                v-model="form.articleLink"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="初始爬取文章数量">
            <el-input
                v-model.number="form.initNum"
                autocomplete="off"
                placeholder="建议100篇以内"
            />
          </el-form-item>
          <el-form-item label="公众号key值">
            <el-input
                v-model="form.key"
                autocomplete="off"
                placeholder="字段暂时无用可以随意填"
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


    <el-dialog
        v-model="dialogWxFormVisible"
        :before-close="closeWxDialog"
        title="微信公众号token"
    >
      <el-scrollbar height="500px">
        <el-form
            :model="wxForm"
            label-position="right"
            label-width="90px"
        >
          <el-form-item label="token">
            <el-input
                v-model="wxForm.token"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="slaveSid">
            <el-input
                v-model="wxForm.slaveSid"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="bizUin">
            <el-input
                v-model="wxForm.bizUin"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="dataTicket">
            <el-input
                v-model="wxForm.dataTicket"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="randInfo">
            <el-input
                v-model="wxForm.randInfo"
                autocomplete="off"
            />
          </el-form-item>
        </el-form>
      </el-scrollbar>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeWxDialog">取 消</el-button>
          <el-button
              type="primary"
              @click="submitWxDialog"
          >修改
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createBenchmark,
  deleteBenchmark,
  getBenchmarkAccount,
  getBenchmarkAccountList,
  updateBenchmark,
  updateWxToken
} from '@/api/benchmarkAccount'
import {ref} from 'vue'
import {ElMessage} from 'element-plus'
import {formatDate} from '@/utils/format'
import {getTopicList} from "@/api/topic";

defineOptions({
  name: 'Benchmark'
})


const form = ref({
  accountName: '',
  topic: '',
  initNum: 10,
  articleLink: '',
  key: ''
})

const wxForm = ref({
  slaveSid: '',
  bizUin: '',
  dataTicket: '',
  randInfo: '',
  token: '',
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const topicArr = ref([])

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
  const topicSelect = await getTopicList({page: page.value, pageSize: pageSize.value})
  if (topicSelect.code === 0) {
    topicArr.value = topicSelect.data.list
  }

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
const dialogWxFormVisible = ref(false)
const type = ref('')


const deleteWechatBenchmarkAccount = async (row) => {
  row.visible = false
  const res = await deleteBenchmark({ID: row.ID})
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

const updateWechatBenchmark = async (row) => {
  const res = await getBenchmarkAccount({ID: row.ID})
  type.value = 'update'
  if (res.code === 0) {
    form.value = res.data.benchmarkAccount
    dialogFormVisible.value = true
  }
}


const enterDialog = async () => {
  let res
  switch (type.value) {
    case 'create':
      res = await createBenchmark(form.value)
      break
    case 'update':
      res = await updateBenchmark(form.value)
      break
    default:
      res = await createBenchmark(form.value)
      break
  }

  if (res.code === 0) {
    closeDialog()
    getTableData()
  }
}

const submitWxDialog = async () => {
  let res = await updateWxToken(wxForm.value)

  if (res.code === 0) {
    closeWxDialog()
    getTableData()
  }
}

const closeDialog = () => {
  dialogFormVisible.value = false
  form.value = {
    accountName: '',
    topic: '',
    initNum: 10,
    articleLink: '',
    key: '',
  }
}

const closeWxDialog = () => {
  dialogWxFormVisible.value = false
  form.value = {
    slaveSid: '',
    bizUin: '',
    dataTicket: '',
    randInfo: '',
    token: '',
  }
}

const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}


const openWxDialog = () => {
  type.value = 'create'
  dialogWxFormVisible.value = true
}


</script>


<style></style>
