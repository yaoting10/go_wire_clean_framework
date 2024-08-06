package response

type MinUser struct {
	Id     uint   `json:"id"`
	Avatar string `json:"avatar"`
}

type CodedUser struct {
	MinUser
	Code string
}

type NamedUser struct {
	CodedUser
	NickName string
}

type UserParent struct {
	MinUser
	Email   string `json:"email"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type LanguageResp struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
