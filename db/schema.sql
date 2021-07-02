DROP DATABASE IF EXISTS bookstore;
CREATE DATABASE bookstore;

\c bookstore;
DROP TABLE IF EXISTS books;

CREATE TABLE books (
    id INT GENERATED ALWAYS AS IDENTITY,
    title varchar(100) NOT NULL,
    author varchar(100) NOT NULL,
    publisher varchar(100) NOT NULL,
    publishDate date NOT NULL,
    rating int NOT NULL, -- 1-3
    status varchar(10) NOT NULL -- (CheckedIn, CheckedOut)
);

INSERT INTO books(title, author, publisher, publishDate, rating, status) VALUES ('test', 'test', 'test', NOW(), 1, 'good');