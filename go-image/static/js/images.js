document.addEventListener('DOMContentLoaded', function() {
    // 获取所有复制按钮
    const copyButtons = document.querySelectorAll('.copy-btn');
    const modalCopyUrlBtn = document.getElementById('modal-copy-url');
    const modalCopyMarkdownBtn = document.getElementById('modal-copy-markdown');
    const imageModal = document.getElementById('image-modal');
    const modalImagePreview = document.getElementById('modal-image-preview');
    const modalImageName = document.getElementById('modal-image-name');
    const modalImageDate = document.getElementById('modal-image-date');
    const modalImageUrl = document.getElementById('modal-image-url');
    const modalMarkdownUrl = document.getElementById('modal-markdown-url');
    const deleteModal = document.getElementById('delete-modal');
    const closeButtons = document.querySelectorAll('.close');
    const confirmDeleteBtn = document.getElementById('confirm-delete');
    const cancelDeleteBtn = document.getElementById('cancel-delete');
    let currentImageId = null;

    // 复制链接功能
    function copyToClipboard(text) {
        // 创建临时输入框
        const input = document.createElement('input');
        input.style.position = 'fixed';
        input.style.opacity = 0;
        input.value = text;
        document.body.appendChild(input);
        input.select();
        input.setSelectionRange(0, 99999);
        
        try {
            // 执行复制命令
            document.execCommand('copy');
            alert('链接已复制到剪贴板！');
        } catch (err) {
            console.error('复制失败:', err);
            alert('复制失败，请手动复制链接。');
        }
        
        document.body.removeChild(input);
    }

    // 获取完整的URL
    function getFullUrl(path) {
        const baseUrl = window.location.protocol + '//' + window.location.host;
        return baseUrl + path;
    }

    // 为所有复制按钮添加点击事件
    copyButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            e.preventDefault();
            const path = this.getAttribute('data-url');
            const url = getFullUrl(path);
            copyToClipboard(url);
        });
    });

    // 查看图片详情
    const viewButtons = document.querySelectorAll('.view-btn');
    viewButtons.forEach(button => {
        button.addEventListener('click', function() {
            const imageId = this.getAttribute('data-id');
            const imageCard = this.closest('.image-card');
            const imageName = imageCard.querySelector('.image-name').textContent;
            const imageDate = imageCard.querySelector('.image-date').textContent;
            
            currentImageId = imageId;
            modalImagePreview.src = `/i/${imageId}`;
            modalImageName.textContent = imageName;
            modalImageDate.textContent = imageDate;
            
            const imageUrl = getFullUrl(`/i/${imageId}`);
            modalImageUrl.value = imageUrl;
            modalMarkdownUrl.value = `![${imageName}](${imageUrl})`;
            
            imageModal.style.display = 'block';
        });
    });

    // 模态框中的复制功能
    modalCopyUrlBtn.addEventListener('click', function() {
        copyToClipboard(modalImageUrl.value);
    });

    modalCopyMarkdownBtn.addEventListener('click', function() {
        copyToClipboard(modalMarkdownUrl.value);
    });

    // 删除图片功能
    const deleteButtons = document.querySelectorAll('.delete-btn');
    deleteButtons.forEach(button => {
        button.addEventListener('click', function() {
            currentImageId = this.getAttribute('data-id');
            deleteModal.style.display = 'block';
        });
    });

    // 确认删除
    confirmDeleteBtn.addEventListener('click', function() {
        if (currentImageId) {
            fetch(`/images/${currentImageId}`, {
                method: 'DELETE',
            })
            .then(response => response.json())
            .then(data => {
                if (data.message === '删除成功') {
                    const imageCard = document.querySelector(`.image-card[data-id="${currentImageId}"]`);
                    if (imageCard) {
                        imageCard.remove();
                    }
                    deleteModal.style.display = 'none';
                    // 检查是否还有图片
                    const remainingImages = document.querySelectorAll('.image-card');
                    if (remainingImages.length === 0) {
                        location.reload(); // 刷新页面显示"暂无图片"提示
                    }
                } else {
                    alert('删除失败：' + data.error);
                }
            })
            .catch(error => {
                console.error('删除请求失败:', error);
                alert('删除失败，请稍后重试');
            });
        }
    });

    // 关闭模态框
    closeButtons.forEach(button => {
        button.addEventListener('click', function() {
            this.closest('.modal').style.display = 'none';
        });
    });

    cancelDeleteBtn.addEventListener('click', function() {
        deleteModal.style.display = 'none';
    });

    // 点击模态框外部关闭
    window.addEventListener('click', function(event) {
        if (event.target.classList.contains('modal')) {
            event.target.style.display = 'none';
        }
    });
});