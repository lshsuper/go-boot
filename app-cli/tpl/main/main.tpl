package main

import(

   "gitee.com/lshsuper/go-boot/web"

)

func main(){


   //服务实例化
   s:=web.NewWebServer(web.WebConf{})

   //启动服务
   s.Run(":10086")


}
