// GENERATED CODE - DO NOT EDIT
package routes

import "github.com/revel/revel"


type tAdmin struct {}
var Admin tAdmin


func (_ tAdmin) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Admin.Index", args).Url
}

func (_ tAdmin) Posted(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Admin.Posted", args).Url
}

func (_ tAdmin) Post(
		content string,
		snum string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "content", content)
	revel.Unbind(args, "snum", snum)
	return revel.MainRouter.Reverse("Admin.Post", args).Url
}

func (_ tAdmin) Login(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Admin.Login", args).Url
}

func (_ tAdmin) LoginInternal(
		inputEmail string,
		inputPassword string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "inputEmail", inputEmail)
	revel.Unbind(args, "inputPassword", inputPassword)
	return revel.MainRouter.Reverse("Admin.LoginInternal", args).Url
}

func (_ tAdmin) Logout(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Admin.Logout", args).Url
}

func (_ tAdmin) ChangePassword(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Admin.ChangePassword", args).Url
}

func (_ tAdmin) UpdatePassword(
		oldpsw string,
		newpsw string,
		newpswConfirm string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "oldpsw", oldpsw)
	revel.Unbind(args, "newpsw", newpsw)
	revel.Unbind(args, "newpswConfirm", newpswConfirm)
	return revel.MainRouter.Reverse("Admin.UpdatePassword", args).Url
}

func (_ tAdmin) Register(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Admin.Register", args).Url
}

func (_ tAdmin) AddUser(
		email string,
		psw string,
		role string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "email", email)
	revel.Unbind(args, "psw", psw)
	revel.Unbind(args, "role", role)
	return revel.MainRouter.Reverse("Admin.AddUser", args).Url
}

func (_ tAdmin) ManageAccounts(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Admin.ManageAccounts", args).Url
}

func (_ tAdmin) DeleteAccount(
		email string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "email", email)
	return revel.MainRouter.Reverse("Admin.DeleteAccount", args).Url
}


type tApp struct {}
var App tApp


func (_ tApp) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("App.Index", args).Url
}

func (_ tApp) Post(
		answer string,
		message string,
		qnum int,
		snum string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "answer", answer)
	revel.Unbind(args, "message", message)
	revel.Unbind(args, "qnum", qnum)
	revel.Unbind(args, "snum", snum)
	return revel.MainRouter.Reverse("App.Post", args).Url
}


type tStatic struct {}
var Static tStatic


func (_ tStatic) Serve(
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.Serve", args).Url
}

func (_ tStatic) ServeModule(
		moduleName string,
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "moduleName", moduleName)
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.ServeModule", args).Url
}


type tTestRunner struct {}
var TestRunner tTestRunner


func (_ tTestRunner) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.Index", args).Url
}

func (_ tTestRunner) Suite(
		suite string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	return revel.MainRouter.Reverse("TestRunner.Suite", args).Url
}

func (_ tTestRunner) Run(
		suite string,
		test string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	revel.Unbind(args, "test", test)
	return revel.MainRouter.Reverse("TestRunner.Run", args).Url
}

func (_ tTestRunner) List(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.List", args).Url
}


