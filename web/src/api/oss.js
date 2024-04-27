import OSS from 'ali-oss';


export let client = new OSS({
    region: 'oss-cn-shenzhen',  // 填你的oss所在区域，例如oss-cn-shenzhen
    accessKeyId: '', // 填你的oss的accessKeyId
    accessKeySecret: '', // 填你的oss的accessSecret
    bucket: '' // 你创建的路径名称
})