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
UHA count [-C,--Cumulative|-F,--failure-only|--firstF float|--secondF float|--thirdF float]
```



現在のディレクトリにあるCSVファイルを見つけて、不良数を数え上げます。
- `-C,--Cumulative`
  - 不良数を累積和します
- `-F,--failure-only`
  - 不良数だけを出力します
- `--firstF,--secondF,--thirdF float`
  - 1カラム目、2カラム目、3カラム目がfloat以上だった時、不良とします。すべてAndです


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
UHA make [--SEED int|-N,--VtnVoltage float|-P,--VtpVoltage float|-D,--default|-h,--help|-o,--out string|--sigma float|-y,--yes]
```



シミュレーションタスクを作るコマンドです。オプション無しで実行すると、インタラクティブにタスクを生成できます。タスクは `~/.config/UHA/task`以下に保存されます


- `--SEED int`
  - SEEDを指定します
- `-N,--VtnVoltage|-P,--VtpVoltage float`
  - Vtn,Vtpのしきい値電圧です。設定ファイルより優先度が高いです
- `-o,--out`
  - シミュレーションの書き出し先です
- `-D,--default`
  - `.UHA.json`の設定値をそのままでタスクを発行します
- `--sigma`
  - シグマの値を引数で指定できます。`-D,--default`の値を上書きします
- `-y,--yes`
  - 最後の確認を飛ばすことが出来ます



### push

```sh
UHA push [--Id string|-R,--RangeSEED|-h,--help|-G,--ignore-sigma|-n,--name string]
```

Google SpreadSheetにデータを書き込みます

- `--Id string`
  - スプレッドシートのIDです
- `-R,--RangeSEED`
  - RangeSEEDシミュレーションの結果をPushします
- `-F,--ignore-sigma`
  - Sigmaを省いてPushします
- `-n,--name`
  - シートの名前です



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
UHA run [--GC|--all|-C,--continue|-h,--help|--no-notify|-n,number [num]]
```

発行済のタスクを実行します。


- `--GC`
  - 結果のCSV以外を削除します
- `--no-notify`
  - Slackに通知しません
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
UHA srun [--GC|--all|-C,--continue|-h,--help|--no-notify|-n,--number [num]|-P,--parallel [num]|--start int|-S,--summary]
```



`smake`で作ったタスクを実行します



- `--GC`
  - 結果のCSV以外を削除します
- `--no-notify`
  - Slackに通知しません
- `--start int`
  - SEEDの初期値を指定します
- `-S,--summary`
  - サマリーを出力します
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
UHA upgrade
```

UHAをアップグレードします

### version

```sh
UHA version
```

UHAのバージョン情報を出力して終了します
