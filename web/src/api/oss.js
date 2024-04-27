import OSS from 'ali-oss';


export let client = new OSS({
    region: 'oss-cn-shenzhen',  // 填你的oss所在区域，例如oss-cn-shenzhen
    accessKeyId: 'LTAI5t6zY9ckKVK6pejYexRu', // 填你的oss的accessKeyId
    accessKeySecret: 'PwCVGIrOmwMScYfeDpBBzFu2z8TTee', // 填你的oss的accessSecret
    bucket: 'ai4media' // 你创建的路径名称
})