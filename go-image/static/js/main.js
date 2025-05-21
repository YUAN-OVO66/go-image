// 主要JavaScript功能
document.addEventListener('DOMContentLoaded', function() {
    // 图片列表页面功能
    initImageListPage();
    // 初始化存储使用情况图表
    initStorageChart();
    
    // 初始化存储使用情况图表
    function initStorageChart() {
        const storageChartContainer = document.querySelector('.storage-chart-container');
        if (!storageChartContainer) return;

        // 移除旧的canvas元素
        const oldCanvas = document.getElementById('storage-chart');
        if (oldCanvas) {
            oldCanvas.remove();
        }

        const usedStorage = parseFloat(document.getElementById('used-storage').textContent);
        const totalStorage = parseFloat(document.getElementById('total-storage').textContent);
        
        // 计算使用百分比
        const usagePercentage = totalStorage > 0 ? (usedStorage / totalStorage) * 100 : 0;
        
        // 创建进度条容器
        const progressContainer = document.createElement('div');
        progressContainer.className = 'progress-container';
        
        // 创建进度条
        const progressBar = document.createElement('div');
        progressBar.className = 'progress-bar';
        progressBar.style.width = `${usagePercentage}%`;
        
        // 创建百分比文本
        const percentageText = document.createElement('div');
        percentageText.className = 'percentage-text';
        percentageText.textContent = `${usagePercentage.toFixed(1)}%`;
        
        // 组装进度条
        progressContainer.appendChild(progressBar);
        progressContainer.appendChild(percentageText);
        
        // 添加到容器
        storageChartContainer.appendChild(progressContainer);


        // 更新显示的存储信息
        document.getElementById('used-storage').textContent = formatFileSize(usedStorage);
        document.getElementById('total-storage').textContent = formatFileSize(totalStorage);
    }

    // 初始化图片列表页面功能
    function initImageListPage() {
        // 获取DOM元素
        const imageGrid = document.getElementById('image-grid');
        if (!imageGrid) return; // 不在图片列表页面
        
        const searchInput = document.getElementById('search-input');
        const searchButton = document.getElementById('search-button');
        const imageModal = document.getElementById('image-modal');
        const deleteModal = document.getElementById('delete-modal');
        const modalImagePreview = document.getElementById('modal-image-preview');
        const modalImageName = document.getElementById('modal-image-name');
        const modalImageDate = document.getElementById('modal-image-date');
        const modalImageSize = document.getElementById('modal-image-size');
        const modalImageUrl = document.getElementById('modal-image-url');
        const modalMarkdownUrl = document.getElementById('modal-markdown-url');
        const modalCopyUrl = document.getElementById('modal-copy-url');
        const modalCopyMarkdown = document.getElementById('modal-copy-markdown');
        const confirmDelete = document.getElementById('confirm-delete');
        const cancelDelete = document.getElementById('cancel-delete');
        
        // 当前选中的图片ID
        let currentImageId = null;
        
        // 关闭模态框
        const closeButtons = document.querySelectorAll('.close');
        closeButtons.forEach(button => {
            button.addEventListener('click', function() {
                imageModal.style.display = 'none';
                deleteModal.style.display = 'none';
            });
        });
        
        // 点击模态框外部关闭
        window.addEventListener('click', function(e) {
            if (e.target === imageModal) {
                imageModal.style.display = 'none';
            }
            if (e.target === deleteModal) {
                deleteModal.style.display = 'none';
            }
        });
        
        // 查看图片
        const viewButtons = document.querySelectorAll('.view-btn');
        viewButtons.forEach(button => {
            button.addEventListener('click', function() {
                const imageId = this.getAttribute('data-id');
                showImageDetails(imageId);
            });
        });
        
        // 复制链接
        const copyButtons = document.querySelectorAll('.copy-btn');
        copyButtons.forEach(button => {
            button.addEventListener('click', function() {
                const url = this.getAttribute('data-url');
                copyToClipboard(url, this);
            });
        });
        
        // 删除图片
        const deleteButtons = document.querySelectorAll('.delete-btn');
        deleteButtons.forEach(button => {
            button.addEventListener('click', function() {
                const imageId = this.getAttribute('data-id');
                showDeleteConfirmation(imageId);
            });
        });
        
        // 确认删除
        confirmDelete.addEventListener('click', function() {
            if (currentImageId) {
                deleteImage(currentImageId);
            }
        });
        
        // 取消删除
        cancelDelete.addEventListener('click', function() {
            deleteModal.style.display = 'none';
        });
        
        // 搜索功能
        if (searchButton && searchInput) {
            searchButton.addEventListener('click', function() {
                searchImages(searchInput.value);
            });
            
            searchInput.addEventListener('keypress', function(e) {
                if (e.key === 'Enter') {
                    searchImages(searchInput.value);
                }
            });
        }
        
        // 模态框中的复制按钮
        if (modalCopyUrl) {
            modalCopyUrl.addEventListener('click', function() {
                copyToClipboard(modalImageUrl.value, this);
            });
        }
        
        if (modalCopyMarkdown) {
            modalCopyMarkdown.addEventListener('click', function() {
                copyToClipboard(modalMarkdownUrl.value, this);
            });
        }
        
        // 显示图片详情
        function showImageDetails(imageId) {
            fetch(`/images/${imageId}`)
                .then(response => {
                    if (!response.ok) {
                        throw new Error('获取图片详情失败');
                    }
                    return response.json();
                })
                .then(data => {
                    // 设置模态框内容
                    const baseUrl = window.location.origin;
                    const imageFullUrl = `${baseUrl}/i/${data.id}`;
                    
                    modalImagePreview.src = imageFullUrl;
                    modalImageName.textContent = data.filename;
                    modalImageDate.textContent = `上传时间: ${formatDate(data.uploaded_at)}`;
                    modalImageSize.textContent = `文件大小: ${formatFileSize(data.size)}`;
                    modalImageUrl.value = imageFullUrl;
                    modalMarkdownUrl.value = `![${data.filename}](${imageFullUrl})`;
                    
                    // 显示模态框
                    imageModal.style.display = 'block';
                })
                .catch(error => {
                    alert(error.message);
                });
        }
        
        // 显示删除确认
        function showDeleteConfirmation(imageId) {
            currentImageId = imageId;
            deleteModal.style.display = 'block';
        }
        
        // 删除图片
        function deleteImage(imageId) {
            fetch(`/images/${imageId}`, {
                method: 'DELETE'
            })
                .then(response => response.json())
                .then(data => {
                    if (data.error) {
                        throw new Error(data.error);
                    }
                    
                    // 关闭模态框
                    deleteModal.style.display = 'none';
                    
                    // 从页面移除图片卡片
                    const imageCard = document.querySelector(`.image-card[data-id="${imageId}"]`);
                    if (imageCard) {
                        imageCard.remove();
                    }
                    
                    // 如果没有图片了，显示空状态
                    if (imageGrid.children.length === 0) {
                        const noImages = document.createElement('div');
                        noImages.className = 'no-images';
                        noImages.innerHTML = `
                            <p>暂无图片，去上传一些吧！</p>
                            <a href="/upload" class="btn primary">上传图片</a>
                        `;
                        imageGrid.appendChild(noImages);
                    }
                })
                .catch(error => {
                    alert(error.message || '删除图片失败');
                });
        }
        
        // 搜索图片
        function searchImages(query) {
            // 简单的前端过滤，实际项目中可以改为后端搜索
            const imageCards = document.querySelectorAll('.image-card');
            const normalizedQuery = query.toLowerCase().trim();
            
            let hasVisibleImages = false;
            
            imageCards.forEach(card => {
                const imageName = card.querySelector('.image-name').textContent.toLowerCase();
                if (normalizedQuery === '' || imageName.includes(normalizedQuery)) {
                    card.style.display = 'block';
                    hasVisibleImages = true;
                } else {
                    card.style.display = 'none';
                }
            });
            
            // 显示或隐藏空状态
            let noImagesElement = document.querySelector('.no-images');
            
            if (!hasVisibleImages) {
                if (!noImagesElement) {
                    noImagesElement = document.createElement('div');
                    noImagesElement.className = 'no-images';
                    noImagesElement.innerHTML = `
                        <p>没有找到匹配的图片</p>
                        <button class="btn secondary" id="reset-search">重置搜索</button>
                    `;
                    imageGrid.appendChild(noImagesElement);
                    
                    // 添加重置搜索按钮事件
                    document.getElementById('reset-search').addEventListener('click', function() {
                        searchInput.value = '';
                        searchImages('');
                    });
                }
            } else if (noImagesElement) {
                noImagesElement.remove();
            }
        }
    }
    
    // 复制到剪贴板
    function copyToClipboard(text, button) {
        // 创建临时输入框
        const tempInput = document.createElement('input');
        tempInput.value = text;
        document.body.appendChild(tempInput);
        tempInput.select();
        document.execCommand('copy');
        document.body.removeChild(tempInput);
        
        // 显示复制成功提示
        const originalText = button.textContent;
        button.textContent = '已复制';
        
        setTimeout(() => {
            button.textContent = originalText;
        }, 1500);
    }
    
    // 格式化日期
    function formatDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
        });
    }
    
    // 格式化文件大小
    function formatFileSize(bytes) {
        if (bytes < 1024) {
            return bytes + ' B';
        } else if (bytes < 1024 * 1024) {
            return (bytes / 1024).toFixed(2) + ' KB';
        } else {
            return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
        }
    }
});