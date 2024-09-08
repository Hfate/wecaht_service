<template>
  <div>

    <div class="gva-search-box">
      <el-form
          :inline="true"
          :model="searchInfo"
      >
        <el-form-item label="公众号">
          <el-input
              v-model="searchInfo.Title"
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
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">

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
            label="账号名称"
            prop="accountName"
            width="150"
        >
        </el-table-column>
        <el-table-column
            align="left"
            label="时间区间"
            prop="zone"
            width="300"
        />
        <el-table-column
            align="left"
            label="区间内结算收入"
            prop="settledRevenue"
            width="150"
        >
          <template #default="scope">
            <span>{{ divideBy100(scope.row.settledRevenue) }}</span>
          </template>
        </el-table-column>
        <el-table-column
            align="left"
            label="结算状态"
            prop="settStatus"
            width="150"
        >
          <template #default="scope">
            {{ formatSettStatus(scope.row.settStatus) }}
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
import {getSettlementList} from '@/api/settle'
import {ref} from 'vue'
import axios from "axios";
import {saveAs} from 'file-saver'

defineOptions({
  name: 'CssFormat'
})

const form = ref({
  formatName: '',
  cssCode: '',
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


// 定义过滤器
const divideBy100 = (value) => {
  return value / 100;
};

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async () => {
  const table = await getSettlementList({page: page.value, pageSize: pageSize.value})
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }

}

getTableData()

// 定义计算属性
const formatSettStatus = (status) => {
  switch (status) {
    case 1:
      return '结算中';
    case 2:
    case 3:
      return '已结算';
    case 4:
      return '付款中';
    case 5:
      return '已付款';
    default:
      return '-';
  }
};

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

// 条件搜索前端看此方法
const onDownload = () => {
  page.value = 1
  pageSize.value = 200
  downloadSettlement()
}

const downloadSettlement = async () => {
  // 发起 GET 请求并传递参数
  axios.get(import.meta.env.VITE_BASE_API + '/settlement/download', {
    params: {page: page.value, pageSize: pageSize.value, ...searchInfo.value,},
    responseType: 'blob'
  }).then(response => {
    // 将服务器返回的二进制数据保存为 Blob 对象
    const blob = new Blob([response.data], {type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'});

    // 使用 file-saver 库将 Blob 对象保存为本地文件
    saveAs(blob, 'settlement.xlsx');
  }).catch(error => {
    console.error('Failed to download Excel file:', error);
  });
}

</script>

<style>
</style>
