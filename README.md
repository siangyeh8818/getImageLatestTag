# gdeyamlOperator
主要操作 :
- gettag : 返回該image 最新的tag並儲存到"getImageLatestTag_result.txt"這個檔案
- snapshot : 將k8s環境上的deploy,statefulset,daemonset,cronjob 輸出成gdeyaml格式
- nexus_api : 對Nexus repository 的API操作
- promote : 對Nexus repository 的資產進行搬移或是複製類型的搬移, 支援將gdeyaml格式的image 進行promote
- gitclone : 對於git clone的動作 , 參數以branch優先 , 若branch不存在會找tag 做clone
- git : git 相關操作 , 類似clone,checkoot,commit 等等
- replace : 將environment格式的replace內容取代gdeyaml的指定內容
- new-release : 為了開新的branch (gdeyaml & base的repo) , 並把new branch打上去gdeyaml文件
- imagedump : dump出k8s上的image 並可支援push


# 使用前須安裝
須能使用jq指令和docker指令 <br>
Centos:<br>
  yum -y install jq

# 用法
./getImageLatestTag --imagename dockerhub.pentium.network/grafana

| flag      | 說明    | 預設值     |
| ---------- | :-----------:  | :-----------: |
|  imagename    | docker image , such as dockerhub.pentium.network/grafana     | dockerhub.pentium.network/grafana    |
|  latest_mode    |  push or build , choose one mode to identify latest tag to you    | push    |
|  pull-pattern    |  (pull)pattern for imagename , ex: cr-{{stage}}.pentium.network/{{image}}:{{tag}}"    | null    |
|  push-pattern    |  (push)pattern for imagename , ex: cr-{{stage}}.pentium.network/{{image}}:{{tag}}    | null    |
|  stage    |  replace stage , new stage content    | null    |
|  inputfile    |  input file name , such as deploy.yml    | null    |
|  ouputfile    |  output file name , such as deploy-out.yml    | tmp_out.yml    |
|  user   |  user for docker login    | null   |
|  password   |  password for docker login    | null   |
|  push   |  push this image , default is false    | false   |
|  version   |  prints current binary version    | false   |
|  list    |  After sort tag list , we only deal with these top'number tags    | 5    |
