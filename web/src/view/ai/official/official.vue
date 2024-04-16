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
            label="微信公众号"
            prop="accountName"
            width="120"
        >
        </el-table-column>
        <el-table-column
            align="left"
            label="主题"
            prop="topic"
            width="70"
        />
        <el-table-column
            align="left"
            label="账号邮箱"
            prop="userEmail"
            width="250"
        />
        <el-table-column
            align="left"
            label="默认作者"
            prop="defaultAuthorName"
            width="80"
        />
        <el-table-column
            align="left"
            label="每日默认生成文章数"
            prop="targetNum"
            width="150"
        />
        <el-table-column
            align="left"
            label="创作方式"
            min-width="180"
        >
          <template #default="scope">
            <el-cascader
                v-model="scope.row.createTypeList"
                :options="authOptions"
                :show-all-levels="false"
                collapse-tags
                :props="{ multiple:true,checkStrictly: true,label:'createName',value:'createType',disabled:'disabled',emitPath:false}"
                :clearable="false"
                @visible-change="(flag)=>{changeCreateType(scope.row,flag,0)}"
                @remove-tag="(removeAuth)=>{changeCreateType(scope.row,false,removeAuth)}"
            />
          </template>
        </el-table-column>

        <el-table-column
            align="left"
            label="创建时间"
            width="160"
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
                @click="updateWechatOfficialAccount(scope.row)"
            >更新
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
                    @click="deleteWechatOfficialAccount(scope.row)"
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
            <el-popover
                v-model="scope.row.visible"
                placement="top"
                width="160"
            >
              <p>确定要创作新文章吗？</p>
              <div style="text-align: right; margin-top: 8px;">
                <el-button
                    type="primary"
                    link
                    @click="scope.row.visible = false"
                >取消
                </el-button>
                <el-button
                    type="primary"
                    @click="wechatOfficialAccountCreate(scope.row)"
                >确定
                </el-button>
              </div>
              <template #reference>
                <el-button
                    type="primary"
                    link
                    icon="delete"
                    @click="scope.row.visible = true"
                >创作
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
        title="微信公众号"
    >
      <el-scrollbar height="500px">
        <el-form
            :model="form"
            label-position="right"
            label-width="150px"
        >
          <el-form-item label="微信公众号">
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
          <el-form-item label="账号邮箱">
            <el-input
                v-model="form.userEmail"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="默认作者">
            <el-input
                v-model="form.defaultAuthorName"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="AppId">
            <el-input
                v-model="form.appId"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="AppSecret">
            <el-input
                v-model="form.appSecret"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="Token">
            <el-input
                v-model="form.token"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="EncodingAESKey">
            <el-input
                v-model="form.encodingAesKey"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="备注">
            <el-input
                v-model="form.remark"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="每日默认生成文章数">
            <el-input
                v-model.number="form.targetNum"
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
import {
  createOfficialAccount,
  deleteOfficialAccount,
  getOfficialAccount,
  getOfficialAccountList,
  officialAccountCreate,
  setCreateTypes,
  updateOfficialAccount
} from '@/api/officialAccount'
import {nextTick, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {formatDate} from '@/utils/format'
import {getTopicList} from "@/api/topic";

defineOptions({
  name: 'OfficialAccount'
})

const form = ref({
  accountName: '',
  topic: '',
  userEmail: '',
  defaultAuthorName: '',
  remark: '',
  targetNum: 1,
  createTypeList: [],
  appId: '',
  appSecret: '',
  token: '',
  encodingAesKey: ''
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
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

  const table = await getOfficialAccountList({page: page.value, pageSize: pageSize.value})
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }

  setOptions()
}

getTableData()


const dialogFormVisible = ref(false)
const type = ref('')
const updateWechatOfficialAccount = async (row) => {
  const res = await getOfficialAccount({ID: row.ID})
  type.value = 'update'
  if (res.code === 0) {
    form.value = res.data
    dialogFormVisible.value = true
  }
}
const closeDialog = () => {
  dialogFormVisible.value = false
  form.value = {
    accountName: '',
    topic: '',
    userEmail: '',
    defaultAuthorName: '',
    remark: '',
    appId: '',
    appSecret: '',
    targetNum: 1,
    token: '',
    encodingAesKey: ''
  }
}


const deleteWechatOfficialAccount = async (row) => {
  row.visible = false
  const res = await deleteOfficialAccount({ID: row.ID})
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

const wechatOfficialAccountCreate = async (row) => {
  row.visible = false
  const res = await officialAccountCreate({ID: row.ID})
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '创作成功'
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
      res = await createOfficialAccount(form.value)
      break
    case 'update':
      res = await updateOfficialAccount(form.value)
      break
    default:
      res = await createOfficialAccount(form.value)
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


const authOptions = ref([])

const setOptions = () => {
  authOptions.value = []
  setAuthorityOptions(authOptions.value)
}


const setAuthorityOptions = (optionsData) => {

  const option1 = {
    createType: 1,
    createName: "热点创作",
  }
  optionsData.push(option1)


  const option2 = {
    createType: 2,
    createName: "洗稿",
  }
  optionsData.push(option2)
}


const tempAuth = {}
const changeCreateType = async (row, flag, removeAuth) => {
  if (flag) {
    if (!removeAuth) {
      tempAuth[row.ID] = [...row.createTypeList]
    }
    return
  }
  await nextTick()
  const res = await setCreateTypes({
    ID: row.ID,
    createTypeList: row.createTypeList
  })
  if (res.code === 0) {
    ElMessage({type: 'success', message: '创作类型设置成功'})
  } else {
    if (!removeAuth) {
      row.createTypeList = [...tempAuth[row.ID]]
      delete tempAuth[row.ID]
    } else {
      row.createTypeList = [removeAuth, ...row.createTypeList]
    }
  }
}

</script>


<style></style>
