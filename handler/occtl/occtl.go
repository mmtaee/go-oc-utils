package occtl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// Oc occtl command
type Oc struct{}

// OcInterface occtl command methods
type OcInterface interface {
	Reload(c context.Context) error
	OnlineUsers(c context.Context) (*[]OcUser, error)
	Disconnect(c context.Context, username string) error
	ShowIPBans(c context.Context) (*[]IPBan, error)
	ShowIPBansPoints(c context.Context) (*[]IPBanPoints, error)
	UnBanIP(c context.Context, ip string) error
	ShowStatus(c context.Context) (string, error)
	ShowIRoutes(c context.Context) (*[]IRoute, error)
	ShowUser(c context.Context, username string) (*[]OcUser, error)
}

// NewOcctl Create New Occtl command obj
func NewOcctl() *Oc {
	return &Oc{}
}

// Reload server configuration reload
func (o *Oc) Reload(c context.Context) error {
	_, err := Exec(c, "reload")
	if err != nil {
		return err
	}
	return nil
}

// OnlineUsers list of online users
func (o *Oc) OnlineUsers(c context.Context) (*[]OcUser, error) {
	var users []OcUser
	result, err := Exec(c, "-j show users")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, &users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

// Disconnect expire user session. On disconnected users raise error
func (o *Oc) Disconnect(c context.Context, username string) error {
	_, err := Exec(c, fmt.Sprintf("disconnect user %s", username))
	if err != nil {
		return errors.New("failed to disconnect user " + username)
	}
	return nil
}

// ShowIPBans List of banned IPs
func (o *Oc) ShowIPBans(c context.Context) (*[]IPBan, error) {
	var ipBans []IPBan
	result, err := Exec(c, "-j show ip bans")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, &ipBans)
	if err != nil {
		return nil, err
	}
	return &ipBans, nil
}

// ShowIPBansPoints List of baned IPs with points
func (o *Oc) ShowIPBansPoints(c context.Context) (*[]IPBanPoints, error) {
	var ipBansPoint []IPBanPoints
	result, err := Exec(c, "-j show ip bans points")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, &ipBansPoint)
	if err != nil {
		return nil, err
	}
	return &ipBansPoint, nil
}

// UnBanIP unban banned IP
func (o *Oc) UnBanIP(c context.Context, ip string) error {
	_, err := Exec(c, fmt.Sprintf("unban ip %s", ip))
	if err != nil {
		return err
	}
	return nil
}

// ShowStatus server status
func (o *Oc) ShowStatus(c context.Context) (string, error) {
	result, err := Exec(c, "show status")
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// ShowIRoutes list user IP routes
func (o *Oc) ShowIRoutes(c context.Context) (*[]IRoute, error) {
	result, err := Exec(c, "-j show iroutes")
	if err != nil {
		return nil, err
	}
	var routes []IRoute
	err = json.Unmarshal(result, &routes)
	if err != nil {
		return nil, err
	}
	return &routes, nil
}

// ShowUser show user info with extra data
func (o *Oc) ShowUser(c context.Context, username string) (*[]OcUser, error) {
	var user *[]OcUser
	result, err := Exec(c, fmt.Sprintf("-j show user %s", username))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
