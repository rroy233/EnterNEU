

### 功能

* 自定义姓名、学号、校门名称等
* 上传头像
* 加密用户数据
* 用户自行管理凭证(删除、设置过期时间)
* 支持配合Shadowrocket使用
* 支持TelegramBot



### 使用

请阅读：[常见问题及使用教程](https://github.com/rroy233/EnterNEU/blob/main/assets/tips.md)



### 自行搭建

#### 环境需求

- go1.17（自行编译需要）
- redis（存储用户数据）

#### 步骤

1. 获取可执行文件

   1. 克隆仓库

      ```shell
      git clone https://github.com/rroy233/EnterNEU.git
      ```

   2. 自行编译或前往release下载

      ```shell
      # 自行编译
      # go版本要求：go1.17
      go build -o enterneu
      # 或者使用make交叉编译
      make
      ```

      下载已编译的[可执行文件](https://github.com/rroy233/EnterNEU/releases)，重新命名为`enterneu`，放于项目文件夹内。

2. 编辑config.yaml

   ```shell
   # 新建config.yaml,编辑内容
   cp config.example.yaml config.yaml
   vim config.yaml
   ```

3. 运行

   ```shell
   ./enterneu
   # 或使用后台运行脚本
   bash run.sh
   ```

4. 访问`http://localhost:8994`使用。

   可以使用nginx进行反代



### License

GPL-3.0 license
