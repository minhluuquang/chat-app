package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
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
	email, ok := c.userData["email"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	emailStr, ok := email.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}
	m := md5.New()
	io.WriteString(m, strings.ToLower(emailStr))
	return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
}
