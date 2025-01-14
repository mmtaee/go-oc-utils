package ocuser

// Sync type of user with group to sync ocpasswd file with custom strategies
type Sync struct {
	Username string `json:"username"`
	Group    string `json:"group"`
}
