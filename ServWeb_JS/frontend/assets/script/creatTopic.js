function previewImage(event) {
    const file = event.target.files[0];
    const preview = document.getElementById('preview-image');

    const reader = new FileReader();
    reader.onload = function(event) {
        preview.src = event.target.result;
    };

    if (file) {
        reader.readAsDataURL(file);
    } else {
        preview.src = "#";
    }
}