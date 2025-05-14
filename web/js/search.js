function searchBox() {
    document.addEventListener('DOMContentLoaded', function () {
        // 文件查询表单提交事件
        const fromInput = document.getElementById('from');
        const prevPageButton = document.getElementById('prevPage');
        const nextPageButton = document.getElementById('nextPage');
        const pageInfo = document.getElementById('pageInfo');
        const pagination = document.getElementById('pagination');

        let currentFrom = 1;
        let totalSize = 0;
        const pageSize = 10;
        const resetButton = document.querySelector('#searchForm button[type="reset"]');
        resetButton.addEventListener('click', function () {
            // Clear the table body content
            const tableBody = document.getElementById('fileTable').getElementsByTagName('tbody')[0];
            tableBody.innerHTML = '';
            pagination.style.display = 'none';
            currentFrom = 1;
        });
        const searchForm = document.getElementById('searchForm');
        searchForm.addEventListener('submit', function (e) {
            e.preventDefault();
            const sessionId = localStorage.getItem('session_id');
            if (!sessionId) {
                window.location.href = 'login.html';
                alert('请先登录');
                return;
            }

            const formData = new FormData(this);
            const name = formData.get('name');
            const url = `http://localhost/v1/file/search?name=${name}&from=${currentFrom}&size=${pageSize}`;
            fetch(url, {
                headers: {
                    'X-Session-ID': sessionId
                }
            })
                .then(handleResponseStatus)
                .then(data => {
                    const baseResp = data.baseResp;
                    if (!handleResponseError(baseResp, '文件查询')) {
                        return;
                    }

                    // 清空表格现有数据
                    const tableBody = document.getElementById('fileTable').getElementsByTagName('tbody')[0];
                    tableBody.innerHTML = '';
                    const realData = data.data.files;
                    totalSize = data.data.total;
                    if (realData.length === 0) {
                        pagination.style.display = 'none';
                        alert('没有查询到相关文件');
                    }
                    // 填充表格数据
                    realData.forEach(file => {
                        const row = document.createElement('tr');
                        const nameCell = document.createElement('td');
                        const sourceCell = document.createElement('td');
                        const sizeCell = document.createElement('td');
                        const hashCell = document.createElement('td');
                        const operationCell = document.createElement('td');

                        nameCell.textContent = file.name;
                        if (file.name.length > 20) {
                            nameCell.textContent = file.name.substring(0, 20) + '...';
                        }
                        sizeCell.textContent = formatFileSize(file.size);
                        hashCell.textContent = file.hash.substring(0, 5) + '...'+ file.hash.substring(file.hash.length - 5);
                        sourceCell.textContent = file.source;

                        const downloadButton = document.createElement('button');
                        downloadButton.textContent = '下载';
                        downloadButton.classList.add('btn', 'top-bar-color-btn');
                        downloadButton.addEventListener('click', function () {
                            const downloadUrl = `http://localhost/v1/file/get?hash=${file.hash}`;
                            fetch(downloadUrl,{
                                headers: {
                                    'X-Session-ID': sessionId
                                }
                            })
                                .then(handleResponseStatus)
                                .then(data => {
                                    console.log('data.type=',data.type, data);
                                    if (typeof data === 'object' &&!data.type) {
                                        handleResponseError(data, '文件下载')
                                    } else {
                                        // 处理文件下载
                                        const url = window.URL.createObjectURL(data);
                                        const a = document.createElement('a');
                                        a.href = url;
                                        a.download = file.name;
                                        a.click();
                                        window.URL.revokeObjectURL(url);
                                    }
                                })
                                .catch(error => {
                                    console.error('文件下载失败', error);
                                });
                        });

                        operationCell.appendChild(downloadButton);

                        row.appendChild(nameCell);
                        row.appendChild(sourceCell);
                        row.appendChild(sizeCell);
                        row.appendChild(hashCell);
                        row.appendChild(operationCell);

                        tableBody.appendChild(row);
                    });

                    // 更新分页按钮状态，第一页时禁用上一页按钮
                    prevPageButton.disabled = currentFrom === 1;
                    nextPageButton.disabled = currentFrom * pageSize >= totalSize;
                    const cur1 = (currentFrom - 1) * pageSize + 1;
                    let cur2 = currentFrom * pageSize;
                    if (currentFrom * pageSize >= totalSize) {
                        cur2 = totalSize;
                    }
                    pageInfo.textContent = `找到相关结果 约  ${totalSize}  个， 当前显示： ${cur1} - ${cur2} `;
                    pagination.style.display = 'flex';
                })
                .catch(error => {
                    console.error('文件查询失败', error);
                });
        });

        prevPageButton.addEventListener('click', function () {
            if (currentFrom > 1) {
                currentFrom --;
                fromInput.value = currentFrom;
                searchForm.dispatchEvent(new Event('submit'));
            }
        });

        nextPageButton.addEventListener('click', function () {
            if (currentFrom * pageSize >= totalSize) {
                alert('没有更多数据了');
                return;
            }
            currentFrom ++;
            fromInput.value = currentFrom;
            searchForm.dispatchEvent(new Event('submit'));
        });
    });
}