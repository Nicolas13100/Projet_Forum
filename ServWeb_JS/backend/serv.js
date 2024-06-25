const path = require('path');
require('dotenv').config({ path: path.resolve(__dirname, '.env') });
const express = require('express');
const cors = require('cors');
const cookieParser = require('cookie-parser');
const bodyParser = require('body-parser');
const displayRoutes = require('./src/routes/routes');

const app = express();

app.use(cors());
app.use('/static', express.static(path.resolve(__dirname, '../frontEnd/assets')));
app.set('view engine', 'ejs');
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

const viewsDirectories = [
    path.resolve(__dirname, '../frontend/templates'),
];
app.set('views', viewsDirectories);

app.use(cookieParser());
app.use(express.json());
app.use('/', displayRoutes);

const port = process.env.PORT || 3300;
app.listen(port, () => {
    console.log(`Server is running on port ${port}`);
});
