-- Fake data for users_table
INSERT INTO users_table (username, email, password, biography, profile_pic)
VALUES
    ('john_doe', 'john@example.com', 'password123', 'Hello, I am John Doe.', '/assets/img/john_profile.jpg'),
    ('jane_smith', 'jane@example.com', 'pass123word', 'Hi there! I am Jane Smith.', '/assets/img/jane_profile.jpg'),
    ('sam_green', 'sam@example.com', 'green123', 'Nice to meet you all! I am Sam Green.', '/assets/img/sam_profile.jpg');

-- Fake data for Topics_Table
INSERT INTO Topics_Table (title, body, status, user_id)
VALUES
    ('Introduction', 'This is the introduction topic.', 1, 1),
    ('Technology', 'Discussing the latest in tech.', 1, 2),
    ('Travel', 'Share your travel experiences here.', 1, 3);

-- Fake data for Tags_Table
INSERT INTO Tags_Table (tag_name)
VALUES
    ('Discussion'),
    ('Tech'),
    ('Travel');

-- Fake data for Messages_Table
INSERT INTO Messages_Table (body, topic_id, user_id)
VALUES
    ('Welcome everyone!', 1, 1),
    ('What are your thoughts on the new iPhone?', 2, 2),
    ('I just came back from a trip to Europe!', 3, 3);

-- Fake data for Admin_Logs_Table
INSERT INTO Admin_Logs_Table (action_descrbition)
VALUES
    ('Created a new topic.'),
    ('Deleted a spam message.'),
    ('Updated user permissions.');

-- Fake data for images_Table
INSERT INTO images_Table (image_origin_name, image_serv_name, message_id, topic_id)
VALUES
    ('profile_pic1.jpg', '/assets/img/profile1.jpg', NULL, NULL),
    ('profile_pic2.jpg', '/assets/img/profile2.jpg', NULL, NULL),
    ('profile_pic3.jpg', '/assets/img/profile3.jpg', NULL, NULL);

-- Fake data for have
INSERT INTO have (topic_id, tag_id)
VALUES
    (1, 1),
    (2, 2),
    (3, 3);

-- Fake data for admin
INSERT INTO admin (topic_id, message_id, log_id)
VALUES
    (1, 1, 1),
    (2, 2, 2),
    (3, 3, 3);

-- Fake data for react_topic
INSERT INTO react_topic (user_id, topic_id, status)
VALUES
    (1, 2, 1),
    (2, 1, 2),
    (3, 3, 1);

-- Fake data for friendship
INSERT INTO friendship (sender_id, reciver_id, status, created_at, updated_at)
VALUES
    (1, 2, 1, NOW(), NOW()),
    (2, 3, 1, NOW(), NOW()),
    (3, 1, 1, NOW(), NOW());

-- Fake data for react_message
INSERT INTO react_message (user_id, message_id, status)
VALUES
    (1, 1, 1),
    (2, 2, 2),
    (3, 3, 1);

-- Fake data for follow
INSERT INTO follow (user_id, topic_id, followed_date)
VALUES
    (1, 2, NOW()),
    (2, 1, NOW()),
    (3, 3, NOW());
