// Code generated by pegomock. DO NOT EDIT.
package matchers

import (
	"reflect"
	"github.com/petergtz/pegomock"
	ch360 "github.com/CloudHub360/ch360.go/ch360"
)

func AnyCh360Document() ch360.Document {
	pegomock.RegisterMatcher(pegomock.NewAnyMatcher(reflect.TypeOf((*(ch360.Document))(nil)).Elem()))
	var nullValue ch360.Document
	return nullValue
}

func EqCh360Document(value ch360.Document) ch360.Document {
	pegomock.RegisterMatcher(&pegomock.EqMatcher{Value: value})
	var nullValue ch360.Document
	return nullValue
}