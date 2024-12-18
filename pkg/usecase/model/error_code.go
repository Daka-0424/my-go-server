package model

const (
	// E /   00    /  00
	// E / カテゴリ / 連番

	// DB系
	E0001 = "E0001" // エラーが発生しました。しばらくしてから、もう一度お試しください。
	E0002 = "E0002" // データの取得に失敗しました。しばらくしてから、もう一度お試しください。

	// ユーザ系
	E0101 = "E0101" // 認証に失敗しました。しばらくしてから、もう一度お試しください。
	E0102 = "E0102" // 作成済みのアカウントです。
	E0103 = "E0103" // ユーザ作成に失敗しました。しばらくしてから、もう一度お試しください。
	E0104 = "E0104" // ユーザー取得に失敗しました。しばらくしてから、もう一度お試しください。
	E0105 = "E0105" // アカウントが停止されています。
	E0106 = "E0106" // ユーザが存在しています。

	// Admin系
	E2001 = "E2001" // 管理者登録の有効期限が切れました。もう一度登録をやり直してください。
	E2002 = "E2002" // 管理者登録に失敗しました。しばらくしてから、もう一度お試しください。
	E2003 = "E2003" // 管理者ログインに失敗しました。しばらくしてから、もう一度お試しください。
	E2004 = "E2004" // セッションの有効期限が切れました。もう一度ログインしてください。
	E2005 = "E2005" // 管理者が存在していません。
	E2006 = "E2006" // 管理者取得に失敗しました。しばらくしてから、もう一度お試しください。
	E2007 = "E2007" // 管理者更新に失敗しました。しばらくしてから、もう一度お試しください。
	E2008 = "E2008" // 管理者削除に失敗しました。しばらくしてから、もう一度お試しください。

	// PlatformProduct系
	E3001 = "E3001" // 商品の取得に失敗しました。しばらくしてから、もう一度お試しください。

	// 課金系
	E9001 = "E9001" // 署名の検証に失敗しました。しばらくしてから、もう一度お試しください。
	E9002 = "E9002" // 購入情報の取得に失敗しました。しばらくしてから、もう一度お試しください。
	E9003 = "E9003" // 購入情報の保存に失敗しました。しばらくしてから、もう一度お試しください。
	E9004 = "E9004" // 購入情報がキャンセルか保留になっています。
	E9005 = "E9005" // 購入情報が一致しません。
	E9006 = "E9006" // 購入情報が正しくありません。
	E9007 = "E9007" // 既に購入済みです。

	// EarnedPoint系
	E9101 = "E9101" // 通貨の取得に失敗しました。しばらくしてから、もう一度お試しください。
	E9102 = "E9102" // 通貨の作成に失敗しました。しばらくしてから、もう一度お試しください。
	E9103 = "E9103" // 通貨の更新に失敗しました。しばらくしてから、もう一度お試しください。

	// メンテナンス
	E9801 = "E9801" // メンテナンス中です。しばらくしてから、もう一度お試しください。

	// その他
	E9901 = "E9901" // 不正なリクエストです。
	E9999 = "E9999" // エラーが発生しました。しばらくしてから、もう一度お試しください。
)
