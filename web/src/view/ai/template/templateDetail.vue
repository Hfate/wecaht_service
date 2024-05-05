<template>
  <div>
    <div class="gva-form-box">
      <el-form
          label-position="right"
          label-width="100px"
          :model="form"
      >
        <el-form-item label="目标公众号">
          <el-select
              v-model="form.accountId"
              class="w-56"
          >
            <el-option
                v-for="item in accountArr"
                :value="item.appId"
                :label="item.accountName"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="模板">
          <Toolbar
              style="border-bottom: 1px solid #ccc"
              :editor="editorRef"
              :defaultConfig="toolbarConfig"
              :mode="simple"
          />
          <Editor
              style="height: 400px; overflow-y: hidden;"
              v-model="form.templateValue"
              :defaultConfig="editorConfig"
              :mode="simple"
              @onCreated="onCreated"
          />
        </el-form-item>
      </el-form>


      <div class="dialog-footer">
        <el-button
            type="primary"
            @click="enterDialog"
        >保存
        </el-button>
      </div>

    </div>
  </div>

</template>


<style scoped>
/* Flexbox布局样式 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 30px; /* 根据需要调整间距 */
  margin-right: 130px; /* 根据需要调整边距 */
}
</style>

<style src="@wangeditor/editor/dist/css/style.css"></style>


<script setup>
import {getTemplate, updateTemplate} from "@/api/template";
import {useRoute} from 'vue-router';
import {onBeforeUnmount, onMounted, ref, watch} from 'vue';
import {getOfficialAccountList} from "@/api/officialAccount";
import {ElMessage} from "element-plus";
import {Editor, Toolbar} from '@wangeditor/editor-for-vue'


defineOptions({
  name: 'TemplateDetail'
})

const form = ref({
  accountName: '',
  accountId: '',
  templateValue: '',
})


const accountArr = ref([])
const articleId = ref('');

const editorRef = ref(null)
const toolbarConfig = ref({})
const editorConfig = ref({placeholder: '请输入内容...', MENU_CONF: {}})


onMounted(() => {
  getWechatTemplate();
  getAccountData();
});

onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor) {
    editor.destroy() // 组件销毁时，及时销毁编辑器
  }
})

// 方法
const onCreated = (editor) => {
  editorRef.value = Object.seal(editor) // 使用 Object.seal() 封装

  console.log(basePath)

}


//
editorConfig.value.MENU_CONF['uploadImage'] = {
  server: '/api/file/upload',
  // form-data fieldName ，默认值 'wangeditor-uploaded-image'
  fieldName: 'file',

  // 单个文件的最大体积限制，默认为 2M
  maxFileSize: 1 * 1024 * 1024, // 1M

  // 最多可上传几个文件，默认为 100
  maxNumberOfFiles: 1,

  // 选择文件时的类型限制，默认为 ['image/*'] 。如不想限制，则设置为 []
  allowedFileTypes: ['image/*'],


  // 将 meta 拼接到 url 参数中，默认 false
  metaWithUrl: false,


  // 跨域是否传递 cookie ，默认为 false
  withCredentials: true,

  // 超时时间，默认为 10 秒
  timeout: 5 * 1000, // 5 秒
}


const enterDialog = async () => {
  let res = await updateTemplate(form.value)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: res.msg
    })
    closeDialog()
    getTableData()
  }
}


// 监视器，保持双向数据绑定的一致性
watch(() => form.value.templateValue, (newValue) => {
  // 这里可以进行一些操作，比如将内容同步到其他状态管理中
});


const getWechatTemplate = async () => {
  // 解析哈希中的查询参数
  const route = useRoute();
  articleId.value = route.query.id || ''
  const res = await getTemplate({ID: articleId.value})
  if (res.code === 0) {
    form.value = res.data.template
  }
}


// 查询所有公众号
const getAccountData = async () => {
  const accountSelect = await getOfficialAccountList({page: 1, pageSize: 1000})
  if (accountSelect.code === 0) {
    accountArr.value = accountSelect.data.list
  }
}


</script>

