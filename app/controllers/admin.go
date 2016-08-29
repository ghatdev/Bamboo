package controllers

import (
	"strings"

	"crypto/sha256"
	"encoding/base64"

	"github.com/revel/revel"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Admin Controller
type Admin struct {
	*revel.Controller
}

// Account struct
// 관리자 계정 구조
type Account struct {
	Email    string
	ID       string
	Password string
	Role     string
}

// Content struct
// 제보내용 구조
type Content struct {
	Message string
	Time    string
	Posted  string
	Snum    string
	Ipaddr  string
}

// Index func
// 관리자 페이지
func (c Admin) Index() revel.Result {
	if c.Session["id"] == "" { // 로그인 안했으면 로그인 페이지로 Redirect
		return c.Redirect(Admin.Login)
	}

	session, err := mgo.Dial("localhost") // DB 연결. 여기서는 localhost의 MongoDB에 연결. 오류 발생하면 err에 오류값이 저장된다.
	if err != nil {                       // 오류 발생한 경우
		c.Flash.Error("오류로 인해 내용을 불러올 수 없습니다. 오류가 지속될 경우 관리자에게 문의 바랍니다.") // 에러 메세지
		return c.Redirect(Admin.Login)                                    // 내용 작성하던 페이지로 Redirect
	}
	defer session.Close() // panic 이 호출된 경우 세션을 Close 하고 즉시 리턴

	collection := session.DB("bamboo").C("content") // MongoDB의 bamboo 데이터베이스에서 content 컬렉션 선택

	var contents []Content                                          // 가져올 content 형식. 여러개의 content임으로 배열로 정의
	err = collection.Find(bson.M{"posted": "false"}).All(&contents) // {"posted":"false"}인, 즉 게시되지 않은 모든 게시물 정보 가져와서 contents에 배열로 저장
	if err != nil {                                                 // 못찾았다면
		c.Flash.Error("게시되지 않은 내용이 없습니다") //게시되지 않은 내용 없음. 열일하시네요
	}

	session.Close() // DB 쓰고 나면 꼭 session 닫기

	return c.Render(contents) // contens 값 View 로 전달한 후 View 렌더링
}

// Posted func
// 게시된 게시물들 List 페이지
func (c Admin) Posted() revel.Result {
	if c.Session["id"] == "" { // 로그인 되어있는지 체크. 안되있다면 로그인 페이지로 이동
		c.Flash.Error("로그인이 필요합니다")
		return c.Redirect(Admin.Login)
	}

	session, err := mgo.Dial("localhost") // DB 연결. 여기서는 localhost의 MongoDB에 연결. 오류 발생하면 err에 오류값이 저장된다.
	if err != nil {                       // 오류 발생한 경우
		c.Flash.Error("오류로 인해 내용을 불러올 수 없습니다. 오류가 지속될 경우 관리자에게 문의 바랍니다.") // 에러 메세지
		return c.Redirect(Admin.Login)                                    // 내용 작성하던 페이지로 Redirect
	}
	defer session.Close() // panic 이 호출된 경우 세션을 Close 하고 즉시 리턴

	collection := session.DB("bamboo").C("content") // MongoDB의 bamboo 데이터베이스에서 content 컬렉션 선택

	var contents []Content                                         // 가져올 content 형식. 여러개의 content임으로 배열로 정의
	err = collection.Find(bson.M{"posted": "true"}).All(&contents) // {"posted":"true"}인, 즉 게시된 모든 게시물 정보 가져와서 contents에 배열로 저장
	if err != nil {                                                // 못찾았다면
		c.Flash.Error("게시된 내용이 없습니다") //게시된 내용 없음. 열일하세요
	}

	session.Close() // DB 쓰고 나면 꼭 session 닫기

	return c.Render(contents) // contens 값 View 로 전달한 후 View 렌더링
}

// Post func
// 게시 처리 MethodName
func (c Admin) Post(content string, snum string) revel.Result {
	if c.Session["id"] == "" { // 로그인 되어있는지 체크. 안되있다면 로그인 페이지로 이동
		c.Flash.Error("로그인이 필요합니다")
		return c.Redirect(Admin.Login)
	}

	session, err := mgo.Dial("localhost") // DB 연결. 여기서는 localhost의 MongoDB에 연결. 오류 발생하면 err에 오류값이 저장된다.
	if err != nil {                       // 오류 발생한 경우
		c.Flash.Error("오류로 인해 내용을 게시할 수 없습니다. 오류가 지속될 경우 관리자에게 문의 바랍니다.") // 에러 메세지
		return c.Redirect(Admin.Index)                                    // 내용 작성하던 페이지로 Redirect
	}
	defer session.Close() // panic 이 호출된 경우 세션을 Close 하고 즉시 리턴

	collection := session.DB("bamboo").C("content") // MongoDB의 bamboo 데이터베이스에서 content 컬렉션 선택

	result := Content{}                                      // 가져올 content 형식. 한개임으로 Content 형식을 가진 하나의 변수
	err = collection.Find(bson.M{"snum": snum}).One(&result) // snum이 전달된 snum 이랑 같은 게시물 검색.
	if err != nil {                                          // snum인 게시물 없다면
		c.Flash.Error("올바르지 않은 접근입니다") // 말도 안되는거지 이게
		return c.Redirect(Admin.Index) // 오류 출력하고 Index로 복귀
	}

	err = collection.Update(bson.M{"snum": snum}, bson.M{"message": content, "time": result.Time, "posted": "true", "snum": snum, "ipaddr": result.Ipaddr}) // 찾은 내용에서 posted 값만 ture로 바꾸고 나머지는 그대로
	if err != nil {                                                                                                                                         //실패했다면
		c.Flash.Error("올바르지 않은 접근입니다") // 이것도 말도안되는거지.
		return c.Redirect(Admin.Index) // 오류 출력하고 Index 로 복귀
	}

	session.Close() // DB 쓰고나선 꼭 연결 닫자.

	return c.Redirect(Admin.Index) // 앞에서 아무문제 없었다면 아무 메세지 없이  Index로 복귀
}

// Login func
// 관리자 로그인 페이지
func (c Admin) Login() revel.Result {
	return c.Render() // 단순히 페이지 렌더링
}

// LoginInternal func
// 관리자 로컬계정 로그인
// POST 전용 func
func (c Admin) LoginInternal(inputEmail string, inputPassword string) revel.Result {
	session, err := mgo.Dial("localhost") // DB 연결. 여기서는 localhost의 MongoDB에 연결. 오류 발생하면 err에 오류값이 저장된다.
	if err != nil {                       // 오류 발생한 경우
		c.Flash.Error("오류로 인해 로그인에 실패하였습니다. 오류가 지속될 경우 관리자에게 문의 바랍니다.") // 에러 메세지
		return c.Redirect(Admin.Login)                                  // 내용 작성하던 페이지로 Redirect
	}
	defer session.Close() // panic 이 호출된 경우 세션을 Close 하고 즉시 리턴

	//session.SetMode(mgo.Monotonic, true) // 모드 설정.

	collection := session.DB("bamboo").C("accounts")

	result := Account{}
	err = collection.Find(bson.M{"email": inputEmail}).One(&result)
	if err != nil {
		c.Flash.Error("존재하지 않는 아이디 입니다")
		return c.Redirect(Admin.Login)
	}

	indexOfAt := strings.LastIndex(inputEmail, "@")
	id := inputEmail[:indexOfAt]

	h := sha256.New()
	h.Write([]byte(id + inputPassword))
	bs := h.Sum([]byte{})

	psw := base64.StdEncoding.EncodeToString(bs)

	if strings.Compare(result.Password, psw) != 0 {
		c.Flash.Error("비밀번호가 일치하지 않습니다")
		return c.Redirect(Admin.Login)
	}

	session.Close()

	c.Session["id"] = result.Email
	c.Session["role"] = result.Role
	return c.Redirect(Admin.Index)
}

// Logout func
// 로그아웃. id와 role 세션을 초기화
// delete 사용해도 될듯?
func (c Admin) Logout() revel.Result {
	if c.Session["id"] == "" {
		c.Flash.Error("로그인이 필요합니다")
		return c.Redirect(Admin.Login)
	}

	c.Session["id"] = ""
	c.Session["role"] = ""

	return c.Redirect(Admin.Login)
}

// ChangePassword func
// 비밀번호 변경
func (c Admin) ChangePassword() revel.Result {
	if c.Session["id"] == "" {
		c.Flash.Error("로그인이 필요합니다")
		return c.Redirect(Admin.Login)
	}

	return c.Render()
}

// UpdatePassword func
// ChangePassword POST func
func (c Admin) UpdatePassword(oldpsw string, newpsw string, newpswConfirm string) revel.Result {
	if c.Session["id"] == "" {
		c.Flash.Error("로그인이 필요합니다")
		return c.Redirect(Admin.Login)
	}

	if len(newpsw) < 8 {
		c.Flash.Error("패스워드는 8자리 이상이여야 합니다")
		return c.Redirect(Admin.ChangePassword)
	}

	if strings.Compare(newpsw, newpswConfirm) != 0 {
		c.Flash.Error("패스워드가 일치하지 않습니다")
		return c.Redirect(Admin.ChangePassword)
	}

	email := c.Session["id"] // email 변수에 로그인된 사용자의 email 주소 저장

	session, err := mgo.Dial("localhost") // DB 연결. 여기서는 localhost의 MongoDB에 연결. 오류 발생하면 err에 오류값이 저장된다.
	if err != nil {                       // 오류 발생한 경우
		c.Flash.Error("내용 저장에 실패하였습니다. 오류가 지속될 경우 관리자에게 문의 바랍니다.") // 에러 메세지
		c.FlashParams()                                            // 기존에 작성했던 내용 저장
		c.Redirect(App.Index)                                      // 내용 작성하던 페이지로 Redirect
	}
	defer session.Close() // panic 이 호출된 경우 세션을 Close 하고 즉시 리턴

	collection := session.DB("bamboo").C("accounts") // accounts collection 선택

	// ---------------------------- 기존 패스워드 재확인 과정 ---------------------------------
	result := Account{}
	err = collection.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		c.Flash.Error("존재하지 않는 아이디 입니다")
		return c.Redirect(Admin.Login)
	}

	indexOfAt := strings.LastIndex(email, "@")
	id := email[:indexOfAt]

	h := sha256.New()
	h.Write([]byte(id + oldpsw))
	bs := h.Sum([]byte{})

	psw := base64.StdEncoding.EncodeToString(bs)

	if strings.Compare(result.Password, psw) != 0 {
		c.Flash.Error("현재 패스워드가 일치하지 않습니다")
		return c.Redirect(Admin.ChangePassword)
	}

	// --------------------------------- 새로운 비밀번호로 업데이트 하는 과정 ---------------------------------

	h = sha256.New()
	h.Write([]byte(id + newpsw))
	bs = h.Sum([]byte{})

	HashednewPsw := base64.StdEncoding.EncodeToString(bs)

	err = collection.Update(bson.M{"email": email}, bson.M{"email": email, "password": HashednewPsw, "role": c.Session["role"]})
	if err != nil {
		c.Flash.Error("패스워드 변경에 실패하였습니다. " + err.Error())
		return c.Redirect(Admin.ChangePassword)
	}

	c.Flash.Success("패스워드가 정상적으로 변경되었습니다")

	session.Close()

	return c.Redirect(Admin.ChangePassword)
}

// Register func
// 관리자 등록 func
func (c Admin) Register() revel.Result {
	if c.Session["id"] == "" {
		c.Flash.Error("로그인이 필요합니다")
		return c.Redirect(Admin.Login)
	}

	if c.Session["role"] != "Master" {
		c.Flash.Error("권한이 부족합니다")
		return c.Redirect(Admin.Index)
	}

	return c.Render()
}

// AddUser func
// 관리자 등록 POST func
func (c Admin) AddUser(email string, psw string, role string) revel.Result {
	if c.Session["id"] == "" {
		c.Flash.Error("로그인이 필요합니다")
		return c.Redirect(Admin.Login)
	}

	if c.Session["role"] != "Master" {
		c.Flash.Error("권한이 부족합니다")
		return c.Redirect(Admin.Index)
	}

	session, err := mgo.Dial("localhost") // DB 연결. 여기서는 localhost의 MongoDB에 연결. 오류 발생하면 err에 오류값이 저장된다.
	if err != nil {                       // 오류 발생한 경우
		c.Flash.Error("계정 등록에 실패하였습니다. 오류가 지속될 경우 관리자에게 문의 바랍니다.") // 에러 메세지
		c.FlashParams()                                            // 기존에 작성했던 내용 저장
		return c.Redirect(Admin.Register)                          // 내용 작성하던 페이지로 Redirect
	}
	defer session.Close() // panic 이 호출된 경우 세션을 Close 하고 즉시 리턴

	//session.SetMode(mgo.Monotonic, true) // 모드 설정.

	indexOfAt := strings.LastIndex(email, "@")
	id := email[:indexOfAt]

	h := sha256.New()
	h.Write([]byte(id + psw))
	bs := h.Sum([]byte{})

	hpsw := base64.StdEncoding.EncodeToString(bs)

	if role == "" {
		role = "Admin"
	}

	collection := session.DB("bamboo").C("accounts")         // MongoDB에서 DB와 collection 설정
	err = collection.Insert(&Account{email, id, hpsw, role}) // 선택된 DB, collection 에 전달받은 message와 저장되는 시간 구조화하여 MongoDB에 저장.
	// 오류가 발생할경우 err에 에러정보가 입력되며 에러가 없을경우 err은 nil(null)이 된다.

	if err != nil { // 오류 발생한 경우
		c.Flash.Error("계정 등록에 실패하였습니다. 오류가 지속될 경우 관리자에게 문의 바랍니다.") // 오류 메세지
		c.FlashParams()                                            // 기존에 작성했던 내용 저장
		return c.Redirect(App.Index)                               // 내용 작성하던 페이지로 Redirect
	}

	// 위의 과정에서 아무 오류가 없다면 데이터 저장에 성공한 것임으로
	session.Close() // DB와 연결된 세션을 닫고

	c.Flash.Success("계정이 정상적으로 등록되었습니다.") // 성공했다는 메세지
	return c.Redirect(Admin.Register)
}

// ManageAccounts func
// 계정 관리 func
func (c Admin) ManageAccounts() revel.Result {
	if c.Session["id"] == "" {
		c.Flash.Error("로그인이 필요합니다")
		return c.Redirect(Admin.Login)
	}

	if c.Session["role"] != "Master" {
		c.Flash.Error("권한이 부족합니다")
		return c.Redirect(Admin.Index)
	}

	session, err := mgo.Dial("localhost") // DB 연결. 여기서는 localhost의 MongoDB에 연결. 오류 발생하면 err에 오류값이 저장된다.
	if err != nil {                       // 오류 발생한 경우
		c.Flash.Error("DB에 접속할 수 없습니다. 해당 오류가 지속될 경우 관리자에게 문의 바랍니다.") // 에러 메세지
		c.FlashParams()                                               // 기존에 작성했던 내용 저장
		return c.Redirect(Admin.Index)                                // 내용 작성하던 페이지로 Redirect
	}
	defer session.Close() // panic 이 호출된 경우 세션을 Close 하고 즉시 리턴

	var results []Account
	collection := session.DB("bamboo").C("accounts")             // MongoDB에서 DB와 collection 설정
	err = collection.Find(bson.M{"role": "Admin"}).All(&results) // 선택된 DB, collection 에 전달받은 message와 저장되는 시간 구조화하여 MongoDB에 저장.
	// 오류가 발생할경우 err에 에러정보가 입력되며 에러가 없을경우 err은 nil(null)이 된다.

	if err != nil { // 오류 발생한 경우
		c.Flash.Error("등록된 계정을 찾을 수 없습니다.") // 오류 메세지
		session.Close()
		return c.Redirect(Admin.ManageAccounts) // 내용 작성하던 페이지로 Redirect
	}

	session.Close() // DB와 연결된 세션을 닫고

	return c.Render(results)
}

// DeleteAccount func
// 계정 삭제 func
// Post only
func (c Admin) DeleteAccount(email string) revel.Result {
	if c.Session["id"] == "" {
		c.Flash.Error("로그인이 필요합니다")
		return c.Redirect(Admin.Login)
	}

	if c.Session["role"] != "Master" {
		c.Flash.Error("권한이 부족합니다")
		return c.Redirect(Admin.Index)
	}

	session, err := mgo.Dial("localhost") // DB 연결. 여기서는 localhost의 MongoDB에 연결. 오류 발생하면 err에 오류값이 저장된다.
	if err != nil {                       // 오류 발생한 경우
		c.Flash.Error("DB에 접속할 수 없습니다. 해당 오류가 지속될 경우 관리자에게 문의 바랍니다.") // 에러 메세지
		c.FlashParams()                                               // 기존에 작성했던 내용 저장
		return c.Redirect(Admin.Index)                                // 내용 작성하던 페이지로 Redirect
	}
	defer session.Close() // panic 이 호출된 경우 세션을 Close 하고 즉시 리턴

	collection := session.DB("bamboo").C("accounts") // MongoDB에서 DB와 collection 설정
	err = collection.Remove(bson.M{"email": email})

	// 오류가 발생할경우 err에 에러정보가 입력되며 에러가 없을경우 err은 nil(null)이 된다.

	if err != nil { // 오류 발생한 경우
		c.Flash.Error("계정을 삭제하는 중 오류가 발생하였습니다.") // 오류 메세지
		session.Close()
		return c.Redirect(Admin.ManageAccounts) // 내용 작성하던 페이지로 Redirect
	}

	session.Close() // DB와 연결된 세션을 닫고
	c.Flash.Success("계정을 성공적으로 삭제하였습니다")
	return c.Redirect(Admin.ManageAccounts)
}
