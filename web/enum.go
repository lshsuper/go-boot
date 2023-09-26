package web

import (
	"net/http"
	"strings"
)

type MethodMode int

const (
	_ MethodMode = iota
	RichMethod
	NormalMethod
)

func (e MethodMode) GetMethod(actName string) string {
	switch e {
	case RichMethod:
		return richMethod(actName)
	case NormalMethod:
		return normalMethod(actName)
	default:
		panic("method mode not define")
	}
}

func richMethod(actName string) string {

	actName = strings.ToLower(actName)
	if strings.Index(actName, "get") >= 0 || strings.Index(actName, "query") >= 0 {
		return http.MethodGet
	}
	if strings.Index(actName, "put") >= 0 || strings.Index(actName, "edit") >= 0 {
		return http.MethodPut
	}
	if strings.Index(actName, "add") >= 0 || strings.Index(actName, "create") >= 0 {
		return http.MethodPost
	}
	if strings.Index(actName, "del") >= 0 || strings.Index(actName, "delete") >= 0 || strings.Index(actName, "remove") >= 0 {
		return http.MethodDelete
	}
	return http.MethodGet

}

func normalMethod(actName string) string {

	if strings.Index(actName, "del") >= 0 || strings.Index(actName, "delete") >= 0 || strings.Index(actName, "remove") >= 0 ||
		strings.Index(actName, "put") >= 0 || strings.Index(actName, "edit") >= 0 ||
		strings.Index(actName, "add") >= 0 || strings.Index(actName, "create") >= 0 {
		return http.MethodPost
	}
	return http.MethodGet

}
