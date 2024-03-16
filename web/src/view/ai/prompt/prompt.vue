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
            label="主题"
            prop="topic"
            width="180"
        >
        </el-table-column>
        <el-table-column
            align="left"
            label="类型"
            prop="promptType"
            width="120"
        >
          <template #default="scope">
            <span>{{ scope.row.promptType == 1 ? "内容改写" : "标题改写" }}</span>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="Prompt"
            prop="prompt"
            width="300"
        />
        <el-table-column
            align="left"
            label="语言"
            prop="language"
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
            label="更新时间"
            width="180"
        >
          <template #default="scope">
            <span>{{ formatDate(scope.row.UpdatedAt) }}</span>
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
                @click="updateWechatPrompt(scope.row)"
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
                    @click="deleteWechatPrompt(scope.row)"
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
        title="Prompt"
    >
      <el-scrollbar height="500px">
        <el-form
            :model="form"
            label-position="right"
            label-width="90px"
        >
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
          <el-form-item label="类型" prop="promptType">
            <el-select
                v-model="form.promptType"
                style="width:100%"
            >
              <el-option
                  :value="1"
                  label="内容改写"
              />
              <el-option
                  :value="2"
                  label="标题改写"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="prompt">
            <el-input
                v-model="form.prompt"
                type="textarea"
                :rows="20"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="语言">
            <el-input
                v-model="form.language"
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
import {createPrompt, deletePrompt, getPrompt, getPromptList, updatePrompt} from '@/api/prompt'
import {getTopicList} from '@/api/topic'
import {ref} from 'vue'
import {ElMessage} from 'element-plus'
import {formatDate} from '@/utils/format'

defineOptions({
  name: 'Prompt'
})

const form = ref({
  topic: '',
  promptType: '',
  prompt: '',
  language: '',
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

  const table = await getPromptList({page: page.value, pageSize: pageSize.value})
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
const updateWechatPrompt = async (row) => {
  const res = await getPrompt({ID: row.ID})
  type.value = 'update'
  if (res.code === 0) {
    form.value = res.data.prompt
    dialogFormVisible.value = true
  }
}
const closeDialog = () => {
  dialogFormVisible.value = false
  form.value = {
    topic: '',
    promptType: '',
    prompt: '',
    language: '',
  }
}
const deleteWechatPrompt = async (row) => {
  row.visible = false
  const res = await deletePrompt({ID: row.ID})
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
      res = await createPrompt(form.value)
      break
    case 'update':
      res = await updatePrompt(form.value)
      break
    default:
      res = await createPrompt(form.value)
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
