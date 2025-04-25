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