
    document.addEventListener('DOMContentLoaded', (event) => {
    const posts = document.querySelectorAll('.post');

    posts.forEach(post => {
    const creationDate = post.getAttribute('data-creation-date');
    const creationMoment = moment(creationDate);
    const now = moment();
    const timeDifference = moment.duration(now.diff(creationMoment));

    let timeAgo = '';
    switch (true) {
    case timeDifference.asSeconds() < 60:
    timeAgo = `${Math.floor(timeDifference.asSeconds())} seconds ago`;
    break;
    case timeDifference.asMinutes() < 60:
    timeAgo = `${Math.floor(timeDifference.asMinutes())} minutes ago`;
    break;
    case timeDifference.asHours() < 24:
    timeAgo = `${Math.floor(timeDifference.asHours())} hours ago`;
    break;
    case timeDifference.asDays() < 30:
    timeAgo = `${Math.floor(timeDifference.asDays())} days ago`;
    break;
    case timeDifference.asMonths() < 12:
    timeAgo = `${Math.floor(timeDifference.asMonths())} months ago`;
    break;
    default:
    timeAgo = `${Math.floor(timeDifference.asYears())} years ago`;
    break;
    }


    post.querySelector('.post__date-time').innerText = timeAgo;
});
});