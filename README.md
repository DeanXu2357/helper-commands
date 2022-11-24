# helper-commands

## arangodb-repo-gen 

arangodb 的 curd repo 生成器，根據傳入 entity struct 生成對應的 repository 。

### 使用方法

* 安裝套件  
`go install github.com/DeanXu2357/arangodb-repo-gen@latest`
* 執行  
`arangodb-repo-gen create [entity import] [entity name]`  
ex:  
`go run main.go  create github.com/DeanXu2357/arangodb-repo-gen/entity User`

