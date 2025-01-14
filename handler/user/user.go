package ocservUser

import (
	"context"
	"fmt"
	"os/exec"
)

// OcUser ocserv user
type OcUser struct{}

// OcUserInterface ocserv user methods
type OcUserInterface interface {
	Create(c context.Context, username, password, group string) error
	Update(c context.Context, username, password, group string) error
	Lock(c context.Context, username string) error
	UnLock(c context.Context, username string) error
	Delete(c context.Context, username string) error
}

var (
	ocpasswdCMD = "/usr/bin/ocpasswd"    // ocpasswd os command path
	passwdFile  = "/etc/ocserv/ocpasswd" // ocpasswd file path
)

// NewOcUser create new ocserv user obj
func NewOcUser() *OcUser {
	return &OcUser{}
}

// Create  ocserv user creation with password and group
func (u *OcUser) Create(c context.Context, username, password, group string) error {
	if group == "defaults" || group == "" {
		group = ""
	} else {
		group = fmt.Sprintf("-g %s", group)
	}
	command := fmt.Sprintf("/usr/bin/echo -e \"%s\\n%s\\n\" | %s %s -c %s %s",
		password,
		password,
		ocpasswdCMD,
		group,
		passwdFile,
		username,
	)
	return exec.CommandContext(c, "sh", "-c", command).Run()
}

// Update  ocserv user updating with password and group
func (u *OcUser) Update(c context.Context, username, password, group string) error {
	return u.Create(c, username, password, group)
}

// Lock disable ocserv user to connect to server(Ocserv User Locked)
func (u *OcUser) Lock(c context.Context, username string) error {
	command := fmt.Sprintf("%s %s -c %s %s", ocpasswdCMD, "-l", passwdFile, username)
	return exec.CommandContext(c, "sh", "-c", command).Run()
}

// UnLock enable ocserv user to connect to server(Ocserv User UnLocked)
func (u *OcUser) UnLock(c context.Context, username string) error {
	command := fmt.Sprintf("%s %s -c %s %s", ocpasswdCMD, "-u", passwdFile, username)
	return exec.CommandContext(c, "sh", "-c", command).Run()
}

// Delete ocserv user deleting account
func (u *OcUser) Delete(c context.Context, username string) error {
	command := fmt.Sprintf("%s -c %s -d %s", ocpasswdCMD, passwdFile, username)
	return exec.CommandContext(c, "sh", "-c", command).Run()
}
