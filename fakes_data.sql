-- Insert Users Data
INSERT INTO users_table (username, email, password, biography, isAdmin, isModerator, profile_pic)
VALUES ('filmlover123', 'filmlover123@example.com', 'password123', 'Passionate about movies!', 0, 0,
        'profile_pic1.jpg'),
       ('cinemafanatic', 'cinemafanatic@example.com', 'securepassword', 'Obsessed with all things cinema.', 0, 0,
        'profile_pic2.jpg'),
       ('moviebuff456', 'moviebuff456@example.com', 'moviesarelife', 'Movie enthusiast exploring new genres.', 0, 0,
        'profile_pic3.jpg');

-- Insert Topics Data
INSERT INTO Topics_Table (title, body, status, is_private, user_id)
VALUES ('Favorite Directors', 'Discuss your favorite film directors and their works here!', 1, 0, 1),
       ('Latest Releases', 'Share your thoughts on the latest movies hitting theaters.', 1, 0, 2),
       ('Classic Films', 'Rediscover and analyze timeless classics from the golden era of cinema.', 1, 0, 3);

-- Insert Tags Data
INSERT INTO Tags_Table (tag_name)
VALUES ('Action'),
       ('Drama'),
       ('Comedy'),
       ('Sci-Fi'),
       ('Horror');

-- Insert Messages Data
INSERT INTO Messages_Table (body, topic_id, user_id)
VALUES ('I just watched the latest superhero movie and it was amazing!', 2, 1),
       ('Has anyone seen the new indie film everyone is talking about?', 2, 2),
       ('Let''s discuss the cinematography in this classic film!', 3, 3);

-- Insert Admin Logs Data
INSERT INTO Admin_Logs_Table (action_descrbition)
VALUES ('Banned user filmlover123 for violating forum rules.'),
       ('Moderated inappropriate content in topic Latest Releases.');

-- Insert Images Data
INSERT INTO images_Table (image_origin_name, image_serv_name, topic_id)
VALUES ('movie_poster.jpg', 'movie_poster_1.jpg', 2),
       ('cinema_snapshot.png', 'cinema_snapshot_1.png', 3);

-- Insert Have Data
INSERT INTO have (topic_id, tag_id)
VALUES (2, 1),
       (2, 4),
       (3, 2),
       (3, 3);

-- Insert Admin Data
INSERT INTO admin (topic_id, message_id, log_id)
VALUES (2, 1, 1),
       (2, 2, 2),
       (3, 3, 3);

-- Insert React Topic Data
INSERT INTO react_topic (user_id, topic_id, status)
VALUES (1, 2, 1),
       (2, 2, 2),
       (3, 3, 1);

-- Insert Friendship Data
INSERT INTO friendship (sender_id, reciver_id, status, created_at, updated_at)
VALUES (1, 2, 1, NOW(), NOW()),
       (2, 3, 1, NOW(), NOW());

-- Insert React Message Data
INSERT INTO react_message (user_id, message_id, status)
VALUES (1, 1, 1),
       (2, 1, 2),
       (3, 3, 1);

-- Insert Follow Data
INSERT INTO follow (user_id, topic_id, followed_date)
VALUES (1, 2, NOW()),
       (2, 3, NOW()),
       (3, 1, NOW());