package models

import (
	"log"
	"strconv"
)

type Post struct {
	By          string
	Descendants int
	ID          int64
	Score       int
	Time        int64
	Title       string
	Type        string
	URL         string
}

func (post *Post) Update() {
	err := db.Cmd("HMSET", "hackernews:post:"+strconv.FormatInt(post.ID, 10), "id", post.ID, "title", post.Title, "url", post.URL).Err
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Post with %d is saved ", post.ID)
}

func UpdateTopStories(ids []int64) error {

	// getting new connection from pool
	conn, err := db.Get()
	if err != nil {
		log.Fatal(err)
	}
	// defering connection close
	defer db.Put(conn)

	// starting multi query
	err = conn.Cmd("MULTI").Err
	if err != nil {
		log.Fatal(err)
	}

	// deleting existing top stories
	err = conn.Cmd("DEL", "hackernews:topstories").Err
	if err != nil {
		log.Fatal(err)
	}
	// adding top stories to sorted set

	for i, id := range ids {
		err = conn.Cmd("ZADD", "hackernews:topstories", i, id).Err
		if err != nil {
			log.Fatal(err)
		}
	}

	// executing multiple commands at one
	err = conn.Cmd("EXEC").Err
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Top stories updated!")
	return nil
}

func FindTopFive() ([]*Post, error) {
	// getting new connection from pool
	conn, err := db.Get()
	if err != nil {
		log.Fatal(err)
	}
	// defering connection close
	defer db.Put(conn)

	reply, err := conn.Cmd("ZRANGE", "hackernews:topstories", 0, 4).List()
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Cmd("MULTI").Err
	if err != nil {
		log.Fatal(err)
	}

	for _, id := range reply {
		err := conn.Cmd("HGETALL", "hackernews:post:"+id).Err
		if err != nil {
			return nil, err
		}
	}

	ereply := conn.Cmd("EXEC")
	if ereply.Err != nil {
		return nil, err
	}

	areply, err := ereply.Array()
	if err != nil {
		return nil, err
	}
	// log.Println(areply)
	posts := make([]*Post, 5)

	for i := 0; i < 5; i++ {
		reply := areply[i]
		mreply, err := reply.Map()
		if err != nil {
			return nil, err
		}

		post, err := populatePost(mreply)
		if err != nil {
			return nil, err
		}

		posts[i] = post
	}

	return posts, nil
}

func populatePost(reply map[string]string) (*Post, error) {
	var err error
	post := new(Post)
	post.ID, err = strconv.ParseInt(reply["id"], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	post.Title = reply["title"]
	post.URL = reply["url"]
	if err != nil {
		return nil, err
	}
	return post, nil
}
