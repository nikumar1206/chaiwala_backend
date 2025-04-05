CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    bio TEXT,
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT NOW ()
);

CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    instructions TEXT NOT NULL,
    image_url TEXT,
    prep_time_minutes INTEGER,
    brew_time_minutes INTEGER,
    servings INTEGER,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW (),
    updated_at TIMESTAMP DEFAULT NOW ()
);

CREATE TABLE recipe_steps (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER REFERENCES recipes (id) ON DELETE CASCADE,
    step_number INTEGER NOT NULL,
    description TEXT NOT NULL,
    media_url TEXT
);

CREATE TABLE ingredients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

-- Join table for recipe ingredients
CREATE TABLE recipe_ingredients (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER REFERENCES recipes (id) ON DELETE CASCADE,
    ingredient_id INTEGER REFERENCES ingredients (id),
    quantity VARCHAR(50) -- e.g. "2 tsp", "1 cup"
);

-- Tags table
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Join table for recipe tags
CREATE TABLE recipe_tags (
    recipe_id INTEGER REFERENCES recipes (id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (recipe_id, tag_id)
);

CREATE TABLE recipe_comments (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER REFERENCES recipes (id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    comment TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW ()
);

-- Favorites (likes/bookmarks)
CREATE TABLE favorites (
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    recipe_id INTEGER REFERENCES recipes (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW (),
    PRIMARY KEY (user_id, recipe_id)
);

-- Indexes for performance
CREATE INDEX idx_recipes_user_id ON recipes (user_id);

CREATE INDEX idx_favorites_user_id ON favorites (user_id);

CREATE INDEX idx_recipe_tags_recipe_id ON recipe_tags (recipe_id);

CREATE INDEX idx_recipe_ingredients_recipe_id ON recipe_ingredients (recipe_id);

CREATE INDEX idx_recipe_comments_recipe_id ON recipe_comments (recipe_id);

CREATE INDEX idx_recipe_steps_recipe_id ON recipe_steps (recipe_id);
