function main() {
    document.addEventListener('DOMContentLoaded', function () {
        // 设置main的高度
        const topBar = document.getElementById('top-bar');
        const mainContent = document.getElementById('main-content');
        const topBarHeight = topBar.offsetHeight;
        mainContent.style.marginTop = topBarHeight + 'px';

        const uploadLink = document.getElementById('uploadLink');
        const searchLink = document.getElementById('searchLink');
        const deleteLink = document.getElementById('deleteLink');
        const uploadSection = document.getElementById('uploadSection');
        const searchSection = document.getElementById('searchSection');
        const deleteSection = document.getElementById('deleteSection');

        searchLink.classList.add('active');
        searchSection.classList.add('active');

        uploadLink.addEventListener('click', function (e) {
            e.preventDefault();
            removeActiveClass();
            this.classList.add('active');
            hideAllSections();
            uploadSection.classList.add('active');
        });
        deleteLink.addEventListener('click', function (e) {
            e.preventDefault();
            removeActiveClass();
            this.classList.add('active');
            hideAllSections();
            deleteSection.classList.add('active');
        })

        searchLink.addEventListener('click', function (e) {
            e.preventDefault();
            removeActiveClass();
            this.classList.add('active');
            hideAllSections();
            searchSection.classList.add('active');
        });
        function removeActiveClass() {
            const sidebarItems = document.querySelectorAll('#sidebar ul li a');
            sidebarItems.forEach(item => {
                item.classList.remove('active');
            });
        }
        function hideAllSections() {
            const sections = document.querySelectorAll('.content-section');
            sections.forEach(section => {
                section.classList.remove('active');
            });
        }
    });

}


function logout() {
    document.getElementById('logoutButton').addEventListener('click', function (){
        const sessionId = localStorage.getItem('session_id');
        if (!sessionId) {
            window.location.href = 'login.html';
            alert('未登录');
            return;
        }

        fetch('http://localhost/v1/user/logout', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Session-ID': sessionId
            }
        })
            .then(handleResponseStatus)
            .then(data => {
                if (data.code === 0) {
                    // 退出成功，跳转到登录页面
                    window.location.href = 'login.html';
                    console.log('Logout successful');
                    // 删除 localStorage 中的 session_id
                    localStorage.removeItem('session_id');
                } else {
                    alert('退出失败: ' + data.message);
                }
            })
            .catch(error => {
                console.error('退出请求出错:', error);
                alert('退出时发生错误，请稍后重试');
            });
    });
}