# go-wab-app

# したこと

* tinygo を使う。
* dockerが推奨されていたが動かない。

* とりあえず非推奨のwindows版を利用した。
	* 下のをdlして/c/tinygoに展開。
		* https://github.com/tinygo-org/tinygo/releases/download/v0.9.0/tinygo0.9.0.windows-amd64.zip
	* `export PATH=$PATH:/c/tinygo/bin`する。
	* `git clone https://github.com/talentlessguy/go-web-app` する。
	* `go build && mv go-web-app gwa && cp gwa $GOBIN` する。
	* 好きなところで`gwa init 名前`する。
	* `cd 名前 && gwa build` する。
		* これで`build`ディレクトリにwasmのファイルが生成される。
	* `gwa dev --port 8080`してブラウザでlocalhost:8080する。
	* 見えた。
	* これを参考にした。 https://dev.to/talentlessguy/create-frontend-go-apps-using-go-web-app-56cb

# 感想

* small is beautiful.
* よいと思う
* まだ44コミットしかない出来立てに近いが、将来性が大きいと判断した。
* ドキュメントもないがまだソースコードが把握できる範囲だと考えよう。


