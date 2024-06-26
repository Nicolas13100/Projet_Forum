const express = require('express');
const router = express.Router();
const axios = require('axios');
const FormData = require('form-data');
const multer = require('multer');
const path = require('path');
const uploadDir = path.join(__dirname, '../../../../ServWeb_JS/frontend/assets/images/TopicsImg');
// Multer configuration
const storage = multer.diskStorage({
    destination: function (req, file, cb) {
        cb(null, uploadDir);
    },
    filename: function (req, file, cb) {
        cb(null, file.fieldname + '-' + Date.now() + path.extname(file.originalname));
    }
});

const upload = multer({ storage: storage });

router.get('/', (req, res) => {
    res.redirect('/home');
});

router.get('/category', (req, res) => {
    res.render('category', { title: 'Category' });
});

router.get('/createTopic', (req, res) => {

    res.render('createTopic', { title: 'CreateTopic' });
});

router.post('/createTopic', upload.single('image'), async (req, res) => {
    const token = req.cookies.token;
    const logged = token !== undefined;
    if (!logged) {
        return res.redirect('/login');
    }

    const userIDByToken = 'http://localhost:8080/api/getUserIDByToken';
    let userID;
    try {
        const response = await axios.get(`${userIDByToken}/${token}`);
        userID = response.data.UserID;
    } catch (e) {
        console.log(e.message);
        return res.status(500).send('Error fetching user ID');
    }

    const { title, description, categories } = req.body;
    const imagePath = req.file ? `/static/images/TopicsImg/${req.file.filename}` : null;

    try {
        const createTopicUrl = "http://localhost:8080/api/createTopic";
        const response = await axios.post(createTopicUrl, {
            title,
            description,
            categories,
            imagePath,
            userID
        });

         const topicID = response.data.topic_id;
         res.redirect(`/topic/${topicID}`);
    } catch (err) {
        console.error(err);
        res.status(500).send('Error creating topic');
    }
});

router.get('/home', async (req, res) => {
    const url = 'http://localhost:8080/api/getHome/1/10';
    const TagUrl = 'http://localhost:8080/api/getTopicTag';
    const ownerUrl = 'http://localhost:8080/api/getTopicOwner';
    const likeUrl = 'http://localhost:8080/api/getLikeTopicNumber';
    const topicImgUrl = 'http://localhost:8080/api/getTopicImg'
    const randomUserUrl = 'http://localhost:8080/api/getRandomUser';
    const followUrl = 'http://localhost:8080/api/getFollowers';
    const userIDByToken = 'http://localhost:8080/api/getUserIDByToken';
    const userDataUrl = 'http://localhost:8080/api/getUser';
    const token = req.cookies.token;
    const logged = token !== undefined;
    let response;
    try {
        // Fetch the main response
        response = await axios.get(url, {});

        // Arrays to store promises for fetching tag and owner data
        const tagPromises = [];
        const ownerPromises = [];
        const numberOfLikePromises =[];
        const topicImgPromises =[];

        // Iterate through topics and create promises for fetching tag and owner data
        for (const topic of response.data.resp.topics) {
            const topic_id = topic.topic_id;

            // Create promises for fetching tag and owner data
            const fetchTagPromise = axios.get(`${TagUrl}/${topic_id}`);
            const fetchOwnerPromise = axios.get(`${ownerUrl}/${topic_id}`);
            const fetchNumberOfLikePromise = axios.get(`${likeUrl}/${topic_id}`);
            const fetchTopicImgPromise=axios.get(`${topicImgUrl}/${topic_id}`)

            // Store the promises
            tagPromises.push(fetchTagPromise);
            ownerPromises.push(fetchOwnerPromise);
            numberOfLikePromises.push(fetchNumberOfLikePromise);
            topicImgPromises.push(fetchTopicImgPromise)
        }

        // Wait for all tag and owner data requests to resolve
        const tagResponses = await Promise.all(tagPromises);
        const ownerResponses = await Promise.all(ownerPromises);
        const numberOfLikeResponses = await Promise.all(numberOfLikePromises)
        const topicImgResponses=await Promise.all(topicImgPromises)

        // Merge tag and owner data into respective topics
        tagResponses.forEach((tagResponse, index) => {
            response.data.resp.topics[index].tags = tagResponse.data.data;
        });

        ownerResponses.forEach((ownerResponse, index) => {
            response.data.resp.topics[index].owner = ownerResponse.data.UserData;
        });

        numberOfLikeResponses.forEach((likeResponse, index) => {
           response.data.resp.topics[index].numberOfLike = likeResponse.data.NumberOfLike.like_count;
        })

        topicImgResponses.forEach((imgResponse, index) => {
            response.data.resp.topics[index].imgPath = imgResponse.data.ImagePath;
        })

    } catch (error) {
        console.log(error);
    }

    let forYou
    const followerPromises = [];
    try{
        forYou = await axios.get(randomUserUrl, {});
        for (const user of forYou.data.UsersData) {
            const user_id = user.user_id
            const fetchFollowersPromise = axios.get(`${followUrl}/${user_id}`);
            followerPromises.push(fetchFollowersPromise);
        }

        const followerResponses = await Promise.all(followerPromises);

        followerResponses.forEach((followerResponses, index) => {
            forYou.data.UsersData[index].followers = followerResponses.data.FollowerData.follower_count;
        });
    }catch(e){
        console.log(e);
    }
    let user
if (logged){
    try{
        const userID = await axios.get(`${userIDByToken}/${token}`, {});
        const userData = await axios.get(`${userDataUrl}/${userID.data.UserID}`, {});
        user = userData.data.user
    }catch(e){
        console.log(e.data);
    }
}
    const data = {
        topics: response.data.resp.topics,
        logged: logged,
        user: user,
        forYou : forYou.data.UsersData
    };

    res.render('home', data);
});

router.get('/login', (req, res) => {
    res.render('login', { messageLogin: null, messageRegister : null });
});

router.get('/logout', async (req, res) => {
    const token = req.cookies.token;  // Retrieves token from cookies
    const logged = token !== undefined;  // Checks if token exists (user is logged in)
    const deleteTokenUrl = `http://localhost:8080/api/logout/${token}`;  // URL to delete token on the backend

    if (!logged) {
        res.redirect('/login');  // If not logged in, redirects to login page
    } else {
        try {
            // Send DELETE request to backend to delete token
            await axios.delete(deleteTokenUrl, {
                headers: {
                    // Add any necessary headers (e.g., authorization headers)
                }
            });
            console.log('Token deleted successfully from backend');
        } catch (e) {
            console.error('Error deleting token:', e);
            // Handle error gracefully, e.g., notify the user
        }
    }

    // Clear the token cookie by setting an expired cookie
    res.clearCookie('token');

    // Redirects to the login page after logging out
    res.redirect('/login');
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
        if (error.response) {
            // The request was made and the server responded with a status code that falls out of the range of 2xx
            // console.error('Error response data:', error.response.data);
            // console.error('Error response status:', error.response.status);
            // console.error('Error response headers:', error.response.headers);
            if (error.response.data){
                res.render('login', { messageLogin: error.response.data.message, messageRegister:null});
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

// POST route to handle form submission
router.post('/register', async (req, res) => {
    const { username, mail, password } = req.body;
    const url = 'http://localhost:8080/api/register';
    const data = new FormData();
    data.append('username', username);
    data.append('password', password);
    data.append('mail', mail);

    try {
        await axios.post(url, data, {
            headers: {
                ...data.getHeaders()
            }
        });

        res.redirect('/login');

    } catch (error) {
        if (error.response) {
            // The request was made and the server responded with a status code that falls out of the range of 2xx
            // console.error('Error response data:', error.response.data);
            // console.error('Error response status:', error.response.status);
            // console.error('Error response headers:', error.response.headers);
            if (error.response.data){
                res.render('login', { messageRegister: error.response.data.message , messageLogin : null});
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

    
//afficher un utilisateur par son id
router.get('/user/:id', async (req, res) => {
    const userId = req.params.id;
    const ThisUrl = `http://localhost:8080/api/getUser/${userId}`;
    const randomUserUrl = 'http://localhost:8080/api/getRandomUser';
    const userTopicsUrl = `http://localhost:8080/api/getUserTopics/${userId}`;
    const tagUrl = 'http://localhost:8080/api/getTopicTag';
    const ownerUrl = 'http://localhost:8080/api/getTopicOwner';
    const likeUrl = 'http://localhost:8080/api/getLikeTopicNumber';
    const topicImgUrl = 'http://localhost:8080/api/getTopicImg';
    const followUrl = 'http://localhost:8080/api/getFollowers';
    const userFollowersUrl = `http://localhost:8080/api/getUserFollow/${userId}`;
    const userFollowingUrl = `http://localhost:8080/api/getUserFollowings/${userId}`;
    const token = req.cookies.token;
    const logged = token !== undefined;

    if (!logged) {
        return res.redirect('/login');
    }

    let topics = [];
    try {
        // Fetch user topics
        const response = await axios.get(userTopicsUrl);
        if (response.data && response.data.TopicList) {
            topics = response.data.TopicList;
        }

        // Create promises for fetching additional topic data
        const tagPromises = topics.map(topic => axios.get(`${tagUrl}/${topic.topic_id}`));
        const ownerPromises = topics.map(topic => axios.get(`${ownerUrl}/${topic.topic_id}`));
        const likePromises = topics.map(topic => axios.get(`${likeUrl}/${topic.topic_id}`));
        const imgPromises = topics.map(topic => axios.get(`${topicImgUrl}/${topic.topic_id}`));

        // Resolve all promises
        const [tagResponses, ownerResponses, likeResponses, imgResponses] = await Promise.all([
            Promise.all(tagPromises),
            Promise.all(ownerPromises),
            Promise.all(likePromises),
            Promise.all(imgPromises)
        ]);

        // Merge fetched data into topics
        topics.forEach((topic, index) => {
            topic.tags = tagResponses[index]?.data?.data || [];
            topic.owner = ownerResponses[index]?.data?.UserData || {};
            topic.numberOfLike = likeResponses[index]?.data?.NumberOfLike?.like_count || 0;
            topic.imgPath = imgResponses[index]?.data?.ImagePath || '';
        });

    } catch (error) {
        if (error.response && error.response.data && error.response.data.message === "No topics found for this user") {
            topics = [];
        } else {
            console.error('Error fetching topics:', error);
            return res.status(500).send('An error occurred while fetching topics');
        }
    }

    let forYou = [];
    try {
        // Fetch random users
        const forYouResponse = await axios.get(randomUserUrl);
        const users = forYouResponse.data.UsersData || [];

        // Create promises for fetching followers
        const followerPromises = users.map(user => axios.get(`${followUrl}/${user.user_id}`));
        const followerResponses = await Promise.all(followerPromises);

        // Merge follower data into users
        users.forEach((user, index) => {
            user.followers = followerResponses[index]?.data?.FollowerData?.follower_count || 0;
        });

        forYou = users;

    } catch (error) {
        console.error('Error fetching random users:', error);
    }

    let user = {};
    try {
        const userResponse = await axios.get(ThisUrl);
        user = userResponse.data.user || {};
    } catch (error) {
        console.error('Error fetching user data:', error);
        return res.status(500).send('An error occurred while fetching user data');
    }

    let followersNumber ={}
    let followingNumber ={}

    try {
        const followers = await axios.get(userFollowersUrl);
        const following= await axios.get(userFollowingUrl);
        followersNumber= followers.data.FollowerData.follower_count;
        followingNumber= following.data.FollowerData.following_count;
    }catch (e) {
        console.log(e)
    }

    const topicsLength = topics.length;

    const renderOptions = {
        user,
        logged,
        topics,
        forYou,
        topicsLength,
        followersNumber,
        followingNumber
    };

    res.render('profil', renderOptions);
});

router.get('/search', (req, res) => {
    const token = req.cookies.token;
    const logged = token !== undefined;

    res.render('search', { logged });
});

module.exports = router;
