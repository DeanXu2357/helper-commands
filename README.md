# helper-commands

## 安裝套件  
`go install github.com/DeanXu2357/helper-commands@latest`

## arangodb-repo-gen 

arangodb 的 curd repo 生成器，根據傳入 entity struct 生成對應的 repository 。

### 使用方法

#### 新增 repo 程式碼
* 執行  
`helper-commands create [entity import] [entity name]`  
ex:  
`go run main.go  create github.com/DeanXu2357/arangodb-repo-gen/User User`

#### 新增 migration 創建 collection
* 執行
`helper-commands collectionMigration [collection name]`
