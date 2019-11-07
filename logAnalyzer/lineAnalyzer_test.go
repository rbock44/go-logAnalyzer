package logAnalyzer

import (
	"reflect"
	"regexp"
	"testing"
)

var resultOK bool
var resultHit string
var resultName *NamedRegEx

func BenchmarkIsLineOKWithEmail(b *testing.B) {
	compiled, _ := regexp.Compile(".*test.com")
	emailRegEx := NamedRegEx{Level: REGEX_LEVEL_LINE, Name: "email", RegEx: compiled}

	named := []NamedRegEx{emailRegEx}
	ignore := []IgnoreRegEx{}
	text := "blabla@test.com"
	var ok bool
	var hit string
	var name *NamedRegEx

	// run the IsLineOK function b.N times
	for n := 0; n < b.N; n++ {
		ok, hit, name = IsLineOK(named, ignore, text)
	}

	resultOK = ok
	resultHit = hit
	resultName = name
}

func TestIsLineOK(t *testing.T) {
	emailTestAddress := "mytest@test.com"
	emailPattern, _ := regexp.Compile(".*test.com")
	emptyHitString := ""
	emailRegEx := NamedRegEx{Level: REGEX_LEVEL_LINE, Name: "email", RegEx: emailPattern}
	ignoreEmailRegEx := IgnoreRegEx{Level: REGEX_LEVEL_SOURCEFILE, RegEx: emailPattern}

	type args struct {
		regularExpressions          []NamedRegEx
		regexToIdentifyIgnoredParts []IgnoreRegEx
		stringToAnalyze             string
	}
	tests := []struct {
		name          string
		args          args
		wantResult    bool
		wantHitString string
		wantHit       *NamedRegEx
	}{
		{name: "empty line",
			args: args{
				regularExpressions: []NamedRegEx{
					emailRegEx,
				},
				regexToIdentifyIgnoredParts: []IgnoreRegEx{},
				stringToAnalyze:             emptyHitString,
			},
			wantResult:    true,
			wantHitString: emptyHitString,
			wantHit:       nil,
		},
		{name: "mail address match",
			args: args{
				regularExpressions: []NamedRegEx{
					emailRegEx,
				},
				regexToIdentifyIgnoredParts: []IgnoreRegEx{},
				stringToAnalyze:             emailTestAddress,
			},
			wantResult:    false,
			wantHitString: emailTestAddress,
			wantHit:       &emailRegEx,
		},
		{name: "mail address ignore",
			args: args{
				regularExpressions: []NamedRegEx{
					emailRegEx,
				},
				regexToIdentifyIgnoredParts: []IgnoreRegEx{
					ignoreEmailRegEx,
				},
				stringToAnalyze: emailTestAddress,
			},
			wantResult:    true,
			wantHitString: emptyHitString,
			wantHit:       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotHitString, gotHit := IsLineOK(tt.args.regularExpressions, tt.args.regexToIdentifyIgnoredParts, tt.args.stringToAnalyze)
			if gotResult != tt.wantResult {
				t.Errorf("IsLineOK() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotHitString != tt.wantHitString {
				t.Errorf("IsLineOK() gotHitString = %v, want %v", gotHitString, tt.wantHitString)
			}
			if !reflect.DeepEqual(gotHit, tt.wantHit) {
				t.Errorf("IsLineOK() gotHit = %v, want %v", gotHit, tt.wantHit)
			}
		})
	}
}
