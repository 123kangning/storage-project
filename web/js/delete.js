function deleteBox() {
    document.addEventListener('DOMContentLoaded', function () {

        // 文件查询表单提交事件
        const fromInput = document.getElementById('from2');
        const prevPageButton = document.getElementById('prevPage2');
        const nextPageButton = document.getElementById('nextPage2');
        const pageInfo = document.getElementById('pageInfo2');
        const pagination = document.getElementById('pagination2');
        let currentFrom = 1;
        let totalSize = 0;
        const pageSize = 10;
        const resetButton = document.querySelector('#deleteForm button[type="reset"]');
        resetButton.addEventListener('click', function () {
            // Clear the table body content
            const tableBody = document.getElementById('fileTable2').getElementsByTagName('tbody')[0];
            tableBody.innerHTML = '';
            pagination.style.display = 'none';
            currentFrom = 1;
        });

        const deleteForm = document.getElementById('deleteForm');
        deleteForm.addEventListener('submit', function (e) {
            e.preventDefault();
            const sessionId = localStorage.getItem('session_id');
            if (!sessionId) {
                window.location.href = 'login.html';
                alert('请先登录');
                return;
            }

            const formData = new FormData(this);
            const name = formData.get('name');
            const url = `http://localhost/v1/file/list?name=${name}&from=${currentFrom}&size=${pageSize}`;

            // 获取文件列表（与查询页面相同）
            fetch(url, {
                headers: {'X-Session-ID': sessionId}
            }).then(handleResponseStatus)
                .then(data => {
                    if (!handleResponseError(data.baseResp, '文件查询')) return;

                    const tableBody = document.getElementById('fileTable2').getElementsByTagName('tbody')[0];
                    tableBody.innerHTML = '';
                    const realData = data.data.files;
                    totalSize = data.data.total;
                    if (realData.length === 0) {
                        pagination.style.display = 'none';
                        alert('没有查询到相关文件');
                    }

                    realData.forEach(file => {
                        const row = document.createElement('tr');
                        const nameCell = document.createElement('td');
                        const sizeCell = document.createElement('td');
                        const hashCell = document.createElement('td');
                        const operationCell = document.createElement('td');

                        nameCell.textContent = file.name;
                        if (file.name.length > 20) {
                            nameCell.textContent = file.name.substring(0, 20) + '...';
                        }
                        sizeCell.textContent = formatFileSize(file.size);
                        hashCell.textContent = file.hash.substring(0, 5) + '...'+ file.hash.substring(file.hash.length - 5);

                        // 删除按钮
                        const deleteBtn = document.createElement('button');
                        deleteBtn.textContent = '删除';
                        deleteBtn.classList.add('btn', 'btn-danger');
                        deleteBtn.addEventListener('click', function() {
                            // 两步确认
                            if (!confirm('此操作不可撤销！再次确认删除？')) return;

                            // 调用删除接口
                            fetch(`http://localhost/v1/file/del?hash=${file.hash}`, {
                                method: 'DELETE',
                                headers: {'X-Session-ID': sessionId}
                            })
                                .then(handleResponseStatus)
                                .then(response => {
                                    if (handleResponseError(response, '文件删除')) {
                                        row.remove(); // 删除表格行
                                        alert('文件删除成功');
                                    }
                                })
                                .catch(error => console.error('删除失败:', error));
                        });

                        operationCell.appendChild(deleteBtn);

                        row.appendChild(nameCell);
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
                });
        });
        prevPageButton.addEventListener('click', function () {
            if (currentFrom > 1) {
                currentFrom --;
                fromInput.value = currentFrom;
                deleteForm.dispatchEvent(new Event('submit'));
            }
        });


        nextPageButton.addEventListener('click', function () {
            if (currentFrom * pageSize >= totalSize) {
                alert('没有更多数据了');
                return;
            }
            currentFrom ++;
            fromInput.value = currentFrom;
            deleteForm.dispatchEvent(new Event('submit'));
        });
    });
}