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
            label="排版名称"
            prop="formatName"
            width="150"
        >
        </el-table-column>
        <el-table-column
            align="left"
            label="css代码"
            prop="cssCode"
            width="500"
            show-overflow-tooltip="true"
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
            <el-button
                type="primary"
                link
                icon="edit"
                @click="updateWechatCssFormat(scope.row)"
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
                    @click="deleteWechatCssFormat(scope.row)"
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
        title="css排版"
    >
      <el-scrollbar height="500px">
        <el-form
            :model="form"
            label-position="right"
            label-width="150px"
        >
          <el-form-item label="排版名称">
            <el-input
                v-model="form.formatName"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="css代码">
            <el-input
                v-model.number="form.cssCode"
                type="textarea"
                :rows="30"
                autocomplete="off"
            />
          </el-form-item>
          <el-form-item label="Rendered Content in iframe">
            <iframe ref="iframe" style="width: 100%; height: 500px; border: 1px solid #ccc;"></iframe>
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
import {createCssFormat, deleteCssFormat, getCssFormat, getCssFormatList, updateCssFormat} from '@/api/cssFormat'
import {onMounted, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import {formatDate} from '@/utils/format'

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

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async () => {
  const table = await getCssFormatList({page: page.value, pageSize: pageSize.value})
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
const updateWechatCssFormat = async (row) => {
  const res = await getCssFormat({ID: row.ID})
  type.value = 'update'
  if (res.code === 0) {
    form.value = res.data
    dialogFormVisible.value = true
  }
}
const closeDialog = () => {
  dialogFormVisible.value = false
  form.value = {
    formatName: '',
    cssCode: '',
  }
}


const deleteWechatCssFormat = async (row) => {
  row.visible = false
  const res = await deleteCssFormat({ID: row.ID})
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
      res = await createCssFormat(form.value)
      break
    case 'update':
      res = await updateCssFormat(form.value)
      break
    default:
      res = await createCssFormat(form.value)
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


const htmlContent = ref('<!DOCTYPE html>\n' +
    '<html lang="en">\n' +
    '<head>\n' +
    '    <meta charset="UTF-8">\n' +
    '    <title>测试页面</title>\n' +
    '</head>\n' +
    '<body>\n' +
    '    <h1>这是一级标题</h1>\n' +
    '    <p>这是一段普通文本。<span class="colorful">这部分文本将应用自定义颜色。</span></p>\n' +
    '    <ul>\n' +
    '        <li>无序列表项一</li>\n' +
    '        <li>无序列表项二</li>\n' +
    '        <li>无序列表项三</li>\n' +
    '    </ul>\n' +
    '    <ol>\n' +
    '        <li>有序列表项一</li>\n' +
    '        <li>有序列表项二</li>\n' +
    '        <li>有序列表项三</li>\n' +
    '    </ol>\n' +
    '    <table>\n' +
    '        <tr>\n' +
    '            <th>表头单元格一</th>\n' +
    '            <th>表头单元格二</th>\n' +
    '        </tr>\n' +
    '        <tr>\n' +
    '            <td>表格单元格一</td>\n' +
    '            <td>表格单元格二</td>\n' +
    '        </tr>\n' +
    '        <tr>\n' +
    '            <td>表格单元格三</td>\n' +
    '            <td>表格单元格四</td>\n' +
    '        </tr>\n' +
    '    </table>\n' +
    '    <button>这是一个按钮</button>\n' +
    '    <br>\n' +
    '    <hr>\n' +
    '    <form action="">\n' +
    '        <input type="text" placeholder="输入一些文本">\n' +
    '        <input type="submit" value="提交">\n' +
    '    </form>\n' +
    '    <div>\n' +
    '        内联元素，紧接着是另一个<br>\n' +
    '        内联元素，它们位于同一行。\n' +
    '    </div>\n' +
    '    <pre>这是一个预格式化的文本块。它保留了空白和\\n换行符。</pre>\n' +
    '    <blockquote>这是一个引用块，通常用于引用其他文本。</blockquote>\n' +
    '    <code>这是行内代码元素。</code>\n' +
    '    <address>这是作者或拥有者的联系信息。</address>\n' +
    '    <small>这是副标题或者法律条款的小字体文本。</small>\n' +
    '</body>\n' +
    '</html>');  // 默认HTML内容

const iframe = ref(null);  // Ref for the iframe element

// Function to update iframe content
const updateIframeContent = () => {
  if (!iframe.value) return; // Check if the iframe is correctly referenced

  const doc = iframe.value.contentDocument || iframe.value.contentWindow.document;
  doc.open();
  doc.write(`<style>${form.value.cssCode}</style>${htmlContent.value}`);
  doc.close();
};


// Watch form.cssFormat and form.content for changes
watch(() => [form.value.cssCode, htmlContent.value], () => {
  updateIframeContent();
}, {deep: true});

// Update the iframe on component mount
onMounted(updateIframeContent);

</script>

<style>
</style>
