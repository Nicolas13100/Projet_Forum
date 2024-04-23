-- Fake data for users_table
INSERT INTO users_table (username, email, password, biography, profile_pic)
VALUES
    ('user1', 'user1@example.com', 'password123', 'Movie lover and reviewer.', '/assets/img/user1_profile.jpg'),
    ('user2', 'user2@example.com', 'pass123word', 'Film enthusiast exploring new releases.', '/assets/img/user2_profile.jpg'),
    ('user3', 'user3@example.com', 'filmlover456', 'Cinema fan always looking for hidden gems.', '/assets/img/user3_profile.jpg');

-- Fake data for Films_Table
INSERT INTO Films_Table (title, release_year, director_id, genre_id, description, poster)
VALUES
    ('Inception', 2010, 1, 1, 'A mind-bending heist movie set within the architecture of the mind.', '/assets/img/inception_poster.jpg'),
    ('The Shawshank Redemption', 1994, 2, 2, 'Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.', '/assets/img/shawshank_redemption_poster.jpg'),
    ('Pulp Fiction', 1994, 3, 3, 'The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.', '/assets/img/pulp_fiction_poster.jpg');

-- Fake data for Directors_Table
INSERT INTO Directors_Table (director_name)
VALUES
    ('Christopher Nolan'),
    ('Frank Darabont'),
    ('Quentin Tarantino');

-- Fake data for Genres_Table
INSERT INTO Genres_Table (genre_name)
VALUES
    ('Action'),
    ('Drama'),
    ('Crime');

-- Fake data for Actors_Table
INSERT INTO Actors_Table (actor_name)
VALUES
    ('Leonardo DiCaprio'),
    ('Morgan Freeman'),
    ('Tim Robbins'),
    ('Samuel L. Jackson'),
    ('John Travolta'),
    ('Uma Thurman');

-- Fake data for Reviews_Table
INSERT INTO Reviews_Table (film_id, user_id, rating, review_text, review_date)
VALUES
    (1, 1, 5, 'Absolutely mind-blowing! A must-watch for any film buff.', '2024-04-20'),
    (2, 2, 5, 'An emotional rollercoaster. One of the best movies ever made.', '2024-04-21'),
    (3, 3, 4, 'Pulp Fiction is a masterpiece of storytelling and filmmaking.', '2024-04-22');

-- Fake data for Film_Actors_Table (many-to-many relationship between Films_Table and Actors_Table)
INSERT INTO Film_Actors_Table (film_id, actor_id)
VALUES
    (1, 1), -- Inception - Leonardo DiCaprio
    (2, 2), -- The Shawshank Redemption - Morgan Freeman
    (2, 3), -- The Shawshank Redemption - Tim Robbins
    (3, 4), -- Pulp Fiction - Samuel L. Jackson
    (3, 5), -- Pulp Fiction - John Travolta
    (3, 6); -- Pulp Fiction - Uma Thurman
