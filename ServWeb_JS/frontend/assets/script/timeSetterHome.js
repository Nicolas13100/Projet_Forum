
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
        if (timeDifference.asSeconds() < 1) {
        timeAgo = 'Just now';
    }
    timeAgo = `${Math.floor(timeDifference.asSeconds())} seconds ago`;
    break;
    case timeDifference.asMinutes() < 60:
        if (timeDifference.asMinutes() < 2) {
        timeAgo = '1 minute ago';
    }
    timeAgo = `${Math.floor(timeDifference.asMinutes())} minutes ago`;
    break;
    case timeDifference.asHours() < 24:
        if (Math.floor(timeDifference.asHours()) < 2) {
        timeAgo = '1 hour ago';
    }
    timeAgo = `${Math.floor(timeDifference.asHours())} hours ago`;
    break;
    case timeDifference.asDays() < 30:
        if (timeDifference.asDays() < 2) {
        timeAgo = '1 day ago';
    }
    timeAgo = `${Math.floor(timeDifference.asDays())} days ago`;
    break;
    case timeDifference.asMonths() < 12:
        if (Math.floor(timeDifference.asMonths()) < 2) {
        timeAgo = '1 month ago';
    }
    timeAgo = `${Math.floor(timeDifference.asMonths())} months ago`;
    break;
    default:
        if (timeDifference.asYears() < 1) {
        timeAgo = '1 year ago';
    } else {
    timeAgo = `${Math.floor(timeDifference.asYears())} years ago`;
    break;
    }
    }
    post.querySelector('.post__date-time').innerText = timeAgo;
});
});