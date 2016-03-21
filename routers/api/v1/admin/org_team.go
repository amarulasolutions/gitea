// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package admin

import (
	api "github.com/gogits/go-gogs-client"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/modules/context"
	"github.com/gogits/gogs/routers/api/v1/convert"
	"github.com/gogits/gogs/routers/api/v1/user"
)

func ListTeams(ctx *context.APIContext) {
	org := user.GetUserByParamsName(ctx, ":orgname")
	if ctx.Written() {
		return
	}

	if err := org.GetTeams(); err != nil {
		ctx.Error(500, "GetTeams", err)
		return
	}

	apiTeams := make([]*api.Team, len(org.Teams))
	for i := range org.Teams {
		apiTeams[i] = convert.ToTeam(org.Teams[i])
	}
	ctx.JSON(200, apiTeams)
}

func CreateTeam(ctx *context.APIContext, form api.CreateTeamOption) {
	org := user.GetUserByParamsName(ctx, ":orgname")
	if ctx.Written() {
		return
	}

	team := &models.Team{
		OrgID:       org.Id,
		Name:        form.Name,
		Description: form.Description,
		Authorize:   models.ParseAccessMode(form.Permission),
	}
	if err := models.NewTeam(team); err != nil {
		if models.IsErrTeamAlreadyExist(err) {
			ctx.Error(422, "NewTeam", err)
		} else {
			ctx.Error(500, "NewTeam", err)
		}
		return
	}

	ctx.JSON(200, convert.ToTeam(team))
}