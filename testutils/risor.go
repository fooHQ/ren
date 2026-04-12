package testutils

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/foohq/ren"
)

type MockOS struct {
	mock.Mock
}

func (m *MockOS) Create(name string) (ren.File, error) {
	args := m.Called(name)
	return args.Get(0).(ren.File), args.Error(1)
}

func (m *MockOS) Mkdir(name string, perm ren.FileMode) error {
	args := m.Called(name, perm)
	return args.Error(0)
}

func (m *MockOS) MkdirAll(path string, perm ren.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

func (m *MockOS) MkdirTemp(dir, pattern string) (string, error) {
	args := m.Called(dir, pattern)
	return args.String(0), args.Error(1)
}

func (m *MockOS) Open(name string) (ren.File, error) {
	args := m.Called(name)
	return args.Get(0).(ren.File), args.Error(1)
}

func (m *MockOS) OpenFile(name string, flag int, perm ren.FileMode) (ren.File, error) {
	args := m.Called(name, flag, perm)
	return args.Get(0).(ren.File), args.Error(1)
}

func (m *MockOS) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockOS) Remove(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockOS) RemoveAll(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *MockOS) Rename(oldpath, newpath string) error {
	args := m.Called(oldpath, newpath)
	return args.Error(0)
}

func (m *MockOS) Stat(name string) (ren.FileInfo, error) {
	args := m.Called(name)
	return args.Get(0).(ren.FileInfo), args.Error(1)
}

func (m *MockOS) Symlink(oldname, newname string) error {
	args := m.Called(oldname, newname)
	return args.Error(0)
}

func (m *MockOS) TempDir() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockOS) WriteFile(name string, data []byte, perm ren.FileMode) error {
	args := m.Called(name, data, perm)
	return args.Error(0)
}

func (m *MockOS) ReadDir(name string) ([]ren.DirEntry, error) {
	args := m.Called(name)
	return args.Get(0).([]ren.DirEntry), args.Error(1)
}

func (m *MockOS) PathSeparator() rune {
	args := m.Called()
	return rune(args.Int(0))
}

func (m *MockOS) PathListSeparator() rune {
	args := m.Called()
	return rune(args.Int(0))
}

func (m *MockOS) Args() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockOS) Chdir(dir string) error {
	args := m.Called(dir)
	return args.Error(0)
}

func (m *MockOS) Environ() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockOS) Exit(code int) {
	m.Called(code)
}

func (m *MockOS) Getpid() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockOS) Getuid() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockOS) Getwd() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockOS) Hostname() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockOS) Setenv(key, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockOS) Getenv(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *MockOS) Unsetenv(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockOS) LookupEnv(key string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

func (m *MockOS) UserCacheDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockOS) UserConfigDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockOS) UserHomeDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockOS) Stdin() ren.File {
	args := m.Called()
	return args.Get(0).(ren.File)
}

func (m *MockOS) Stdout() ren.File {
	args := m.Called()
	return args.Get(0).(ren.File)
}

func (m *MockOS) CurrentUser() (ren.User, error) {
	args := m.Called()
	return args.Get(0).(ren.User), args.Error(1)
}

func (m *MockOS) LookupUser(name string) (ren.User, error) {
	args := m.Called(name)
	return args.Get(0).(ren.User), args.Error(1)
}

func (m *MockOS) LookupUid(uid string) (ren.User, error) {
	args := m.Called(uid)
	return args.Get(0).(ren.User), args.Error(1)
}

func (m *MockOS) LookupGroup(name string) (ren.Group, error) {
	args := m.Called(name)
	return args.Get(0).(ren.Group), args.Error(1)
}

func (m *MockOS) LookupGid(gid string) (ren.Group, error) {
	args := m.Called(gid)
	return args.Get(0).(ren.Group), args.Error(1)
}

type MockFile struct {
	mock.Mock
}

func (m *MockFile) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *MockFile) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *MockFile) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockFile) Stat() (ren.FileInfo, error) {
	args := m.Called()
	return args.Get(0).(ren.FileInfo), args.Error(1)
}

type MockFileInfo struct {
	mock.Mock
}

func (m *MockFileInfo) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockFileInfo) Size() int64 {
	args := m.Called()
	return int64(args.Int(0))
}

func (m *MockFileInfo) Mode() ren.FileMode {
	args := m.Called()
	return args.Get(0).(ren.FileMode)
}

func (m *MockFileInfo) ModTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockFileInfo) IsDir() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockFileInfo) Sys() any {
	args := m.Called()
	return args.Get(0)
}

type MockDirEntry struct {
	mock.Mock
}

func (m *MockDirEntry) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockDirEntry) IsDir() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockDirEntry) Type() ren.FileMode {
	args := m.Called()
	return args.Get(0).(ren.FileMode)
}

func (m *MockDirEntry) Info() (ren.FileInfo, error) {
	args := m.Called()
	return args.Get(0).(ren.FileInfo), args.Error(1)
}

type MockUser struct {
	uid      string
	gid      string
	username string
	name     string
	homeDir  string
}

func NewMockUser(uid, gid, username, name, homeDir string) *MockUser {
	return &MockUser{
		uid:      uid,
		gid:      gid,
		username: username,
		name:     name,
		homeDir:  homeDir,
	}
}

func (u *MockUser) Uid() string {
	return u.uid
}

func (u *MockUser) Gid() string {
	return u.gid
}

func (u *MockUser) Username() string {
	return u.username
}

func (u *MockUser) Name() string {
	return u.name
}

func (u *MockUser) HomeDir() string {
	return u.homeDir
}

type MockGroup struct {
	gid  string
	name string
}

func NewMockGroup(gid, name string) *MockGroup {
	return &MockGroup{
		gid:  gid,
		name: name,
	}
}

func (g *MockGroup) Gid() string {
	return g.gid
}

func (g *MockGroup) Name() string {
	return g.name
}
