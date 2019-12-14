package generalUsecase

import (
	"github.com/suggest-go/suggest/pkg/dictionary"
	"github.com/suggest-go/suggest/pkg/metric"
	"github.com/suggest-go/suggest/pkg/suggest"
	"log"
)

func formFreelancerDictionary() (dictionary.Dictionary, error){
	dict := dictionary.NewInMemoryDictionary([]string{
		"dasqha efimova",
		"daqqsha efimdddova",
		"daswha efimoddva",
		"dashwa efimovfa",
		"dashea efffimova",
		"dasrutha easdfimova",
		"dashera aaefimova",
	})

	return dict, nil
}

func formJobDictionary() (dictionary.Dictionary, error){
	dict := dictionary.NewInMemoryDictionary([]string{
		"golang",
		"java",
		"python",
		"server",
		"my server",
		"server hello",
		"server golang",
	})
	return dict, nil
}

func NewSuggestService() (*suggest.Service, error) {
	jDict, err := formJobDictionary()
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	fDict, err := formFreelancerDictionary()
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	jobIndex := suggest.IndexDescription {
		Name:      "jobs",
		NGramSize: 2,
		Wrap:      [2]string{"$", "$"},
		Pad:       "$",
		Alphabet:  []string{"english", "$"},
	}

	freelancerIndex := suggest.IndexDescription {
		Name:      "freelancers",
		NGramSize: 3,
		Wrap:      [2]string{"$", "$"},
		Pad:       "$",
		Alphabet:  []string{"english", "$"},
	}

	builder1, err := suggest.NewRAMBuilder(jDict, jobIndex)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	builder2, err := suggest.NewRAMBuilder(fDict, freelancerIndex)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	service := suggest.NewService()
	if err := service.AddIndex(jobIndex.Name, jDict, builder1); err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	if err := service.AddIndex(freelancerIndex.Name, fDict, builder2); err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	return service, nil
}

func GetSuggest(service *suggest.Service, query string, dict string) ([]string, error) {
	log.Println(query)
	searchConf, err := suggest.NewSearchConfig(query, 5, metric.CosineMetric(), 0.4)
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

	return values, nil
}
