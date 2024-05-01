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
          <el-button
              icon="download"
              @click="onDownload"
          >下载
          </el-button>
          <el-button
              type="primary"
              icon="refresh"
              @click="openDialog"
          >素材余量统计
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
                :disabled="!articles.length"
                @click="deleteVisible = true"
            >删除
            </el-button>
          </template>
        </el-popover>

        <el-upload
            accept=".xlsx"
            ref="upload"
            auto-upload=auto-upload
            :action="`/api/article/upload`"
            :limit="1"
            :on-success="uploadSuccess"
            :on-error="uploadError"
            :show-file-list="true"
            :file-list="files"
            class="upload-btn"
        >
          <el-button type="primary">上传文件</el-button>
        </el-upload>
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
            label="文章标题"
            prop="title"
            width="300"
            show-overflow-tooltip="true"
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
            label="作者"
            prop="authorName"
            width="200"
        />
        <el-table-column
            align="left"
            label="发布时间"
            width="160"
        >
          <template #default="scope">
            <span>{{ formatDate(scope.row.publishTime) }}</span>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="阅读量"
            prop="readNum"
            width="70"
        />
        <el-table-column
            align="left"
            label="评论量"
            prop="commentNum"
            width="70"
        />
        <el-table-column
            align="left"
            label="点赞量"
            prop="likeNum"
            width="70"
        />
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
                    @click="deleteWechatArticle(scope.row)"
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
              <p>确定要改写吗？</p>
              <div style="text-align: right; margin-top: 8px;">
                <el-button
                    type="primary"
                    link
                    @click="scope.row.visible = false"
                >取消
                </el-button>
                <el-button
                    type="primary"
                    @click="recreationWechatArticle(scope.row)"
                >确定
                </el-button>
              </div>
              <template #reference>
                <el-button
                    type="primary"
                    link
                    icon="delete"
                    @click="scope.row.visible = true"
                >改写
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
        title="剩余可用素材"
        width="300"
    >
      <el-table
          ref="multipleTable"
          :data="articleStat"
          style="width:100%"
      >
        <el-table-column
            type="selection"
            width="55"
            v-if="false"
        />
        <el-table-column
            align="left"
            label="主题"
            prop="topic"
            width="100"
        />
        <el-table-column
            align="left"
            label="剩余素材"
            prop="count"
            width="160"
        />
      </el-table>
    </el-dialog>


  </div>

</template>


<script setup>
import {deleteArticle, deleteArticlesByIds, getArticleList, getArticleStats, recreationArticle} from '@/api/article'
import {ref} from 'vue'
import {ElMessage} from 'element-plus'
import {formatDate} from '@/utils/format'
import axios from 'axios'
import {saveAs} from 'file-saver'


defineOptions({
  name: 'Article'
})

const articles = ref([])
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const files = ref([])
const articleStat = ref([])


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
const onDownload = () => {
  page.value = 1
  pageSize.value = 200
  downloadArticle()
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


// 查询
const getTableData = async () => {
  const table = await getArticleList({page: page.value, pageSize: pageSize.value, ...searchInfo.value,})
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }

  const stats = await getArticleStats()
  if (stats.code == 0) {
    articleStat.value = stats.data.list
  }
}

getTableData()

const dialogFormVisible = ref(false)


const deleteVisible = ref(false)
const onDelete = async () => {
  const ids = articles.value.map(item => item.ID)
  const res = await deleteArticlesByIds({ids})
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


const closeDialog = () => {
  dialogFormVisible.value = false
}


const downloadArticle = async () => {
  // 发起 GET 请求并传递参数
  axios.get(import.meta.env.VITE_BASE_API + '/article/download', {
    params: {page: page.value, pageSize: pageSize.value, ...searchInfo.value,},
    responseType: 'blob'
  }).then(response => {
    // 将服务器返回的二进制数据保存为 Blob 对象
    const blob = new Blob([response.data], {type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'});

    // 使用 file-saver 库将 Blob 对象保存为本地文件
    saveAs(blob, 'article.xlsx');
  }).catch(error => {
    console.error('Failed to download Excel file:', error);
  });
}

const deleteWechatArticle = async (row) => {
  row.visible = false
  const res = await deleteArticle({ID: row.ID})
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


const recreationWechatArticle = async (row) => {
  row.visible = false
  const res = await recreationArticle({ID: row.ID})
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


const uploadSuccess = (res) => {
  files.value = []
  getTableData()
  const {data} = res
  if (data.file) {
    emit('on-success', data.file.url)
  }
}

const uploadError = () => {
  files.value = []
  ElMessage({
    type: 'error',
    message: '上传失败'
  })
}


const openDialog = () => {
  dialogFormVisible.value = true
}

</script>


<style></style>
