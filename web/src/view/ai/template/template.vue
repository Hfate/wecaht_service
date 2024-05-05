<template>
  <div>
    <div class="gva-table-box">
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
            label="公众号"
            prop="accountName"
            width="200"
        >
        </el-table-column>


        <el-table-column
            align="left"
            label="创建时间"
            width="200"
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
                @click="openNewTab('#/layout/ai/templateDetail?id='+scope.row.ID)"
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
                    @click="deleteWechatTemplate(scope.row)"
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
              <p>确定要克隆吗？</p>
              <div style="text-align: right; margin-top: 8px;">
                <el-button
                    type="primary"
                    link
                    @click="scope.row.visible = false"
                >取消
                </el-button>
                <el-button
                    type="primary"
                    @click="cloneWechatTemplate(scope.row)"
                >确定
                </el-button>
              </div>
              <template #reference>
                <el-button
                    type="primary"
                    link
                    icon="delete"
                    @click="scope.row.visible = true"
                >克隆
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

import {cloneTemplate, deleteTemplate, getTemplateList} from '@/api/template'

import {ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import {formatDate} from '@/utils/format'
import {getTopicList} from "@/api/topic";


import {getOfficialAccountList,} from '@/api/officialAccount'


defineOptions({
  name: 'Template'
})

const aiTemplates = ref([])

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
  const table = await getTemplateList({page: page.value, pageSize: pageSize.value, ...searchInfo.value,})
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
  aiTemplates.value = val
}


// 监视器，保持双向数据绑定的一致性
watch(() => form.value.content, (newValue) => {
  // 这里可以进行一些操作，比如将内容同步到其他状态管理中
});


const deleteWechatTemplate = async (row) => {
  row.visible = false
  const res = await deleteTemplate({ID: row.ID})
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

const cloneWechatTemplate = async (row) => {
  row.visible = false
  const res = await cloneTemplate({ID: row.ID})
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '克隆成功'
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}


</script>


<style></style>
