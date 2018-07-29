## UHA setting

`.UHA.json`について記述します



---



## .UHA.json

`UHA` の設定ファイルです。デフォルトでは `~/.config/UHA/.UHA.json` が読み込まれます。



```sh
$ UHA --config /path/to/.UHA.json command
```



と記述すれば別の設定を読み込むことが出来ます。



## 設定項目

割とあります



```json
{
  "Simulation":{
    "Monte":["100", "200", "500", "1000", "2000", "5000", "10000", "20000", "50000"],
    "Range":{
      "Start":"2.5ns",
      "Stop":"17.5ns",
      "Step":"7.5ns"
    },
    "Signal":"signalname",
    "SimDir":"/path/to/netlist",
    "DstDir":"/path/to/result",
    "Vtp":{
      "Voltage":0.8,
      "Sigma":0.0,
      "Deviation":1.0
    },
    "Vtn":{
      "Voltage":0.8,
      "Sigma":0.0,
      "Deviation":1.0
    }
  },
  "TaskDir":"/path/to/.config/UHA/task",
  "Repository":[
    {
      "Type":"Git",
      "Path":"/path/to/repository"
    }
  ],
  "SpreadSheet":{
    "Id":"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "CSPath":"/path/to/client_secret.json",
    "TokenPath":"/path/to/token.json"
  },
  "SlackConfig":{
    "Token":"xxxxxxxxxxxxxxxxxxxxxxxxx",
    "User":"user",
    "Channel":"channel"
  }
}
```

---

### `Simulation` クループ

#### Monte

モンテカルロの回数のリストです。文字列のリストで記述します。



#### Range

WaveViewでデータを取るときに使います。

- `Start` : データの始点です
- `Stop` : データの終点です
- `Step` : データを取る刻み幅です



#### Signal

データを取りたい信号線の名前です



#### SimDir

`netlist/`ディレクトリまでの絶対パスです



#### DstDir

結果を書き出すディレクトリまでの絶対パスです



#### Vtp,Vtn

VtnとVtpのデフォルト値を決めます

- `Voltage` : しきい値電圧です
- `Sigma` : シグマです
- `Deviation` : なにこれ



### TaskDir

`~/.config/UHA/task`までの絶対パスです

---

### `Repository`グループ

データを集めるときに使う設定です。主に`UHA get`コマンドで使用します。`UHA get`を使わない場合、この設定はしなくてもいいです。

#### Type

リポジトリのタイプを指定します。Gitリポジトリの場合は`Git`を指定します。それ以外の場合は`Dir`を指定します



#### Path

データを格納するディレクトリへの絶対パスです。`UHA get`を実行するとこのディレクトリ以下にデータが集まります。



---

### `SpreadSheed`グループ

データをGoogle SpreadSheetに書き込むときに使う設定です。Google SpreadSheetを使わない場合は設定しなくてもいいです。

#### ID

スプレッドシートのIDです



#### CSPath

`client_secret.json`への絶対パスです



#### TokenPath

`token.json` への絶対パスです


---

### `SlackConfig`グループ
シミュレーション終了後Slackへ通知するための設定です

#### Token

BotのTokenです

#### Channel

投稿するチャンネルです

#### User

リプライを飛ばすUser名です


