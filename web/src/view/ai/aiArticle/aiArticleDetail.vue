<template>
  <div>
    <div class="gva-form-box">
      <el-form
          label-position="right"
          label-width="100px"
          :model="form"
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
        <el-form-item label="标签">
          <el-input
              v-model="form.tags"
              autocomplete="off"
          />
        </el-form-item>
        <el-form-item label="排版文本">
          <editor v-model="form.content"
                  api-key="nujoppecgjjvzjg005s5fie8aqqwbsu8f28aya3no68zzi4v"
                  :init="init"
                  :disabled="false"
                  @onClick="onClick">
          </editor>
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

<script setup>
import Editor from "@tinymce/tinymce-vue";
import {client} from '@/api/oss';
import {getAIArticle, updateArticle} from "@/api/aiArticle";
import md5 from 'blueimp-md5';
import {useRoute} from 'vue-router';
import {onMounted, ref, watch} from 'vue';
import {getOfficialAccountList} from "@/api/officialAccount";
import {getTopicList} from "@/api/topic";
import {ElMessage} from "element-plus";

defineOptions({
  name: 'AIArticleDetail'
})

const form = ref({
  title: '',
  topic: '',
  targetAccountName: '',
  targetAccountId: '',
  originalContent: '',
  content: '',
  tags: '',
})

const page = ref(1);
const pageSize = ref(100);
const topicArr = ref([])
const accountArr = ref([])
const articleId = ref('');

onMounted(() => {
  getWechatArticle();
  getTopicData();
  getAccountData();
});


// 初始化配置
const init = {
  //language_url: '/static/tinymce/langs/zh_CN.js',
  //language: 'zh_CN',
  //skin_url: '/static/tinymce/skins/ui/oxide',
  height: 500,
  width: 1200,
  plugins: 'lists image media table wordcount preview code',
  toolbar: 'code | undo redo | formatselect | bold italic | alignleft aligncenter alignright alignjustify | bullist numlist outdent indent | lists image media table | removeformat preview',
  branding: false,
  menubar: true,
  images_upload_handler: (blobInfo, success, failure) => {
    const filename = blobInfo.filename();
    const suffix = filename.substring(filename.lastIndexOf('.') + 1);
    const nameWithMd5AndTime = `${md5(blobInfo.base64())}${getTime()}.${suffix}`;
    client.multipartUpload(nameWithMd5AndTime, blobInfo.blob()).then((result) => {
      if (result.res.requestUrls) {
        console.log('返回结果', result);
        success(result.res.requestUrls[0].split('?')[0]);
      }
    }).catch((err) => {
      console.log(err);
    });
  },
};


const enterDialog = async () => {
  let res = await updateArticle(form.value)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: res.msg
    })
    closeDialog()
    getTableData()
  }
}

// 时间戳方法
const getTime = () => {
  const time = new Date();
  // ... 生成时间字符串的逻辑 ...
  return time;
};

// onClick 事件处理
const onClick = (e) => {
  // 触发 onClick 事件，传递当前事件和 tinymce 实例
  // 这里可以根据需要进行事件处理
};

// 监视器，保持双向数据绑定的一致性
watch(() => form.value.content, (newValue) => {
  // 这里可以进行一些操作，比如将内容同步到其他状态管理中
});


const getWechatArticle = async () => {
  // 解析哈希中的查询参数
  const route = useRoute();
  articleId.value = route.query.id || ''
  const res = await getAIArticle({ID: articleId.value})
  if (res.code === 0) {
    form.value = res.data.article
  }
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


</script>

