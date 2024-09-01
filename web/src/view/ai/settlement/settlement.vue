<template>
  <div>

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

</script>

<style>
</style>
