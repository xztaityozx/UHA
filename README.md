# UHA

[![img](https://travis-ci.org/xztaityozx/UHA.svg?branch=master)](https://travis-ci.org/xztaityozx/UHA) [MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)



Ultra H_SPICE Attacker



H_PICEのモンテカルロをうまくするツール

---

- [インストール](doc/Install.md)
- [使い方](doc/usage.md)
- [設定](doc/setting.md)

  

## Require

- go (>=1.10.2)



## Example

シグマを0.00、それ以外をデフォルト値にしてタスクを作ったあと、実行

```sh
# タスクを作る
$ UHA make -D --sigma 0.000 -y
# タスクの実行
$ UHA run 
```





始点、刻み幅、回数を指定してシグマを動かしながらタスクを連続生成した後、すべて実行

```sh
# タスクを作る 0から2刻みで10まで
$ UHA rmake -n 0 2 10
# 実行 失敗しても継続
$ UHA run --all -C
```



50000回のモンテカルロを2回するタスクを発行して実行

```sh
# タスク生成
$ UHA smake -T 50000 -n 2 --sigma 0.00
# 実行
$ UHA srun
```

