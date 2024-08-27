CREATE TABLE comments
(
    id SERIAL PRIMARY KEY,
    news_id INT DEFAULT NULL, 
    parent_comment_id INT DEFAULT NULL, 
    content TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE comments owner to admin;

INSERT INTO comments (news_id, parent_comment_id, content)
VALUES 
    (1, NULL, 'This is the first comment.'),
    (2, NULL, 'This is the second comment.'),
    (1, 1, 'This is a reply to the first comment.'),
    (3, NULL, 'This is the third comment.'),
    (2, 2, 'This is a reply to the second comment.');

-- Проверка добавленных записей
SELECT * FROM comments;
