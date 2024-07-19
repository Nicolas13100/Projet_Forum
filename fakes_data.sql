INSERT INTO users_table (username, email, password, biography, isAdmin, isModerator, profile_pic)
VALUES ('john_doe', 'john.doe@example.com', 'hashed_password_1', 'Film enthusiast and critic.', 1, 0,
        '/static/images/userAvatar/default-user.png'),
       ('jane_smith', 'jane.smith@example.com', 'hashed_password_2', 'Avid moviegoer and blogger.', 0, 1,
        '/static/images/userAvatar/default-user.png'),
       ('alice_jones', 'alice.jones@example.com', 'hashed_password_3', 'Movie director and producer.', 0, 0,
        '/static/images/userAvatar/default-user.png'),
       ('bob_brown', 'bob.brown@example.com', 'hashed_password_4', 'Film student and reviewer.', 0, 0,
        '/static/images/userAvatar/default-user.png');
INSERT INTO Topics_Table (title, body, status, user_id)
VALUES ('Top 10 Sci-Fi Movies of All Time', 'In this list, we explore the greatest science fiction films ever made.', 1,
        1),
       ('The Evolution of Action Films', 'A deep dive into how action films have evolved over the decades.', 1, 2),
       ('Classic Horror Films to Watch', 'A selection of classic horror films that every fan should see.', 1, 3),
       ('The Impact of Streaming Services on Cinema',
        'An analysis of how streaming platforms have changed the film industry.', 1, 4);
INSERT INTO Tags_Table (tag_name)
VALUES ('Sci-Fi'),
       ('Action'),
       ('Horror'),
       ('Streaming'),
       ('Classic'),
       ('Modern');
INSERT INTO TopicTags (topic_id, tag_id)
VALUES (1, 1), -- Top 10 Sci-Fi Movies of All Time - Sci-Fi
       (2, 2), -- The Evolution of Action Films - Action
       (3, 3), -- Classic Horror Films to Watch - Horror
       (4, 4); -- The Impact of Streaming Services on Cinema - Streaming
-- Original Messages
INSERT INTO Messages_Table (body, topic_id, user_id)
VALUES ('I completely agree with the list of sci-fi movies! Especially love "Blade Runner".', 1, 2),
       ('Great points on how action films have changed. The evolution is fascinating.', 2, 3),
       ('The classic horror films are a must-watch. "The Shining" is my all-time favorite!', 3, 4),
       ('Streaming services have definitely changed the game. More content than ever!', 4, 1);
-- Replies
INSERT INTO Messages_Table (body, base_message_id, user_id)
VALUES ('I totally agree with your point about "Blade Runner"! It’s a classic.', 1, 3),
       ('The evolution of action films is fascinating. I think "Die Hard" was a game-changer.',  2, 4),
       ('Have you seen "Psycho"? It’s a great classic horror film that’s not on the list.', 3, 1),
       ('I agree. Streaming services have made it easier to find rare films.', 4, 2);
INSERT INTO Admin_Logs_Table (action_descrbition)
VALUES ('Added new topic about Sci-Fi movies.'),
       ('Updated tags for action films topic.'),
       ('Deleted a message related to horror films.'),
       ('Created new user for film reviews.');
INSERT INTO images_Table (image_link, message_id, topic_id)
VALUES ('/static/images/TopicsImg/movie_poster_1.jpg', 1, 1),
       ('/static/images/TopicsImg/cinema_snapshot_1.png', 2, 2),
       ('/static/images/TopicsImg/clap.jpg', 3, 3),
       ('/static/images/TopicsImg/flms.webp', 4, 4);
INSERT INTO admin (topic_id, message_id, log_id)
VALUES (1, 1, 1),
       (2, 2, 2),
       (3, 3, 3),
       (4, 4, 4);
INSERT INTO react_topic (user_id, topic_id, status)
VALUES (1, 1, 1), -- User 1 likes Topic 1
       (2, 2, 0), -- User 2 has no comment on Topic 2
       (3, 3, 2), -- User 3 dislikes Topic 3
       (4, 4, 1); -- User 4 likes Topic 4
INSERT INTO friendship (sender_id, reciver_id, status, created_at, updated_at)
VALUES (1, 2, 1, '2024-01-01 10:00:00', '2024-01-02 11:00:00'),
       (2, 3, 0, '2024-02-01 14:00:00', '2024-02-01 14:00:00'),
       (3, 4, 1, '2024-03-01 16:00:00', '2024-03-02 17:00:00');
INSERT INTO react_message (user_id, message_id, status)
VALUES (1, 1, 1), -- User 1 likes Message 1
       (2, 2, 0), -- User 2 has no comment on Message 2
       (3, 3, 2), -- User 3 dislikes Message 3
       (4, 4, 1); -- User 4 likes Message 4
INSERT INTO follow (user_id, topic_id)
VALUES (1, 1), -- User 1 follows Topic 1
       (2, 2), -- User 2 follows Topic 2
       (3, 3), -- User 3 follows Topic 3
       (4, 4); -- User 4 follows Topic 4
INSERT INTO followUser (user_id, follower_id)
VALUES (1, 2), -- User 1 follows User 2
       (2, 3), -- User 2 follows User 3
       (3, 4), -- User 3 follows User 4
       (4, 1); -- User 4 follows User 1
INSERT INTO tokens (user_id, end_date, token)
VALUES (1, '2024-12-31 23:59:59', 'abcdef1234567890'),
       (2, '2024-12-31 23:59:59', '123456abcdef7890');
INSERT INTO password_reset_tokens (user_id, end_date, token)
VALUES (1, '2024-12-31 23:59:59', 'resetabcdef123456'),
       (2, '2024-12-31 23:59:59', 'reset123456abcdef');
