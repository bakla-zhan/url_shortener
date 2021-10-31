package link

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Link struct {
	Long  string
	Short string
}

type LinkStore interface {
	CreateLink(ctx context.Context, link Link) error
	ReadLink(ctx context.Context, shortLink string) (longLink string, err error)
}

type Links struct {
	lstore LinkStore
}

func NewLinks(lstore LinkStore) *Links {
	return &Links{
		lstore: lstore,
	}
}

func (l *Links) Create(ctx context.Context, longLink string) (shortLink string, err error) {
	shortLink = shortLinkGenerator()
	link := Link{
		Long:  longLink,
		Short: shortLink,
	}
	err = l.lstore.CreateLink(ctx, link)
	if err != nil {
		return "", fmt.Errorf("create link error: %w", err)
	}
	return shortLink, nil
}

func (l *Links) Read(ctx context.Context, shortLink string) (longLink string, err error) {
	return l.lstore.ReadLink(ctx, shortLink)
}

var chars = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func shortLinkGenerator() (shortLink string) {
	shortLinkArr := make([]byte, 8)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < cap(shortLinkArr); i++ {
		shortLinkArr[i] = chars[rand.Intn(len(chars)-1)]
	}

	for _, val := range shortLinkArr {
		shortLink = fmt.Sprint(shortLink, string(val))
	}
	return shortLink
}
