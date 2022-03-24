package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/go-github/v43/github"
	"golang.org/x/time/rate"
	"gopkg.in/rx.v0"
)

func main() {
	client := github.NewClient(nil)
	rl := rate.NewLimiter(rate.Every(time.Second), 1)

	fetchTopContributors := func(owner string, repository string) rx.Observable[*github.Contributor] {
		return rx.Defer(func(subscriber rx.Lifecycle) rx.Observable[*github.Contributor] {
			ctx := subscriber.Context(nil)
			rl.Wait(ctx)
			contribs, _, err := client.Repositories.ListContributors(ctx, owner, repository, nil)
			if err != nil {
				return rx.Error[*github.Contributor](err)
			}
			return rx.List(contribs)
		})
	}

	fetchUserFollowings := func(contrib *github.Contributor) rx.Observable[*github.User] {
		return rx.Func(func(subscriber rx.Writer[*github.User]) (err error) {
			if contrib.Login == nil {
				return
			}
			ctx := subscriber.Context(nil)
			rl.Wait(ctx)
			followings, _, err := client.Users.ListFollowing(ctx, *contrib.Login, nil)
			if err != nil {
				return
			}
			for _, following := range followings {
				if !subscriber.Write(following) {
					return
				}
			}
			return
		})
	}

	downloadUserProfileImages := func(user *github.User) rx.Observable[[]byte] {
		return rx.Func(func(subscriber rx.Writer[[]byte]) (err error) {
			if user.AvatarURL == nil {
				return
			}
			resp, err := http.Get(*user.AvatarURL)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			image, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}
			if !subscriber.Write(image) {
				return
			}
			return
		})
	}

	contribs := fetchTopContributors("ReactiveX", "rxjs")
	followings := rx.MergeMap(contribs, fetchUserFollowings, rx.MergeMapWithStandbyConcurrency(5))
	images := rx.MergeMap(followings, downloadUserProfileImages, rx.MergeMapWithStandbyConcurrency(10))

	writer, reader := rx.Pipe[[]byte]()
	images.Subscribe(writer)

	for {
		if image, ok := reader.Read(); ok {
			fmt.Printf("Image Size: %d Hash: %s\n", len(image), hex.EncodeToString(SHA256(image)))
		} else {
			break
		}
	}

	if err := reader.Wait(); err != nil {
		panic(err)
	}
}

func SHA256(b []byte) []byte {
	m := sha256.Sum256(b)
	return m[:]
}
