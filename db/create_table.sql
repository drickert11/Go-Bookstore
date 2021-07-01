use Bookstore;
DROP TABLE IF EXISTS Books;

CREATE TABLE Books (
    [ID] SERIAL PRIMARY KEY NOT NULL,
    [Title] varchar(100) NOT NULL,
    [Author] varchar(100) NOT NULL,
    [Publisher] varchar(100) NOT NULL,
    [PublishDate] date NOT NULL,
    [Rating] int NOT NULL, -- 1-3
    [Status] varchar(10) NOT NULL -- (CheckedIn, CheckedOut)
);