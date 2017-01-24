+++
title = "ドライバ | ドキュメント"
type = "docs"
category = "docs"
lang = "ja"
basename = "drivers.html"
+++

# ドライバ

Esshのドライバとは、タスク実行時にシェルスクリプトを構築するためのテンプレートシステムです。ドライバを使用してタスクの動作を変更することができます。

## 例

~~~lua
driver "custom_driver" { 
    engine = [=[
    
    
    ]=],
    
    foo = "foo",
    
    bar = "bar",
}
~~~

WIP...