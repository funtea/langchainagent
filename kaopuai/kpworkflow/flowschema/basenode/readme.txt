test.json为测试json

literal 的 content 改为  literalContent

list 的 schema 改为  listSchema

object 的 schema 改为  objectSchema

javascript   替换  typescript格式

datasetList  改为  StringArrayContent

topK  minScore  strategy的literalContent改为字符串类型

code,llm,condition,plugins  输出的outputMap在scriptResultJson字段里，是json字符串。

json字符串进行解析时对列表需要进行：  key4.0.key45.key451