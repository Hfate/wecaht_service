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
            label="门户"
            prop="portalName"
            width="180"
        >
        </el-table-column>
        <el-table-column
            align="left"
            label="门户Key"
            prop="portalKey"
            width="120"
        />
        <el-table-column
            align="left"
            label="文章key"
            prop="articleKey"
            width="120"
        />
        <el-table-column
            align="left"
            label="门户主页"
            prop="link"
            width="200"
        />
        <el-table-column
            align="left"
            label="主题"
            prop="topic"
            width="120"
        />
        <el-table-column
            align="left"
            label="目标爬取网页数"
            prop="targetNum"
            width="160"
        />
        <el-table-column
            align="left"
            label="备注"
            prop="remark"
            width="120"
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
                @click="updateWechatPortal(scope.row)"
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
                    @click="deleteWechatPortal(scope.row)"
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
        title="门户"
    >
      <el-scrollbar height="500px">
        <el-form
            :model="form"
            label-position="right"
            label-width="90px"
        >
          <el-form-item label="门户名字">
            <el-input
                v-model="form.portalName"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="门户Key">
            <el-input
                v-model="form.portalKey"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="文章Key">
            <el-input
                v-model="form.articleKey"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="门户主页">
            <el-input
                v-model="form.link"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="主题">
            <el-input
                v-model="form.topic"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="GraphQuery">
            <el-input
                v-model="form.graphQuery"
                type="textarea"
                :rows="10"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="目标网页数">
            <el-input
                v-model="form.targetNum"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="备注">
            <el-input
                v-model="form.remark"
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
import {createPortal, deletePortal, getPortal, getPortalList, updatePortal} from '@/api/portal'
import {ref} from 'vue'
import {ElMessage} from 'element-plus'
import {formatDate} from '@/utils/format'

defineOptions({
  name: 'Portal'
})

const form = ref({
  portalName: '',
  portalKey: '',
  articleKey: '',
  link: '',
  topic: '',
  graphQuery: '',
  targetNum: 0,
  remark: '',
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

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
  const table = await getPortalList({page: page.value, pageSize: pageSize.value})
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
const updateWechatPortal = async (row) => {
  const res = await getPortal({ID: row.ID})
  type.value = 'update'
  if (res.code === 0) {
    form.value = res.data.portal
    dialogFormVisible.value = true
  }
}
const closeDialog = () => {
  dialogFormVisible.value = false
  form.value = {
    portalName: '',
    portalKey: '',
    articleKey: '',
    link: '',
    topic: '',
    graphQuery: '',
    targetNum: 0,
    remark: '',
  }
}
const deleteWechatPortal = async (row) => {
  row.visible = false
  const res = await deletePortal({ID: row.ID})
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
      res = await createPortal(form.value)
      break
    case 'update':
      res = await updatePortal(form.value)
      break
    default:
      res = await createPortal(form.value)
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

</script>


<style></style>
