// Package morfeusz provides a morphological analyser for Polish.
// See the documentation of the C++ API at
// http://download.sgjp.pl/morfeusz/Morfeusz2.pdf
//
// The names of most methods correspond to the names in the C++ API
// in an obvious way, with two major exceptions:
//  * hasNext() and next() are renamed to Next() and TokenInfo(),
//  * all getFoo() methods are renamed to Foo().
package morfeusz

/*
#cgo LDFLAGS: -lmorfeusz2
#include "morfeusz-cgo.h"

static struct String makeStructString(_GoString_ s) {
  struct String ret = { _GoStringPtr(s), _GoStringLen(s) };
  return ret;
}
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

// Morfeusz is the type of a struct capable of morphological
// analysis and/or generation.
type Morfeusz struct {
	morf C.Morf
}

// Result is the type of a struct representing the result
// of morphological analysis of a text.
type Result struct {
	res C.Res
}

// TokenInfo is the type of a struct representing the morphological
// interpretation of a token in the result of morphological analysis.
type TokenInfo struct {
	info C.struct_TokenInfo
}

type (
	// Charset determines the encoding that Morfeusz uses
	// in its input and output.
	Charset C.enum_Charset
	// TokenNumbering determines when the token index is reset.
	TokenNumbering C.enum_TokenNumbering
	// CaseHandling determines how to deal with tokens
	// whose letter case does not match dictionary entries.
	CaseHandling C.enum_CaseHandling
	// WhitespaceHandling determines whether whitespace
	// appears in the result of morphological analysis.
	WhitespaceHandling C.enum_WhitespaceHandling
	// Usage determines whether Morfeusz is capable of
	// morphological analysis and/or generation.
	Usage C.enum_Usage
)

const (
	// UTF8 (the default) makes Morfeusz use the UTF-8
	// encoding in its input and output.
	UTF8 Charset = C.UTF8
	// ISO8859_2 makes Morfeusz use the ISO/IEC 8859-2
	// encoding in its input and output.
	ISO8859_2 = C.ISO8859_2
	// CP1250 makes Morfeusz use the Windows-1250
	// code page in its input and output.
	CP1250 = C.CP1250
	// CP852 makes Morfeusz use the DOS code page 852
	// in its input and output.
	CP852 = C.CP852
)

const (
	// SeparateNumbering (the default) makes Morfeusz
	// reset the token index on every call to Analyse.
	SeparateNumbering TokenNumbering = C.SEPARATE_NUMBERING
	// ContinuousNumbering makes Morfeusz reset the token index
	// only on calls to SetTokenNumbering.
	ContinuousNumbering = C.CONTINUOUS_NUMBERING
)

const (
	// ConditionallyCaseSensitive (the default) makes Morfeusz allow
	// the interpretations whose letter case does not match any
	// dictionary entry when there is no choice with matching case.
	ConditionallyCaseSensitive CaseHandling = C.CONDITIONALLY_CASE_SENSITIVE
	// StrictlyCaseSensitive makes Morfeusz reject the interpretations
	// whose letter case does not match any dictionary entry.
	StrictlyCaseSensitive = C.STRICTLY_CASE_SENSITIVE
	// IgnoreCase makes Morfeusz ignore the letter case when
	// looking up a form in the dictionary.
	IgnoreCase = C.IGNORE_CASE
)

const (
	// SkipWhitespaces (the default) makes Morfeusz ignore whitespace.
	SkipWhitespaces WhitespaceHandling = C.SKIP_WHITESPACES
	// AppendWhitespaces makes Morfeusz append whitespace
	// to the preceding token.
	AppendWhitespaces = C.APPEND_WHITESPACES
	// KeepWhitespaces makes Morfeusz emit whitespace tokens.
	KeepWhitespaces = C.KEEP_WHITESPACES
)

const (
	// BothAnalyseAndGenerate (the default) tells New to create
	// an instance of Morfeusz capable both of morphological
	// analysis and generation.
	BothAnalyseAndGenerate Usage = C.BOTH_ANALYSE_AND_GENERATE
	// AnalyseOnly tells New to create an instance of Morfeusz
	// capable only of morphological analysis.
	AnalyseOnly = C.ANALYSE_ONLY
	// GenerateOnly tells New to create an instance of Morfeusz
	// capable only of morphological generation.
	GenerateOnly = C.GENERATE_ONLY
)

// Config informs New about the parameters
// of the instance of Morfeusz to be created.
type Config struct {
	DictName           string
	Aggl               string
	Praet              string
	Charset            Charset
	TokenNumbering     TokenNumbering
	CaseHandling       CaseHandling
	WhitespaceHandling WhitespaceHandling
	Usage              Usage
}

var (
	errInvalidCharset = errors.New("Invalid charset")
	errInvalidUsage   = errors.New("Invalid usage option")
)

// New returns a fresh instance of Morfeusz. New(nil), equivalent
// to New(&Config{}), creates the instance with default parameters.
func New(c *Config) (*Morfeusz, error) {
	if c == nil {
		c = &Config{}
	}
	// Workaround for the lack of usage checks.
	if c.Usage < BothAnalyseAndGenerate || c.Usage > GenerateOnly {
		return nil, errInvalidUsage
	}
	m := C.createInstance(
		C.makeStructString(c.DictName), C.enum_Usage(c.Usage))
	if m == nil {
		return nil, fmt.Errorf(
			"Failed to load dictionary \"%s\"", c.DictName)
	}
	// Make sure that the associated C++ object m
	// will be freed when ret is garbage-collected.
	ret := gcMorfeusz(m)
	if c.Aggl != "" {
		if err := ret.SetAggl(c.Aggl); err != nil {
			return nil, err
		}
	}
	if c.Praet != "" {
		if err := ret.SetPraet(c.Praet); err != nil {
			return nil, err
		}
	}
	if err := ret.SetCharset(c.Charset); err != nil {
		return nil, err
	}
	if err := ret.SetCaseHandling(c.CaseHandling); err != nil {
		return nil, err
	}
	if err := ret.SetTokenNumbering(c.TokenNumbering); err != nil {
		return nil, err
	}
	if err := ret.SetWhitespaceHandling(c.WhitespaceHandling); err != nil {
		return nil, err
	}
	return ret, nil
}

// Analyse returns the result of morphological analysis
// of a byte slice. Use the Next and TokenInfo functions
// of the result to get the interpretation of the tokens.
func (m Morfeusz) Analyse(text []byte) *Result {
	return m.AnalyseString(string(text))
}

// AnalyseString returns the result of morphological
// analysis of a string. Use the Next and TokenInfo
// functions of the result to get the interpretation
// of the tokens.
func (m Morfeusz) AnalyseString(text string) *Result {
	r := C.analyseString(m.morf, C.makeStructString(text))
	if r == nil {
		return nil
	}
	// Make sure that the associated C++ object r
	// will be freed when the returned *Result
	// is garbage-collected.
	return gcResult(r)
}

// Next returns true when there is more information
// available in the result of the analysis. It does
// not modify the internals of the result.
func (r Result) Next() bool {
	return C.hasNext(r.res) != 0
}

// TokenInfo returns the next *TokenInfo, or nil if the analysis is done.
// It modifies the internals of the Result so that the next call will return
// another piece of information.
func (r Result) TokenInfo() *TokenInfo {
	t := C.next(r.res)
	if t.orth.n == 0 {
		return nil
	}
	// Make sure that the associated C++ object t and its
	// character arrays will be freed when the returned
	// *TokenInfo is garbage-collected.
	return gcTokenInfo(t)
}

// StartNode returns the index of the node where a token starts.
func (t *TokenInfo) StartNode() int {
	return int(t.info.startNode)
}

// EndNode returns the index of the node where a token ends.
func (t *TokenInfo) EndNode() int {
	return int(t.info.endNode)
}

// Orth returns the spelling of a token.
func (t *TokenInfo) Orth() string {
	return goString(t.info.orth)
}

// Lemma returns the lemma of a token.
func (t *TokenInfo) Lemma() string {
	return goString(t.info.lemma)
}

// IsIgn returns true only when a token is an unknown word.
func (t *TokenInfo) IsIgn() bool {
	return t.info.tagID == 0
}

// IsWhitespace returns true when a token represents whitespace.
func (t *TokenInfo) IsWhitespace() bool {
	return t.info.tagID == 1
}

// Tag returns the tag for a token.
func (t *TokenInfo) Tag(morf *Morfeusz) string {
	return morf.Tag(int(t.info.tagID))
}

// Name returns the named entity for a token.
func (t *TokenInfo) Name(morf *Morfeusz) string {
	return morf.Name(int(t.info.nameID))
}

// LabelsAsAtring returns the string form of the labels for a token.
func (t *TokenInfo) LabelsAsString(morf *Morfeusz) string {
	return morf.LabelsAsString(int(t.info.labelsID))
}

// Labels returns the slice-of-strings form of labels for a token.
func (t *TokenInfo) Labels(morf *Morfeusz) []string {
	return morf.Labels(int(t.info.labelsID))
}

// TagestID returns the current tagset ID, as specified
// in the first line of the tagset file.
func (m Morfeusz) TagsetID() string {
	return goStringFree(C.tagsetId(m.morf))
}

// Tag returns the inflectional tag for a given ID,
// or an empty string when the ID is invalid.
func (m Morfeusz) Tag(tagID int) string {
	return goStringFree(C.tag(m.morf, C.int(tagID)))
}

// TagID returns the ID for a given inflectional tag,
// or -1 when the tag is invalid.
func (m Morfeusz) TagID(tag string) int {
	return int(C.tagId(m.morf, C.makeStructString(tag)))
}

// Name returns the named entity for a given ID,
// or an empty string when the ID is invalid.
func (m Morfeusz) Name(nameID int) string {
	return goStringFree(C.name(m.morf, C.int(nameID)))
}

// NameID returns the ID for a given named entity,
// or -1 when the name is invalid.
func (m Morfeusz) NameID(name string) int {
	return int(C.nameId(m.morf, C.makeStructString(name)))
}

// LabelsAsString returns the string form of the labels for a given ID,
// or an empty string when the ID is invalid.
func (m Morfeusz) LabelsAsString(labelsID int) string {
	return goStringFree(C.labelsAsString(m.morf, C.int(labelsID)))
}

// Labels returns the slice-of-strings form of labels for a given ID,
// or an empty slice when the ID is invalid.
func (m Morfeusz) Labels(labelsID int) []string {
	return fromStringArray(C.labels(m.morf, C.int(labelsID)))
}

// LabelsID returns the ID for given labels,
// or -1 when the labels are invalid.
func (m Morfeusz) LabelsID(labels string) int {
	return int(C.labelsId(m.morf, C.makeStructString(labels)))
}

// TagsCount returns the number of tags in the current dictionary.
func (m Morfeusz) TagsCount() int {
	return int(C.tagsCount(m.morf))
}

// NamesCount returns the number of named entity types
// in the current dictionary.
func (m Morfeusz) NamesCount() int {
	return int(C.namesCount(m.morf))
}

// LabelsCount returns the number of different labels
// in the current dictionary.
func (m Morfeusz) LabelsCount() int {
	return int(C.labelsCount(m.morf))
}

// Generate returns a list of all inflected forms for a given lemma.
func (m Morfeusz) Generate(lemma string) ([]*TokenInfo, error) {
	return fromTokenInfoArray(C.generate(
		m.morf, C.makeStructString(lemma)))
}

// GenerateWithTagID returns a list of inflected forms for a given lemma
// that have a specific inflectional tag.
func (m Morfeusz) GenerateWithTagID(
	tagID int, lemma string) ([]*TokenInfo, error) {
	return fromTokenInfoArray(C.generateWithTagID(
		m.morf, C.int(tagID), C.makeStructString(lemma)))
}

// DictID returns the ID of the current dictionary.
func (m Morfeusz) DictID() string {
	return goStringFree(C.dictId(m.morf))
}

// DictCopyright returns the copyright text of the current dictionary.
func (m Morfeusz) DictCopyright() string {
	return goStringFree(C.dictCopyright(m.morf))
}

// SetAggl sets the kind of agglutination rules.
func (m Morfeusz) SetAggl(aggl string) error {
	return newError(C.setAggl(m.morf, C.makeStructString(aggl)))
}

// SetPraet sets the kind of past tense segmentation.
func (m Morfeusz) SetPraet(praet string) error {
	return newError(C.setPraet(m.morf, C.makeStructString(praet)))
}

// SetCharset sets the input and output charset.
func (m Morfeusz) SetCharset(encoding Charset) error {
	// Workaround for the lack of charset checks.
	if encoding < UTF8 || encoding > CP852 {
		return errInvalidCharset
	}
	return newError(C.setCharset(m.morf, C.enum_Charset(encoding)))
}

// SetCaseHandling sets the kind of case handling.
func (m Morfeusz) SetCaseHandling(caseHandling CaseHandling) error {
	return newError(C.setCaseHandling(
		m.morf, C.enum_CaseHandling(caseHandling)))
}

// SetTokenNumbering sets the kind of token numbering.
func (m Morfeusz) SetTokenNumbering(numbering TokenNumbering) error {
	return newError(C.setTokenNumbering(
		m.morf, C.enum_TokenNumbering(numbering)))
}

// SetWhitespaceHandling sets the kind of whitespace handling.
func (m Morfeusz) SetWhitespaceHandling(handling WhitespaceHandling) error {
	return newError(C.setWhitespaceHandling(
		m.morf, C.enum_WhitespaceHandling(handling)))
}

// SetDictionary sets current dictionary to the one with given name.
func (m Morfeusz) SetDictionary(dictName string) error {
	return newError(C.setDictionary(m.morf, C.makeStructString(dictName)))
}

// SetDebug turns debugging output on and off.
func (m Morfeusz) SetDebug(debug bool) {
	intDebug := C.int(0)
	if debug {
		intDebug = 1
	}
	C.setDebug(m.morf, intDebug)
}

// Charset returns the current input and output charset.
func (m Morfeusz) Charset() Charset {
	return Charset(C.charset(m.morf))
}

// Aggl returns the current kind of agglutination rules.
func (m Morfeusz) Aggl() string {
	return goStringFree(C.aggl(m.morf))
}

// Praet returns the current kind of past tense segmentation.
func (m Morfeusz) Praet() string {
	return goStringFree(C.praet(m.morf))
}

// CaseHandling returns the current kind of case handling.
func (m Morfeusz) CaseHandling() CaseHandling {
	return CaseHandling(C.caseHandling(m.morf))
}

// TokenNumbering returns the current kind of token numbering.
func (m Morfeusz) TokenNumbering() TokenNumbering {
	return TokenNumbering(C.tokenNumbering(m.morf))
}

// WhitespaceHandling returns the current kind of whitespace handling.
func (m Morfeusz) WhitespaceHandling() WhitespaceHandling {
	return WhitespaceHandling(C.whitespaceHandling(m.morf))
}

// AvailableAgglOptions returns the allowed values for the argument
// of SetAgglOptions.
func (m Morfeusz) AvailableAgglOptions() []string {
	return fromStringArray(C.availableAgglOptions(m.morf))
}

// AvailablePraetOptions returns the allowed values for the argument
// of SetPraetOptions.
func (m Morfeusz) AvailablePraetOptions() []string {
	return fromStringArray(C.availablePraetOptions(m.morf))
}

// DictionarySearchPaths returns the paths where the Morfeusz instance
// looks for dictionaries.
func (m Morfeusz) DictionarySearchPaths() []string {
	return fromStringArray(C.dictionarySearchPaths(m.morf))
}

// PrependToDictionarySearchPaths inserts path at the beginning of the list.
func (m Morfeusz) PrependToDictionarySearchPaths(path string) {
	C.prependToDictionarySearchPaths(m.morf, C.makeStructString(path))
}

// AppendToDictionarySearchPaths adds path at the end of the list.
func (m Morfeusz) AppendToDictionarySearchPaths(path string) {
	C.appendToDictionarySearchPaths(m.morf, C.makeStructString(path))
}

// RemoveFromDictionarySearchPaths removes from dictionary search paths
// elements equal to path. It returns the number of removed elements.
func (m Morfeusz) RemoveFromDictionarySearchPaths(path string) int {
	return int(C.removeFromDictionarySearchPaths(
		m.morf, C.makeStructString(path)))
}

// ClearDictionarySearchPaths removes all dictionary search paths.
func (m Morfeusz) ClearDictionarySearchPaths() {
	C.clearDictionarySearchPaths(m.morf)
}

// Clone copies an instance of Morfeusz. Beware: as of Morfeusz 1.9.16,
// the copy and the original share the charset, token numbering, case
// handling, whitespace handling, and dictionary search paths.
func (m Morfeusz) Clone() *Morfeusz {
	// Make sure that the associated C++ object will be freed
	// when the returned *Morfeusz is garbage-collected.
	return gcMorfeusz(C.cloneMorf(m.morf))
}

// Version returns the version of the underlying Morfeusz 2 library.
func Version() string {
	return goStringFree(C.version())
}

// DefaultDictName returns the default dictionary name.
func DefaultDictName() string {
	return goStringFree(C.defaultDictName())
}

// Copyright returns the copyright text of the underlying Morfeusz 2 library.
func Copyright() string {
	return goStringFree(C.copyright())
}

func gcMorfeusz(m C.Morf) *Morfeusz {
	ret := &Morfeusz{m}
	runtime.SetFinalizer(ret, freeMorfeusz)
	return ret
}

func gcResult(r C.Res) *Result {
	ret := &Result{r}
	runtime.SetFinalizer(ret, freeResult)
	return ret
}

func gcTokenInfo(t C.struct_TokenInfo) *TokenInfo {
	ret := &TokenInfo{t}
	runtime.SetFinalizer(ret, freeTokenInfo)
	return ret
}

func freeMorfeusz(m *Morfeusz) {
	C.freeMorf(m.morf)
}

func freeResult(r *Result) {
	C.freeRes(r.res)
}

func freeTokenInfo(t *TokenInfo) {
	C.freeTokenInfo(&t.info)
}

func fromStringArray(arr C.struct_StringArray) []string {
	sliceView := (*[1 << 28]C.struct_String)(
		unsafe.Pointer(arr.strings))[:arr.length:arr.length]
	ret := make([]string, 0, arr.length)
	for _, s := range sliceView {
		ret = append(ret, goStringFree(s))
	}
	C.freeStringArray(&arr)
	return ret
}

func fromTokenInfoArray(arr C.struct_TokenInfoArray) ([]*TokenInfo, error) {
	if arr.error.p != nil {
		return nil, newError(arr.error)
	}
	sliceView := (*[1 << 28]C.struct_TokenInfo)(
		unsafe.Pointer(arr.tokens))[:arr.length:arr.length]
	ret := make([]*TokenInfo, 0, arr.length)
	for _, t := range sliceView {
		ret = append(ret, gcTokenInfo(t))
	}
	C.freeTokenInfoArray(&arr)
	return ret, nil
}

func newError(s C.struct_String) error {
	if s.n == 0 {
		return nil
	}
	return errors.New(goStringFree(s))
}

func goStringFree(s C.struct_String) string {
	ret := goString(s)
	C.freeCharArray(s.p)
	return ret
}

func goString(s C.struct_String) string {
	return C.GoStringN(s.p, s.n)
}
