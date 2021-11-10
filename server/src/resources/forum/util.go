package forum

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/forPelevin/gomoji"
	"github.com/openmultiplayer/web/server/src/db"
)

const MaxShortBodyLength = 128

func MakeShortBody(long string) string {
	full := math.Min(float64(len(long)), MaxShortBodyLength)
	firstPara := strings.Index(long, "\n")
	if firstPara == -1 {
		firstPara = 999999
	}

	end := int(math.Min(full, float64(firstPara)))

	return long[:end]
}

func timeOrNil(x bool) *time.Time {
	if x {
		t := time.Now()
		return &t
	}
	return nil
}

func CanUserMutatePost(ctx context.Context, d *db.PrismaClient, authorID, id string) error {
	post, err := d.Post.
		FindUnique(db.Post.ID.Equals(id)).
		With(db.Post.Author.Fetch()).
		Exec(ctx)
	if err != nil {
		return err
	}
	if post.Author().Admin {
		return nil
	}
	if post.Author().ID == authorID {
		return nil
	}
	return ErrUnauthorised
}

func IsValidEmoji(e string) (string, bool) {
	if len(e) == 0 {
		return "", false
	}
	if e[1] == ':' {
		return "", false
	}
	return e, gomoji.ContainsEmoji(e)
}
