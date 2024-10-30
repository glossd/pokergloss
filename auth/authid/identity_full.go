package authid

type IdentityFull struct {
	Identity
	Provider string `json:"provider"`
	EmailVerified bool `json:"emailVerified"`
}

func (id IdentityFull) IsAnonymous() bool {
	return id.Provider == "anonymous"
}
