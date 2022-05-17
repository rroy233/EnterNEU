<p align="center">
    <img src="https://socialify.git.ci/rroy233/EnterNEU/image?description=1&descriptionEditable=%E4%B8%9CB%E5%A4%A7%E5%AD%A6e%E7%A0%81%E9%80%9A%E7%94%9F%E6%88%90(%E5%90%8E%E7%AB%AF)&language=1&logo=https%3A%2F%2Fs2.loli.net%2F2022%2F05%2F05%2FDx8lpyk7mcuZ3Tr.png&name=1&owner=1&pattern=Circuit%20Board&theme=Light">
</p>


<p align="center">
   <a href="https://github.com/rroy233/EnterNEU">
      <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/rroy233/EnterNEU?style=flat-square">
   </a>
   <a href="https://github.com/rroy233/EnterNEU/releases">
      <img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/rroy233/EnterNEU?style=flat-square">
   </a>
   <a href="https://github.com/rroy233/EnterNEU/blob/main/LICENSE">
      <img alt="GitHub license" src="https://img.shields.io/github/license/rroy233/EnterNEU?style=flat-square">
   </a>
   <a href="https://github.com/rroy233/EnterNEU">
      <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/rroy233/EnterNEU/Go?style=flat-square">
   </a>
   <a href="https://github.com/rroy233/EnterNEU/commits/main">
      <img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/rroy233/EnterNEU?style=flat-square">
   </a>
   <a href="https://enterneu.icu">
      <img alt="Demo" src="https://img.shields.io/endpoint?url=https%3A%2F%2Fupload.enterneu.icu%2Fshields.json">
   </a>
</p>


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

1. 克隆仓库

   ```shell
   git clone https://github.com/rroy233/EnterNEU.git
   ```
   
2. 获取可执行文件

   1. 自行编译

      ```shell
      cd EnterNEU/
      # 自行编译
      # go版本要求：go1.17
      go build -o enterneu
      # 或者使用make交叉编译
      make
      ```

   2. 前往release下载

      下载已编译的[可执行文件](https://github.com/rroy233/EnterNEU/releases)，重新命名为`enterneu`，放于项目文件夹内。

3. 编辑config.yaml

   ```shell
   # 新建config.yaml,编辑内容
   cp config.example.yaml config.yaml
   vim config.yaml
   ```

4. 运行

   ```shell
   ./enterneu
   # 或使用后台运行脚本
   bash run.sh
   ```

5. 访问`http://localhost:8994`使用。

   可以使用nginx进行反代



### License

GPL-3.0 license
