package static

const (
	ScutJwApi             = "http://jw.scut.edu.cn/zhinan/cms/article/v2/findInformNotice.do"
	ScutJwContentType     = "application/x-www-form-urlencoded;charset=UTF-8"
	ScutJwCookie          = "JSESSIONID=4658DB1262F616F7404775101777E62A"
	ScutJwOrigin          = "http://jw.scut.edu.cn"
	ScutJwReferer         = "http://jw.scut.edu.cn/zhinan/cms/toPosts.do?category=0"
	ScutJwNoticeUrl       = "http://jw.scut.edu.cn/zhinan/cms/article/view.do?type=posts&id=%s"
	ScutJwNoticeMobileUrl = "http://jw.scut.edu.cn/dist/#/detail/index?type=posts&id=%s"

	ScutSeWebUrl       = "http://www2.scut.edu.cn/sse/%s/list.htm"
	ScutSeNoticeWebUrl = "http://www2.scut.edu.cn/%s"

	ScutGrWebUrl       = "http://www2.scut.edu.cn/graduate/14562/list%d.htm"
	ScutGrNoticeWebUrl = "http://www2.scut.edu.cn/%s"

	ScutGzicWebUrl       = "http://www2.scut.edu.cn/gzic/%s/list.htm"
	ScutGzicNoticeWebUrl = "http://www2.scut.edu.cn/%s"
)

var (
	ScutJwTagNames = map[int]string{1: "选课", 2: "考试", 3: "实践", 4: "交流", 5: "教师", 6: "信息"}

	ScutSeTagParts = []string{"xyjd_17232", "17235", "17236", "gwtz", "kytz"}
	ScutSeTagNames = []string{"学院焦点", "本科生通知", "研究生通知", "公务通知", "科研通知"}

	ScutGrTagName = "通知公告"

	ScutGzicTagParts = []string{"30284", "30307", "30283"}
	ScutGzicTagNames = []string{"学术预告", "教研通知", "事务通知"}
)
