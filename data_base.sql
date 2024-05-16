CREATE TABLE users_table
(
    user_id           INT AUTO_INCREMENT,
    username          VARCHAR(26)  NOT NULL,
    email             VARCHAR(50)  NOT NULL,
    password          VARCHAR(50)  NOT NULL,
    registration_date DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    last_login_date   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
    biography         VARCHAR(200),
    isAdmin           TINYINT      NOT NULL DEFAULT 0 CHECK (isAdmin = 0 OR isAdmin = 1),         -- 0 false, 1 true
    isModerator       TINYINT      NOT NULL DEFAULT 0 CHECK (isModerator = 0 OR isModerator = 1), -- 0 false, 1 true
    is_deleted        TINYINT      NOT NULL DEFAULT 0 CHECK ( is_deleted = 0 OR is_deleted = 1 ), -- 0 false, 1 true
    profile_pic       VARCHAR(100) NOT NULL,
    PRIMARY KEY (user_id),
    UNIQUE (username),
    UNIQUE (email)
)ENGINE = INNODB;

CREATE TABLE Topics_Table
(
    topic_id      INT AUTO_INCREMENT,
    title         VARCHAR(50)   NOT NULL,
    body          VARCHAR(1000) NOT NULL,
    creation_date DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    status        INT           NOT NULL,
    is_private    TINYINT       NOT NULL DEFAULT 0 CHECK (is_private = 0 OR is_private = 1), -- 0 false, 1 true
    user_id       INT           NOT NULL,
    PRIMARY KEY (topic_id),
    FOREIGN KEY (user_id) REFERENCES users_table (user_id)
)ENGINE = INNODB;

CREATE TABLE Tags_Table
(
    tag_id   INT AUTO_INCREMENT,
    tag_name VARCHAR(50) NOT NULL,
    PRIMARY KEY (tag_id),
    UNIQUE (tag_name)
)ENGINE = INNODB;

CREATE TABLE Messages_Table
(
    message_id      INT AUTO_INCREMENT,
    body            VARCHAR(500) NOT NULL,
    date_sent       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    topic_id        INT          NOT NULL,
    base_message_id INT,
    user_id         INT          NOT NULL,
    PRIMARY KEY (message_id),
    FOREIGN KEY (topic_id) REFERENCES Topics_Table (topic_id),
    FOREIGN KEY (base_message_id) REFERENCES Messages_Table (message_id),
    FOREIGN KEY (user_id) REFERENCES users_Table (user_id)
)ENGINE = INNODB;

CREATE TABLE Admin_Logs_Table
(
    log_id             int AUTO_INCREMENT,
    action_descrbition TEXT     NOT NULL,
    date_performed     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    PRIMARY KEY (log_id)
)ENGINE = INNODB;

CREATE TABLE images_Table
(
    image_id          INT AUTO_INCREMENT,
    image_origin_name VARCHAR(1000) NOT NULL,
    image_serv_name   VARCHAR(1000) NOT NULL,
    image_link        VARCHAR(1000) NOT NULL DEFAULT '/assets/img/default.png',
    message_id        INT,
    topic_id          INT,
    PRIMARY KEY (image_id),
    UNIQUE (image_serv_name),
    UNIQUE (image_link),
    FOREIGN KEY (message_id) REFERENCES Messages_Table (message_id),
    FOREIGN KEY (topic_id) REFERENCES Topics_Table (topic_id)
)ENGINE = INNODB;

CREATE TABLE have
(
    topic_id INT,
    tag_id   INT,
    PRIMARY KEY (topic_id, tag_id),
    FOREIGN KEY (topic_id) REFERENCES Topics_Table (topic_id),
    FOREIGN KEY (tag_id) REFERENCES Tags_Table (tag_id)
)ENGINE = INNODB;

CREATE TABLE admin
(
    topic_id   INT,
    message_id INT,
    log_id     INT,
    PRIMARY KEY (topic_id, message_id, log_id),
    FOREIGN KEY (topic_id) REFERENCES Topics_Table (topic_id),
    FOREIGN KEY (message_id) REFERENCES Messages_Table (message_id),
    FOREIGN KEY (log_id) REFERENCES Admin_Logs_Table (log_id)
)ENGINE = INNODB;

CREATE TABLE react_topic
(
    user_id  INT,
    topic_id INT,
    status   TINYINT NOT NULL DEFAULT 0 CHECK (status >= 0 AND status <= 2), -- 0 no comment, 1 like, 2 dislike
    PRIMARY KEY (user_id, topic_id),
    FOREIGN KEY (user_id) REFERENCES users_Table (user_id),
    FOREIGN KEY (topic_id) REFERENCES Topics_Table (topic_id)
)ENGINE = INNODB;

CREATE TABLE friendship
(
    sender_id  INT,
    reciver_id INT,
    status     TINYINT  NOT NULL DEFAULT 0 CHECK (status >= 0 AND status <= 2), -- 0 Pending, 1 accepted, 2 refused
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    PRIMARY KEY (sender_id, reciver_id),
    FOREIGN KEY (sender_id) REFERENCES users_Table (user_id),
    FOREIGN KEY (reciver_id) REFERENCES users_Table (user_id)
)ENGINE = INNODB;

CREATE TABLE react_message
(
    user_id    INT,
    message_id INT,
    status     TINYINT NOT NULL DEFAULT 0 CHECK (status >= 0 AND status <= 2), -- 0 nothing, 1 like, 2 dislike
    PRIMARY KEY (user_id, message_id),
    FOREIGN KEY (user_id) REFERENCES users_Table (user_id),
    FOREIGN KEY (message_id) REFERENCES Messages_Table (message_id)
)ENGINE = INNODB;

CREATE TABLE follow
(
    user_id       INT,
    topic_id      INT,
    followed_date DATETIME NOT NULL,
    PRIMARY KEY (user_id, topic_id),
    FOREIGN KEY (user_id) REFERENCES users_Table (user_id),
    FOREIGN KEY (topic_id) REFERENCES Topics_Table (topic_id)
)ENGINE = INNODB;
