document.addEventListener('DOMContentLoaded', function () {
    const likeButtons = document.querySelectorAll('.like-button');

    likeButtons.forEach(button => {
        button.addEventListener('click', async function () {
            const topicId = this.dataset.topicId;
            const userId = this.dataset.userId;
            const isLiked = this.classList.contains('liked');
            const url = isLiked ? `/api/dislike/${topicId}/${userId}` : `/api/like/${topicId}/${userId}`;

            try {
                const response = await fetch(url, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });

                if (response.ok) {
                    this.classList.toggle('liked'); // Toggle the 'liked' class
                    if (this.classList.contains('liked')) {
                        this.querySelector('svg').style.fill = 'red';
                    } else {
                        this.querySelector('svg').style.fill = '';
                    }
                } else {
                    console.error('Failed to like/dislike the topic', response.statusText);
                }
            } catch (error) {
                console.error('An error occurred', error);
            }
        });
    });
});
