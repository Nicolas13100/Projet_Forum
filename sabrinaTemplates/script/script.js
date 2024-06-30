// Elements
//recuperer l'html
const html = document.querySelector('html');
const toggleThemeBtn = document.querySelector('.header__theme-button');
const posts = document.querySelectorAll('.post');
const postsContent = document.querySelectorAll('.post__content');

// ===================================
// Parcourir chaque POSTS
posts.forEach((post) => {
  // Sélection bouton à l'intérieur de cet article
  const button = post.querySelector('.comment_btn');
  const commentContainer = post.querySelector('.comment_section');
const hiddenBtn = post.querySelector('.hidden_comment');
const textPost =post.querySelector('.post-text');
const mediaPost=post.querySelector('.post__medias')
const categoryPost =post.querySelector('.post-category');
const Btn_moreOption=post.querySelector('.post__more-options');

//fonction affcher les details d'un post


  // Fonction pour afficher la section de commentaires
  function displayCommentSection() {
    commentContainer.style.display = 'flex';
    //desactiber l'overflow de la page;
    html.style.overflow = 'hidden';
    categoryPost.style.display = 'none';
    textPost.style.display = 'none';
    mediaPost.style.display = 'none';
  }

  // Fonction pour masquer la section de commentaires
  function hideCommentSection() {
    html.style.overflow = '';
    commentContainer.style.display = 'none';
    categoryPost.style.display = 'flex';
    textPost.style.display = 'flex';
    mediaPost.style.display = 'flex';
  }

  // gestionnaire pour afficher/masquerla section de commentaires 
  button.addEventListener('click', displayCommentSection);
  hiddenBtn.addEventListener('click', hideCommentSection);
});

// ===================================
// DARK/LIGHT THEME
// Set initial theme from LocalStorage
document.onload = setInitialTheme(localStorage.getItem('theme'));
function setInitialTheme(themeKey) {
  if (themeKey === 'dark') {
    document.documentElement.classList.add('darkTheme');
  } else {
    document.documentElement.classList.remove('darkTheme');
  }
}

// Toggle theme button
toggleThemeBtn.addEventListener('click', () => {
  // Toggle root class
  document.documentElement.classList.toggle('darkTheme');

  // Saving current theme on LocalStorage
  if (document.documentElement.classList.contains('darkTheme')) {
    localStorage.setItem('theme', 'dark');
  } else {
    localStorage.setItem('theme', 'light');
  }
});


//repondre a un message
document.addEventListener('DOMContentLoaded', function () {
  const replyButtons = document.querySelectorAll('.reply_btn');
  const commentInput = document.getElementById('commentInput');

  replyButtons.forEach(button => {
    button.addEventListener('click', function () {
      const userName = button.closest('.comment_user').querySelector('.comment_user_name a').textContent;
      commentInput.value = `@${userName.trim()} `;
      commentInput.focus();
    });
  });
});

