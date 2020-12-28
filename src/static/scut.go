package static

var (
	SCUT_JW_API_URL         = "http://jw.scut.edu.cn/zhinan/cms/article/v2/findInformNotice.do"
	SCUT_JW_USER_AGENT      = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36"
	SCUT_JW_REFERER         = "http://jw.scut.edu.cn/zhinan/cms/toPosts.do"
	SCUT_JW_TAG_NAMES       = []string{"选课", "考试", "实践", "交流", "教师", "信息"}
	SCUT_JW_ITEM_URL        = "http://jw.scut.edu.cn/zhinan/cms/article/view.do?type=posts&id=%s"
	SCUT_JW_ITEM_MOBILE_URL = "http://jw.scut.edu.cn/dist/#/detail/index?id=%s&type=notice"

	SCUT_SE_WEB_URL    = "http://www2.scut.edu.cn/sse/%s/list.htm"
	SCUT_SE_USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36"
	SCUT_SE_TAG_PARTS  = []string{"xyjd_17232", "17235", "17236", "gwtz", "kytz"}
	SCUT_SE_TAG_NAMES  = []string{"学院焦点", "本科生通知", "研究生通知", "公务通知", "科研通知"}
	SCUT_SE_ITEM_URL   = "http://www2.scut.edu.cn/%s"
)
