## 使用介绍  

db-rest，无需编写后端代码，即可将mysql,postgresql,sqlite,oracle,mssql等关系型数据库快速暴露为功能强大的 REST API。让数据共享和集成变得前所未有的简单。

```
初始化工作目录  

db-rest init 工作目录

在工作目录中配置db to api。

抽取数据

db-rest extract 工作目录


启动服务

db-rest serve 工作目录



``` 


##  接口详情  

开启验证时获取token地址 POST /api/internal/token  

```
curl --location 'http://localhost:8080/api/internal/token' \
--form 'name="admin"' \
--form 'pass="admin"'

```
查看工作空间详情 GET /api/internal/token/workspace/list
```

curl --location 'http://localhost:8080/api/internal/workspace/list'

```
查询数据  POST /api/集合名/表配置名/list  

```

curl --location 'http://localhost:8080/api/question/question/list' \
--header 'Content-Type: application/json' \
--data '{
    "pageSize": 100,
    "page": 1,
    "query": {
       "query_type":"and",
       "queries":[
            {
            "field": "input",
            "op": "like",
            "value":"大模型"
          },
            {
            "field": "instruct",
            "op": "notnull",
            "value":""
          },
          {
            "field": "output",
            "op": "eq",
            "value":"应受善保存放置，做好衬垫、绑扎工作，防止翻倒、振动、碰撞。"
          }

       ]

    }
}'
查询不建议使用and模式，该实现采用子查询方式设计，性能会比较差,小量数据可以使用。
curl --location 'http://localhost:8080/api/question/question/list' \
--header 'Content-Type: application/json' \
--data '{
    "pageSize": 100,
    "page": 1,
    "query": {
       "queries":[
            {
            "field": "input",
            "op": "like",
            "value":"大模型"
          }

       ]

    }
}'
默认使用or模式。
如需要实现更多复杂的查询，可以自行不采用sql模式实现查询引擎。
```  

