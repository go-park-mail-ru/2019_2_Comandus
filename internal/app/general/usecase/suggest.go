package generalUsecase

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients/interfaces"
	"github.com/suggest-go/suggest/pkg/dictionary"
	"github.com/suggest-go/suggest/pkg/metric"
	"github.com/suggest-go/suggest/pkg/suggest"
	"strings"
)

func formFreelancerDictionary(client clients.ClientUser) (dictionary.Dictionary, suggest.IndexDescription, error) {
	users, err := client.GetNamesFromServer()
	if err != nil {
		return nil, suggest.IndexDescription{}, err
	}

	dict := dictionary.NewInMemoryDictionary(users.Names)
	index := suggest.IndexDescription{
		Name:      "freelancers",
		NGramSize: 3,
		Wrap:      [2]string{"$", "$"},
		Pad:       "$",
		Alphabet:  []string{"russian", "$"},
	}

	return dict, index, nil
}

func formJobDictionary(client clients.ClientJob) (dictionary.Dictionary, suggest.IndexDescription, error) {
	tags, err := client.GetTags()
	if err != nil {
		return nil, suggest.IndexDescription{}, err
	}

	tags = unique(tags)

	index := suggest.IndexDescription{
		Name:      "jobs",
		NGramSize: 2,
		Wrap:      [2]string{"$", "$"},
		Pad:       "$",
		Alphabet:  []string{"english", "$"},
	}

	dict := dictionary.NewInMemoryDictionary(tags)
	return dict, index, nil
}

func NewSuggestService(uClient clients.ClientUser, jClient clients.ClientJob) (*suggest.Service, error) {
	jDict, jIndex, err := formJobDictionary(jClient)
	if err != nil {
		return nil, err
	}

	builder1, err := suggest.NewRAMBuilder(jDict, jIndex)
	if err != nil {
		return nil, err
	}

	fDict, fIndex, err := formFreelancerDictionary(uClient)
	if err != nil {
		return nil, err
	}

	builder2, err := suggest.NewRAMBuilder(fDict, fIndex)
	if err != nil {
		return nil, err
	}

	service := suggest.NewService()
	if err := service.AddIndex(jIndex.Name, jDict, builder1); err != nil {
		return nil, err
	}

	if err := service.AddIndex(fIndex.Name, fDict, builder2); err != nil {
		return nil, err
	}

	return service, nil
}

func GetSuggest(service *suggest.Service, query string, dict string) ([]string, error) {
	words := strings.Fields(query)
	lastWord := words[len(words)-1]
	searchConf, err := suggest.NewSearchConfig(lastWord, 5, metric.CosineMetric(), 0.4)
	if err != nil {
		return nil, err
	}

	result, err := service.Suggest(dict, searchConf)
	if err != nil {
		return nil, err
	}

	values := make([]string, 0, len(result))
	for _, item := range result {
		values = append(values, item.Value)
	}

	values = unique(values)

	return values, nil
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, tagLine := range stringSlice {
		for _, entry := range strings.Split(tagLine, ",") {
			if _, value := keys[entry]; !value {
				keys[entry] = true
				list = append(list, entry)
			}
		}
	}
	return list
}