package model

type Chat struct {
	ID         int64  `json:"id"`
	Freelancer int64  `json:"freelancerId"`
	Manager    int64  `json:"managerId"`
	Name       string `json:"name"`
	ProposalId int64  `json:"proposalId"`
	UserId	   int64  `json:"userId"`
}
