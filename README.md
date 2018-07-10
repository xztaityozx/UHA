# UHA
Ultra H○PICE Attacker

できるだけ楽したい！

H○SPICEでモンテカルロ・シミュレーションをうまくやるツールです

## Require
- go (>=1.10.2)

## Install
Goが無い環境ならGoをインストールしたほうが更新が楽になります

### Ubuntu
```sh
$ sudo apt update && sudo apt upgrade
$ sudo apt install golang
$ export GOPATH=~/.go
$ echo 'export GOPATH=~/.go' >> ~/.bashrc
$ export PATH=$PATH:$GOPATH
$ echo 'export PATH=$PATH:$GOPATH' >> ~/.bashrc
$ go get github.com/xztaityozx/UHA
$ cd $GOPATH/src/github.com/xztaityozx/UHA
$ go install
```

## Usage
```sh
UHA [--config FILE|-h,--help]  [command] OPTIONS
```

### command
#### make
シミュレーションセットを作成します

```sh
UHA [--config FILE] make [flags]
```

Flags:
  -D, --default       設定ファイルをそのままタスクにします。オプションで値をしているとそちらが優先されます
  -h, --help          help for make
      --sigma float   Vtp,VtnのSigmaを設定します

#### run
makeで作ったシミュレーションセットを実行します。

```sh
UHA [--config FILE] run [flags]
```
Flags:
  -C, --continue     連続して実行する時、どれかがコケても次のシミュレーションを行います
  -h, --help         help for run
  -n, --number int   実行するシミュレーションセットの個数です (default 1)


## Configuration
UHAを快適に使うには設定が必要です。設定ファイルは`.UHA.json`で、通常`~/.config/UHA/`以下に置かれており、これがデフォルトの設定として扱われます。

各コマンドの`--config`オプションの引数に設定ファイルを設定すると、デフォルトの代わりに設定ファイルとして読み込みます。

### フォーマット
```json
{
  "Simulation":{
    "Monte":["100", "200", "500", "1000", "2000", "5000", "10000", "20000", "50000"],
    "Range":{
      "Start":"2.5ns",
      "Stop":"17.5ns",
      "Step":"7.5ns"
    },
    "Signal":"N2",
    "SimDir":"",
    "DstDir":"",
    "Vtp":{
      "Voltage":-0.6,
      "Sigma":0.0,
      "Deviation":1.0
    },
    "Vtn":{
      "Voltage":0.6,
      "Sigma":0.0,
      "Deviation":1.0
    }
  },
  "TaskDir":"~/.config/UHA/task",
  "Repository":[
    {
      "Path":"",
      "Type",
    }
  ],
  "SpreadSheet":{
    "Id":"xxxxxxxxxxxxxxxxxxxxxxx",
    "CSPath":"/path/to/client_secret.json",
    "TokenPath":"/path/to/token.json"
  }
}
```

### Simuration
|項目名|型|説明|
|:--:|:--:|:--:|
|Monte|stringのリスト|一個のシミュレーションセットに含めたいモンテカルロの回数のリスト|
|Range|Range型|CSVの書き出しに使うデータ|
|Start|string|書き出しの開始時間|
|Stop|string|書き出しの終了時間|
|Step|string|書き出しの刻み幅|
|Signal|string|プロットする信号線名|
|SimDir|string|netlistというディレクトリがH○SPICEによって生成されるので、そこへのパス|
|DstDir|string|結果を書き出したいディレクトリへのパス|
|Vtp|Node型|Vtpの設定値|
|Vtn|Node型|Vtnの設定値|
|Voltage|float|AGAUSS(ここ,s,m)|
|Sigma|float|AGAUSS(v,ここ,m)|
|Deviation|float|AGAUSS(v,s,ここ)|


### TaskDir
シミュレーションセットが保存されるディレクトリです。デフォルトでは`$HOME/.config/UHA/task`です

### Repository
リポジトリを複数指定できます。

|項目名|型|説明|
|:--:|:--:|:--:|
|Type|string|リポジトリタイプです。Git,Dir,AWSが指定できます。GitとDirはほぼ同じです|
|Path|string|リポジトリへのパスです。結果を集積したい具体的なパスを書くようにしてください|

例えば`Path`が`~/WorkSpace`であると、データは`~/WorkSpace`以下に蓄積されます。

### SpreadSheet
Google SpreadSheetのAPIを使って結果を書き込むことが出来ます。それ用の設定です

|項目名|型|説明|
|:--:|:--:|:--:|
|Id|string|SpreadSheetのIDです|
|CSPath|string|client_secret.jsonへのパスです|
|TokenPath|string|token.jsonへのパスです|


## input.spiのテンプレートを準備する
この項目をクリアすれば使えるようになります。

まず、`netlist`というディレクトリがH○SPICEによって作られているはずなので、そこへ移動します。

そのディレクトリにある`input.spi`を、設定ファイルのあるディレクトリに置きます。この時ファイル名を`spitemplate.spi`にします。

次に`spitemplate.spi`の以下の行を編集します。

|行|変更前|変更後|
|だいたい8行目|.param vtn=AGAUSS(xx,yy,zz) vtp=AGAUSS(xx,yy,zz)|.param vtn=AGAUSS( __%.4f,%.4f,%.4f__ ) vtp=AGAUSS( __%.4f,%.4f,%.4f__ )|
|だいたい53行目|.tran 10p 20n start=0 uic sweep monte=xxx firstrun=1|.tran 10p 20n start=0 uic sweep __monte=%s__ firstrun=1|

編集したら保存して終わりです



