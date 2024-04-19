CREATE TABLE Users_Table(
                            user_id COUNTER,
                            username VARCHAR(26) NOT NULL,
                            email VARCHAR(50) NOT NULL,
                            password VARCHAR(50) NOT NULL mininmum 8 / hashed 256,
                            registration_date DATETIME NOT NULL,
                            last_login_date DATETIME NOT NULL,
                            biography VARCHAR(200),
                            isAdmin LOGICAL NOT NULL,
                            isModerator LOGICAL NOT NULL,
                            profile_pic VARCHAR(100) NOT NULL,
                            PRIMARY KEY(user_id),
                            UNIQUE(username),
                            UNIQUE(email)
);

CREATE TABLE Topics_Table(
                             topic_id COUNTER,
                             title VARCHAR(50) NOT NULL,
                             body VARCHAR(1000) NOT NULL,
                             creation_date DATETIME NOT NULL,
                             status INT NOT NULL,
                             is_private LOGICAL NOT NULL,
                             user_id INT NOT NULL,
                             PRIMARY KEY(topic_id),
                             FOREIGN KEY(user_id) REFERENCES Users_Table(user_id)
);

CREATE TABLE Tags_Table(
                           tag_id INT,
                           tag_name VARCHAR(50) NOT NULL,
                           PRIMARY KEY(tag_id),
                           UNIQUE(tag_name)
);

CREATE TABLE Messages_Table(
                               message_id COUNTER,
                               body VARCHAR(500) NOT NULL,
                               date_sent DATETIME NOT NULL,
                               topic_id INT NOT NULL,
                               message_id_1 INT NOT NULL,
                               user_id INT NOT NULL,
                               PRIMARY KEY(message_id),
                               FOREIGN KEY(topic_id) REFERENCES Topics_Table(topic_id),
                               FOREIGN KEY(message_id_1) REFERENCES Messages_Table(message_id),
                               FOREIGN KEY(user_id) REFERENCES Users_Table(user_id)
);

CREATE TABLE Admin_Logs_Table(
                                 log_id COUNTER,
                                 action_descrbition TEXT NOT NULL,
                                 date_performed DATETIME NOT NULL,
                                 PRIMARY KEY(log_id)
);

CREATE TABLE images_Table(
                             image_id COUNTER,
                             image_origin_name VARCHAR(1000) NOT NULL,
                             image_serv_name VARCHAR(1000) NOT NULL,
                             image_link VARCHAR(1000) NOT NULL,
                             message_id INT,
                             topic_id INT,
                             PRIMARY KEY(image_id),
                             UNIQUE(image_serv_name),
                             UNIQUE(image_link),
                             FOREIGN KEY(message_id) REFERENCES Messages_Table(message_id),
                             FOREIGN KEY(topic_id) REFERENCES Topics_Table(topic_id)
);

CREATE TABLE have(
                     topic_id INT,
                     tag_id INT,
                     PRIMARY KEY(topic_id, tag_id),
                     FOREIGN KEY(topic_id) REFERENCES Topics_Table(topic_id),
                     FOREIGN KEY(tag_id) REFERENCES Tags_Table(tag_id)
);

CREATE TABLE admin(
                      topic_id INT,
                      message_id INT,
                      log_id INT,
                      PRIMARY KEY(topic_id, message_id, log_id),
                      FOREIGN KEY(topic_id) REFERENCES Topics_Table(topic_id),
                      FOREIGN KEY(message_id) REFERENCES Messages_Table(message_id),
                      FOREIGN KEY(log_id) REFERENCES Admin_Logs_Table(log_id)
);

CREATE TABLE react(
                      user_id INT,
                      topic_id INT,
                      status INT NOT NULL,
                      PRIMARY KEY(user_id, topic_id),
                      FOREIGN KEY(user_id) REFERENCES Users_Table(user_id),
                      FOREIGN KEY(topic_id) REFERENCES Topics_Table(topic_id)
);

CREATE TABLE friendship(
                           user_id INT,
                           user_id_1 INT,
                           status INT NOT NULL,
                           created_at DATETIME NOT NULL,
                           updated_at DATETIME NOT NULL,
                           PRIMARY KEY(user_id, user_id_1),
                           FOREIGN KEY(user_id) REFERENCES Users_Table(user_id),
                           FOREIGN KEY(user_id_1) REFERENCES Users_Table(user_id)
);

CREATE TABLE react_1(
                        user_id INT,
                        message_id INT,
                        status INT NOT NULL,
                        PRIMARY KEY(user_id, message_id),
                        FOREIGN KEY(user_id) REFERENCES Users_Table(user_id),
                        FOREIGN KEY(message_id) REFERENCES Messages_Table(message_id)
);

CREATE TABLE follow(
                       user_id INT,
                       topic_id INT,
                       followed_date DATETIME NOT NULL,
                       PRIMARY KEY(user_id, topic_id),
                       FOREIGN KEY(user_id) REFERENCES Users_Table(user_id),
                       FOREIGN KEY(topic_id) REFERENCES Topics_Table(topic_id)
);
