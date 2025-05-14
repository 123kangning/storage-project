function handleFileUpload(file) {
    const sessionId = localStorage.getItem('session_id');
    if (!sessionId) {
        window.location.href = 'login.html';
        alert('请先登录');
        return;
    }

    const formData = new FormData();
    formData.append('file', file);

    fetch('http://localhost/v1/file/post/', {
        method: 'post',
        headers: {
            'X-Session-ID': sessionId
        },
        body: formData
    })
        .then(response => response.json())
        .then(data => {
            if (handleResponseError(data, '文件上传')) {
                console.log('文件上传成功', data);
                alert('文件上传成功');
            }
        })
        .catch(error => {
            console.error('文件上传失败', error);
        });
}

function uploadBox() {
    const uploadDropzone = document.getElementById('uploadDropzone');
    const fileInput = document.getElementById('fileInput');

    // 点击上传区域触发文件选择
    uploadDropzone.addEventListener('click', () => {
        fileInput.click();
    });

    // 处理文件选择
    fileInput.addEventListener('change', (e) => {
        const file = e.target.files[0];
        if (file) {
            handleFileUpload(file);
            // 处理完文件后，将 input 的 value 置为空
            e.target.value = '';
        }
    });

    // 拖拽相关事件
    uploadDropzone.addEventListener('dragover', (e) => {
        e.preventDefault();
        uploadDropzone.style.borderColor = '#2196F3';
    });

    uploadDropzone.addEventListener('dragleave', () => {
        uploadDropzone.style.borderColor = '#63b3ed';
    });

    uploadDropzone.addEventListener('drop', (e) => {
        e.preventDefault();
        uploadDropzone.style.borderColor = '#63b3ed';
        const file = e.dataTransfer.files[0];
        if (file) {
            handleFileUpload(file);
        }
    });
}