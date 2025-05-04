package recipes

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
