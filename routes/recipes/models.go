package recipes

type Step struct {
	StepNumber  int    `json:"stepNumber,omitempty"`
	Description string `json:"description,omitempty"`
	AssetId     string `json:"assetId,omitempty"`
}

type CreateRecipeBody struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Instructions    string `json:"instructions"`
	Steps           []Step `json:"steps"`
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
