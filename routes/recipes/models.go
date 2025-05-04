package recipes

import "ChaiwalaBackend/db"

type TeaType int32

const (
	TTBlack TeaType = iota
	TTGreen
	TTWhite
	TTOolong
	TTPuErh
	TTYellow
	TTHerbal
	TTRooibos
	TTYerbaMate
	TTMatcha
	TTChai
	TTFlavored
	TTBlooming
)

var TEANAMES = []string{
	"Black",
	"Green",
	"White",
	"Oolong",
	"Pu-erh",
	"Yellow",
	"Herbal",
	"Rooibos",
	"Yerba Mate",
	"Matcha",
	"Chai",
	"Flavored",
	"Blooming",
}

func (t TeaType) String() string {
	if t < 0 || int(t) >= len(TEANAMES) {
		return "Unknown"
	}

	return TEANAMES[t]
}

type Step struct {
	StepNumber  int    `json:"stepNumber,omitempty"`
	Description string `json:"description,omitempty"`
	AssetId     string `json:"assetId,omitempty"`
}

type SavedStep struct {
	ID int `json:"id"`
	Step
}

type CreateRecipeBody struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Instructions    string `json:"instructions"`
	TeaType         int    `json:"teaType"`
	Steps           []Step `json:"steps"`
	AssetId         string `json:"assetId"`
	PrepTimeMinutes int32  `json:"prepTimeMinutes"`
	Servings        int32  `json:"servings"`
	IsPublic        bool   `json:"isPublic"`
}

type UpdateRecipeBody struct {
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Instructions    string          `json:"instructions"`
	Steps           []db.RecipeStep `json:"steps"`
	TeaType         int             `json:"teaType"`
	AssetID         string          `json:"assetId"`
	PrepTimeMinutes int32           `json:"prepTimeMinutes"`
	Servings        int32           `json:"servings"`
	IsPublic        bool            `json:"isPublic"`
}

type GetRecipe struct {
	ID             int32           `json:"id"`
	Recipe         db.Recipe       `json:"recipe"`
	CreatedBy      db.User         `json:"createdBy"`
	Steps          []db.RecipeStep `json:"steps"`
	CommentsCount  int32           `json:"commentsCount"`
	FavoritesCount int32           `json:"favoritesCount"`
}
