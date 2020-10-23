package service

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wzru/gitran-server/config"
	"github.com/wzru/gitran-server/model"
)

var (
	urlNameReg = "^[A-Za-z0-9-]{1,32}$"
)

//checkURLName checks if a name is legal (to be in URL)
func checkURLName(name string) bool {
	ok, _ := regexp.Match(urlNameReg, []byte(name))
	return ok
}

func checkLang(lang *config.Lang) bool {
	for _, cfgLang := range config.Langs {
		if cfgLang.Code == lang.Code {
			return true
		}
	}
	return false
}

func checkLangs(langs []config.Lang) bool {
	for _, lang := range langs {
		ok := checkLang(&lang)
		if !ok {
			return false
		}
	}
	return true
}

func checkGitURLBranch(url string, br string) bool {
	//TODO
	return true
}

func createGitProj(ctx *gin.Context) error {
	return nil
}

//GetUserProjByName get a user project info
func GetUserProjByName(ctx *gin.Context, owner string, name string) *model.Project {
	user := model.GetUserByLogin(owner)
	if user == nil {
		return nil
	}
	self := (hasUserPermission(ctx, owner))
	return model.GetProjByOIDName(user.ID, name, self)
}

//GetOrgProjByName get an org project info
func GetOrgProjByName(ctx *gin.Context, owner string, name string) *model.Project {
	//TODO
	return nil
}

//GetProj get a project info
func GetProj(ctx *gin.Context) {
	owner := ctx.Param("owner")
	name := ctx.Param("project")
	proj := GetUserProjByName(ctx, owner, name)
	if proj == nil {
		proj = GetOrgProjByName(ctx, owner, name)
	}
	if proj == nil {
		ctx.JSON(http.StatusNotFound, model.Result{
			Success: false,
			Msg:     "Not Found",
			Data:    nil,
		})
		return
	}
	projInfo := model.ProjInfo{
		ID:        proj.ID,
		OwnerID:   proj.OwnerID,
		IsUsers:   proj.IsUsers,
		Name:      proj.Name,
		Desc:      proj.Desc,
		IsPrivate: proj.IsPrivate,
		IsGit:     proj.IsGit,
		GitURL:    proj.GitURL,
		GitBranch: proj.GitBranch,
		SyncTime:  proj.SyncTime,
		SrcLangs:  model.GetLangsFromString(proj.SrcLangs),
		TgtLangs:  model.GetLangsFromString(proj.TgtLangs),
		CreatedAt: proj.CreatedAt,
		UpdatedAt: proj.UpdatedAt,
	}
	ctx.JSON(http.StatusNotFound, model.Result{
		Success: true,
		Msg:     "",
		Data: gin.H{
			"proj_info": projInfo,
		},
	})
}

//CreateUserProj create a new user project
func CreateUserProj(ctx *gin.Context) {
	name := ctx.PostForm("name")
	isUsers := true
	desc := ctx.PostForm("desc")
	userID, _ := strconv.Atoi(ctx.GetString("user-id"))
	isPrvt := ctx.PostForm("is_private") == "true"
	isGit := ctx.PostForm("is_git") == "true"
	gitURL := ctx.PostForm("git_url")
	gitBranch := ctx.PostForm("git_branch")
	syncTime, _ := strconv.Atoi(ctx.PostForm("sync_time"))
	src := ctx.PostForm("src_langs")
	tgt := ctx.PostForm("tgt_langs")
	srcLangs := model.GetLangsFromString(src)
	tgtLangs := model.GetLangsFromString(tgt)
	if !checkURLName(name) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "名字不合法",
			Data:    nil,
		})
		return
	}
	if len(srcLangs) == 0 {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "源语言不能为空",
			Data:    nil,
		})
		return
	}
	if !checkLangs(srcLangs) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "源语言不合法",
			Data:    nil,
		})
		return
	}
	if !checkLangs(tgtLangs) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "目标语言不合法",
			Data:    nil,
		})
		return
	}
	if !checkGitURLBranch(gitURL, gitBranch) {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "Git URL不合法",
			Data:    nil,
		})
		return
	}
	proj := &model.Project{
		Name:      name,
		Desc:      desc,
		OwnerID:   uint64(userID),
		IsUsers:   isUsers,
		IsPrivate: isPrvt,
		IsGit:     isGit,
		GitURL:    gitURL,
		GitBranch: gitBranch,
		SyncTime:  uint64(syncTime),
		SrcLangs:  src,
		TgtLangs:  tgt,
	}
	if findProj := model.GetProjByOIDName(proj.OwnerID, proj.Name, true); findProj != nil {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     "项目已存在",
			Data:    nil,
		})
		return
	}
	proj, err := model.NewProj(proj)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.Result{
			Success: false,
			Msg:     err.Error(),
			Data:    nil,
		})
		return
	} else {
		ctx.JSON(http.StatusCreated, model.Result{
			Success: true,
			Msg:     "创建项目成功",
			Data:    nil,
		})
	}
}

//CreateOrgProj create a new organization project
func CreateOrgProj(ctx *gin.Context) {
	//TODO
}
