package morfeusz_test

import (
	"fmt"
	"testing"

	"github.com/go-morfeusz/morfeusz"
)

type tokenInfo struct {
	start        int
	end          int
	orth         string
	lemma        string
	isIgn        bool
	isWhitespace bool
	tag          string
	name         string
	labels       string
}

func ExampleMorfeusz() {
	m, _ := morfeusz.New(nil)
	r := m.AnalyseString("Ala ma kota.")
	for r.Next() {
		t := r.TokenInfo()
		fmt.Println(
			t.StartNode(), t.EndNode(), t.Orth(), t.Lemma(),
			t.Tag(m), t.Name(m), t.LabelsAsString(m), ";")
	}
	// Output:
	// 0 1 Ala Ala subst:sg:nom:f imię  ;
	// 0 1 Ala Al subst:sg:gen.acc:m1 imię  ;
	// 0 1 Ala Alo subst:sg:gen.acc:m1 imię  ;
	// 1 2 ma mieć fin:sg:ter:imperf   ;
	// 1 2 ma mój:A adj:sg:nom.voc:f:pos   ;
	// 2 3 kota kota subst:sg:nom:f nazwa_pospolita  ;
	// 2 3 kota kot:Sm1 subst:sg:gen.acc:m1 nazwa_pospolita pot.,środ. ;
	// 2 3 kota kot:Sm2 subst:sg:gen.acc:m2 nazwa_pospolita  ;
	// 3 4 . . interp   ;
}

func TestNew(t *testing.T) {
	tests := []struct {
		conf morfeusz.Config
		give string
	}{
		{morfeusz.Config{DictName: "xyz"}, "DictName"},
		{morfeusz.Config{Aggl: "xyz"}, "Aggl"},
		{morfeusz.Config{Praet: "xyz"}, "Praet"},
		{morfeusz.Config{Charset: 4}, "Charset"},
		{morfeusz.Config{TokenNumbering: 2}, "TokenNumbering"},
		{morfeusz.Config{CaseHandling: 3}, "CaseHandling"},
		{morfeusz.Config{WhitespaceHandling: 3}, "WhitespaceHandling"},
		{morfeusz.Config{Usage: 3}, "Usage"},
	}
	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			_, err := morfeusz.New(&tt.conf)
			assertError(t, err)
		})
	}
	t.Run("DefaultConfig", func(t *testing.T) {
		_, err := morfeusz.New(&morfeusz.Config{})
		assertNoError(t, err)
	})
}

func TestMorfeuszMethods(t *testing.T) {
	m, _ := morfeusz.New(nil)
	validIntTests := []struct {
		got  int
		give string
	}{
		{m.TagID("subst:sg:nom:f"), "TagID"},
		{m.NameID("imię"), "NameID"},
		{m.LabelsID("pot."), "LabelsID"},
		{m.TagsCount(), "TagsCount"},
		{m.NamesCount(), "NamesCount"},
		{m.LabelsCount(), "LabelsCount"},
	}
	for _, tt := range validIntTests {
		t.Run(tt.give, func(t *testing.T) {
			if tt.got < 0 {
				t.Errorf("got %d; want >= 0", tt.got)
			}
		})
	}
	invalidIntTests := []struct {
		got  int
		give string
	}{
		{m.TagID("xyz"), "TagID"},
		{m.NameID("xyz"), "NameID"},
		{m.LabelsID("xyz"), "LabelsID"},
	}
	for _, tt := range invalidIntTests {
		t.Run(tt.give, func(t *testing.T) {
			assertEqualInt(t, tt.got, -1)
		})
	}

	validStringTests := []struct {
		got  string
		give string
	}{
		{m.TagsetID(), "TagsetID"},
		{m.Tag(611), "Tag"},
		{m.Name(12), "Name"},
		{m.LabelsAsString(466), "LabelsAsString"},
		{m.DictID(), "DictID"},
		{m.DictCopyright(), "DictCopyright"},
	}
	for _, tt := range validStringTests {
		t.Run(tt.give, func(t *testing.T) {
			assertNotEqualString(t, tt.got, "")
		})
	}
	invalidStringTests := []struct {
		got  string
		give string
	}{
		{m.Tag(-1), "Tag"},
		{m.Name(-1), "Name"},
		{m.LabelsAsString(-1), "LabelsAsString"},
	}
	for _, tt := range invalidStringTests {
		t.Run(tt.give, func(t *testing.T) {
			assertEqualString(t, tt.got, "")
		})
	}

	t.Run("Labels", func(t *testing.T) {
		assertNonEmpty(t, len(m.Labels(466)))
	})
}

func TestMorfeuszSettersAndGetters(t *testing.T) {
	m, _ := morfeusz.New(nil)
	setCharset := func(x int) error {
		return m.SetCharset(morfeusz.Charset(x))
	}
	setTokenNumbering := func(x int) error {
		return m.SetTokenNumbering(morfeusz.TokenNumbering(x))
	}
	setCaseHandling := func(x int) error {
		return m.SetCaseHandling(morfeusz.CaseHandling(x))
	}
	setWhitespaceHandling := func(x int) error {
		return m.SetWhitespaceHandling(morfeusz.WhitespaceHandling(x))
	}
	charset := func() int { return int(m.Charset()) }
	tokenNumbering := func() int { return int(m.TokenNumbering()) }
	caseHandling := func() int { return int(m.CaseHandling()) }
	whitespaceHandling := func() int { return int(m.WhitespaceHandling()) }

	validIntSetterTests := []struct {
		set  func(int) error
		get  func() int
		want int
		give string
	}{
		{setCharset, charset, morfeusz.CP852, "SetCharset"},
		{setTokenNumbering, tokenNumbering,
			morfeusz.ContinuousNumbering, "SetTokenNumbering"},
		{setCaseHandling, caseHandling,
			morfeusz.IgnoreCase, "SetCaseHandling"},
		{setWhitespaceHandling, whitespaceHandling,
			morfeusz.KeepWhitespaces, "SetWhitespaceHandling"},
	}
	for _, tt := range validIntSetterTests {
		t.Run(tt.give, func(t *testing.T) {
			assertNoError(t, tt.set(tt.want))
			assertEqualInt(t, tt.get(), tt.want)
			assertError(t, tt.set(42))
			assertEqualInt(t, tt.get(), tt.want)
		})
	}

	stringSetterTests := []struct {
		set  func(string) error
		get  func() string
		want string
		give string
	}{
		{m.SetAggl, m.Aggl, "permissive", "Aggl"},
		{m.SetPraet, m.Praet, "composite", "Praet"},
	}
	for _, tt := range stringSetterTests {
		t.Run(tt.give, func(t *testing.T) {
			assertNoError(t, tt.set(tt.want))
			assertEqualString(t, tt.get(), tt.want)
			assertError(t, tt.set("xyz"))
			assertEqualString(t, tt.get(), tt.want)
		})
	}

	stringSliceGetterTests := []struct {
		get  func() []string
		give string
	}{
		{m.AvailableAgglOptions, "AvailableAgglOptions"},
		{m.AvailablePraetOptions, "AvailablePraetOptions"},
	}
	for _, tt := range stringSliceGetterTests {
		t.Run(tt.give, func(t *testing.T) {
			assertNonEmpty(t, len(tt.get()))
		})
	}
}

func TestAnalyse(t *testing.T) {
	m, _ := morfeusz.New(nil)
	m.SetWhitespaceHandling(morfeusz.KeepWhitespaces)
	got := analyseToTokenInfoSlice(t, m, "bez xyz")
	want := []tokenInfo{
		{0, 1, "bez", "bez:P", false, false,
			"prep:gen:nwok", "", ""},
		{0, 1, "bez", "bez:S", false, false,
			"subst:sg:nom.acc:m3", "nazwa_pospolita", "bot."},
		{0, 1, "bez", "beza", false, false,
			"subst:pl:gen:f", "nazwa_pospolita", ""},
		{1, 2, " ", " ", false, true, "sp", "", ""},
		{2, 3, "xyz", "xyz", true, false, "ign", "", ""},
	}
	assertEqualTokenInfoSlices(t, got, want)
}

func TestGenerate(t *testing.T) {
	m, _ := morfeusz.New(nil)
	np := "nazwa_pospolita"
	bot := "bot."
	generateTests := []struct {
		want []tokenInfo
		give string
	}{
		{[]tokenInfo{
			{0, 1, "bez", "bez:S", false, false,
				"subst:sg:nom.acc:m3", np, bot},
			{0, 1, "bzu", "bez:S", false, false,
				"subst:sg:gen:m3", np, bot},
			{0, 1, "bzowi", "bez:S", false, false,
				"subst:sg:dat:m3", np, bot},
			{0, 1, "bzem", "bez:S", false, false,
				"subst:sg:inst:m3", np, bot},
			{0, 1, "bzie", "bez:S", false, false,
				"subst:sg:loc:m3", np, bot},
			{0, 1, "bzie", "bez:S", false, false,
				"subst:sg:voc:m3", np, bot},
			{0, 1, "bzy", "bez:S", false, false,
				"subst:pl:nom.acc.voc:m3", np, bot},
			{0, 1, "bzów", "bez:S", false, false,
				"subst:pl:gen:m3", np, bot},
			{0, 1, "bzom", "bez:S", false, false,
				"subst:pl:dat:m3", np, bot},
			{0, 1, "bzami", "bez:S", false, false,
				"subst:pl:inst:m3", np, bot},
			{0, 1, "bzach", "bez:S", false, false,
				"subst:pl:loc:m3", np, bot},
			{0, 1, "beze", "bez:P", false, false,
				"prep:gen:wok", "", ""},
			{0, 1, "bez", "bez:P", false, false,
				"prep:gen:nwok", "", ""},
			{0, 1, "b", "bez", false, false, "brev:pun", "", ""},
		}, "bez"},
		{[]tokenInfo{
			{0, 1, "beze", "bez:P", false, false,
				"prep:gen:wok", "", ""},
			{0, 1, "bez", "bez:P", false, false,
				"prep:gen:nwok", "", ""},
		}, "bez:P"},
		{[]tokenInfo{
			{0, 1, "xyz", "xyz", true, false, "ign", "", ""},
		}, "xyz"},
	}
	for _, tt := range generateTests {
		t.Run(tt.give, func(t *testing.T) {
			got := generateToTokenInfoSlice(t, m, tt.give)
			assertEqualTokenInfoSlices(t, got, tt.want)
		})
	}

	generateWithTagIDTests := []struct {
		want []tokenInfo
		tag  string
		give string
	}{
		{[]tokenInfo{
			{0, 1, "bzu", "bez:S", false, false,
				"subst:sg:gen:m3", np, bot},
		}, "subst:sg:gen:m3", "bez:S"},
		{[]tokenInfo{
			{0, 1, "bzu", "bez:S", false, false,
				"subst:sg:gen:m3", np, bot},
		}, "subst:sg:gen:m3", "bez"},
		{[]tokenInfo{
			{0, 1, "xyz", "xyz", true, false, "ign", "", ""},
		}, "ign", "xyz"},
	}
	for _, tt := range generateWithTagIDTests {
		t.Run(tt.give, func(t *testing.T) {
			got := generateToTokenInfoSliceWithTag(
				t, m, tt.tag, tt.give)
			assertEqualTokenInfoSlices(t, got, tt.want)
		})
	}
	generateWithTagIDEmptyTests := []struct {
		tag  string
		give string
	}{
		{"subst:sg:gen:m1", "bez"},
		{"", "xyz"},
	}
	for _, tt := range generateWithTagIDEmptyTests {
		t.Run(tt.give, func(t *testing.T) {
			got := generateToTokenInfoSliceWithTag(
				t, m, tt.tag, tt.give)
			assertEmpty(t, len(got))
		})
	}
}

func TestUsage(t *testing.T) {
	ma, _ := morfeusz.New(&morfeusz.Config{Usage: morfeusz.AnalyseOnly})
	_, err := ma.Generate("dom")
	assertError(t, err)

	mg, _ := morfeusz.New(&morfeusz.Config{Usage: morfeusz.GenerateOnly})
	if mg.AnalyseString("dom") != nil {
		t.Error("got Analyse() != nil; want Analyse() == nil")
	}
}

func TestDictionarySearchPaths(t *testing.T) {
	m, _ := morfeusz.New(nil)
	paths := m.DictionarySearchPaths()

	m.PrependToDictionarySearchPaths("first_path")
	want := append([]string{"first_path"}, paths...)
	assertEqualStringSlices(t, m.DictionarySearchPaths(), want)

	m.PrependToDictionarySearchPaths("first_path")
	want = append([]string{"first_path"}, want...)
	assertEqualStringSlices(t, m.DictionarySearchPaths(), want)

	m.AppendToDictionarySearchPaths("last_path")
	want = append(want, "last_path")
	assertEqualStringSlices(t, m.DictionarySearchPaths(), want)

	assertEqualInt(t, m.RemoveFromDictionarySearchPaths("first_path"), 2)
	want = want[2:]
	assertEqualStringSlices(t, m.DictionarySearchPaths(), want)

	assertEqualInt(t, m.RemoveFromDictionarySearchPaths("xyz"), 0)
	assertEqualStringSlices(t, m.DictionarySearchPaths(), want)

	m.ClearDictionarySearchPaths()
	assertEmpty(t, len(m.DictionarySearchPaths()))
}

func TestClone(t *testing.T) {
	m, _ := morfeusz.New(nil)
	c := m.Clone()

	aWant := analyseToTokenInfoSlice(t, m, "dom")
	aGot := analyseToTokenInfoSlice(t, c, "dom")
	assertEqualTokenInfoSlices(t, aGot, aWant)

	gWant := generateToTokenInfoSlice(t, m, "dom")
	gGot := generateToTokenInfoSlice(t, c, "dom")
	assertEqualTokenInfoSlices(t, gGot, gWant)

	tWant := generateToTokenInfoSliceWithTag(
		t, m, "subst:sg:dat:m3", "dom")
	tGot := generateToTokenInfoSliceWithTag(
		t, c, "subst:sg:dat:m3", "dom")
	assertEqualTokenInfoSlices(t, tGot, tWant)
}

func expandTokenInfo(
	t *morfeusz.TokenInfo, m *morfeusz.Morfeusz) tokenInfo {
	// Check against double freeing of the underlying C.struct_String.
	t.Orth()
	t.Lemma()
	return tokenInfo{
		t.StartNode(), t.EndNode(),
		t.Orth(), t.Lemma(),
		t.IsIgn(), t.IsWhitespace(),
		t.Tag(m), t.Name(m), t.LabelsAsString(m),
	}
}

func analyseToTokenInfoSlice(
	t *testing.T, m *morfeusz.Morfeusz, text string) []tokenInfo {
	r := m.AnalyseString(text)
	var ret []tokenInfo
	for r.Next() {
		ret = append(ret, expandTokenInfo(r.TokenInfo(), m))
	}
	// Check that nil is returned when there is no more information.
	if r.TokenInfo() != nil {
		t.Error("got TokenInfo() != nil; want TokenInfo() == nil")
	}
	return ret
}

func generateToTokenInfoSlice(
	t *testing.T, m *morfeusz.Morfeusz, lemma string) []tokenInfo {
	ts, err := m.Generate(lemma)
	assertNoError(t, err)
	return makeTokenInfoSlice(ts, m)
}

func generateToTokenInfoSliceWithTag(
	t *testing.T, m *morfeusz.Morfeusz, tag, lemma string) []tokenInfo {
	tagID := m.TagID(tag)
	assertNotEqualInt(t, tagID, -1)
	ts, err := m.GenerateWithTagID(tagID, lemma)
	assertNoError(t, err)
	return makeTokenInfoSlice(ts, m)
}

func makeTokenInfoSlice(
	ts []*morfeusz.TokenInfo, m *morfeusz.Morfeusz) []tokenInfo {
	var ret []tokenInfo
	for _, t := range ts {
		ret = append(ret, expandTokenInfo(t, m))
	}
	return ret
}

func makeMultiset(tis []tokenInfo) map[tokenInfo]int {
	ret := map[tokenInfo]int{}
	for _, ti := range tis {
		ret[ti] += 1
	}
	return ret
}

func assertEqualTokenInfoSlices(t *testing.T, got, want []tokenInfo) {
	assertNonEmpty(t, len(got))
	assertEqualInt(t, len(got), len(want))
	gotSet := makeMultiset(got)
	wantSet := makeMultiset(want)
	for g, gn := range gotSet {
		wn := wantSet[g]
		if gn != wn {
			t.Errorf("got %v %d times; want it %d times", g, gn, wn)
		}
	}
	for w, wn := range wantSet {
		gn := gotSet[w]
		if gn == 0 {
			t.Errorf("got %v %d times; want it %d times", w, 0, wn)
		}
	}
}

func assertEqualStringSlices(t *testing.T, got, want []string) {
	assertEqualInt(t, len(got), len(want))
	for i, g := range got {
		assertEqualString(t, g, want[i])
	}
}

func assertEmpty(t *testing.T, length int) {
	if length != 0 {
		t.Errorf("got len() = %d; want len() == 0", length)
	}
}

func assertNonEmpty(t *testing.T, length int) {
	if length == 0 {
		t.Error("got len() = 0; want len() != 0")
	}
}

func assertEqualInt(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got == %d; want == %d", got, want)
	}
}

func assertNotEqualInt(t *testing.T, got, dontWant int) {
	if got == dontWant {
		t.Errorf("got == %d; want != %d", got, dontWant)
	}
}

func assertEqualString(t *testing.T, got, want string) {
	if got != want {
		t.Errorf(`got == "%s"; want == "%s"`, got, want)
	}
}

func assertNotEqualString(t *testing.T, got, dontWant string) {
	if got == dontWant {
		t.Errorf(`got == "%s"; want != "%s"`, got, dontWant)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf(`got err == "%s"; want err == nil`, err)
	}
}

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Error("got err == nil; want err != nil")
	}
}
