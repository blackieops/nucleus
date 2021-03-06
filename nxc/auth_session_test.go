package nxc

import (
	"testing"

	"go.b8s.dev/nucleus/data"
	testUtils "go.b8s.dev/nucleus/internal/testing"
)

func TestCreateNextcloudAuthSession(t *testing.T) {
	testUtils.WithData(func(c *data.Context) {
		s, err := CreateNextcloudAuthSession(c)
		if err != nil {
			t.Errorf("Failed to create a NextcloudAuthSession: %v", err)
		}
		if s.LoginToken == s.PollToken {
			t.Errorf("NextcloudAuthSession's Poll and Login tokens were identical.")
		}
		if len(s.PollToken) != 128 {
			t.Errorf("NextcloudAuthSession poll token was wrong size: %d", len(s.PollToken))
		}
		if len(s.LoginToken) != 128 {
			t.Errorf("NextcloudAuthSession login token was wrong size: %d", len(s.LoginToken))
		}
	})
}

func TestFindNextcloudAuthSessionByPollToken(t *testing.T) {
	testUtils.WithData(func(c *data.Context) {
		s, err := CreateNextcloudAuthSession(c)
		if err != nil {
			t.Errorf("Test setup failed to create a NextcloudAuthSession: %v", err)
		}
		found, err := FindNextcloudAuthSessionByPollToken(c, s.PollToken)
		if err != nil {
			t.Errorf("Failed to find NextcloudAuthSession: %v", err)
		}
		if found.ID != s.ID {
			t.Errorf("Found wrong NextcloudAuthSession: %d instead of %d", found.ID, s.ID)
		}
		_, err = FindNextcloudAuthSessionByPollToken(c, "thisisnotatoken")
		if err == nil {
			t.Errorf("Should not have found a NextcloudAuthSession!")
		}
	})
}

func TestFindNextcloudAuthSessionByLoginToken(t *testing.T) {
	testUtils.WithData(func(c *data.Context) {
		s, err := CreateNextcloudAuthSession(c)
		if err != nil {
			t.Errorf("Test setup failed to create a NextcloudAuthSession: %v", err)
		}
		found, err := FindNextcloudAuthSessionByLoginToken(c, s.LoginToken)
		if err != nil {
			t.Errorf("Failed to find NextcloudAuthSession: %v", err)
		}
		if found.ID != s.ID {
			t.Errorf("Found wrong NextcloudAuthSession: %d instead of %d", found.ID, s.ID)
		}
		_, err = FindNextcloudAuthSessionByLoginToken(c, "thisisnotatoken")
		if err == nil {
			t.Errorf("Should not have found a NextcloudAuthSession!")
		}
	})
}

func TestDestroyNextcloudAuthSession(t *testing.T) {
	testUtils.WithData(func(c *data.Context) {
		s, err := CreateNextcloudAuthSession(c)
		if err != nil {
			t.Errorf("Test setup failed to create a NextcloudAuthSession: %v", err)
		}
		err = DestroyNextcloudAuthSession(c, s)
		if err != nil {
			t.Errorf("Failed to destroy NextcloudAuthSession: %v", err)
		}
		err = DestroyNextcloudAuthSession(c, &NextcloudAuthSession{})
		if err == nil {
			t.Errorf("Should not have destroyed invalid NextcloudAuthSession: %v", err)
		}
	})
}
