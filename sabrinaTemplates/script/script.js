// Elements
const toggleThemeBtn = document.querySelector('.header__theme-button');
const storiesContent = document.querySelector('.stories__content');
const storiesLeftButton = document.querySelector('.stories__left-button');
const storiesRightButton = document.querySelector('.stories__right-button');
const posts = document.querySelectorAll('.post');
const postsContent = document.querySelectorAll('.post__content');

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


// Parcourez chaque POSTS
posts.forEach((post) => {
  // Sélection bouton à l'intérieur de cet article
  const button = post.querySelector('.comment_btn');
  const commentContainer = post.querySelector('.comment_section');
const hiddenBtn = post.querySelector('.hidden_comment');
const textPost =post.querySelector('.post-text');
const mediaPost=post.querySelector('.post__medias')
const categoryPost =post.querySelector('.post-category');
  // Fonction pour afficher la section de commentaires
  function displayCommentSection() {
    commentContainer.style.display = 'flex';
    categoryPost.style.display = 'none';
    textPost.style.display = 'none';
    mediaPost.style.display = 'none';
  }

  // Fonction pour masquer la section de commentaires
  function hideCommentSection() {
    commentContainer.style.display = 'none';
    categoryPost.style.display = 'flex';
    textPost.style.display = 'flex';
    mediaPost.style.display = 'flex';
  }

  // gestionnaire pour afficher/masquerla section de commentaires 
  button.addEventListener('click', displayCommentSection);
  hiddenBtn.addEventListener('click', hideCommentSection);
});

