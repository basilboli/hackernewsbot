package models

import "log"

// User represents user
type User struct {
	ChatID       string
	LastModified int64
}

// SubscribeForTopStories subscribes user for notification of top five stories
func SubscribeForTopStories(id int64) error {
	err := db.Cmd("SADD", "hackernews:subscription:topstories", id).Err
	if err != nil {
		return err
	}
	log.Printf("User %d is subscribed for topstories!", id)
	return nil
}

// UnSubscribeForTopStories subscribes user for notification of top five stories
func UnSubscribeForTopStories(id int64) error {
	err := db.Cmd("SREM", "hackernews:subscription:topstories", id).Err
	if err != nil {
		return err
	}
	log.Printf("User %d is unsubscribed!", id)
	return nil
}

// GetAllSubscribedUsers subscribes user for notification of top five stories
func GetAllSubscribedUsers() ([]string, error) {
	ids, err := db.Cmd("SMEMBERS", "hackernews:subscription:topstories").List()
	if err != nil {
		log.Fatal(err)
	}
	return ids, nil
}

// IsSubscribed checks whether the user is subscribed
func IsSubscribed(id int64) (bool, error) {
	isMember, err := db.Cmd("SISMEMBER", "hackernews:subscription:topstories", id).Int()
	if err != nil {
		return false, err
	}

	if isMember == 1 {
		return true, nil
	}
	return false, nil
}
