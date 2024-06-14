const express = require('express');
const axios = require('axios');
const router = express.Router();

router.get('/', (req, res) => {
    const data = { title: 'Home' };
    res.render('home', data);

});

router.get('/category', (req, res) => {
    const data = { title: 'Category' };
    res.render('category', data);
});

router.get('/createTopic', (req, res) => {
    const data = { title: 'CreateTopic' };
    res.render('createTopic', data);
});

router.get('/landingPage', (req, res) => {
    const data = { title: 'LandingPage' };
    res.render('landingPage', data);
});

router.get('/login', (req, res) => {
    const data = { title: 'Login' };
    res.render('login', data);
});

router.get('/message', (req, res) => {
    const data = { title: 'Message' };
    res.render('message', data);

});

router.get('/profil', (req, res) => {
    const data = { title: 'Profil' };
    res.render('profil', data);
});

router.get('/search', (req, res) => {
    const data = { title: 'Search' };
    res.render('search', data);
});






