-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (
  email, password_hash, bio, avatar_url
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: ListPublicRecipes :many
SELECT * FROM recipes
WHERE is_public = true
ORDER BY created_at DESC;

-- name: GetRecipe :one
SELECT * FROM recipes
WHERE id = $1;

-- name: ListUserRecipes :many
SELECT * FROM recipes
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: CreateRecipe :one
INSERT INTO recipes (
  user_id, title, description, type, asset_id,
  prep_time_minutes, servings, is_public
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateRecipe :exec
UPDATE recipes SET
  title = $2,
  description = $3,
  type = $4,
  asset_id = $5,
  prep_time_minutes = $6,
  servings = $7,
  is_public = $8,
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteRecipe :exec
DELETE FROM recipes
WHERE id = $1;

-- name: AddRecipeStep :one
INSERT INTO recipe_steps (
  recipe_id, step_number, description, asset_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateRecipeStep :exec
UPDATE recipe_steps
SET
  step_number = $2,
  description = $3,
  asset_id = $4
WHERE id = $1;


-- name: DeleteRecipeStep :exec
DELETE FROM recipe_steps
WHERE id = $1;

-- name: GetRecipeStep :one
SELECT * FROM recipe_steps
WHERE id = $1;

-- name: GetRecipeStepByNumber :one
SELECT * FROM recipe_steps
WHERE recipe_id = $1 AND step_number = $2;

-- name: GetRecipeStepsByRecipe :many
SELECT * FROM recipe_steps
WHERE recipe_id = $1
ORDER BY step_number;

-- name: ListRecipeSteps :many
SELECT * FROM recipe_steps
WHERE recipe_id = $1
ORDER BY step_number;

-- name: AddComment :one
INSERT INTO recipe_comments (
  recipe_id, user_id, comment
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: ListComments :many
SELECT rc.comment, rc.created_at, u.email, u.avatar_url
FROM recipe_comments rc
JOIN users u ON rc.user_id = u.id
WHERE rc.recipe_id = $1
ORDER BY rc.created_at DESC;

-- name: ListCommentsByUser :many
SELECT * FROM recipe_comments
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateComment :exec
UPDATE recipe_comments
SET comment = $2
WHERE id = $1;

-- name: DeleteComment :exec
DELETE FROM recipe_comments
WHERE id = $1;

-- name: FavoriteRecipe :exec
INSERT INTO favorites (user_id, recipe_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: UnfavoriteRecipe :exec
DELETE FROM favorites
WHERE user_id = $1 AND recipe_id = $2;

-- name: ListUserFavorites :many
SELECT r.*
FROM favorites f
JOIN recipes r ON f.recipe_id = r.id
WHERE f.user_id = $1
ORDER BY f.created_at DESC;

-- name: IsRecipeFavorited :one
SELECT EXISTS (
  SELECT 1 FROM favorites
  WHERE user_id = $1 AND recipe_id = $2
) AS favorited;
