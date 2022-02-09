# CeVIGO

これはCeVIO AIのGo用の薄いAPIラッパーです。

CeVIO AIをサポート対象としていますが、CSでも多分動く気がします。

CeVIO AIのCOMコンポーネントを移植しただけなのでドキュメントは公式のを見てください。

https://cevio.jp/guide/cevio_ai/interface/com/

# How to use

参考実装としてcmd以下にgRPCでテキストを音声に変換するTTSサーバとそのサーバにリクエストを投げるクライアントの実装があります。

簡単なトークンを付けた認証しか行なっていないので外部に晒したり本番環境では利用しないことをおすすめします。リバースプロキシか何かで暗号化さえかけちゃえば大丈夫かとは思いますが。

# 実装状況

CeVIO AIのAPIを見て実装しましたが、型が違うだけでほぼAPIが同じなのでCeVIO CSでも動きます。

## CeVIO AI

- [x] interface ITalker2V40
  - トーク機能を提供します。
- [x] interface ITalkerComponentArray2
  - キャストの感情パラメータマップを表すオブジェクト。
- [x] interface ITalkerComponent2
  - 感情パラメータの単位オブジェクト。
- [x] interface ISpeakingState2
  - 再生状態を表すオブジェクト。
- [ ] interface IPhonemeDataArray2
  - 音素データの配列を表すオブジェクト。
  - 実装をサボっています。
- [ ] interface IPhonemeData2
  - 音素データの単位オブジェクト。
  - 実装をサボっています。
- [x] interface IStringArray2
  - 文字列の配列を表すオブジェクト。
  - AvailableCastsの返り値です、たぶん。
- [ ] interface IServiceControl2V40
  - 【CeVIO Creative Studio】制御機能を提供します。
  - やります。

## CeVIO CS
- [x] interface ITalkerV40
  - 動きました。
- [x] interface ITalkerComponentArray
  - 動きました。
- [x] interface ITalkerComponent
  - 動きました。
- [x] interface ISpeakingState
  - 未確認、たぶん動きます
- [ ] interface IPhonemeDataArray
- [ ] interface IPhonemeData
- [x] interface IStringArray
  - たぶん動きます。
- [ ] interface IServiceControlV40
  - やります。
