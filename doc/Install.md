# UHA install

`UHA` をインストールする手順です。



----



OSによってGolangのインストール方法が違うのでここでは [linuxbrew](http://linuxbrew.sh/) を使う方法を書きます



### Goのインストール

```sh
$ brew install go
$ export GOPATH=~/.go
$ echo 'export GOPATH=~/.go' >> ~/.bashrc
$ export PATH=$PATH:$GOPATH/bin
$ echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
```



### UHAのインストール

```sh
$ go get github.com/xztaityozx/UHA
```



### UHAの準備

#### いるもの



|      名前       |                           説明                           |       場所       | Require |
| :-------------: | :------------------------------------------------------: | :--------------: | :-----: |
| spitemplate.spi |    `hspice` に入力するSPIスクリプトのテンプレートです    | `~/.config/UHA/` |  必要   |
|   addfile.txt   |     SPIにインクルードするファイルのテンプレートです      | `~/.config/UHA/` |  必要   |
|   .UHA.config   |                 `UHA` の設定ファイルです                 | `~/.config/UHA`  |  必要   |
|    next.json    | データの数え上げとスプレッドシートに書き込むのに使います | `~/.config/UHA/` |  必要   |



#### spitemplate.spi

1. シミュレーション対象のディレクトリから、`input.spi`を探してきます。
   1. だいたいは `~/simulation/User/...../monte_caro/netlist`にある
2. `input.spi` を `~/.config/UHA` 以下へコピーします
3. `input.spi` を `spitemplate.spi` にリネームします
4. `input.spi` を以下のように編集します



|      場所      |                         変更前                         |                            変更後                            |
| :------------: | :----------------------------------------------------: | :----------------------------------------------------------: |
| だいたい8行目  |   `.param vtn=AGAUSS(xx,xx,xx) vtp=AGAUSS(yy,yy,yy)`   | `.param vtn=AGAUSS(%.4f,%.4f,%.4f) vtp=AGAUSS(%.4f,%.4f,%.4f)` |
| だいたい10行目 |           `.include '/path/to/addfile.txt'`            |                       `.include '%s'`                        |
| だいたい53行目 | `.tran 10p 20n start=0 uic sweep monte=xxx firstrun=1` |    `.tran 10p 20n start=0 uic sweep monte=%s firstrun=1`     |

コレ以外はそのままでOKです



#### addfile.txt

1. シミュレーションに使っている `addfile.txt` を `~/.config/UHA/` へコピーします
2. `.option SEED=xxxx` と書いている行を `.option SEED=%d`に変更します



#### .UHA.json

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
  }
}
```





上記の`json`をコピーして `~/.config/UHA` に設置してください。その後適切な値を書き込んでください。詳細は[設定](setting.md)を見てください



#### next.json

以下のコマンドを実行します



```sh
$ touch ~/.config/UHA/next.json
```



----

