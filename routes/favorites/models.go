package favorites

type Favorite struct {
	RecipeID int32 `json:"recipeId"`
	UserID   int32 `json:"userId"`
}
