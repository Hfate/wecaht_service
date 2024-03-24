<template>
  <div>
    <el-form>
      <el-select v-model="script_model" style="float: left">
        <el-option label="other脚本(低性能)" value="other"></el-option>
        <el-option label="python脚本(中性能)" value="python"></el-option>
        <el-option label="go脚本(高性能)" value="go"></el-option>
      </el-select>
      <br><br><br>
      <el-upload
          :data="{'script_model':script_model}"
          accept=".jpg,.png"
          ref="upload"
          :action="`${path}/fileUploadAndDownload/upload`"
          :limit="1"
          :on-success="uploadSuccess"
          :on-error="uploadError"
          :show-file-list="false"
          class="upload-btn"
      >
        <el-button type="primary">上传脚本</el-button>
      </el-upload>
    </el-form>

  </div>
</template>

<script setup>

import {ref} from 'vue'
import {ElMessage} from 'element-plus'
import {useUserStore} from '@/pinia/modules/user'
import {isImageMime, isVideoMime} from '@/utils/image'

defineOptions({
  name: 'UploadCommon',
})

const emit = defineEmits(['on-success'])
const path = ref(import.meta.env.VITE_BASE_API)

const userStore = useUserStore()
const fullscreenLoading = ref(false)

const script_model = ref({
  inputValue: '',
})

const uploadData = ref({
  extraField: '',
})

const checkFile = (file) => {
  fullscreenLoading.value = true
  const isLt500K = file.size / 1024 / 1024 < 0.5 // 500K, @todo 应支持在项目中设置
  const isLt5M = file.size / 1024 / 1024 < 5 // 5MB, @todo 应支持项目中设置
  const isVideo = isVideoMime(file.type)
  const isImage = isImageMime(file.type)
  let pass = true
  if (!isVideo && !isImage) {
    ElMessage.error('上传图片只能是 jpg,png,svg,webp 格式, 上传视频只能是 mp4,webm 格式!')
    fullscreenLoading.value = false
    pass = false
  }
  if (!isLt5M && isVideo) {
    ElMessage.error('上传视频大小不能超过 5MB')
    fullscreenLoading.value = false
    pass = false
  }
  if (!isLt500K && isImage) {
    ElMessage.error('未压缩的上传图片大小不能超过 500KB，请使用压缩上传')
    fullscreenLoading.value = false
    pass = false
  }

  console.log('upload file check result: ', pass)

  return pass
}


const submitUpload = (res) => {
  this.uploadData.extraField = this.form.inputValue; // 在上传前将输入框的值赋给data对象的属性
  this.$refs.upload.submit(); // 调用 el-upload 的 submit 方法来上传文件
}

const uploadSuccess = (res) => {
  const {data} = res
  if (data.file) {
    emit('on-success', data.file.url)
  }
}

const uploadError = () => {
  ElMessage({
    type: 'error',
    message: '上传失败'
  })
  fullscreenLoading.value = false
}

</script>

