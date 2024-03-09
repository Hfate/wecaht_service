import service from '@/utils/request'
import axios from 'axios';

// @Tags FileUploadAndDownload
// @Summary 分页文件列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "分页获取文件户列表"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /fileUploadAndDownload/getFileList [post]
export const getFileList = (data) => {
    return service({
        url: '/fileUploadAndDownload/getFileList',
        method: 'post',
        data
    })
}

// @Tags FileUploadAndDownload
// @Summary 删除文件
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body dbModel.FileUploadAndDownload true "传入文件里面id即可"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /fileUploadAndDownload/deleteFile [post]
export const deleteFile = (data) => {
    return service({
        url: '/fileUploadAndDownload/deleteFile',
        method: 'post',
        data
    })
}

/**
 * 编辑文件名或者备注
 * @param data
 * @returns {*}
 */
export const editFileName = (data) => {
    return service({
        url: '/fileUploadAndDownload/editFileName',
        method: 'post',
        data
    })
}


const fetch = axios.create({
    baseURL: import.meta.env.VITE_BASE_API,
    xsrfCookieName: '',
    withCredentials: true,
    timeout: 10000,
});

/**
 * 通过 fetch 请求下载 xlsx 文件并判断状态是否失败
 * @param url 文件地址
 * @param name 文件名称
 */
export function downloadFile(url) {
    fetch
        .get(url, {responseType: 'blob'})
        .then(res => {
            const contentType = res.headers['content-type'];
            const contentDisposition = res.headers['content-disposition'] || '';
            const matchFileName = contentDisposition.match(/filename=([\s\S]*)/)?.[1] || '';
            const fileName = decodeURI(matchFileName) || name || 'filename';
            // 文件流格式
            if (contentType === 'application/octet-stream') {
                // 设置文件格式 xlsx
                const blob = new Blob([res.data], {
                    type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;',
                });
                const downloadUrl = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.style.display = 'none';
                a.href = downloadUrl;
                a.download = fileName;
                document.body.appendChild(a);
                a.click();
                document.body.removeChild(a);
                URL.revokeObjectURL(url);
            }
        })
        .catch(e => {
            console.error('e', e);
        });
}