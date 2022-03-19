package maildoor

// productConfig options allow to customize the productConfig name and logo
// as well as the favicon. These are used in the email that gets
// sent to the user and the login form.
type productConfig struct {
	Name       string
	LogoURL    string
	FaviconURL string
}
