const express = require('express');
const app = express();
const cors = require('cors');
const cookieParser = require('cookie-parser');
const path = require('path');
require('dotenv').config({ path: path.resolve(__dirname, 'data.env') });
const port = process.env.PORT || 3300; // 3300 is the default port if not found 3000
const viewsDirectories = [
    path.resolve(__dirname, '../frontend/templates'),
];
app.set ('view engine', 'ejs');

app.use(cors());
app.use(cookieParser());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.set ('views', viewsDirectories);


app.listen(port, () => {
    console.log(`Server is running on port ${port}`);
});
