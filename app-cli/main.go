package main

import (
	"fmt"
	"gitee.com/lshsuper/go-boot/app-cli/cmd"
	"gitee.com/lshsuper/go-boot/define"
	"github.com/urfave/cli"
	"log"
	"os"
)

const (
	verStr = `
                                                               ___     
                               ,---,                         ,--.'|_   
              ,---.          ,---.'|      ,---.     ,---.    |  | :,'  
  ,----._,.  '   ,'\         |   | :     '   ,'\   '   ,'\   :  : ' :  
 /   /  ' / /   /   |        :   : :    /   /   | /   /   |.;__,'  /   
|   :     |.   ; ,. :        :     |,-..   ; ,. :.   ; ,. :|  |   |    
|   | .\  .'   | |: :        |   : '  |'   | |: :'   | |: ::__,'| :    
.   ; ';  |'   | .; :        |   |  / :'   | .; :'   | .; :  '  : |__  
'   .   . ||   :    |        '   : |: ||   :    ||   :    |  |  | '.'| 
 ----'| | \   \  /         |   | '/ : \   \  /  \   \  /   ;  :    ; 
 .'__/\_: |  '----'          |   :    |  '----'    '----'    |  ,   /  
 |   :    :                  /    \  /                        ----'   
\   \  /                   '-'----'                                  
   ----'                                                              ver:%s`
)

func main() {

	app := &cli.App{
		Name:        "app-cli",
		Description: "快速构建应用相关",
		Commands: []cli.Command{ //命令
			cmd.APICommand,
			cmd.MODELCommand,
		},
		Version: fmt.Sprintf(verStr, define.Version), //定义版本
		Action: func(ctx *cli.Context) {

		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
