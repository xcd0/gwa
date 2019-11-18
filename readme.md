# go-wab-app

## したこと

これを参考にした。 https://dev.to/talentlessguy/create-frontend-go-apps-using-go-web-app-56cb

* tinygo を使う。
* dockerが推奨されていたが動かない。

## インストール
* とりあえず非推奨のwindows版を利用した。
* 下のをdlして/c/tinygoに展開。
	* https://github.com/tinygo-org/tinygo/releases/download/v0.9.0/tinygo0.9.0.windows-amd64.zip
* `export PATH=$PATH:/c/tinygo/bin`する。
* `go get -u -v https://github.com/talentlessguy/go-web-app` する。
* `cd $GOPATH/src/github.com/talentlessguy/go-web-app && go build -o gwa && cp gwa $GOBIN` する。

## アプリ作成
	* 好きなところで`gwa init 名前`する。
	* `cd 名前 && gwa build` する。
		* これで`build`ディレクトリにwasmのファイルが生成される。
	* `gwa dev --port 8080`してブラウザでlocalhost:8080する。
	* 見えた。

# 感想

* small is beautiful.  
	go-web-appの小ささはすごい。
* よいと思う
	* electronはchromiumの更新があると大変らしい。
	* gotronは悪くはない。
* まだ44コミットしかない出来立てに近いが、将来性が大きいと判断した。
* ドキュメントもないがまだソースコードが把握できる範囲だと考えよう。


