package arn

// ConnectGoogle connects the user's account with a Google account.
func (user *User) ConnectGoogle(googleID string) error {
	user.Accounts.Google.ID = googleID

	return DB.Set("GoogleToUser", googleID, &GoogleToUser{
		ID:     googleID,
		UserID: user.ID,
	})
}

// ConnectFacebook connects the user's account with a Facebook account.
func (user *User) ConnectFacebook(facebookID string) error {
	user.Accounts.Facebook.ID = facebookID

	return DB.Set("FacebookToUser", facebookID, &FacebookToUser{
		ID:     facebookID,
		UserID: user.ID,
	})
}

// ConnectTwitter connects the user's account with a Twitter account.
func (user *User) ConnectTwitter(twtterID string) error {
	user.Accounts.Twitter.ID = twtterID

	return DB.Set("TwitterToUser", twtterID, &TwitterToUser{
		ID:     twtterID,
		UserID: user.ID,
	})
}
