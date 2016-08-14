package controllers

import (
	"math/rand"
	"strings"
	"time"

	"github.com/revel/revel"
	"gopkg.in/mgo.v2"
)

// App Controller
type App struct {
	*revel.Controller
}

// Index func
// 제보 패이지. 여기가 메인.
func (c App) Index() revel.Result {
	question := [3]string{
		"2016년 학생회장이 선출된 이후 이 페이지를 개발한 누구는 공식적으로 이 별명을 얻게되었다. 이 별명은?(3글자)",
		"충남삼성고등학교 학생생활부장 선생님의 성함은?",
		"영어선생님 중 노래도 잘하시고 드럼도 잘치시는 선생님이 있다. 이 선생님의 성함은?"} //재학생 확인 질문 목록

	s := rand.NewSource(time.Now().UnixNano()) // 랜덤 시드를 시간 기반으로 설정
	r := rand.New(s)                           // 바뀐 랜덤 시드로 r 생성

	num := r.Intn(len(question)) // 0부터 재학생 확인 질문 갯수까지 무작위 숫자 선정

	q := question[num] // q는 재학생 확인 질문 목록에서 뽑힌 string
	rand.Seed(time.Now().UnixNano())

	snum := RandStringRunes(7)

	return c.Render(num, q, snum) // /App/Index 호출될때 Index.html 렌더링 하면서 num(뽑힌 무작위 숫자), q(뽑힌 재학생 질문) 전달 {{.num}}과 {{.q}}로 받을 수 있다.
}

// Post func
// - answer: 재학생 질문 답
// - message: 내용
// - qnum: 재학생 확인질문 번호
func (c App) Post(answer string, message string, qnum int, snum string) revel.Result {
	c.Validation.MinSize(message, 10).Message("내용이 너무 짧습니다.") // 내용 10자 미만 체크

	if c.Validation.HasErrors() { //내용 길이가 10자 미만이면 오류발생
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(App.Index)
	}

	qanswer := [3]string{"불효자", "손병희", "손재식"} //재학생 질문 확인 답. 질문번호와 답 번호가 같다

	if !strings.Contains(answer, qanswer[qnum]) { // 전달받은 질문 번호의 답을 포함하지 않을경우
		c.Flash.Error("재학생 확인 질문에 대한 답이 잘못되었습니다.") // 오류 메세지
		c.FlashParams()                            // 기존 작성하던 데이터 저장
		return c.Redirect(App.Index)               // 작성하던 페이지로 Redirect
	}

	// ------------------------- 여기서부터 데이터 DB에 저장. -----------------------------
	type Content struct {
		Message string
		Time    string
		Posted  string
		Snum    string
		Ipaddr  string
	} // 저장할 데이터 구조

	session, err := mgo.Dial("localhost") // DB 연결. 여기서는 localhost의 MongoDB에 연결. 오류 발생하면 err에 오류값이 저장된다.
	if err != nil {                       // 오류 발생한 경우
		c.Flash.Error("내용 저장에 실패하였습니다. 오류가 지속될 경우 관리자에게 문의 바랍니다.") // 에러 메세지
		c.FlashParams()                                            // 기존에 작성했던 내용 저장
		return c.Redirect(App.Index)                               // 내용 작성하던 페이지로 Redirect
	}
	defer session.Close() // panic 이 호출된 경우 세션을 Close 하고 즉시 리턴

	//session.SetMode(mgo.Monotonic, true) // 모드 설정.

	s := strings.Split(c.Request.RemoteAddr, ":")
	ip := s[0]

	collection := session.DB("bamboo").C("content")                                    // MongoDB에서 DB와 collection 설정
	err = collection.Insert(&Content{message, time.Now().String(), "false", snum, ip}) // 선택된 DB, collection 에 전달받은 message와 저장되는 시간 구조화하여 MongoDB에 저장.
	// 오류가 발생할경우 err에 에러정보가 입력되며 에러가 없을경우 err은 nil(null)이 된다.

	if err != nil { // 오류 발생한 경우
		c.Flash.Error("내용 저장에 실패하였습니다. 오류가 지속될 경우 관리자에게 문의 바랍니다.") // 오류 메세지
		c.FlashParams()                                            // 기존에 작성했던 내용 저장
		return c.Redirect(App.Index)                               // 내용 작성하던 페이지로 Redirect
	}

	// 위의 과정에서 아무 오류가 없다면 데이터 저장에 성공한 것임으로
	session.Close() // DB와 연결된 세션을 닫고

	c.Flash.Success("내용이 정상적으로 저장되었습니다.") // 성공했다는 메세지
	return c.Redirect(App.Index)          // 이번에는 작성하던 내용 저장하지 않고 Index로 Redirect
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandStringRunes func
// 랜덤 문자열 생성 함수
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
