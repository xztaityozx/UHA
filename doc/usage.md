## UHA usage

`UHA` の使い方です



---



- [基本]()

- [サブコマンド]()
  - [count]()
  - [get]()
  - [help]()
  - [list]()
  - [make]()
  - [push]()
  - [rmake]()
  - [run]()
  - [smake]()
  - [srun]()
  - [upgrade]()
  - [version]()





## 基本

```sh
UHA [--config [cfgFile]] [command] [flag]
```



`UHA` は複数のサブコマンドを持つツールです。それぞれに1つの役割があります





##　サブコマンド

### count

```sh
UHA count [--help]
```



現在のディレクトリにあるCSVファイルを見つけて、不良数を数え上げます。



### get

```sh
UHA get [--from,-f [path]|--help|-P,--push]
```



- `-f,--from`
  - データがあるディレクトリを明示的に指定します
- `-P,--push`
  - 数え上げた後自動で`UHA push`を実行します



### help

```sh
UHA help [--help]
```



ヘルプを表示します

### list
```sh
UHA list [--one,-1|--long,-l]
```

読み込んだコンフィグで設定されている"TaskDir"以下にあるタスクをリストアップします

- `-1,--one` 
  - データを1列に表示します
- `-l,--long`
  - 作成日時も出力します


### make

```sh
UHA make [-D,--default|--sigma [sigma]|-y,--yes|-h,--help]
```



シミュレーションタスクを作るコマンドです。オプション無しで実行すると、インタラクティブにタスクを生成できます。タスクは `~/.config/UHA/task`以下に保存されます



- `-D,--default`
  - `.UHA.json`の設定値をそのままでタスクを発行します
- `--sigma`
  - シグマの値を引数で指定できます。`-D,--default`の値を上書きします
- `-y,--yes`
  - 最後の確認を飛ばすことが出来ます



### push

```sh
UHA push [--help]
```



Google SpreadSheetにデータを書き込みます。書き込むデータは最後に`UHA count`で得たデータです。



### rmake

```sh
UHA rmake [-n,--number] start step stop/number
```



始点、刻み幅、終点をシグマを変えながらタスクを連続発行します。シグマ以外は`~/.UHA.json`の値が使用されます



- `-n,--number`
  - 始点、刻み幅、回数を指定してタスクを発行します
  - `UHA rmake -n 1 2 10`は __1から2刻みで10個タスクを作る__ という意味になります



### run

```sh
UHA run [--all|-C,--continue|-h,--help|-n,number [num]]
```



発行済のタスクを実行します。



- `--all`
  - 発行済のタスクをすべて実行します
- `-C,--continue`
  - 各シミュレーションのうちどれかが、失敗しても他のシミュレーションを継続します
- `-n,--number [num]`
  - `num`個のタスクを実行します



### smake

```sh
UHA smake [-h,--help|-n,--number [num]|-S,--sigma [sigma]|-T,--times [t]]
```



1つのモンテカルロを複数回するタスクを発行します。それぞれのSEED値は連番になります



- `-n,--number [num]`
  - モンテカルロを`num`回行います
- `-S,--sigma [sigma]`
  - シグマの値を指定します
- `-T,--times [t]`
  - `t`回のモンテカルロを行います



### srun

```sh
UHA srun [--all|-C,--continue|-h,--help|-n,--number [num]|-P,--parallel [num]]
```



`smake`で作ったタスクを実行します



- `--all` 
  - すべてのタスクを実行します
- `-C,--continue`
  - 各シミュレーションのうちどれかが、失敗しても他のシミュレーションを継続します
- `-n,--number [num]`
  - `num` 個のタスクを実行します
- `-P,--parallel [num]`
  - シミュレーションを`num`個並列で動かします
  - PCのスペックと相談して __慎重に__ 決めてください
- `--GC`
  - CSV以外の中間データを削除します
- `--no-notify`
  - Slackに通知しないようにします
- `--start num`
  - SEEDの初期値をnumに設定します
- `-S,--summary`
  - Summaryを出力します

### upgrade

```sh
UHA upgrade [-b,--branch [BRANCH]]
```

UHAをアップグレードします

- `-b,--branch [BRANCH]`
  - PullしてくるBranchを指定できます

### version

```sh
UHA version
```

UHAのバージョン情報を出力して終了します
