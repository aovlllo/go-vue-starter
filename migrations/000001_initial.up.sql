CREATE TABLE users (
                       id int NOT NULL AUTO_INCREMENT,
                       email varchar(50),
                       password char(60) binary,
                       name varchar(50),
                       second_name varchar(50),
                       birth date,
                       sex enum('male', 'female', 'non binary'),
                       interests varchar(512),
                       city varchar(50),
                       PRIMARY KEY(id)
);

