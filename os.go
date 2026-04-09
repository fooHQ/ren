// Portions of this file are adapted from Risor (https://github.com/deepnoodle-ai/risor).
// Licensed under the Apache License, Version 2.0.

package ren

import (
	"context"
	"errors"
	"os"
	"os/user"
	"strings"

	"github.com/foohq/urlpath"
)

// User represents a system user.
type User interface {
	Uid() string
	Gid() string
	Username() string
	Name() string
	HomeDir() string
}

// Group represents a system group.
type Group interface {
	Gid() string
	Name() string
}

// OS provides an interface for interacting with the operating system.
type OS interface {
	FS
	Args() []string
	Chdir(dir string) error
	Environ() []string
	Exit(code int)
	Getpid() int
	Getuid() int
	Getwd() (dir string, err error)
	Hostname() (name string, err error)
	Setenv(key, value string) error
	Getenv(key string) string
	Unsetenv(key string) error
	LookupEnv(key string) (string, bool)
	TempDir() string
	UserCacheDir() (string, error)
	UserConfigDir() (string, error)
	UserHomeDir() (string, error)
	Stdin() File
	Stdout() File
	PathSeparator() rune
	PathListSeparator() rune
	CurrentUser() (User, error)
	LookupUser(name string) (User, error)
	LookupUid(uid string) (User, error)
	LookupGroup(name string) (Group, error)
	LookupGid(gid string) (Group, error)
}

type osContextKey struct{}

// WithOS returns a new context with the given OS implementation.
func WithOS(ctx context.Context, o OS) context.Context {
	return context.WithValue(ctx, osContextKey{}, o)
}

// GetOS returns the OS implementation from the context.
func GetOS(ctx context.Context) OS {
	o, _ := ctx.Value(osContextKey{}).(OS)
	return o
}

func isOS(ctx context.Context) bool {
	_, ok := ctx.Value(osContextKey{}).(OS)
	return ok
}

// ExitHandler is a function that handles os.exit calls.
type ExitHandler func(int)

var _ OS = (*osMiddleware)(nil)

type osMiddleware struct {
	wd          string
	fs          fsMiddleware
	stdin       File
	stdout      File
	args        []string
	exitHandler ExitHandler
}

func (o *osMiddleware) Mkdir(name string, perm os.FileMode) error {
	pth, err := urlpath.Abs(name, o.wd)
	if err != nil {
		return err
	}
	return o.fs.Mkdir(pth, perm)
}

func (o *osMiddleware) MkdirAll(path string, perm os.FileMode) error {
	pth, err := urlpath.Abs(path, o.wd)
	if err != nil {
		return err
	}
	return o.fs.MkdirAll(pth, perm)
}

func (o *osMiddleware) MkdirTemp(dir, pattern string) (string, error) {
	pth, err := urlpath.Abs(dir, o.wd)
	if err != nil {
		return "", err
	}
	return o.fs.MkdirTemp(pth, pattern)
}

func (o *osMiddleware) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	pth, err := urlpath.Abs(name, o.wd)
	if err != nil {
		return nil, err
	}
	return o.fs.OpenFile(pth, flag, perm)
}

func (o *osMiddleware) ReadFile(name string) ([]byte, error) {
	pth, err := urlpath.Abs(name, o.wd)
	if err != nil {
		return nil, err
	}
	return o.fs.ReadFile(pth)
}

func (o *osMiddleware) Remove(name string) error {
	pth, err := urlpath.Abs(name, o.wd)
	if err != nil {
		return err
	}
	return o.fs.Remove(pth)
}

func (o *osMiddleware) RemoveAll(path string) error {
	pth, err := urlpath.Abs(path, o.wd)
	if err != nil {
		return err
	}
	return o.fs.RemoveAll(pth)
}

func (o *osMiddleware) Rename(oldpath, newpath string) error {
	oldPth, err := urlpath.Abs(oldpath, o.wd)
	if err != nil {
		return err
	}
	newPth, err := urlpath.Abs(newpath, o.wd)
	if err != nil {
		return err
	}
	return o.fs.Rename(oldPth, newPth)
}

func (o *osMiddleware) Stat(name string) (os.FileInfo, error) {
	pth, err := urlpath.Abs(name, o.wd)
	if err != nil {
		return nil, err
	}
	return o.fs.Stat(pth)
}

func (o *osMiddleware) Symlink(oldname, newname string) error {
	oldPth, err := urlpath.Abs(oldname, o.wd)
	if err != nil {
		return err
	}
	newPth, err := urlpath.Abs(newname, o.wd)
	if err != nil {
		return err
	}
	return o.fs.Symlink(oldPth, newPth)
}

func (o *osMiddleware) TempDir() string {
	return os.TempDir()
}

func (o *osMiddleware) WriteFile(name string, content []byte, perm os.FileMode) error {
	pth, err := urlpath.Abs(name, o.wd)
	if err != nil {
		return err
	}
	return o.fs.WriteFile(pth, content, perm)
}

func (o *osMiddleware) ReadDir(name string) ([]DirEntry, error) {
	pth, err := urlpath.Abs(name, o.wd)
	if err != nil {
		return nil, err
	}
	return o.fs.ReadDir(pth)
}

func (o *osMiddleware) PathSeparator() rune {
	return urlpath.PathSeparator
}

func (o *osMiddleware) PathListSeparator() rune {
	return urlpath.PathListSeparator
}

func (o *osMiddleware) Chdir(dir string) error {
	pth, err := urlpath.Abs(dir, o.wd)
	if err != nil {
		return err
	}
	info, err := o.fs.Stat(pth)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New("chdir " + pth + ": file is not a directory")
	}
	scheme, err := urlpath.Scheme(pth)
	if err != nil {
		return err
	}
	if scheme == "file" {
		pth = strings.TrimPrefix(pth, "file://")
	}
	o.wd = pth
	return nil
}

func (o *osMiddleware) Getwd() (dir string, err error) {
	if o.wd == "" {
		o.wd, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	return o.wd, nil
}

func (o *osMiddleware) Stdout() File {
	return o.stdout
}

func (o *osMiddleware) Stdin() File {
	return o.stdin
}

func (o *osMiddleware) Args() []string {
	return o.args
}

func (o *osMiddleware) Environ() []string {
	return os.Environ()
}

func (o *osMiddleware) Getenv(key string) string {
	return os.Getenv(key)
}

func (o *osMiddleware) Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func (o *osMiddleware) Unsetenv(key string) error {
	return os.Unsetenv(key)
}

func (o *osMiddleware) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (o *osMiddleware) Exit(code int) {
	if o.exitHandler != nil {
		o.exitHandler(code)
	}
}

func (o *osMiddleware) Getpid() int {
	return os.Getpid()
}

func (o *osMiddleware) Getuid() int {
	return os.Getuid()
}

func (o *osMiddleware) Hostname() (string, error) {
	return os.Hostname()
}

func (o *osMiddleware) UserCacheDir() (string, error) {
	return os.UserCacheDir()
}

func (o *osMiddleware) UserConfigDir() (string, error) {
	return os.UserConfigDir()
}

func (o *osMiddleware) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (o *osMiddleware) CurrentUser() (User, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	return &userWrapper{User: u}, nil
}

func (o *osMiddleware) LookupUser(name string) (User, error) {
	u, err := user.Lookup(name)
	if err != nil {
		return nil, err
	}
	return &userWrapper{User: u}, nil
}

func (o *osMiddleware) LookupUid(uid string) (User, error) {
	u, err := user.LookupId(uid)
	if err != nil {
		return nil, err
	}
	return &userWrapper{User: u}, nil
}

func (o *osMiddleware) LookupGroup(name string) (Group, error) {
	g, err := user.LookupGroup(name)
	if err != nil {
		return nil, err
	}
	return &groupWrapper{Group: g}, nil
}

func (o *osMiddleware) LookupGid(gid string) (Group, error) {
	g, err := user.LookupGroupId(gid)
	if err != nil {
		return nil, err
	}
	return &groupWrapper{Group: g}, nil
}

var _ User = (*userWrapper)(nil)

// userWrapper wraps the standard library's user.User type to implement the User interface.
type userWrapper struct {
	*user.User
}

// Uid returns the user ID.
func (u *userWrapper) Uid() string {
	return u.User.Uid
}

// Gid returns the primary group ID.
func (u *userWrapper) Gid() string {
	return u.User.Gid
}

// Username returns the username.
func (u *userWrapper) Username() string {
	return u.User.Username
}

// Name returns the user's name.
func (u *userWrapper) Name() string {
	return u.User.Name
}

// HomeDir returns the user's home directory.
func (u *userWrapper) HomeDir() string {
	return u.User.HomeDir
}

var _ Group = (*groupWrapper)(nil)

// groupWrapper wraps the standard library's user.Group type to implement the Group interface.
type groupWrapper struct {
	*user.Group
}

// Gid returns the group ID.
func (g *groupWrapper) Gid() string {
	return g.Group.Gid
}

// Name returns the group name.
func (g *groupWrapper) Name() string {
	return g.Group.Name
}
