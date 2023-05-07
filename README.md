# 项目架构
采用golang+Vue作为主开发语言

使用技术栈：golang wails作为GUI框架，前端采用Vue+Element plus

> 采用多线程数据抓取，提高程序抓取速度。Http请求采用异常重试和自动代理采集
# 目录结构
- backend：后端golang代码目录
- build：编译后的程序文件
- frontend：前端资源文件
- resource：静态资源文件

# 构建方法
```
 wails build
```

# 开发模式
```
 wails dev
```

# 配置
可修改`wails.json`对程序标题logo等配置，更多请参考： https://wails.io/docs/reference/project-config

# 预览
![app.png](https://raw.githubusercontent.com/auroraruanjian/go_mouyin/main/doc/app.png)