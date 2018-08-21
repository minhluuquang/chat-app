package main

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
)

var (
	ErrNoAvatarURL = errors.New("No Avatar URL")
)

// Avatar is an interface to do anything with avatar of users
type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	url, ok := c.userData["avatar_url"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	urlStr, ok := url.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}
	return urlStr, nil
}

type GravatarAvatar struct{}

var UseGravatarAvatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	userID, ok := c.userData["user_id"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	userIDStr, ok := userID.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}
	return "//www.gravatar.com/avatar/" + userIDStr, nil
}

type UploadAvatar struct{}

var UserUploadAvatar UploadAvatar

func (UploadAvatar) GetAvatarURL(c *client) (string, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(b))
	userID, ok := c.userData["user_id"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	userIDStr, ok := userID.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}
	files, err := ioutil.ReadDir(filepath.Join(basepath, "/avatars"))
	if err != nil {
		return "", ErrNoAvatarURL
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if match, _ := path.Match(userIDStr+"*", file.Name()); match {
			return "/avatars/" + file.Name(), nil
		}
	}
	return "", ErrNoAvatarURL
}
