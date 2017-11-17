# gCrawl
简单爬虫的编写

编译 gCrawTest文件夹下的src.go文件，生成src
配置config.json文件，字段说明:
Mainurl:想爬取的主页面 Header:反爬虫页面将限制访问，有时需加头部，暂未处理
Keyword：想爬取的种子等关键字 RountineNum:协程数目，默认0则为20个

修改gParseLinks/gParseLinks.go 函数parseDetail 定义自定义抓取规则，因每个页面布局不一致
