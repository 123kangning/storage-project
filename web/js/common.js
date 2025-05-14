// 提取公共的错误处理方法
function handleResponseError(data, action) {
    //code==0表示成功
    if (data.code === 0) {
        console.log(`${action}成功`);
        return true;
    }
    console.error(`${action}失败,`, data.message);
    alert(data.message);
    switch (data.code) {
        case 1:
            break;
        case 2:
            //未登录，无权限，跳转到登录页面
            window.location.href = 'login.html';
            break;
        default:
    }
    return false;
}

/**
 * 处理响应状态码
 * @param {Response} response - fetch 请求的响应对象
 * @returns {Promise<any>} - 解析后的响应数据
 */
function handleResponseStatus(response) {
    if (response.status === 401) {
        // 处理 401 状态码，清除 sessionId 并跳转登录页面
        localStorage.removeItem('session_id');
        window.location.href = 'login.html';
        alert('会话已过期，请重新登录');
        return Promise.reject(new Error('会话已过期'));
    } else if (response.status === 502) {
        // 处理 502 状态码，可按需添加更多提示
        alert('服务器暂时不可用，请稍后重试');
        return Promise.reject(new Error('服务器暂时不可用'));
    }
    // 其他状态码按正常逻辑处理
    if (response.ok) {
        const contentType = response.headers.get('Content-Type');
        if (contentType && contentType.includes('application/json')) {
            return response.json();
        }
        return response.blob();
    }
    return Promise.reject(new Error(`请求失败，状态码: ${response.status}`));
}

// 格式化文件大小
function formatFileSize(bytes) {
    if (bytes < 1024) return bytes + ' B';
    const units = ['KB', 'MB', 'GB', 'TB'];
    let unitIndex = -1;
    do {
        bytes /= 1024;
        unitIndex++;
    } while (bytes >= 1024 && unitIndex < units.length - 1);
    return bytes.toFixed(2) + ' ' + units[unitIndex];
}