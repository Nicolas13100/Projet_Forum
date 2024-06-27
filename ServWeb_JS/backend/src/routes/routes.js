const express = require('express');
const router = express.Router();
const axios = require('axios');
const FormData = require('form-data');

router.get('/', (req, res) => {
    res.render('home', { title: 'Home' });
});

router.get('/category', (req, res) => {
    res.render('category', { title: 'Category' });
});

router.get('/createTopic', (req, res) => {
    res.render('createTopic', { title: 'CreateTopic' });
});

router.get('/home', async (req, res) => {
    const url = 'http://localhost:8080/api/getHome/1/10';
    const TagUrl = 'http://localhost:8080/api/getTopicTag';
    let response;
    try {
        response = await axios.get(url, {});

        // Array to store promises for fetching tag data
        const tagPromises = [];

        // Iterate through topics and fetch tag data
        for (const topic of response.data.resp.topics) {
            const topic_id = topic.topic_id;
            const fetchTagPromise = await axios.get(`${TagUrl}/${topic_id}`, {}); // Adjust URL as per your API endpoint
            tagPromises.push(fetchTagPromise.data.data);
        }

        // Wait for all tag data requests to resolve
        const tagResponses = await Promise.all(tagPromises);

        // Merge tag data into respective topics
        tagResponses.forEach((tagResponse, index) => {
            response.data.resp.topics[index].tags = tagPromises[index];
        });

    } catch (error) {
        console.log(error);
    }

    const data = {
        topics: response.data.resp.topics,
        user: {
            profilePic: ""
        }
    };

    res.render('home', data);
});

router.get('/login', (req, res) => {
    res.render('login', { message: null });
});

// POST route to handle form submission
router.post('/loginUser', async (req, res) => {
    const { username, password } = req.body;

    // Axios POST request to another server (http://localhost:8080/api/login)
    const url = 'http://localhost:8080/api/login';
    const data = new URLSearchParams();
    data.append('username', username);
    data.append('password', password);

    try {
        const response = await axios.post(url, data, {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        });

        const token = response.data.token;

        // Set cookie named 'token' with the received token
        res.cookie('token', token, { httpOnly: true });

        // Redirect to /home
        res.redirect('/home');
    } catch (error) {
        console.error('Error sending login request:', error.message);
        res.status(500).send('Error sending login request.'); // Handle error
    }
});

// POST route to handle form submission
router.post('/register', async (req, res) => {
    const { username, mail, password } = req.body;
    const url = 'http://localhost:8080/api/register';
    const data = new FormData();
    data.append('username', username);
    data.append('password', password);
    data.append('mail', mail);
    console.log(username, password, mail);

    try {
        await axios.post(url, data, {
            headers: {
                ...data.getHeaders()
            }
        });

        res.redirect('/home');

    } catch (error) {
        if (error.response) {
            // The request was made and the server responded with a status code that falls out of the range of 2xx
            console.error('Error response data:', error.response.data);
            console.error('Error response status:', error.response.status);
            console.error('Error response headers:', error.response.headers);
            if (error.response.data){
                res.render('login', { message: error.response.data.message});
            }else {
                res.status(error.response.status).send(error.response.data);
            }
        } else if (error.request) {
            // The request was made but no response was received
            console.error('Error request:', error.request);
            res.status(500).send('No response received from the server.');
        } else {
            // Something happened in setting up the request that triggered an Error
            console.error('Error message:', error.message);
            res.status(500).send('Error sending register request.');
        }
    }
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
