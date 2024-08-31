-- name: GetAllComments :many
SELECT * 
FROM comments
ORDER BY id;

-- name: CreateComments :exec
INSERT INTO comments 
(news_id, parent_comment_id, content) 
VALUES ($1, $2, $3);

-- name: DeleteComment :exec
DELETE FROM comments 
WHERE id = $1;