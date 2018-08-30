package user

import (
	api "gitote.com/gitote/go-gitote-client"

	"gitote.com/gitote/gitote/models"
	"gitote.com/gitote/gitote/pkg/context"
)

func responseApiUsers(c *context.APIContext, users []*models.User) {
	apiUsers := make([]*api.User, len(users))
	for i := range users {
		apiUsers[i] = users[i].APIFormat()
	}
	c.JSON(200, &apiUsers)
}

func listUserFollowers(c *context.APIContext, u *models.User) {
	users, err := u.GetFollowers(c.QueryInt("page"))
	if err != nil {
		c.Error(500, "GetUserFollowers", err)
		return
	}
	responseApiUsers(c, users)
}

func ListMyFollowers(c *context.APIContext) {
	listUserFollowers(c, c.User)
}

func ListFollowers(c *context.APIContext) {
	u := GetUserByParams(c)
	if c.Written() {
		return
	}
	listUserFollowers(c, u)
}

func listUserFollowing(c *context.APIContext, u *models.User) {
	users, err := u.GetFollowing(c.QueryInt("page"))
	if err != nil {
		c.Error(500, "GetFollowing", err)
		return
	}
	responseApiUsers(c, users)
}

func ListMyFollowing(c *context.APIContext) {
	listUserFollowing(c, c.User)
}

func ListFollowing(c *context.APIContext) {
	u := GetUserByParams(c)
	if c.Written() {
		return
	}
	listUserFollowing(c, u)
}

func checkUserFollowing(c *context.APIContext, u *models.User, followID int64) {
	if u.IsFollowing(followID) {
		c.Status(204)
	} else {
		c.Status(404)
	}
}

func CheckMyFollowing(c *context.APIContext) {
	target := GetUserByParams(c)
	if c.Written() {
		return
	}
	checkUserFollowing(c, c.User, target.ID)
}

func CheckFollowing(c *context.APIContext) {
	u := GetUserByParams(c)
	if c.Written() {
		return
	}
	target := GetUserByParamsName(c, ":target")
	if c.Written() {
		return
	}
	checkUserFollowing(c, u, target.ID)
}

func Follow(c *context.APIContext) {
	target := GetUserByParams(c)
	if c.Written() {
		return
	}
	if err := models.FollowUser(c.User.ID, target.ID); err != nil {
		c.Error(500, "FollowUser", err)
		return
	}
	c.Status(204)
}

func Unfollow(c *context.APIContext) {
	target := GetUserByParams(c)
	if c.Written() {
		return
	}
	if err := models.UnfollowUser(c.User.ID, target.ID); err != nil {
		c.Error(500, "UnfollowUser", err)
		return
	}
	c.Status(204)
}
