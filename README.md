# ddgen
用于生成DB中数据项说明文件的工具。

## 目标
支持不同类型、版本数据库且能导出多种类型文件。

## 预览

- Markdown

  ![](http://processon.com/chart_image/5d6f3748e4b01080c7333e06.png?_=1567570006592)

- Office Word

  ![](http://processon.com/chart_image/5d6f372fe4b04a195014f76a.png?_=1567569962969)
 
 
## 使用

详细说明见帮助`ddgen -h`

- 直接生成数据项说明文件

  ``` shell
  $ ddgen export --inspector mysql5.7 --source "root:root@tcp(127.0.0.1:3306)/?charset=utf8" --schema db --renderer md --out out.md
  ```

- 仅导出数据项定义文件
  
  ``` shell
  $ ddgen inspect --inspector mysql5.7 --source "root:root@tcp(127.0.0.1:3306)/?charset=utf8" --schema db --out spec.json
  ```

- 将已有数据项定义文件渲染至指定格式

  ``` shell
  $ ddgen render --dat spec.json --renderer md --out spec.md
  ```


## 结构及主要部件

![部件结构](http://processon.com/chart_image/5d6f3296e4b0c5c942b59e78.png?_=1567568950321)

- Inspector 数据项检查器

  对于不同类型、版本的数据库采用不同的inspector实现。
  目前已实现的Inspector为:

  - mysql5.7

- Renderer 数据项渲染器
  
  已实现的渲染器为:
  
  - office-word
  - markdown
