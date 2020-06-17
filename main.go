package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql"
	_ "github.com/GoAdminGroup/themes/adminlte"
	_ "github.com/GoAdminGroup/themes/sword"
	"github.com/moruowait/admin/conn"
	"github.com/moruowait/admin/homestay"

	"github.com/GoAdminGroup/components/echarts"
	ada "github.com/GoAdminGroup/go-admin/adapter/gin"
	adminContext "github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/gin-gonic/gin"
	homestaypages "github.com/moruowait/admin/homestay/pages"
	"github.com/moruowait/admin/login"
	"github.com/moruowait/admin/pages"
	"github.com/moruowait/admin/tables"
)

func main() {
	r := gin.Default()

	eng := engine.Default()

	var loginComps = map[string]template.Component{
		"default": login.Get(),
	}

	template.AddLoginComp(loginComps["default"])
	template.AddComp(chartjs.NewChart())
	template.AddComp(echarts.NewChart())

	rootPath := "/Users/Wushaoqing/data/www/go-admin"
	//rootPath := "."

	cfg := config.ReadFromJson(rootPath + "/config.json")
	cfg.CustomFootHtml = template.HTML(`<div style="display:none;">
    <script type="text/javascript" src="https://s9.cnzz.com/z_stat.php?id=1278156902&web_id=1278156902"></script>
	<!-- Global site tag (gtag.js) - Google Analytics -->
	<script async src="https://www.googletagmanager.com/gtag/js?id=UA-103003647-2"></script>
	<script>
	window.dataLayer = window.dataLayer || [];
	function gtag(){dataLayer.push(arguments);}
	gtag('js', new Date());
	gtag('config', 'UA-103003647-2');
	</script>
</div>`)
	cfg.CustomHeadHtml = template.HTML(`<link rel="icon" type="image/png" sizes="32x32" href="//quick.go-admin.cn/official/assets/imgs/icons.ico/favicon-32x32.png">
        <link rel="icon" type="image/png" sizes="96x96" href="//quick.go-admin.cn/official/assets/imgs/icons.ico/favicon-64x64.png">
        <link rel="icon" type="image/png" sizes="16x16" href="//quick.go-admin.cn/official/assets/imgs/icons.ico/favicon-16x16.png">`)

	cfg.Animation = config.PageAnimation{
		Type:     "fadeInUp",
		Duration: 0.9,
	}
	cfg.AddUpdateProcessFn(func(values form.Values) (values2 form.Values, e error) {
		if values.Get("theme") == "adminlte" && values.Get("asset_url") == "//quick.go-admin.cn/demo/sword" {
			values.Add("asset_url", "//quick.go-admin.cn/demo")
		}
		if values.Get("theme") == "sword" && values.Get("asset_url") == "//quick.go-admin.cn/demo" {
			values.Add("asset_url", "//quick.go-admin.cn/demo/sword")
		}
		if values.Get("site_off") == "true" || values.Get("no_limit_login_ip") == "false" {
			return nil, errors.New("不允许的操作")
		}
		if values.Get("login_title") != "GoAdmin" {
			return nil, errors.New("不允许的操作")
		}
		if values.Get("custom_head_html") != string(cfg.CustomHeadHtml) {
			return nil, errors.New("不允许的操作")
		}
		if values.Get("custom_foot_html") != string(cfg.CustomFootHtml) {
			return nil, errors.New("不允许的操作")
		}
		if values.Get("footer_info") != "" || values.Get("login_logo") != string(cfg.LoginLogo) ||
			values.Get("logo") != string(cfg.Logo) || values.Get("mini_logo") != string(cfg.MiniLogo) {
			return nil, errors.New("不允许的操作")
		}
		if e := values.Get("extra"); e != "" {
			var extra = make(map[string]interface{})
			err := json.Unmarshal([]byte(e), &extra)
			if err != nil && extra["login_theme"] != "" {
				if comp, ok := loginComps[extra["login_theme"].(string)]; ok {
					template.AddLoginComp(comp)
				}
			}
		}
		return values, nil
	})
	//cfg.HideConfigCenterEntrance = true

	if err := eng.AddConfig(cfg).
		AddGenerators(tables.Generators).
		AddGenerators(homestay.Generators).
		AddGenerator("user", tables.GetUserTable).
		// AddPlugins(filemanager.NewFileManagerWithConfig(filemanager.Config{
		// 	Path:          "/data/www/go-admin/fm_example",
		// 	AllowDelete:   false,
		// 	AllowUpload:   true,
		// 	AllowDownload: true,
		// 	AllowRename:   true,
		// 	AllowMove:     true,
		// }), librarian.NewLibrarianWithConfig(librarian.Config{
		// 	Path:      "/data/www/go-admin/fm_example/markdown",
		// 	BuildMenu: false,
		// 	Prefix:    "librarian",
		// })).
		AddNavButtons("网站信息", "", action.PopUp("/website/info", "网站信息",
			func(ctx *adminContext.Context) (success bool, msg string, data interface{}) {
				return true, "ok", `<p>网站由 <a href="https://github.com/chenhg5">cg33<a/> 创造</p>`
			})).
		AddNavButtons("用户管理", "", action.Jump("/admin/info/manager")).
		Use(r); err != nil {
		panic(err)
	}

	r.Static("/uploads", rootPath+"/uploads")

	// you can custom your pages like:

	r.GET("/admin", ada.Content(func(ctx *gin.Context) (panel types.Panel, e error) {
		if config.GetTheme() == "adminlte" {
			return pages.GetDashBoardContent(ctx)
		} else {
			return pages.GetDashBoard2Content(ctx)
		}
	}))
	r.GET("/admin/echarts", ada.Content(pages.GetDashBoard3Content))
	r.GET("/admin/table", ada.Content(pages.GetTableContent))

	r.GET("/admin/homestay/search_dashboard", ada.Content(homestaypages.GetHomestaySearchDashboard))

	r.GET("/admin/form1", ada.Content(pages.GetForm1Content))
	eng.Data("POST", "/admin/form/update", func(ctx *adminContext.Context) {
		fmt.Println("ctx.PostForm()", ctx.PostForm())
		ctx.PjaxUrl("/admin")
	})
	eng.Data("POST", "/admin/popup/form", func(ctx *adminContext.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"code": 200,
			"msg":  "ok",
			"data": map[string]string{
				"url": "/admin/info/profile",
			},
		})
	})

	//eng.DB("mysql").Query()

	eng.Data("GET", "/admin/homestay/item", func(ctx *adminContext.Context) {
		res, err := eng.DB("mysql").Query("select * from homestay_spend_item")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, nil)
			return
		}
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"code": 200,
			"msg":  "ok",
			"data": res,
		})
	})

	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/admin")
	})

	// 设置全局数据库连接
	conn.SetConn(eng.MysqlConnection())

	srv := &http.Server{
		Addr:    ":9033",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

	log.Print("closing database connection")
	eng.MysqlConnection().Close()

}
