-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: CreateUser :one
INSERT INTO users (
  username, email, password_hash, bio, avatar_url
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: ListPublicRecipes :many
SELECT * FROM recipes
WHERE is_public = true
ORDER BY created_at DESC;

-- name: GetRecipe :one
SELECT * FROM recipes
WHERE id = $1 AND is_public = true;

-- name: ListUserRecipes :many
SELECT * FROM recipes
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: CreateRecipe :one
INSERT INTO recipes (
  user_id, title, description, instructions, image_url,
  prep_time_minutes, brew_time_minutes, servings, is_public
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateRecipe :exec
UPDATE recipes SET
  title = $2,
  description = $3,
  instructions = $4,
  image_url = $5,
  prep_time_minutes = $6,
  brew_time_minutes = $7,
  servings = $8,
  is_public = $9,
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteRecipe :exec
DELETE FROM recipes
WHERE id = $1;

-- name: AddRecipeStep :one
INSERT INTO recipe_steps (
  recipe_id, step_number, description, media_url
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: ListRecipeSteps :many
SELECT * FROM recipe_steps
WHERE recipe_id = $1
ORDER BY step_number;

-- name: AddIngredient :one
INSERT INTO ingredients (name)
VALUES ($1)
ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
RETURNING *;

-- name: AddRecipeIngredient :one
INSERT INTO recipe_ingredients (
  recipe_id, ingredient_id, quantity
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: ListRecipeIngredients :many
SELECT ri.quantity, i.name
FROM recipe_ingredients ri
JOIN ingredients i ON ri.ingredient_id = i.id
WHERE ri.recipe_id = $1;

-- name: AddTag :one
INSERT INTO tags (name)
VALUES ($1)
ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
RETURNING *;

-- name: TagRecipe :exec
INSERT INTO recipe_tags (recipe_id, tag_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: ListRecipeTags :many
SELECT t.name
FROM recipe_tags rt
JOIN tags t ON rt.tag_id = t.id
WHERE rt.recipe_id = $1;

-- name: AddComment :one
INSERT INTO recipe_comments (
  recipe_id, user_id, comment
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: ListComments :many
SELECT rc.comment, rc.created_at, u.username, u.avatar_url
FROM recipe_comments rc
JOIN users u ON rc.user_id = u.id
WHERE rc.recipe_id = $1
ORDER BY rc.created_at DESC;

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
