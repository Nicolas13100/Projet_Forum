const express = require('express');
const router = express.Router();

router.get('/', (req, res) => {
    res.render('home', { title: 'Home' });
});

router.get('/category', (req, res) => {
    res.render('category', { title: 'Category' });
});

router.get('/createTopic', (req, res) => {
    res.render('createTopic', { title: 'CreateTopic' });
});

router.get('/landingPage', (req, res) => {
    res.render('landingPage', { title: 'LandingPage' });
});

router.get('/login', (req, res) => {
    res.render('login', { title: 'Login' });
});

router.get('/message', (req, res) => {
    res.render('message', { title: 'Message' });
});

router.get('/profil', (req, res) => {
    res.render('profil', { title: 'Profil' });
});

router.get('/search', (req, res) => {
    res.render('search', { title: 'Search' });
});

module.exports = router;
