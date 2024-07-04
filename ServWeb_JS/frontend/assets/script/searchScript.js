document.getElementById('search-form').addEventListener('submit', function(event) {
    const query = document.getElementById('search-input').value;
    if (!query) {
        event.preventDefault(); // Prevent form submission if the input is empty
    }
});

document.getElementById('search-input').addEventListener('keypress', function(event) {
    if (event.key === 'Enter') {
        const query = document.getElementById('search-input').value;
        if (!query) {
            event.preventDefault(); // Prevent form submission if the input is empty
        }
    }
});