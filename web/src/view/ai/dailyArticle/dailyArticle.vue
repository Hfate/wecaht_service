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
            label="目标公众号"
            width="140"
            prop="targetAccountName"
        >
        </el-table-column>
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
            label="阅读量"
            prop="readNum"
            width="100"
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


  </div>

</template>

<script setup>
import {deleteArticle, deleteArticlesByIds, getArticleList, recreationArticle} from '@/api/dailyArticle'
import {ref} from 'vue'
import {ElMessage} from 'element-plus'
import {formatDate} from '@/utils/format'


defineOptions({
  name: 'Article'
})

const articles = ref([])
const page = ref(1)
const total = ref(0)
const pageSize = ref(50)
const tableData = ref([])
const searchInfo = ref({})


const onReset = () => {
  searchInfo.value = {}
}
// 条件搜索前端看此方法
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10

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


// 查询
const getTableData = async () => {
  const table = await getArticleList({page: page.value, pageSize: pageSize.value, ...searchInfo.value,})
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()


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


</script>


<style></style>
