package comments

type CreateCommentBody struct {
	RecipeID int32  `json:"recipeId"`
	UserID   int32  `json:"userId"`
	Comment  string `json:"comment"`
}

type UpdateCommentBody struct {
	Comment string `json:"comment"`
}
