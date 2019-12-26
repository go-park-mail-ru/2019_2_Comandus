package model

type Chat struct {
	ID 			int64		`json:"id"`
	UserID 		int64		`json:"userId"`
	SupportID 	int64		`json:"supportId"`
	Name		string		`json:"name"`
	ProposalId	int64		`json:"proposalId"`
}
