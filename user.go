package disgopher

import "encoding/json"

//User ...
type User struct {
	state         *clientState
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator int    `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot"`
	System        bool   `json:"system"`
	MFAEnabled    bool   `json:"mfa_enabled"`
	Locale        string `json:"locale"`
	Verified      bool   `json:"verified"`
	Email         string `json:"email"`
	Flags         int    `json:"flags"`
	PremiumType   int    `json:"premium_type"`
	PublicFlags   int    `json:"public_flags"`
}

func newUser(state *clientState, data []byte) *User {
	user := &User{state: state}
	json.Unmarshal(data, user)
	//TODO: gotUser := state.Users[user.ID] if user != nil user.update(data)...
	state.Users[user.ID] = user
	return user
}
