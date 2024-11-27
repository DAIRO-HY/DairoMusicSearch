package form

/**
 * 搜索结果表单
 */
type SearchForm struct {

	/**
	 * 歌名
	 */
	Name string `json:"name"`

	/**
	 * LOGO
	 */
	Logo string `json:"logo"`

	/**
	 * 视频ID
	 */
	VideoId string `json:"videoId"`
}
