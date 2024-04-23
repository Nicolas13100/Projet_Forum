INSERT INTO users_table (username, email, password, biography, isAdmin, isModerator, profile_pic) VALUES
('user1', 'user1@example.com', 'password1', 'Biography of user1', 0, 0, '/profile_pics/user1.jpg'),
('user2', 'user2@example.com', 'password2', 'Biography of user2', 0, 0, '/profile_pics/user2.jpg'),
('admin', 'admin@example.com', 'adminpass', 'Administrator', 1, 0, '/profile_pics/admin.jpg');


INSERT INTO Topics_Table (title, body, status, is_private, user_id) VALUES
('Topic 1', 'Body of topic 4', 1, 0, 2),
('Topic 2', 'Body of topic 5', 1, 1, 1),
('Topic 3', 'Body of topic 6', 1, 0, 3);

INSERT INTO Tags_Table (tag_name) VALUES
('tag1'),
('tag2'),
('tag3');

INSERT INTO Messages_Table (body, topic_id, user_id) VALUES
('Message 1', 1, 1),
('Message 2', 1, 2),
('Message 3', 2, 1);

INSERT INTO Admin_Logs_Table (action_descrbition) VALUES
('Performed action 1'),
('Performed action 2'),
('Performed action 3');

INSERT INTO images_Table (image_origin_name, image_serv_name, message_id, topic_id) VALUES
('image1.jpg', 'image1_serv', 1, 1),
('image2.jpg', 'image2_serv', 2, 1),
('image3.jpg', 'image3_serv', 3, 2);

INSERT INTO have (topic_id, tag_id) VALUES
(1, 1),
(1, 2),
(2, 2);

INSERT INTO admin (topic_id, message_id, log_id) VALUES
(1, 1, 1),
(2, 2, 2),
(3, 3, 3);

INSERT INTO react_topic (user_id, topic_id, status) VALUES
(1, 1, 1),
(1, 2, 2),
(2, 1, 1);

INSERT INTO friendship (sender_id, reciver_id, status, created_at, updated_at) VALUES
(1, 2, 1, NOW(), NOW()),
(2, 1, 1, NOW(), NOW());

INSERT INTO react_message (user_id, message_id, status) VALUES
(1, 1, 1),
(2, 1, 2),
(1, 2, 1);

INSERT INTO follow (user_id, topic_id, followed_date) VALUES
(1, 1, NOW()),
(1, 2, NOW()),
(2, 1, NOW());
