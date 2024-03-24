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
            align="left"
            label="预览"
            width="100"
        >
          <template #default="scope">
            <CustomPic
                pic-type="file"
                :pic-src="scope.row.link"
                preview
            />
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="日期"
            prop="UpdatedAt"
            width="180"
        >
          <template #default="scope">
            <div>{{ formatDate(scope.row.UpdatedAt) }}</div>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="文件名/备注"
            prop="fileName"
            width="180"
        >
        </el-table-column>
        <el-table-column
            align="left"
            label="链接"
            prop="link"
            min-width="300"
        />
        <el-table-column
            align="left"
            label="标签"
            prop="tag"
            width="100"
        >
          <template #default="scope">
            <el-tag
                :type="scope.row.tag === 'jpg' ? 'info' : 'success'"
                disable-transitions
            >{{ scope.row.tag }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="操作"
            width="160"
        >
          <template #default="scope">
            <el-button
                icon="download"
                type="primary"
                link
                @click="downloadFile(scope.row)"
            >下载
            </el-button>
            <el-button
                icon="delete"
                type="primary"
                link
                @click="deleteWechatMedia(scope.row)"
            >删除
            </el-button>
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
        title="素材"
    >
      <el-form>
        <el-select
            v-model="form.targetAccountId"
            class="w-56"
        >
          <el-option
              v-for="item in accountArr"
              :value="item.appId"
              :label="item.accountName"
          />
        </el-select>
        <br><br><br>
        <el-upload
            :data="{'targetAccountId':form.targetAccountId}"
            accept=".jpg,.png"
            ref="upload"
            :action="`/api/media/media`"
            :limit="1"
            :on-success="uploadSuccess"
            :on-error="uploadError"
            :show-file-list="false"
            class="upload-btn"
        >
          <el-button type="primary">上传文件</el-button>
        </el-upload>
      </el-form>
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
import {createMedia, deleteMedia, getMediaPage, updateMedia} from '@/api/media'
import {ref} from 'vue'
import {ElMessage} from 'element-plus'
import {formatDate} from '@/utils/format'
import {downloadImage} from "@/utils/downloadImg";
import {getOfficialAccountList} from "@/api/officialAccount";

defineOptions({
  name: 'Media'
})

const form = ref({
  targetAccountId: '',
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const accountArr = ref([])

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询所有公众号
const getAccountData = async () => {
  const accountSelect = await getOfficialAccountList({page: 1, pageSize: 1000})
  if (accountSelect.code === 0) {
    accountArr.value = accountSelect.data.list
  }
}

// 查询
const getTableData = async () => {
  const table = await getMediaPage({page: page.value, pageSize: pageSize.value})
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getAccountData()
getTableData()

const dialogFormVisible = ref(false)
const type = ref('')

const closeDialog = () => {
  dialogFormVisible.value = false
  form.value = {
    targetAccountId: '',
  }
}

const deleteWechatMedia = async (row) => {
  row.visible = false
  const res = await deleteMedia({ID: row.ID})
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

const downloadFile = (row) => {
  if (row.url.indexOf('http://') > -1 || row.url.indexOf('https://') > -1) {
    downloadImage(row.url, row.name)
  } else {
    debugger
    downloadImage(path.value + '/' + row.url, row.name)
  }
}


const enterDialog = async () => {
  let res
  switch (type.value) {
    case 'create':
      res = await createMedia(form.value)
      break
    case 'update':
      res = await updateMedia(form.value)
      break
    default:
      res = await createMedia(form.value)
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


const uploadSuccess = (res) => {
  dialogFormVisible.value = false
  form.value = {
    targetAccountId: '',
  }
  getTableData()
  const {data} = res
  if (data.file) {
    emit('on-success', data.file.url)
  }
}

const uploadError = () => {
  ElMessage({
    type: 'error',
    message: '上传失败'
  })
  fullscreenLoading.value = false
}

</script>


<style></style>
