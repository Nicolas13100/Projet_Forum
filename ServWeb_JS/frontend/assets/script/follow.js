document.addEventListener('DOMContentLoaded', function () {
    const followButtons = document.querySelectorAll('.side-menu__suggestion-button');

    followButtons.forEach(button => {
        button.addEventListener('click', function () {
            const userIdToFollow = this.getAttribute('data-user-id');
            const action = this.getAttribute('data-action');
            const loggedUserId = button.getAttribute('data-logged-user-id');

            // Check if the user is logged in
            if (loggedUserId) {
                // Define the URL based on the action (follow or unfollow)
                const url = action === 'follow' ? `/follow/${userIdToFollow}/${loggedUserId}` : `/unfollow/${userIdToFollow}/${loggedUserId}`;

                axios.post(url)
                    .then(response => {
                        // Handle success response
                        console.log(response.data);
                        if (response.data.success) {
                            // Update button content and action based on the current action
                            if (action === 'follow') {
                                this.textContent = 'Unfollow';
                                this.setAttribute('data-action', 'unfollow');
                            } else {
                                this.textContent = 'Follow';
                                this.setAttribute('data-action', 'follow');
                            }
                        }
                    })
                    .catch(error => {
                        // Handle error response
                        console.error('Error during follow/unfollow action:', error);
                    });
            } else {
                window.location.href = '/login';
            }
        });
    });
});
