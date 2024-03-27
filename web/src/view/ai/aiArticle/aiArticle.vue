<template>
  <div>
    <div class="gva-search-box">
      <el-form
          :inline="true"
          :model="searchInfo"
      >
        <el-form-item label="文章标题">
          <el-input
              v-model="searchInfo.Title"
          />
        </el-form-item>
        <el-form-item label="公众号">
          <el-input
              v-model="searchInfo.TargetAccountName"
              placeholder=""
          />
        </el-form-item>
        <el-form-item label="门户">
          <el-input
              v-model="searchInfo.PortalName"
          />
        </el-form-item>
        <el-form-item label="主题">
          <el-input
              v-model="searchInfo.Topic"
              placeholder=""
          />
        </el-form-item>
        <el-form-item>
          <el-button
              icon="refresh"
              @click="onGenerate"
          >生成今日文章
          </el-button>
          <el-button
              type="primary"
              icon="search"
              @click="onSubmit"
          >查询
          </el-button>
          <el-button
              icon="refresh"
              @click="onReset"
          >重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-popover
            v-model="deleteVisible"
            placement="top"
            width="160"
        >
          <p>确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button
                type="primary"
                link
                @click="deleteVisible = false"
            >取消
            </el-button>
            <el-button
                type="primary"
                @click="onDelete"
            >确定
            </el-button>
          </div>
          <template #reference>
            <el-button
                icon="delete"
                :disabled="!aiAIArticles.length"
                @click="deleteVisible = true"
            >删除
            </el-button>
          </template>
        </el-popover>

        <el-popover
            v-model="publishVisible"
            placement="top"
            width="160"
        >
          <p>确定要发布吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button
                type="primary"
                link
                @click="publishVisible = false"
            >取消
            </el-button>
            <el-button
                type="primary"
                @click="onPublish"
            >确定
            </el-button>
          </div>
          <template #reference>
            <el-button
                icon="delete"
                :disabled="!aiAIArticles.length"
                @click="publishVisible = true"
            >发布
            </el-button>
          </template>
        </el-popover>
      </div>
      <el-table
          ref="multipleTable"
          :data="tableData"
          style="width: 100%"
          tooltip-effect="dark"
          row-key="ID"
          @selection-change="handleSelectionChange"
      >
        <el-table-column
            type="selection"
            width="55"
        />
        <el-table-column
            align="left"
            label="文章标题"
            prop="title"
            width="300"
        >
          <template #default="scope">
            <a :href="scope.row.link" target="_blank">{{ scope.row.title }}</a>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="门户"
            prop="portalName"
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
            label="目标发送公众号"
            prop="targetAccountName"
            width="250"
        />
        <el-table-column
            align="left"
            label="发布时间"
            width="180"
        >
          <template #default="scope">
            <span>{{ formatDate(scope.row.publishTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="阅读量"
            prop="readNum"
            width="100"
        />
        <el-table-column
            align="left"
            label="评论量"
            prop="commentNum"
            width="100"
        />
        <el-table-column
            align="left"
            label="点赞量"
            prop="likeNum"
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
                @click="updateWechatArticle(scope.row)"
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
                    @click="deleteWechatAIArticle(scope.row)"
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
              <p>确定要重写吗？</p>
              <div style="text-align: right; margin-top: 8px;">
                <el-button
                    type="primary"
                    link
                    @click="scope.row.visible = false"
                >取消
                </el-button>
                <el-button
                    type="primary"
                    @click="recreationWechatAIArticle(scope.row)"
                >确定
                </el-button>
              </div>
              <template #reference>
                <el-button
                    type="primary"
                    link
                    icon="delete"
                    @click="scope.row.visible = true"
                >重写
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
        title="新创作的文章"
    >
      <el-scrollbar height="600px">
        <el-form
            :model="form"
            label-position="right"
            label-width="110px"
        >
          <el-form-item label="标题">
            <el-input
                v-model="form.title"
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
          <el-form-item label="目标公众号">
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
          </el-form-item>
          <el-form-item label="内容">
            <el-input
                v-model="form.content"
                type="textarea"
                :rows="30"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="标签">
            <el-input
                v-model="form.tags"
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
          >更新
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  deleteAIArticle,
  publishAIArticles,
  deleteAIArticlesByIds,
  generateAIArticle,
  getAIArticle,
  getAIArticleList,
  recreationAIArticle,
  updateArticle
} from '@/api/aiArticle'
import {ref} from 'vue'
import {ElMessage} from 'element-plus'
import {formatDate} from '@/utils/format'
import {getTopicList} from "@/api/topic";

import {getOfficialAccountList,} from '@/api/officialAccount'


defineOptions({
  name: 'AIArticle'
})

const aiAIArticles = ref([])

const form = ref({
  title: '',
  topic: '',
  targetAccountName: '',
  targetAccountId: '',
  content: '',
  tags: '',
})

const closeDialog = () => {
  dialogFormVisible.value = false
  form.value = {
    title: '',
    topic: '',
    targetAccountName: '',
    targetAccountId: '',
    content: '',
    tags: '',
  }
}

const enterDialog = async () => {
  let res = await updateArticle(form.value)
  if (res.code === 0) {
    closeDialog()
    getTableData()
  }
}

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const topicArr = ref([])
const accountArr = ref([])

const onReset = () => {
  searchInfo.value = {}
}
// 条件搜索前端看此方法
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10

  getTableData()
}


// 条件搜索前端看此方法
const onGenerate = () => {
  generateAIArticle()

  getTableData()
}


// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询所有topic
const getTopicData = async () => {
  const topicSelect = await getTopicList({page: page.value, pageSize: pageSize.value})
  if (topicSelect.code === 0) {
    topicArr.value = topicSelect.data.list
  }
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
  const table = await getAIArticleList({page: page.value, pageSize: pageSize.value, ...searchInfo.value,})
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}


getTopicData()
getAccountData()
getTableData()

const dialogFormVisible = ref(false)
const type = ref('')

// 批量操作
const handleSelectionChange = (val) => {
  aiAIArticles.value = val
}

const deleteVisible = ref(false)
const publishVisible = ref(false)

const onDelete = async () => {
  const ids = aiAIArticles.value.map(item => item.ID)
  const res = await deleteAIArticlesByIds({ids})
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: res.msg
    })
    if (tableData.value.length === ids.length && page.value > 1) {
      page.value--
    }
    deleteVisible.value = false
    getTableData()
  }
}


const onPublish = async () => {
  const ids = aiAIArticles.value.map(item => item.ID)
  const res = await publishAIArticles({ids})
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: res.msg
    })
    if (tableData.value.length === ids.length && page.value > 1) {
      page.value--
    }
    publishVisible.value = false
    getTableData()
  }
}

const deleteWechatAIArticle = async (row) => {
  row.visible = false
  const res = await deleteAIArticle({ID: row.ID})
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

const updateWechatArticle = async (row) => {
  const res = await getAIArticle({ID: row.ID})
  type.value = 'update'
  if (res.code === 0) {
    form.value = res.data.article
    dialogFormVisible.value = true
  }
}

const recreationWechatAIArticle = async (row) => {
  row.visible = false
  const res = await recreationAIArticle({ID: row.ID})
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '改写成功'
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}

</script>


<style></style>
