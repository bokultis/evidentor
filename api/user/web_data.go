package user

type UserWO struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Gender    string `json:"gender"`
	Birthday  string `json:"birthday"`
	Email     string `json:"eMail"`
	Address   string `json:"address"`
	Optin     string `json:"optin"`
}

type UserInputWO struct {
	ID        *int    `json:"id"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Gender    *string `json:"gender"`
	Birthday  *string `json:"birthday"`
	Email     *string `json:"eMail"`
	Address   *string `json:"address"`
	Optin     *string `json:"optin"`
	Password  *string `json:"password"`
}

type UserUpdateWO struct {
	ID        *int    `json:"id"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Gender    *string `json:"gender"`
	Birthday  *string `json:"birthday"`
	Email     *string `json:"eMail"`
	Address   *string `json:"address"`
	Optin     *string `json:"optin"`
}

// Create a struct to read the username and password from the request body
type UserCredentials struct {
	Password *string `json:"password"`
	Email    *string `json:"email"`
}
