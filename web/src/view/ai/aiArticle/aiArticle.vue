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
            show-overflow-tooltip="true"
        >
          <template #default="scope">
            <a :href="scope.row.link" target="_blank">{{ scope.row.title }}</a>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="主题"
            prop="topic"
            width="60"
        />
        <el-table-column
            align="left"
            label="目标公众号"
            prop="targetAccountName"
            width="100"
        />
        <el-table-column
            align="left"
            label="发布时间"
            width="120"
        >
          <template #default="scope">
            <span>{{ formatDate(scope.row.publishTime) }}</span>
          </template>
        </el-table-column>
        <!--        <el-table-column-->
        <!--            align="left"-->
        <!--            label="阅读量"-->
        <!--            prop="readNum"-->
        <!--            width="70"-->
        <!--        />-->
        <!--        <el-table-column-->
        <!--            align="left"-->
        <!--            label="评论量"-->
        <!--            prop="commentNum"-->
        <!--            width="70"-->
        <!--        />-->
        <!--        <el-table-column-->
        <!--            align="left"-->
        <!--            label="点赞量"-->
        <!--            prop="likeNum"-->
        <!--            width="70"-->
        <!--        />-->
        <el-table-column
            align="left"
            label="创建时间"
            width="120"
        >
          <template #default="scope">
            <span>{{ formatDate(scope.row.CreatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="创作进度"
            prop="processParams"
            width="320"
        >
          <template #default="scope">
            <span>{{ scope.row.processParams }}</span>
            <Progress
                v-if="scope.row.percent !== 100"
                :width="250"
                :stroke-width="10"
                :stroke-color="{
                  '0%': '#108ee9',
                  '100%': '#87d068',
                  direction: 'right'
              }"
                :percent="scope.row.percent"/>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="发布状态"
            width="100"
        >
          <template #default="scope">
            <span>{{ translatedArticleStatus(scope.row.articleStatus) }}</span>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="创作参数"
            prop="params"
            width="100"
        />
        <el-table-column
            align="left"
            label="相似度"
            prop="similarity"
            width="80"
        />
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
                @click="openNewTab('#/layout/ai/aiArticleDetail?id='+scope.row.ID)"
            >详情
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


  </div>
</template>


<style scoped>
@media (min-width: 1024px) {
  #sample {
    display: flex;
    flex-direction: column;
    place-items: center;
    width: 100%;
  }
}
</style>

<script setup>
import Progress from './Progress.vue'

import {
  deleteAIArticle,
  deleteAIArticlesByIds,
  generateAIArticle,
  getAIArticleList,
  publishAIArticles,
  recreationAIArticle
} from '@/api/aiArticle'

import {onMounted, onUnmounted, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import {formatDate} from '@/utils/format'
import {getTopicList} from "@/api/topic";


import {getOfficialAccountList,} from '@/api/officialAccount'


defineOptions({
  name: 'AIArticle'
})

const aiAIArticles = ref([])

const form = ref({
  ID: '',
  title: '',
  topic: '',
  targetAccountName: '',
  targetAccountId: '',
  originalContent: '',
  content: '',
  tags: '',
})


const openNewTab = (detailPath) => {
  // 使用window.open方法打开新标签页
  const url = window.location.origin + '/' + detailPath; // 构建完整的URL
  window.open(url, '_blank'); // 在新标签页中打开
}


const page = ref(1)
const total = ref(0)
const pageSize = ref(50)
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

// 批量操作
const handleSelectionChange = (val) => {
  aiAIArticles.value = val
}

const deleteVisible = ref(false)
const publishVisible = ref(false)


// onClick 事件处理
const onClick = (e) => {
  // 触发 onClick 事件，传递当前事件和 tinymce 实例
  // 这里可以根据需要进行事件处理
};

// 监视器，保持双向数据绑定的一致性
watch(() => form.value.content, (newValue) => {
  // 这里可以进行一些操作，比如将内容同步到其他状态管理中
});


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

const stop = ref(null); // 用于存储定时器的引用
onMounted(() => {
  // 定义调用getTableData的函数
  const fetchData = () => {
    getTableData();
  };

  // 设置定时器，每3秒调用一次fetchData
  stop.value = setInterval(fetchData, 5000);

  // 当组件卸载时清除定时器
  onUnmounted(() => {
    clearInterval(stop.value);
  });
});

const translatedStatus = ref({
  0: '新生成',
  1: '已发送至草稿箱',
  2: '发布成功',
  3: '群发成功',
  4: '发送草稿失败',
  5: '发布失败',
})

const translatedArticleStatus = (statusValue) => {
  return translatedStatus.value[statusValue]
}

</script>


<style></style>
