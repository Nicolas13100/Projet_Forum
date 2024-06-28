
    document.addEventListener('DOMContentLoaded', (event) => {
    const posts = document.querySelectorAll('.post');

    posts.forEach(post => {
    const creationDate = post.getAttribute('data-creation-date');
    const creationMoment = moment(creationDate);
    const now = moment();
    const timeDifference = moment.duration(now.diff(creationMoment));

    let timeAgo = '';
    if (timeDifference.asMinutes() < 60) {
    timeAgo = `${Math.floor(timeDifference.asMinutes())} minutes ago`;
} else if (timeDifference.asHours() < 24) {
    timeAgo = `${Math.floor(timeDifference.asHours())} hours ago`;
} else {
    timeAgo = `${Math.floor(timeDifference.asDays())} days ago`;
}

    post.querySelector('.post__date-time').innerText = timeAgo;
});
});