// 上传页面的JavaScript功能
document.addEventListener('DOMContentLoaded', function() {
    // 获取DOM元素
    const dropArea = document.getElementById('drop-area');
    const fileInput = document.getElementById('file-input');
    const uploadPreview = document.getElementById('upload-preview');
    const previewImage = document.getElementById('preview-image');
    const fileName = document.getElementById('file-name');
    const fileSize = document.getElementById('file-size');
    const uploadButton = document.getElementById('upload-button');
    const cancelButton = document.getElementById('cancel-button');
    const uploadResult = document.getElementById('upload-result');
    const uploadContainer = document.getElementById('upload-container');
    const imageUrl = document.getElementById('image-url');
    const markdownUrl = document.getElementById('markdown-url');
    const copyUrl = document.getElementById('copy-url');
    const copyMarkdown = document.getElementById('copy-markdown');
    const resultImagePreview = document.getElementById('result-image-preview');
    const uploadAnother = document.getElementById('upload-another');
    const uploadError = document.getElementById('upload-error');
    const errorMessage = document.getElementById('error-message');
    const tryAgain = document.getElementById('try-again');

    // 当前选择的文件
    let currentFile = null;

    // 阻止默认拖放行为
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        dropArea.addEventListener(eventName, preventDefaults, false);
    });

    function preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    // 高亮拖放区域
    ['dragenter', 'dragover'].forEach(eventName => {
        dropArea.addEventListener(eventName, highlight, false);
    });

    ['dragleave', 'drop'].forEach(eventName => {
        dropArea.addEventListener(eventName, unhighlight, false);
    });

    function highlight() {
        dropArea.classList.add('dragover');
    }

    function unhighlight() {
        dropArea.classList.remove('dragover');
    }

    // 处理拖放文件
    dropArea.addEventListener('drop', handleDrop, false);

    function handleDrop(e) {
        const dt = e.dataTransfer;
        const files = dt.files;

        if (files.length > 0) {
            handleFiles(files[0]);
        }
    }

    // 处理文件选择
    fileInput.addEventListener('change', function() {
        if (this.files.length > 0) {
            handleFiles(this.files[0]);
        }
    });

    // 处理选择的文件
    function handleFiles(file) {
        // 检查是否为图片
        if (!file.type.match('image.*')) {
            showError('请选择图片文件');
            return;
        }

        // 检查文件大小
        if (file.size > 10 * 1024 * 1024) { // 10MB
            showError('图片大小不能超过10MB');
            return;
        }

        currentFile = file;

        // 显示预览
        const reader = new FileReader();
        reader.onload = function(e) {
            previewImage.src = e.target.result;
            fileName.textContent = `文件名: ${file.name}`;
            fileSize.textContent = `大小: ${formatFileSize(file.size)}`;
            uploadPreview.style.display = 'block';
            dropArea.style.display = 'none';
        };
        reader.readAsDataURL(file);
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

    // 上传按钮点击事件
    uploadButton.addEventListener('click', uploadFile);

    // 上传文件
    function uploadFile() {
        if (!currentFile) {
            return;
        }

        const formData = new FormData();
        formData.append('image', currentFile);

        // 显示上传中状态
        uploadButton.disabled = true;
        uploadButton.textContent = '上传中...';

        fetch('/upload', {
            method: 'POST',
            body: formData
        })
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || '上传失败');
                });
            }
            return response.json();
        })
        .then(data => {
            // 上传成功
            showResult(data);
        })
        .catch(error => {
            // 上传失败
            showError(error.message);
        })
        .finally(() => {
            // 重置上传按钮状态
            uploadButton.disabled = false;
            uploadButton.textContent = '上传';
        });
    }

    // 显示上传结果
    function showResult(data) {
        // 设置结果页面
        const baseUrl = window.location.origin;
        const imageFullUrl = data.url;
        
        imageUrl.value = imageFullUrl;
        markdownUrl.value = `![${data.filename || 'image'}](${imageFullUrl})`;
        resultImagePreview.src = imageFullUrl;
        
        // 显示结果页面
        uploadContainer.style.display = 'none';
        uploadResult.style.display = 'block';
        uploadError.style.display = 'none';
    }

    // 显示错误信息
    function showError(message) {
        errorMessage.textContent = message;
        uploadContainer.style.display = 'none';
        uploadResult.style.display = 'none';
        uploadError.style.display = 'block';
    }

    // 取消上传
    cancelButton.addEventListener('click', resetUpload);

    // 继续上传
    uploadAnother.addEventListener('click', resetUpload);

    // 重试
    tryAgain.addEventListener('click', resetUpload);

    // 重置上传界面
    function resetUpload() {
        currentFile = null;
        fileInput.value = '';
        previewImage.src = '';
        uploadPreview.style.display = 'none';
        dropArea.style.display = 'block';
        uploadContainer.style.display = 'block';
        uploadResult.style.display = 'none';
        uploadError.style.display = 'none';
    }

    // 复制URL
    copyUrl.addEventListener('click', function() {
        copyToClipboard(imageUrl);
    });

    // 复制Markdown
    copyMarkdown.addEventListener('click', function() {
        copyToClipboard(markdownUrl);
    });

    // 复制到剪贴板
    function copyToClipboard(input) {
        input.select();
        document.execCommand('copy');
        
        // 显示复制成功提示
        const originalText = this.textContent;
        this.textContent = '已复制';
        
        setTimeout(() => {
            this.textContent = originalText;
        }, 1500);
    }

    // 支持粘贴上传
    document.addEventListener('paste', function(e) {
        const items = (e.clipboardData || e.originalEvent.clipboardData).items;
        
        for (let i = 0; i < items.length; i++) {
            if (items[i].type.indexOf('image') !== -1) {
                const file = items[i].getAsFile();
                handleFiles(file);
                break;
            }
        }
    });
});