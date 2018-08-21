package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
	}
	// set a value
	testURL := "http://url-to-gravatar/"
	client.userData = map[string]interface{}{"avatar_url": testURL}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value preset")
	}
	if url != testURL {
		t.Error("AuthAvatar.GetAvatarURL should return correct url")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData = map[string]interface{}{"user_id": "0bc83cb571cd1c50ba6f3e8a78ef1346"}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}

func TestUploadAvatar(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(b))
	filename := filepath.Join(basepath, "avatars", "abc.jpg")
	err := ioutil.WriteFile(filename, []byte{}, os.ModePerm)
	if err != nil {
		t.Errorf("Error when create file test: %s", err.Error())
	}
	// defer os.Remove(filename)
	var fileSystemAvatar UploadAvatar
	client := new(client)
	client.userData = map[string]interface{}{"user_id": "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("UploadAvatar.GetAvatarURL should not return an error")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("UploadAvatar.GetAvatarURL wrongly returned %s", url)
	}
}
