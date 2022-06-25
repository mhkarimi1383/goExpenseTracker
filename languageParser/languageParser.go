package languageParser

import (
	"io/ioutil"

	"github.com/mhkarimi1383/goExpenseTracker/configuration"
	"github.com/mhkarimi1383/goExpenseTracker/logger"
	"github.com/mhkarimi1383/goExpenseTracker/types"
	"gopkg.in/yaml.v3"
)

var (
	defaultLanguage = types.Language{
		LanguageName: "english",
		Direction:    "LTR",
		CurrencySign: "$",
		BalanceTitle: "Your Balance:",
		Expense:      "Expense",
		Income:       "Income",
		Add:          "Add",
		Remove:       "Remove",
		Description:  "Description",
		Amount:       "Amount",
	}
	languages        *types.TranslateFile
	selectedLanguage string
)

func init() {
	cfg, err := configuration.GetConfig()
	if err != nil {
		logger.Fatalf(true, "error in initializing configuration: %v", err)
	}
	selectedLanguage = cfg.Language
	languageFile, err := ioutil.ReadFile("translate.yaml")
	if err != nil {
		logger.Warnf(true, "Error reading translation file: %v, using default language values", err)
		languages.Languages[0] = defaultLanguage
	} else {
		err = yaml.Unmarshal(languageFile, languages)
		if err != nil {
			logger.Warnf(true, "Invalid language file: %v, using default language values", err)
			languages.Languages[0] = defaultLanguage
		}
	}
}

func GetLanguageList() []types.Language {
	return languages.Languages
}

func GetSelectedLanguage() types.Language {
	for _, language := range languages.Languages {
		if language.LanguageName == selectedLanguage {
			return language
		}
	}
	return defaultLanguage
}
