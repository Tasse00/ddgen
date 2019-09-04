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

- 导出数据项定义文件
  
  例子
  ``` shell
  $ ddgen -i mysql5.7 -S "root:root@tcp(127.0.0.1:3306)/?charset=utf8" -s pocket_reading -o dat.json
  ```

  使用说明
  ```
  Usage of ./ddgen:
    -S string
          database source
    -h    show this help
    -i string
          inspector (default "mysql5.7")
    -o string
          output data file (json) (default "dat.json")
    -p string
          params pass to inspector.
    -s string
          schema
  ```

- 渲染数据项至指定格式

  例子
  ``` shell
  $ ddrender -d dat.json -t office-word -o res.docx
  ```

  使用说明
  ```
  Usage of ./ddrender:
    -d string
          data filepath (default "./dat.json")
    -h    show this help
    -o string
          output filepath (default "./out")
    -p string
          additional params that renderer need
    -t string
          render type, one of office-word,md
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
