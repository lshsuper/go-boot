package cmd

import (
	"fmt"
	"gitee.com/lshsuper/go-boot/app-cli/tpl"
	"gitee.com/lshsuper/go-boot/app-cli/tpl/model"
	"gitee.com/lshsuper/go-boot/database"
	"gitee.com/lshsuper/go-boot/utils"
	"github.com/urfave/cli"
	"log"
)

//MODELCommand gorm模型构建命令
var MODELCommand = cli.Command{

	Name:        "model",
	Description: "根据数据库批量构建gorm模型",
	Usage:       "根据数据库批量构建gorm模型",
	ShortName:   "m",
	Action: func(ctx *cli.Context) {
		defer func() {

			if err := recover(); err != nil {
				log.Println(utils.ErrorToString(err))
			}

		}()
		var (
			host   = ctx.String("host")
			port   = ctx.Int("port")
			dbname = ctx.String("dbname")
			uname  = ctx.String("uname")
			pwd    = ctx.String("pwd")
		)
		if len(host) <= 0 {
			panic("host不能为空")
		}
		if port <= 0 {
			panic("port字段不能为0")
		}
		if len(dbname) <= 0 {
			panic("dbname不能为空")
		}
		if len(uname) <= 0 {
			panic("uname不能为空")
		}
		if len(pwd) <= 0 {
			panic("pwd不能为空")
		}

		db := database.Register(database.Mysql, database.DbConfig{
			Host:      host,
			Port:      port,
			UserName:  uname,
			Password:  pwd,
			DefaultDb: dbname,
		})
		tbs := db.GetTables(dbname)

		for _, tb := range tbs {
			tpl := tpl.ModelTpl.Tpl()
			arg := model.ModelArg{
				Descr:  tb.TableComment,
				TName:  tb.TableName,
				Fields: []model.ModelArgField{},
			}
			fields := db.GetColumns(tb.TableName)
			for _, field := range fields {
				arg.Fields = append(arg.Fields, model.ModelArgField{
					Column: field.Field,
					DType:  field.GetGoType(),
					Descr:  field.Comment,
				})
			}

			tpl.SArg(arg)

			//导入到项目下模型层
			tpl.Input(fmt.Sprintf("./%s", modelRootPkgName))
		}

	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "host",
			Value: "localhost",
			Usage: "数据库HOST",
		},
		&cli.IntFlag{
			Name:  "port",
			Value: 3306,
			Usage: "数据库端口",
		},
		&cli.StringFlag{
			Name:  "dbname",
			Value: "db",
			Usage: "数据库名称",
		},
		&cli.StringFlag{
			Name:  "uname",
			Value: "user",
			Usage: "用户名",
		},
		&cli.StringFlag{
			Name:  "pwd",
			Value: "password",
			Usage: "密码",
		},
	},
}
