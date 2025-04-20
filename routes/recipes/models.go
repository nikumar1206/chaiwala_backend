package recipes

type CreateRecipeBody struct {
	UserID          int32  `json:"userId"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Instructions    string `json:"instructions"`
	AssetId         string `json:"assetId"`
	PrepTimeMinutes int32  `json:"prepTimeMinutes"`
	Servings        int32  `json:"servings"`
	IsPublic        bool   `json:"isPublic"`
}

type UpdateRecipeBody struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Instructions    string `json:"instructions"`
	AssetId         string `json:"assetId"`
	PrepTimeMinutes int32  `json:"prepTimeMinutes"`
	Servings        int32  `json:"servings"`
	IsPublic        bool   `json:"isPublic"`
}
