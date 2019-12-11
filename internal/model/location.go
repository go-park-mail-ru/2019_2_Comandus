package model

type Country struct {
	ID			int64
	Name		string
}

type City struct {
	ID			int64
	CountryID	int64
	Name		string

}